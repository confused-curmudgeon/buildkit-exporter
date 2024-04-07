package main

import "github.com/prometheus/client_golang/prometheus"

func fetchCacheSizeTotalBytes(client *Client, ch chan<- prometheus.Metric, desc *prometheus.Desc, valType prometheus.ValueType) error {
	var total int64

	duInfo, err := client.DiskUsage(client.ctx)
	if err != nil {
		return err
	}

	for _, d := range duInfo {
		total += d.Size
	}

	ch <- prometheus.MustNewConstMetric(desc, valType, float64(total))
	return nil
}
