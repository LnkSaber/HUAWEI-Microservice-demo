package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var registries = make(map[string]NewRegistry)
var prometheusRegistry = prometheus.NewRegistry()

//NewRegistry create a registry
type NewRegistry func(opts Options) Registry

//Registry holds all of metrics collectors
//name is a unique ID for different type of metrics
type Registry interface {
	CreateGauge(opts GaugeOpts) error
	CreateCounter(opts CounterOpts) error
	CreateSummary(opts SummaryOpts) error
	GaugeSet(name string, val float64, labels map[string]string) error
	CounterAdd(name string, val float64, labels map[string]string) error
	SummaryObserve(name string, val float64, Labels map[string]string) error
}

var defaultRegistry Registry

//CreateGauge init a new gauge
func CreateGauge(opts GaugeOpts) error {
	return defaultRegistry.CreateGauge(opts)
}

//CreateCounter init a new counter
func CreateCounter(opts CounterOpts) error {
	return defaultRegistry.CreateCounter(opts)
}

//CreateSummary init a new summary
func CreateSummary(opts SummaryOpts) error {
	return defaultRegistry.CreateSummary(opts)
}

//GaugeSet set a new value to a collector
func GaugeSet(name string, val float64, labels map[string]string) error {
	return defaultRegistry.GaugeSet(name, val, labels)
}

//CounterAdd increase value of a collector
func CounterAdd(name string, val float64, labels map[string]string) error {
	return defaultRegistry.CounterAdd(name, val, labels)
}

//SummaryObserve gives a value to summary collector
func SummaryObserve(name string, val float64, labels map[string]string) error {
	return defaultRegistry.SummaryObserve(name, val, labels)
}

//CounterOpts is options to create a counter options
type CounterOpts struct {
	Name   string
	Help   string
	Labels []string
}

//GaugeOpts is options to create a gauge collector
type GaugeOpts struct {
	Name   string
	Help   string
	Labels []string
}

//SummaryOpts is options to create summary collector
type SummaryOpts struct {
	Name       string
	Help       string
	Labels     []string
	Objectives map[float64]float64
}

//Options control config
type Options struct {
	FlushInterval time.Duration
}

//InstallPlugin install metrics registry
func InstallPlugin(name string, f NewRegistry) {
	registries[name] = f
}

//Init load the metrics plugin and initialize it
func Init() error {
	var name string
	name = "prometheus"
	f, ok := registries[name]
	if !ok {
		return fmt.Errorf("can not init metrics registry [%s]", name)
	}
	defaultRegistry = f(Options{
		FlushInterval: 10 * time.Second,
	})
	return nil
}

//GetSystemPrometheusRegistry return prometheus registry which go chassis use
func GetSystemPrometheusRegistry() *prometheus.Registry {
	return prometheusRegistry
}
