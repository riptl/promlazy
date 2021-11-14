// Package promlazy is a replacement for the promauto package.
//
// Metrics defined using this package won't be registered until written to once.
// This is useful to avoid exposing empty metrics about dead code paths and disabled modules.
package promlazy

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Factory creates a batch of "lazy" Prometheus metrics that delay registration until first use.
//
// This package will always panic if registration fails.
type Factory struct {
	r          prometheus.Registerer
	initOnce   sync.Once
	collectors []prometheus.Collector
}

// New creates a factory of lazy metrics targeting prometheus.DefaultRegisterer.
func New() Factory {
	return With(prometheus.DefaultRegisterer)
}

// With creates a factory of lazy metrics that eventually register.
func With(r prometheus.Registerer) Factory { return Factory{r: r} }

// Register imports all collectors into the registry.
//
// This method is idempotent. Subsequent invocations after the first successful call do nothing.
// However, the first call can panic the program on metric naming conflicts, similarly to promauto.
// It is recommended to call Register on the earliest indication a particular module to be monitored is going to be used,
// to panic fast in the event of unsuccessful registration.
//
// For example, if you are trying to define metrics on a network client, it makes sense to Register the factory
// whenever a client is instantiated.
func (f *Factory) Register() {
	f.initOnce.Do(f.init)
}

func (f *Factory) init() {
	f.r.MustRegister(f.collectors...)
}

// NewCounter works like the function of the same name in the prometheus package,
// but it automatically registers the Counter on first use.
func (f *Factory) NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	c := prometheus.NewCounter(opts)
	f.collectors = append(f.collectors, c)
	return lazyCounter{c, f}
}

type lazyCounter struct {
	prometheus.Counter
	*Factory
}

func (l lazyCounter) Inc() {
	l.Factory.Register()
	l.Counter.Inc()
}

func (l lazyCounter) Add(x float64) {
	l.Factory.Register()
	l.Counter.Add(x)
}

// NewGauge works like the function of the same name in the prometheus package,
// but it automatically registers the Gauge on first use.
func (f *Factory) NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	c := prometheus.NewGauge(opts)
	f.collectors = append(f.collectors, c)
	return lazyGauge{c, f}
}

type lazyGauge struct {
	prometheus.Gauge
	*Factory
}

func (l lazyGauge) Set(x float64) {
	l.Factory.Register()
	l.Gauge.Set(x)
}

func (l lazyGauge) Inc() {
	l.Factory.Register()
	l.Gauge.Inc()
}

func (l lazyGauge) Dec() {
	l.Factory.Register()
	l.Gauge.Dec()
}

func (l lazyGauge) Add(x float64) {
	l.Factory.Register()
	l.Gauge.Add(x)
}

func (l lazyGauge) Sub(x float64) {
	l.Factory.Register()
	l.Gauge.Sub(x)
}

func (l lazyGauge) SetToCurrentTime() {
	l.Factory.Register()
	l.Gauge.SetToCurrentTime()
}

// NewSummary works like the function of the same name in the prometheus package,
// but it automatically registers the Summary on first use.
func (f *Factory) NewSummary(opts prometheus.SummaryOpts) prometheus.Summary {
	c := prometheus.NewSummary(opts)
	f.collectors = append(f.collectors, c)
	return lazySummary{c, f}
}

type lazySummary struct {
	prometheus.Summary
	*Factory
}

func (l lazySummary) Observe(x float64) {
	l.Factory.Register()
	l.Summary.Observe(x)
}

// NewHistogram works like the function of the same name in the prometheus package,
// but it automatically registers the Histogram on first use.
func (f *Factory) NewHistogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	c := prometheus.NewHistogram(opts)
	f.collectors = append(f.collectors, c)
	return lazyHistogram{c, f}
}

type lazyHistogram struct {
	prometheus.Histogram
	*Factory
}

func (l lazyHistogram) Observe(x float64) {
	l.Factory.Register()
	l.Histogram.Observe(x)
}
