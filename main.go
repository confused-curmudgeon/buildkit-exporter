package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	defaultBuildkitAddr = "unix:///run/buildkit/buildkitd.sock"
	defaultScrapeAddr   = ":9220"
	namespace           = "buildkit"
)

var (
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()

	buildkitAddr = kingpin.Flag(
		"buildkit.address",
		"Address to use for connecting to Buildkit",
	).Default(defaultBuildkitAddr).String()

	tlsInsecureSkipVerify = kingpin.Flag(
		"tls.insecure-skip-verify",
		"Ignore certificate and server verification when using a tls connection.",
	).Bool()

	toolkitFlags = webflag.AddFlags(kingpin.CommandLine, defaultScrapeAddr)
)

func init() {
	prometheus.MustRegister(version.NewCollector("buildkit_exporter"))
}

func newHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("buildkit_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting buildkit_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	fmt.Println("Scrape Addr:", toolkitFlags.WebListenAddresses)
	fmt.Println("Buildkit Addr:", *buildkitAddr)

	ctx := context.Background()
	client := NewClient(ctx, buildkitAddr)
	exporter := NewExporter(ctx, client, false, 10*time.Second, logger)

	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Buildkit Exporter</title></head>
             <body>
             <h1>Buildkit Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

}
