package main

import (
	"context"
	"fmt"
	"os"

	buildkitclient "github.com/moby/buildkit/client"
)

type Client struct {
	ctx context.Context
	*buildkitclient.Client
}

func NewClient(ctx context.Context, buildkitAddr *string) *Client {
	bkc, err := buildkitclient.New(ctx, *buildkitAddr)
	if err != nil {
		fmt.Println("ERROR Creating buildkit client:", err)
		os.Exit(1)
	}

	return &Client{ctx, bkc}
}
