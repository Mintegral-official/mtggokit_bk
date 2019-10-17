// Package metricslog provides Prometheus implementations for metrics.
// Individual metrics are mapped to their Prometheus counterparts, and
// (depending on the constructor used) may be automatically registered in the
// global Prometheus metrics registry.
package metricslog

import (
    "github.com/Mintegral-official/mtggokit/metrics/metrics"
    "github.com/Mintegral-official/mtggokit/metrics/internal/lv"

    "github.com/Mintegral-official/mtggokit/metrics/metricslog/client_golang/metricslog"
)

// Counter implements Counter, via a Prometheus CounterVec.
type Counter struct {
    cv  *metricslog.CounterVec
    lvs lv.LabelValues
}

// NewCounterFrom constructs and registers a Prometheus CounterVec,
// and returns a usable Counter object.
func NewCounterFrom(opts metricslog.CounterOpts, logOpts metricslog.CounterLogOpts, labelNames []string) *Counter {
    cv := metricslog.NewCounterVec(opts, logOpts, labelNames)
    metricslog.MustRegister(cv)
    return NewCounter(cv)
}

// NewCounter wraps the CounterVec and returns a usable Counter object.
func NewCounter(cv *metricslog.CounterVec) *Counter {
    return &Counter{
        cv: cv,
    }
}

// With implements Counter.
func (c *Counter) With(labelValues ...string) metrics.Counter {
    return &Counter{
        cv:  c.cv,
        lvs: c.lvs.With(labelValues...),
    }
}

// Add implements Counter.
func (c *Counter) Add(delta float64) {
    c.cv.With(makeLabels(c.lvs...)).Add(delta)
}

// Gauge implements Gauge, via a Prometheus GaugeVec.
type Gauge struct {
    gv  *metricslog.GaugeVec
    lvs lv.LabelValues
}

// NewGaugeFrom construts and registers a Prometheus GaugeVec,
// and returns a usable Gauge object.
func NewGaugeFrom(opts metricslog.GaugeOpts, logOpts metricslog.GaugeLogOpts, labelNames []string) *Gauge {
    gv := metricslog.NewGaugeVec(opts, logOpts, labelNames)
    metricslog.MustRegister(gv)
    return NewGauge(gv)
}

// NewGauge wraps the GaugeVec and returns a usable Gauge object.
func NewGauge(gv *metricslog.GaugeVec) *Gauge {
    return &Gauge{
        gv: gv,
    }
}

// With implements Gauge.
func (g *Gauge) With(labelValues ...string) metrics.Gauge {
    return &Gauge{
        gv:  g.gv,
        lvs: g.lvs.With(labelValues...),
    }
}

// Set implements Gauge.
func (g *Gauge) Set(value float64) {
    g.gv.With(makeLabels(g.lvs...)).Set(value)
}

// Add is supported by Prometheus GaugeVecs.
func (g *Gauge) Add(delta float64) {
    g.gv.With(makeLabels(g.lvs...)).Add(delta)
}

// Summary implements Histogram, via a Prometheus SummaryVec. The difference
// between a Summary and a Histogram is that Summaries don't require predefined
// quantile buckets, but cannot be statistically aggregated.
type Summary struct {
    sv  *metricslog.SummaryVec
    lvs lv.LabelValues
}

// NewSummaryFrom constructs and registers a Prometheus SummaryVec,
// and returns a usable Summary object.
func NewSummaryFrom(opts metricslog.SummaryOpts, logOpts metricslog.SummaryLogOpts, labelNames []string) *Summary {
    sv := metricslog.NewSummaryVec(opts, logOpts, labelNames)
    metricslog.MustRegister(sv)
    return NewSummary(sv)
}

// NewSummary wraps the SummaryVec and returns a usable Summary object.
func NewSummary(sv *metricslog.SummaryVec) *Summary {
    return &Summary{
        sv: sv,
    }
}

// With implements Histogram.
func (s *Summary) With(labelValues ...string) metrics.Histogram {
    return &Summary{
        sv:  s.sv,
        lvs: s.lvs.With(labelValues...),
    }
}

// Observe implements Histogram.
func (s *Summary) Observe(value float64) {
    s.sv.With(makeLabels(s.lvs...)).Observe(value)
}

// Histogram implements Histogram via a Prometheus HistogramVec. The difference
// between a Histogram and a Summary is that Histograms require predefined
// quantile buckets, and can be statistically aggregated.
type Histogram struct {
    hv  *metricslog.HistogramVec
    lvs lv.LabelValues
}

// NewHistogramFrom constructs and registers a Prometheus HistogramVec,
// and returns a usable Histogram object.
func NewHistogramFrom(opts metricslog.HistogramOpts, logOpts metricslog.HistogramLogOpts, labelNames []string) *Histogram {
    hv := metricslog.NewHistogramVec(opts, logOpts, labelNames)
    metricslog.MustRegister(hv)
    return NewHistogram(hv)
}

// NewHistogram wraps the HistogramVec and returns a usable Histogram object.
func NewHistogram(hv *metricslog.HistogramVec) *Histogram {
    return &Histogram{
        hv: hv,
    }
}

// With implements Histogram.
func (h *Histogram) With(labelValues ...string) metrics.Histogram {
    return &Histogram{
        hv:  h.hv,
        lvs: h.lvs.With(labelValues...),
    }
}

// Observe implements Histogram.
func (h *Histogram) Observe(value float64) {
    h.hv.With(makeLabels(h.lvs...)).Observe(value)
}

func makeLabels(labelValues ...string) metricslog.Labels {
    labels := metricslog.Labels{}
    for i := 0; i < len(labelValues); i += 2 {
        labels[labelValues[i]] = labelValues[i+1]
    }
    return labels
}
