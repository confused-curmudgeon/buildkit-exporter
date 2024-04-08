package main

import (
	"fmt"
	"strings"

	buildkitclient "github.com/moby/buildkit/client"
	"github.com/prometheus/client_golang/prometheus"
)

func fetchCacheSizeTotalBytes(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error {
	sizes := make(map[buildkitclient.UsageRecordType]int64)

	duInfo, err := client.DiskUsage(client.ctx)
	if err != nil {
		return err
	}

	for _, d := range duInfo {
		sizes[d.RecordType] += d.Size
	}

	for key, val := range sizes {
		ch <- prometheus.MustNewConstMetric(desc, valType, float64(val), string(key))
	}

	return nil
}

func fetchObjectCounts(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error {
	counts := make(map[buildkitclient.UsageRecordType]int)

	duInfo, err := client.DiskUsage(client.ctx)
	if err != nil {
		return err
	}

	for _, d := range duInfo {
		counts[d.RecordType] += 1
	}

	for key, val := range counts {
		ch <- prometheus.MustNewConstMetric(desc, valType, float64(val), string(key))
	}
	return nil
}

func fetchHistoriesCount(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error {
	totals := make(map[string]int)
	sep := ";"

	events, err := client.getAllHistories()
	if err != nil {
		return err
	}

	for _, event := range events {
		var k string
		if event.Record.Exporters != nil {
			k = fmt.Sprintf("%s%s%s",
				event.Record.Exporters[0].Type,
				sep,
				event.Record.Exporters[0].Attrs["name"],
			)
		} else {
			k = "cache-only;undefined"
		}

		totals[k] += 1
	}

	for key, val := range totals {
		sp := strings.Split(key, sep)
		exporterType, image := sp[0], sp[1]
		ch <- prometheus.MustNewConstMetric(desc, valType, float64(val), exporterType, image)
	}
	return nil
}

func fetchBuildStepCounts(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error {
	totals := make(map[string]int32)
	sep := ";"

	events, err := client.getAllHistories()
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
	}

	for key, val := range totals {
		sp := strings.Split(key, sep)
		exporterType, imageFQN := sp[0], sp[1]
		img, err := splitImageFQN(imageFQN)
		if err != nil {

		}
		ch <- prometheus.MustNewConstMetric(desc, valType, float64(val), exporterType, image)
	}

	return nil
}
