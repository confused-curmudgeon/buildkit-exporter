package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	noopFetcher = func(*Client, chan<- prometheus.Metric, *prometheus.Desc, prometheus.ValueType) error { return nil }

	buildkitMetrics = []metricInfo{
		newBuildkitMetric("build_histories_total",
			"Count of build histories that have not yet been pruned.",
			prometheus.GaugeValue,
			[]string{"exporter_type", "image"},
			fetchHistoriesCount),

		newBuildkitMetric("cache_objects_size_bytes",
			"Total bytes used on by cache objects.",
			prometheus.GaugeValue,
			[]string{"type"},
			fetchCacheSizeTotalBytes),

		newBuildkitMetric("cache_objects_total",
			"Count of cache objects that have not yet been pruned.",
			prometheus.GaugeValue,
			[]string{"type"},
			fetchObjectCounts),
	}
)

type metricFetcher func(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error

type metricInfo struct {
	Desc  *prometheus.Desc
	Name  *string
	Type  prometheus.ValueType
	fetch metricFetcher
}

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
