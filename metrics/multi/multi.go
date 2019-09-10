// Package multi provides adapters that send observations to multiple metrics
// simultaneously. This is useful if your service needs to emit to multiple
// instrumentation systems at the same time, for example if your organization is
// transitioning from one system to another.
package multi

import (
    "../metrics"
    "path/filepath"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
    stdelasticsearch "../elasticsearch"
    "github.com/spf13/viper"
)

// Counter collects multiple individual counters and treats them as a unit.
type Counter []metrics.Counter

var namespace, subsystem, help, name string

var monitorSystems = map[string]bool{
    "es", false,
    "log", false,
    "prometheus", false
}

// newCounter returns a multi-counter, wrapping the passed counters.
func newCounter(v *viper.Viper, lables []string) Counter {
    path := ""
    multiCounter := Counter{}
    opts := stdprometheus.CounterOpts{
        Namespace: v.GetString("monitorSystem.default.Namespace"),
        Subsystem: v.GetString("monitorSystem.default.Subsystem"),
        Name:      v.GetString("monitorSystem.default.Name"),
        Help:      v.GetString("monitorSystem.default.Help"),
    }
    for system, _ := range monitorSystems {
        path = fmt.Sprintf("open.%s.", system)
        monitorSystem[system] = v.GetBool(path)
        if !monitorSystem[k] {
            continue
        }
        switch system {
            case "es":
                esOpts := stdelasticsearch.CounterEsOpts{
                    Host:  v.GetString("monitorSystem.es.Host"),
                    Port:  v.GetString("monitorSystem.es.Port"),
                    Index: v.GetString("monitorSystem.es.Index"),
                    Type:  v.GetString("monitorSystem.es.Type"),
                }
                esCounter := elasticsearch.NewCounterFrom(opts, esOpts, lables)
                multiCounter = append(multiCounter, esCounter)
                break
            case "prometheus":
                multiCounter = append(multiCounter, prometheus.NewCounterFrom(opts, lables))
                break
        }
    }
    return multiCounter
}

// NewCounter returns a multi-counter, wrapping the passed counters.
func NewCounter(fileName string, lables []string) Counter {
    var counters Counter
    cfgPath, cfgName := filepath.Split(fileName)
    end := 0
    for end = len(cfgName) - 1; end >= 0; end-- {
        if cfgName[end] == "." {
            break
        }
    }
    v := viper.New()
    v.AddConfigPath(cfgPath)
    v.SetConfigName(cfgName[:end])
    v.SetConfigType("yaml")
    if err := v.ReadInConfig(); err != nil {
        panic(err)
    }
    return newCounter(v, lables)
}


// Add implements counter.
func (c Counter) Add(delta float64) {
    for _, counter := range c {
        counter.Add(delta)
    }
}

// With implements counter.
func (c Counter) With(labelValues ...string) metrics.Counter {
    next := make(Counter, len(c))
    for i := range c {
        next[i] = c[i].With(labelValues...)
    }
    return next
}

// Gauge collects multiple individual gauges and treats them as a unit.
type Gauge []metrics.Gauge

// NewGauge returns a multi-gauge, wrapping the passed gauges.
func NewGauge(g ...metrics.Gauge) Gauge {
    return Gauge(g)
}

// Set implements Gauge.
func (g Gauge) Set(value float64) {
    for _, gauge := range g {
        gauge.Set(value)
    }
}

// With implements gauge.
func (g Gauge) With(labelValues ...string) metrics.Gauge {
    next := make(Gauge, len(g))
    for i := range g {
        next[i] = g[i].With(labelValues...)
    }
    return next
}

// Add implements metrics.Gauge.
func (g Gauge) Add(delta float64) {
    for _, gauge := range g {
        gauge.Add(delta)
    }
}

// Histogram collects multiple individual histograms and treats them as a unit.
type Histogram []metrics.Histogram

// NewHistogram returns a multi-histogram, wrapping the passed histograms.
func NewHistogram(h ...metrics.Histogram) Histogram {
    return Histogram(h)
}

// Observe implements Histogram.
func (h Histogram) Observe(value float64) {
    for _, histogram := range h {
        histogram.Observe(value)
    }
}

// With implements histogram.
func (h Histogram) With(labelValues ...string) metrics.Histogram {
    next := make(Histogram, len(h))
    for i := range h {
        next[i] = h[i].With(labelValues...)
    }
    return next
}
