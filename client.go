package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	controlapi "github.com/moby/buildkit/api/services/control"
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

func (c *Client) getAllHistories() ([]*controlapi.BuildHistoryEvent, error) {
	events := []*controlapi.BuildHistoryEvent{}

	resp, err := c.ControlClient().ListenBuildHistory(c.ctx, &controlapi.BuildHistoryRequest{EarlyExit: true})
	if err != nil {
		return nil, err
	}

	for {
		ev, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		events = append(events, ev)
	}

	return events, nil
}
