package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
)

const (
	defaultBuildkitAddr = "unix:///run/buildkit/buildkitd.sock"
	defaultScrapeAddr   = ":9220"
	namespace           = "buildkit"
)

var (
	buildkitAddr string
	metricsPath  = "/metrics"
	scrapeAddr   string

	cmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

var Usage = func() {
	fmt.Fprintf(cmd.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&scrapeAddr, "scrape-addr", defaultScrapeAddr, "Scrape address for exporter")
	flag.StringVar(&buildkitAddr, "buildkit-addr", defaultBuildkitAddr, fmt.Sprintf("Buildkit socket address. Defaults to '%s'. Can be remote TCP URL, ex: tcp://my.buildkit.host:1234", defaultBuildkitAddr))

	flag.Parse()
}

func main() {
	fmt.Println("buildkit-exporter started!")
	fmt.Println("Scrape Addr:", scrapeAddr)
	fmt.Println("Buildkit Addr:", buildkitAddr)

	ctx := context.Background()

	c := NewClient(ctx, &buildkitAddr)
	c.TotalDiskUsageBytes(ctx)

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	exporter, err := NewExporter(ctx, false, 10*time.Second, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating an exporter", "err", err)
		os.Exit(1)
	}

	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("buildkit_exporter"))

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Buildkit Exporter</title></head>
             <body>
             <h1>Buildkit Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{Registry: promRegistry}))
	log.Fatal(http.ListenAndServe(scrapeAddr, nil))

	// srv := &http.Server{}
	// if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
	// 	level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
	// 	os.Exit(1)
	// }
}
