// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
    "strings"
    "testing"

    "github.com/Schneizelw/mtggokit/metrics/metricslog/client_golang/metricslog"
)

type untypedCollector struct{}

func (u untypedCollector) Describe(c chan<- *metricslog.Desc) {
    c <- metricslog.NewDesc("name", "help", nil, nil)
}

func (u untypedCollector) Collect(c chan<- metricslog.Metric) {
    c <- metricslog.MustNewConstMetric(
        metricslog.NewDesc("name", "help", nil, nil),
        metricslog.UntypedValue,
        2001,
    )
}

func TestToFloat64(t *testing.T) {
    gaugeWithAValueSet := metricslog.NewGauge(metricslog.GaugeOpts{})
    gaugeWithAValueSet.Set(3.14)

    counterVecWithOneElement := metricslog.NewCounterVec(metricslog.CounterOpts{}, []string{"foo"})
    counterVecWithOneElement.WithLabelValues("bar").Inc()

    counterVecWithTwoElements := metricslog.NewCounterVec(metricslog.CounterOpts{}, []string{"foo"})
    counterVecWithTwoElements.WithLabelValues("bar").Add(42)
    counterVecWithTwoElements.WithLabelValues("baz").Inc()

    histogramVecWithOneElement := metricslog.NewHistogramVec(metricslog.HistogramOpts{}, []string{"foo"})
    histogramVecWithOneElement.WithLabelValues("bar").Observe(2.7)

    scenarios := map[string]struct {
        collector metricslog.Collector
        panics    bool
        want      float64
    }{
        "simple counter": {
            collector: metricslog.NewCounter(metricslog.CounterOpts{}),
            panics:    false,
            want:      0,
        },
        "simple gauge": {
            collector: metricslog.NewGauge(metricslog.GaugeOpts{}),
            panics:    false,
            want:      0,
        },
        "simple untyped": {
            collector: untypedCollector{},
            panics:    false,
            want:      2001,
        },
        "simple histogram": {
            collector: metricslog.NewHistogram(metricslog.HistogramOpts{}),
            panics:    true,
        },
        "simple summary": {
            collector: metricslog.NewSummary(metricslog.SummaryOpts{}),
            panics:    true,
        },
        "simple gauge with an actual value set": {
            collector: gaugeWithAValueSet,
            panics:    false,
            want:      3.14,
        },
        "counter vec with zero elements": {
            collector: metricslog.NewCounterVec(metricslog.CounterOpts{}, nil),
            panics:    true,
        },
        "counter vec with one element": {
            collector: counterVecWithOneElement,
            panics:    false,
            want:      1,
        },
        "counter vec with two elements": {
            collector: counterVecWithTwoElements,
            panics:    true,
        },
        "histogram vec with one element": {
            collector: histogramVecWithOneElement,
            panics:    true,
        },
    }

    for n, s := range scenarios {
        t.Run(n, func(t *testing.T) {
            defer func() {
                r := recover()
                if r == nil && s.panics {
                    t.Error("expected panic")
                } else if r != nil && !s.panics {
                    t.Error("unexpected panic: ", r)
                }
                // Any other combination is the expected outcome.
            }()
            if got := ToFloat64(s.collector); got != s.want {
                t.Errorf("want %f, got %f", s.want, got)
            }
        })
    }
}

func TestCollectAndCompare(t *testing.T) {
    const metadata = `
        # HELP some_total A value that represents a counter.
        # TYPE some_total counter
    `

    c := metricslog.NewCounter(metricslog.CounterOpts{
        Name: "some_total",
        Help: "A value that represents a counter.",
        ConstLabels: metricslog.Labels{
            "label1": "value1",
        },
    })
    c.Inc()

    expected := `

        some_total{ label1 = "value1" } 1
    `

    if err := CollectAndCompare(c, strings.NewReader(metadata+expected), "some_total"); err != nil {
        t.Errorf("unexpected collecting result:\n%s", err)
    }
}

