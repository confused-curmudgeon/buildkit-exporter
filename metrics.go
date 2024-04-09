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

const (
	namespace = "buildkit"

	subsystemBuild        = "build"
	subsystemCacheObjects = "cache_objects"

	gaugeBuildHistories       = "histories"
	gaugeBuildSteps           = "steps"
	gaugeCacheObjectSizeBytes = "size_bytes"
	gaugeCacheObjects         = "count"

	histogramBuildDuration = "duration_seconds"
)

var (
	noopFetcher = func(*Client, chan<- prometheus.Metric, *prometheus.Desc, prometheus.ValueType) error { return nil }
	imageFields = []string{"registry", "path", "name", "tag"}

	buildkitMetrics = map[string]prometheus.Collector{
		gaugeBuildHistories: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystemBuild,
				Name:      gaugeBuildHistories,
				Help:      "Count of build histories that have not yet been pruned.",
			},
			append(imageFields, "exporter_type"),
		),
		gaugeBuildSteps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystemBuild,
				Name:      gaugeBuildSteps,
				Help:      "Count of per-image build steps in histories that have not yet been pruned.",
			},
			append(imageFields, "count"),
		),
		gaugeCacheObjectSizeBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystemCacheObjects,
				Name:      gaugeCacheObjectSizeBytes,
				Help:      "Total bytes used by cache objects of each type that have not yet been pruned.",
			},
			[]string{"type"},
		),
		gaugeCacheObjects: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystemCacheObjects,
				Name:      gaugeCacheObjects,
				Help:      "Count of cache objects of each type that have not yet been pruned.",
			},
			[]string{"type"},
		),
		histogramBuildDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystemBuild,
				Name:      histogramBuildDuration,
				Help:      "Time taken to complete an image build.",
				Buckets: []float64{
					30, 60, 90, 120,
					180, 240, 300, 450,
					600, 750, 900, 1050,
					1200,
				},
			},
			append(imageFields, "exporter_type", "status"),
		),
	}
)
