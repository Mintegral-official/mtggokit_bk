// Package multi provides adapters that send observations to multiple metrics
// simultaneously. This is useful if your service needs to emit to multiple
// instrumentation systems at the same time, for example if your organization is
// transitioning from one system to another.
package multi

import (
	"fmt"
	"strings"
	"path/filepath"
    "github.com/spf13/viper"

    "github.com/Schneizelw/mtggokit/metrics"
    "github.com/Schneizelw/mtggokit/metrics/elasticsearch"
    "github.com/Schneizelw/mtggokit/metrics/prometheus"

    stdprometheus "github.com/prometheus/client_golang/prometheus"
    stdelasticsearch "github.com/Schneizelw/elasticsearch/client_golang/elasticsearch"
)

// Counter collects multiple individual counters and treats them as a unit.
type Counter []metrics.Counter

var monitorSystems = []string{"Elasticsearch", "Log", "Prometheus"}

func newEsCounter(v *viper.Viper, lables []string) metrics.Counter {
    baseOpts := stdelasticsearch.CounterOpts{
        Namespace: v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem: v.GetString("MonitorSystem.Default.Subsystem"),
        Name:      v.GetString("MonitorSystem.Default.Name"),
        Help:      v.GetString("MonitorSystem.Default.Help"),
    }
    esOpts := stdelasticsearch.CounterEsOpts{
        Host:      v.GetString("MonitorSystem.Elasticsearch.Host"),
        Port:      v.GetString("MonitorSystem.Elasticsearch.Port"),
        EsIndex :  v.GetString("MonitorSystem.Elasticsearch.Index"),
        EsType  :  v.GetString("MonitorSystem.Elasticsearch.Type"),
        Interval:  v.GetInt("MonitorSystem.Elasticsearch.Interval"),
    }
    return elasticsearch.NewCounterFrom(baseOpts, esOpts, lables)
}

func newPrometheusCounter(v *viper.Viper, lables []string) metrics.Counter {
    baseOpts := stdprometheus.CounterOpts{
        Namespace: v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem: v.GetString("MonitorSystem.Default.Subsystem"),
        Name:      v.GetString("MonitorSystem.Default.Name"),
        Help:      v.GetString("MonitorSystem.Default.Help"),
    }
    return prometheus.NewCounterFrom(baseOpts, lables)
}

// newCounter returns a multi-counter, wrapping the passed counters.
func newCounter(v *viper.Viper, lables []string) Counter {
    path := ""
	isOpen := false
    multiCounter := Counter{}
    for _, system := range monitorSystems {
        path = fmt.Sprintf("Open.%s", system)
        isOpen = v.GetBool(path)
        if !isOpen {
            continue
        }
        switch system {
            case "Elasticsearch":
                esCounter := newEsCounter(v, lables)
                multiCounter = append(multiCounter, esCounter)
            case "Prometheus":
                prometheusCounter := newPrometheusCounter(v, lables)
                multiCounter = append(multiCounter, prometheusCounter)
        }
    }
    return multiCounter
}

func setViper(fileName string) *viper.Viper {
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
	return v
}

