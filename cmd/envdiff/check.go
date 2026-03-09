package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check .env is in sync with .env.example",
		Long:  "Auto-detects .env and .env.example (or .env.sample, .env.template) in current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString("source")
			example, _ := cmd.Flags().GetString("example")
			return runCheck(cmd, src, example)
		},
		SilenceUsage: true,
	}
	cmd.Flags().String("source", ".env", "source .env file")
	cmd.Flags().String("example", "", "example .env file (auto-detected if empty)")
	return cmd
}

func runCheck(cmd *cobra.Command, source, example string) error {
	if example == "" {
		// Auto-detect example file
		candidates := []string{".env.example", ".env.sample", ".env.template"}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				example = c
				break
			}
		}
		if example == "" {
			fmt.Fprintf(os.Stderr, "no example file found (looked for .env.example, .env.sample, .env.template)\n")
			return &exitError{code: 2}
		}
	}

	// source must exist
	if _, err := os.Stat(source); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "source file not found: %s\n", source)
		return &exitError{code: 2}
	}

	// Run diff with example vs source: keys in example but not source = missing
	// keys in source but not example = extra
	return runDiff(cmd, example, source)
}
