package main

import (
	"fmt"
	"os"

	"github.com/shahadulhaider/envdiff/cmd/envdiff/root"
)

var version = "dev"

func main() {
	cmd := root.NewRootCmd(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