// NewCounter returns a multi-counter, wrapping the passed counters.
func NewCounter(fileName string, lables []string) Counter {
    v := setViper(fileName)
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

func newEsGauge(v *viper.Viper, lables []string) metrics.Gauge {
    baseOpts := stdelasticsearch.GaugeOpts{
        Namespace: v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem: v.GetString("MonitorSystem.Default.Subsystem"),
        Name:      v.GetString("MonitorSystem.Default.Name"),
        Help:      v.GetString("MonitorSystem.Default.Help"),
    }
    esOpts := stdelasticsearch.GaugeEsOpts{
        Host:      v.GetString("MonitorSystem.Elasticsearch.Host"),
        Port:      v.GetString("MonitorSystem.Elasticsearch.Port"),
        EsIndex :  v.GetString("MonitorSystem.Elasticsearch.Index"),
        EsType  :  v.GetString("MonitorSystem.Elasticsearch.Type"),
        Interval:  v.GetInt("MonitorSystem.Elasticsearch.Interval"),
    }
    return elasticsearch.NewGaugeFrom(baseOpts, esOpts, lables)
}

func newPrometheusGauge(v *viper.Viper, lables []string) metrics.Gauge {
    baseOpts := stdprometheus.GaugeOpts{
        Namespace: v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem: v.GetString("MonitorSystem.Default.Subsystem"),
        Name:      v.GetString("MonitorSystem.Default.Name"),
        Help:      v.GetString("MonitorSystem.Default.Help"),
    }
    return prometheus.NewGaugeFrom(baseOpts, lables)
}

func newGauge(v *viper.Viper, lables []string) Gauge {
    path := ""
	isOpen := false
    multiGauge := Gauge{}
    for _, system := range monitorSystems {
        path = fmt.Sprintf("Open.%s", system)
        isOpen = v.GetBool(path)
        if !isOpen {
            continue
        }
        switch system {
            case "Elasticsearch":
                esGauge := newEsGauge(v, lables)
                multiGauge = append(multiGauge, esGauge)
            case "Prometheus":
                prometheusGauge := newPrometheusGauge(v, lables)
                multiGauge = append(multiGauge, prometheusGauge)
        }
    }
    return multiGauge
}

// NewGauge returns a multi-gauge, wrapping the passed gauges.
func NewGauge(fileName string, lables []string) Gauge {
    v := setViper(fileName)
    return newGauge(v, lables)
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

// Summary collects multiple individual summaries and treats them as a unit.
type Summary []metrics.Histogram

func newEsSummary(v *viper.Viper, lables []string) metrics.Histogram {
	objectives := map[float64]float64 {
		0.5 : v.GetFloat64("MonitorSystem.Metrics.Summary.Quantile50"),
		0.9 : v.GetFloat64("MonitorSystem.Metrics.Summary.Quantile90"),
		0.99: v.GetFloat64("MonitorSystem.Metrics.Summary.Quantile99"),
	}
    baseOpts := stdelasticsearch.SummaryOpts{
        Namespace:  v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem:  v.GetString("MonitorSystem.Default.Subsystem"),
        Name:       v.GetString("MonitorSystem.Default.Name"),
        Help:       v.GetString("MonitorSystem.Default.Help"),
		Objectives: objectives,
	}
    esOpts := stdelasticsearch.SummaryEsOpts{
        Host:      v.GetString("MonitorSystem.Elasticsearch.Host"),
        Port:      v.GetString("MonitorSystem.Elasticsearch.Port"),
        EsIndex :  v.GetString("MonitorSystem.Elasticsearch.Index"),
        EsType  :  v.GetString("MonitorSystem.Elasticsearch.Type"),
        Interval:  v.GetInt("MonitorSystem.Elasticsearch.Interval"),
    }
    return elasticsearch.NewSummaryFrom(baseOpts, esOpts, lables)
}

func newPrometheusSummary(v *viper.Viper, lables []string) metrics.Histogram {
    baseOpts := stdprometheus.SummaryOpts{
        Namespace: v.GetString("MonitorSystem.Default.Namespace"),
        Subsystem: v.GetString("MonitorSystem.Default.Subsystem"),
        Name:      v.GetString("MonitorSystem.Default.Name"),
        Help:      v.GetString("MonitorSystem.Default.Help"),
    }
    return prometheus.NewSummaryFrom(baseOpts, lables)
}

func newSummary(v *viper.Viper, lables []string) Summary {
    path := ""
	isOpen := false
    multiSummary := Summary{}
    for _, system := range monitorSystems {
        path = fmt.Sprintf("Open.%s", system)
        isOpen = v.GetBool(path)
        if !isOpen {
            continue
        }
        switch system {
            case "Elasticsearch":
                esSummary := newEsSummary(v, lables)
                multiSummary = append(multiSummary, esSummary)
            case "Prometheus":
                prometheusSummary := newPrometheusSummary(v, lables)
                multiSummary = append(multiSummary, prometheusSummary)
        }
    }
    return multiSummary
}

// NewSummary returns a multi-summary, wrapping the passed summary.
func NewSummary(fileName string, lables []string) Summary {
    v := setViper(fileName)
    return newSummary(v, lables)
}

// Observe implements Histogram.
func (s Summary) Observe(value float64) {
    for _, summary := range s {
        summary.Observe(value)
    }
}

// With implements histogram.
func (s Summary) With(labelValues ...string) metrics.Histogram {
    next := make(Summary, len(s))
    for i := range s {
        next[i] = s[i].With(labelValues...)
    }
    return next
}

