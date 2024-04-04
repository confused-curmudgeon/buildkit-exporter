package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	bkclient "github.com/moby/buildkit/client"
)

const (
	defaultBuildkitAddr = "unix:///run/buildkit/buildkitd.sock"
	defaultScrapePort   = 9220
)

var (
	buildkitAddr string
	scrapePort   int

	cmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

var Usage = func() {
	fmt.Fprintf(cmd.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.IntVar(&scrapePort, "scrape-port", defaultScrapePort, "Scrape port for exporter")
	flag.StringVar(&buildkitAddr, "buildkit-addr", defaultBuildkitAddr, fmt.Sprintf("Buildkit socket address. Defaults to '%s'. Can be remote TCP URL, ex: tcp://my.buildkit.host:1234", defaultBuildkitAddr))

	flag.Parse()
}

func main() {
	fmt.Println("buildkit-exporter started!")
	fmt.Println("Scrape Port:", scrapePort)
	fmt.Println("Buildkit Addr:", buildkitAddr)

	ctx := context.Background()
	client, err := bkclient.New(ctx, buildkitAddr)
	if err != nil {
		fmt.Println("ERROR Creating buildkit client:", err)
		os.Exit(1)
	}

	duInfo, err := client.DiskUsage(ctx)
	if err != nil {
		fmt.Println("ERROR Getting DiskUsage:", err)
		os.Exit(1)
	}

	for _, d := range duInfo {
		fmt.Println("Usage::ID:", d.ID, "Usage::Size:", d.Size)
	}
}
