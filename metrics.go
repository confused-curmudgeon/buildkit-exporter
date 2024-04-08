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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	noopFetcher = func(*Client, chan<- prometheus.Metric, *prometheus.Desc, prometheus.ValueType) error { return nil }
	imageFields = []string{"registry", "path", "name", "tag"}

	buildkitMetrics = []metricInfo{
		newBuildkitMetric("build_histories",
			"Count of build histories that have not yet been pruned.",
			prometheus.GaugeValue,
			append(imageFields, "exporter_type"),
			fetchHistoriesCount),

		newBuildkitMetric("cache_objects_size_bytes",
			"Total bytes used on by cache objects.",
			prometheus.GaugeValue,
			[]string{"type"},
			fetchCacheSizeTotalBytes),

		newBuildkitMetric("cache_objects",
			"Count of cache objects that have not yet been pruned.",
			prometheus.GaugeValue,
			[]string{"type"},
			fetchObjectCounts),

		newBuildkitMetric("build_steps",
			"Count of per-image build steps in histories that have not yet been pruned.",
			prometheus.GaugeValue,
			append(imageFields, "count"),
			fetchBuildStepCounts),
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
