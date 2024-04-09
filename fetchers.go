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
	"errors"
	"fmt"
	"strings"

	buildkitclient "github.com/moby/buildkit/client"
	"github.com/prometheus/client_golang/prometheus"
)

type fetcher func(e *Exporter, c prometheus.Collector) error

var (
	sep = ";"

	fetchers = map[string]fetcher{
		gaugeCacheObjectSizeBytes: fetchGaugeCacheObjectSizeBytes,
		gaugeCacheObjects:         fetchGaugeCacheObjectCounts,
		gaugeBuildHistories:       fetchGaugeBuildHistories,
		gaugeBuildSteps:           fetchGaugeBuildSteps,
		histogramBuildDuration:    fetchHistogramBuildDuration,
	}
)

func fetchGaugeCacheObjectSizeBytes(e *Exporter, c prometheus.Collector) error {
	metric := c.(*prometheus.GaugeVec)
	sizes := make(map[buildkitclient.UsageRecordType]int64)

	duInfo, err := e.client.DiskUsage(e.ctx)
	if err != nil {
		return err
	}

	for _, d := range duInfo {
		sizes[d.RecordType] += d.Size
	}

	for key, val := range sizes {
		metric.WithLabelValues(string(key)).Set(float64(val))
	}

	return nil
}

func fetchGaugeCacheObjectCounts(e *Exporter, c prometheus.Collector) error {
	metric := c.(*prometheus.GaugeVec)
	counts := make(map[buildkitclient.UsageRecordType]int)

	duInfo, err := e.client.DiskUsage(e.ctx)
	if err != nil {
		return err
	}

	for _, d := range duInfo {
		counts[d.RecordType] += 1
	}

	for key, val := range counts {
		metric.WithLabelValues(string(key)).Set(float64(val))

	}
	return nil
}

func fetchGaugeBuildHistories(e *Exporter, c prometheus.Collector) error {
	metric := c.(*prometheus.GaugeVec)
	var collectedErrors error
	totals := make(map[string]int)

	events, err := e.client.getAllHistories()
	if err != nil {
		return err
	}

	for _, event := range events {
		var image string
		if event.Record.Exporters != nil {
			image = fmt.Sprintf("%s%s%s",
				event.Record.Exporters[0].Type,
				sep,
				event.Record.Exporters[0].Attrs["name"],
			)
		} else {
			image = "cache-only;undefined"
		}

		totals[image] += 1
	}

	for key, val := range totals {
		sp := strings.Split(key, sep)
		exporterType, imageFQN := sp[0], sp[1]

		img, err := splitImageFQN(imageFQN)
		if err != nil {
			collectedErrors = errors.Join(collectedErrors, err)
		}

		labelValues := append(img.Values(), exporterType)
		metric.WithLabelValues(labelValues...).Set(float64(val))
	}
	return nil
}

func fetchGaugeBuildSteps(e *Exporter, c prometheus.Collector) error {
	metric := c.(*prometheus.GaugeVec)
	var collectedErrors error
	totals := make(map[string]int32)

	events, err := e.client.getAllHistories()
	if err != nil {
		return err
	}

	for _, event := range events {
		if event.Record.Exporters == nil {
			continue
		}
		total := event.Record.GetNumTotalSteps()
		completed := event.Record.GetNumCompletedSteps()
		cached := event.Record.GetNumCachedSteps()

		failed := total - completed
		uncached := total - failed - cached

		image := event.Record.Exporters[0].Attrs["name"]
		totals[fmt.Sprintf("%s%s%s", "cached", sep, image)] += cached
		totals[fmt.Sprintf("%s%s%s", "completed", sep, image)] += completed
		totals[fmt.Sprintf("%s%s%s", "total", sep, image)] += total
		totals[fmt.Sprintf("%s%s%s", "uncached", sep, image)] += uncached
		totals[fmt.Sprintf("%s%s%s", "failed", sep, image)] += failed
	}

	for key, val := range totals {
		sp := strings.Split(key, sep)
		exporterType, imageFQN := sp[0], sp[1]

		img, err := splitImageFQN(imageFQN)
		if err != nil {
			collectedErrors = errors.Join(collectedErrors, err)
		}

		labelValues := append(img.Values(), exporterType)
		metric.WithLabelValues(labelValues...).Set(float64(val))
	}

	return collectedErrors
}

func fetchHistogramBuildDuration(e *Exporter, c prometheus.Collector) error {
	metric := c.(*prometheus.HistogramVec)
	var collectedErrors error

	events, err := e.client.getAllHistories()
	if err != nil {
		return err
	}

	for _, event := range events {
		exporterType := "cache-only"
		img := &image{}
		status := "success"
		if event.Record.GetError() != nil {
			status = "failed"
		}

		elapsed := event.Record.CompletedAt.Sub(*event.Record.CreatedAt).Seconds()

		if event.Record.Exporters != nil {
			img, err = splitImageFQN(event.Record.Exporters[0].Attrs["name"])
			if err != nil {
				collectedErrors = errors.Join(collectedErrors, err)
				continue
			}

			exporterType = event.Record.Exporters[0].Type
		} else {
			img = &image{name: "undefined"}
		}

		labelValues := append(img.Values(), exporterType, status)
		metric.WithLabelValues(labelValues...).Observe(elapsed)
	}

	return collectedErrors
}
