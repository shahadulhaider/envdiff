package main

import (
	"fmt"
	"os"

	"github.com/shahadulhaider/envdiff/cmd/envdiff/root"
)

var version = "dev"

func main() {
	rootCmd := root.NewRootCmd(version)
	rootCmd.AddCommand(
		newDiffCmd(),
		newCheckCmd(),
		newCompareCmd(),
		newValidateCmd(),
		newInitCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		if !handleExitError(err) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