func TestCollectAndCompareNoLabel(t *testing.T) {
    const metadata = `
        # HELP some_total A value that represents a counter.
        # TYPE some_total counter
    `

    c := metricslog.NewCounter(metricslog.CounterOpts{
        Name: "some_total",
        Help: "A value that represents a counter.",
    })
    c.Inc()

    expected := `

        some_total 1
    `

    if err := CollectAndCompare(c, strings.NewReader(metadata+expected), "some_total"); err != nil {
        t.Errorf("unexpected collecting result:\n%s", err)
    }
}

func TestCollectAndCompareHistogram(t *testing.T) {
    inputs := []struct {
        name        string
        c           metricslog.Collector
        metadata    string
        expect      string
        observation float64
    }{
        {
            name: "Testing Histogram Collector",
            c: metricslog.NewHistogram(metricslog.HistogramOpts{
                Name:    "some_histogram",
                Help:    "An example of a histogram",
                Buckets: []float64{1, 2, 3},
            }),
            metadata: `
                # HELP some_histogram An example of a histogram
                # TYPE some_histogram histogram
            `,
            expect: `
                some_histogram{le="1"} 0
                some_histogram{le="2"} 0
                some_histogram{le="3"} 1
                    some_histogram_bucket{le="+Inf"} 1
                    some_histogram_sum 2.5
                    some_histogram_count 1

            `,
            observation: 2.5,
        },
        {
            name: "Testing HistogramVec Collector",
            c: metricslog.NewHistogramVec(metricslog.HistogramOpts{
                Name:    "some_histogram",
                Help:    "An example of a histogram",
                Buckets: []float64{1, 2, 3},
            }, []string{"test"}),

            metadata: `
                # HELP some_histogram An example of a histogram
                # TYPE some_histogram histogram
            `,
            expect: `
                        some_histogram_bucket{test="test",le="1"} 0
                        some_histogram_bucket{test="test",le="2"} 0
                        some_histogram_bucket{test="test",le="3"} 1
                        some_histogram_bucket{test="test",le="+Inf"} 1
                        some_histogram_sum{test="test"} 2.5
                    some_histogram_count{test="test"} 1
        
            `,
            observation: 2.5,
        },
    }

    for _, input := range inputs {
        switch collector := input.c.(type) {
        case metricslog.Histogram:
            collector.Observe(input.observation)
        case *metricslog.HistogramVec:
            collector.WithLabelValues("test").Observe(input.observation)
        default:
            t.Fatalf("unsuported collector tested")

        }

        t.Run(input.name, func(t *testing.T) {
            if err := CollectAndCompare(input.c, strings.NewReader(input.metadata+input.expect)); err != nil {
                t.Errorf("unexpected collecting result:\n%s", err)
            }
        })

    }
}

func TestNoMetricFilter(t *testing.T) {
    const metadata = `
        # HELP some_total A value that represents a counter.
        # TYPE some_total counter
    `

    c := metricslog.NewCounter(metricslog.CounterOpts{
        Name: "some_total",
        Help: "A value that represents a counter.",
        ConstLabels: metricslog.Labels{
            "label1": "value1",
        },
    })
    c.Inc()

    expected := `
        some_total{label1="value1"} 1
    `

    if err := CollectAndCompare(c, strings.NewReader(metadata+expected)); err != nil {
        t.Errorf("unexpected collecting result:\n%s", err)
    }
}

func TestMetricNotFound(t *testing.T) {
    const metadata = `
        # HELP some_other_metric A value that represents a counter.
        # TYPE some_other_metric counter
    `

    c := metricslog.NewCounter(metricslog.CounterOpts{
        Name: "some_total",
        Help: "A value that represents a counter.",
        ConstLabels: metricslog.Labels{
            "label1": "value1",
        },
    })
    c.Inc()

    expected := `
        some_other_metric{label1="value1"} 1
    `

    expectedError := `
metric output does not match expectation; want:

# HELP some_other_metric A value that represents a counter.
# TYPE some_other_metric counter
some_other_metric{label1="value1"} 1

got:

# HELP some_total A value that represents a counter.
# TYPE some_total counter
some_total{label1="value1"} 1
`

    err := CollectAndCompare(c, strings.NewReader(metadata+expected))
    if err == nil {
        t.Error("Expected error, got no error.")
    }

    if err.Error() != expectedError {
        t.Errorf("Expected\n%#+v\nGot:\n%#+v", expectedError, err.Error())
    }
}
