package main

import (
	"context"
	"fmt"
	"os"

	buildkitclient "github.com/moby/buildkit/client"
)

type Client struct {
	*buildkitclient.Client
}

func NewClient(ctx context.Context, buildkitAddr *string) *Client {
	bkc, err := buildkitclient.New(ctx, *buildkitAddr)
	if err != nil {
		fmt.Println("ERROR Creating buildkit client:", err)
		os.Exit(1)
	}

	return &Client{bkc}
}

func (c *Client) TotalDiskUsageBytes(ctx context.Context) {
	var total int64
	duInfo, err := c.DiskUsage(ctx)
	if err != nil {
		fmt.Println("ERROR Getting DiskUsage:", err)
		os.Exit(1)
	}

	for _, d := range duInfo {
		total += d.Size
	}
	fmt.Println("TotalDiskUsageBytes", total)
}
