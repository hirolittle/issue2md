package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hirolittle/issue2md/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <github-url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]
	ctx := context.Background()

	if err := cli.Run(ctx, url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
