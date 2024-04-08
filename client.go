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
