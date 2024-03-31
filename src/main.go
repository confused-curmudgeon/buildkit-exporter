package main

import (
	"context"
	"fmt"
	"os"

	bkclient "github.com/moby/buildkit/client"
)

func main() {
	fmt.Println("buildkit-exporter started!")
	ctx := context.Background()
	client, err := bkclient.New(ctx, "unix:///run/buildkit.sock")
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
