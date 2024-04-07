package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	noopFetcher = func(*Client, chan<- prometheus.Metric, *prometheus.Desc, prometheus.ValueType) error { return nil }

	buildkitMetrics = []metricInfo{
		newBuildkitMetric("cache_size_total_bytes", "Total bytes used on this node's disk by the entire local Buildkit cache.", prometheus.GaugeValue, []string{}, fetchCacheSizeTotalBytes),
		newBuildkitMetric("build_histories_current", "Count of build histories that have not yet been pruned.", prometheus.GaugeValue, []string{}, noopFetcher),
		newBuildkitMetric("build_histories_total", "Count of all build histories, pruned or not, seen since exporter startup.", prometheus.CounterValue, []string{}, noopFetcher),
	}
)

type metricInfo struct {
	Desc  *prometheus.Desc
	Name  *string
	Type  prometheus.ValueType
	fetch metricFetcher
}

type metricFetcher func(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error

func newBuildkitMetric(metricName string, docString string, t prometheus.ValueType, labelNames []string, fetcher metricFetcher) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cache", metricName),
			docString,
			labelNames,
			prometheus.Labels{},
		),
		Name:  &metricName,
		Type:  t,
		fetch: fetcher,
	}
}
