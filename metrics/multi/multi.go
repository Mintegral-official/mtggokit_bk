// Package multi provides adapters that send observations to multiple metrics
// simultaneously. This is useful if your service needs to emit to multiple
// instrumentation systems at the same time, for example if your organization is
// transitioning from one system to another.
package multi

import (
	"fmt"
	"strings"
    "../../metrics"
    "../elasticsearch"
    "../prometheus"
	"path/filepath"
    "github.com/spf13/viper"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
    stdelasticsearch "github.com/Schneizelw/prometheus/client_golang/prometheus"
)

// Counter collects multiple individual counters and treats them as a unit.
type Counter []metrics.Counter

var monitorSystems = map[string]bool{
    "es": false,
    "log": false,
    "prometheus": false,
}

func newEsCounter(v *viper.Viper, lables []string) metrics.Counter {
    baseOpts := stdelasticsearch.CounterOpts{
        Namespace: v.GetString("monitorSystem.default.Namespace"),
        Subsystem: v.GetString("monitorSystem.default.Subsystem"),
        Name:      v.GetString("monitorSystem.default.Name"),
        Help:      v.GetString("monitorSystem.default.Help"),
    }
    esOpts := stdelasticsearch.CounterEsOpts{
        Host:  v.GetString("monitorSystem.es.Host"),
        Port:  v.GetString("monitorSystem.es.Port"),
        EsIndex: v.GetString("monitorSystem.es.Index"),
        EsType:  v.GetString("monitorSystem.es.Type"),
    }
    return elasticsearch.NewCounterFrom(baseOpts, esOpts, lables)
}

func newPrometheusCounter(v *viper.Viper, lables []string) metrics.Counter {
    baseOpts := stdprometheus.CounterOpts{
        Namespace: v.GetString("monitorSystem.default.Namespace"),
        Subsystem: v.GetString("monitorSystem.default.Subsystem"),
        Name:      v.GetString("monitorSystem.default.Name"),
        Help:      v.GetString("monitorSystem.default.Help"),
    }
    return prometheus.NewCounterFrom(baseOpts, lables)
}

// newCounter returns a multi-counter, wrapping the passed counters.
func newCounter(v *viper.Viper, lables []string) Counter {
    path := ""
    multiCounter := Counter{}
    for system, _ := range monitorSystems {
        path = fmt.Sprintf("open.%s.", system)
        monitorSystems[system] = v.GetBool(path)
        if !monitorSystems[system] {
            continue
        }
        switch system {
            case "es":
                esCounter := newEsCounter(v, lables)
                multiCounter = append(multiCounter, esCounter)
                break
            case "prometheus":
                prometheusCounter := newPrometheusCounter(v, lables)
                multiCounter = append(multiCounter, prometheusCounter)
                break
        }
    }
    return multiCounter
}

// NewCounter returns a multi-counter, wrapping the passed counters.
func NewCounter(fileName string, lables []string) Counter {
    configPath, configName := filepath.Split(fileName)
	dotIndex := strings.LastIndex(configName, ".")
	if dotIndex == -1 || configName[dotIndex:] != ".yaml" {
        panic("config file format must be yaml")
	}
    v := viper.New()
    v.AddConfigPath(configPath)
    v.SetConfigName(configName[:dotIndex])
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
