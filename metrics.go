package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	buildkitMetrics = []metricInfo{
		newBuildkitMetric("cache_size_total_bytes", "Total bytes used on this node's disk by the entire local Buildkit cache.", prometheus.GaugeValue, []string{}),
		newBuildkitMetric("build_histories_current", "Count of build histories that have not yet been pruned.", prometheus.GaugeValue, nil),
		newBuildkitMetric("build_histories_total", "Count of all build histories, pruned or not, seen since exporter startup.", prometheus.CounterValue, nil),
	}
)

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

func initRegistry(e *Exporter) *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		stageDurationHist,
	)

	return reg
}

func newBuildkitMetric(metricName string, docString string, t prometheus.ValueType, labelNames []string) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cache", metricName),
			docString,
			labelNames,
			prometheus.Labels{},
		),
		Type: t,
	}
}
