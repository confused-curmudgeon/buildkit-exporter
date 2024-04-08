// Copyright 2024 Cody Boggs
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

package main

import (
	"context"
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
	mutex           sync.RWMutex
	client          *Client
	up              prometheus.Gauge
	totalScrapes    prometheus.Counter
	logger          log.Logger
	buildkitMetrics []metricInfo
}

// Verify if Exporter implements prometheus.Collector
var _ prometheus.Collector = (*Exporter)(nil)

// NewExporter returns an initialized Exporter.
func NewExporter(ctx context.Context, client *Client, sslVerify bool, timeout time.Duration, logger log.Logger) *Exporter {
	return &Exporter{
		client: client,
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
		logger:          logger,
		buildkitMetrics: buildkitMetrics,
	}
}

// Describe describes all the metrics ever exported by the Buildkit exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- buildkitInfo
	ch <- buildkitUp
	ch <- e.totalScrapes.Desc()
	for _, m := range e.buildkitMetrics {
		ch <- m.Desc
	}
}

// Collect fetches the stats from configured Buildkit location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	up := e.scrape(ch)

	ch <- prometheus.MustNewConstMetric(buildkitUp, prometheus.GaugeValue, up)
	ch <- e.totalScrapes
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) (up float64) {
	e.totalScrapes.Inc()
	for _, m := range e.buildkitMetrics {
		if err := m.fetch(e.client, ch, m.Desc, m.Type); err != nil {
			level.Error(e.logger).Log("msg", "Error fetching metric", "name", *m.Name, "err", err)
		}
	}
	return 1
}
