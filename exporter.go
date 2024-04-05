package main

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	buildkitInfo = prometheus.NewDesc(prometheus.BuildFQName(namespace, "version", "info"), "Buildkit version info.", []string{"release_date", "version"}, nil)
	buildkitUp   = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of Buildkit successful.", nil, nil)
)

// Exporter collects Buildkit stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	mutex        sync.RWMutex
	fetchStat    func() (io.ReadCloser, error)
	up           prometheus.Gauge
	totalScrapes prometheus.Counter
	logger       log.Logger
}

// NewExporter returns an initialized Exporter.
func NewExporter(ctx context.Context, sslVerify bool, timeout time.Duration, logger log.Logger) (*Exporter, error) {

	var fetchStat func() (io.ReadCloser, error)

	return &Exporter{
		fetchStat: fetchStat,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of Buildkit successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total Buildkit scrapes.",
		}),
		logger: logger,
	}, nil
}

// Describe describes all the metrics ever exported by the Buildkit exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- buildkitInfo
	ch <- buildkitUp
	ch <- e.totalScrapes.Desc()
}

// Collect fetches the stats from configured Buildkit location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	// up := e.scrape(ch)
	up := e.scrape()

	ch <- prometheus.MustNewConstMetric(buildkitUp, prometheus.GaugeValue, up)
	ch <- e.totalScrapes
}

// func (e *Exporter) scrape(ch chan<- prometheus.Metric) (up float64) {
func (e *Exporter) scrape() (up float64) {
	e.totalScrapes.Inc()
	var err error

	body, err := e.fetchStat()
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't scrape Buildkit", "err", err)
		return 0
	}
	defer body.Close()

	return 1
}
