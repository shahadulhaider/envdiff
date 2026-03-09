package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// isCI returns true if running in a CI environment.
func isCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" ||
		os.Getenv("JENKINS_URL") != ""
}

// isGitHubActions returns true if running in GitHub Actions.
func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") != ""
}

func newCICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Run checks in CI mode",
		Long:  "Run envdiff checks suitable for CI pipelines. Auto-detects GitHub Actions format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			requires, _ := cmd.Flags().GetStringArray("require")
			if len(requires) == 0 {
				fmt.Fprintf(os.Stderr, "error: --require flag is required (specify .env.example or .env.schema.toml)\n")
				return &exitError{code: 2}
			}
			return runCI(cmd, requires)
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringArray("require", nil, "file to check against (.env.example or .env.schema.toml); can be specified multiple times")
	return cmd
}

func runCI(cmd *cobra.Command, requires []string) error {
	// Force CI-appropriate settings
	// Override color to never in CI
	if isCI() {
		// Set color to never by overriding the flag value
		if err := cmd.Root().Flags().Set("color", "never"); err != nil {
			// Persistent flag — try parent
			_ = err
		}
		// Try persistent flags
		if f := cmd.Root().PersistentFlags().Lookup("color"); f != nil {
			f.Value.Set("never")
		}
	}

	// Auto-select github format when in GitHub Actions and format not explicitly set
	if isGitHubActions() {
		if f := cmd.Root().PersistentFlags().Lookup("format"); f != nil && !f.Changed {
			f.Value.Set("github")
		}
	}

	var totalErrors, totalWarnings int
	anyFailure := false

	for _, req := range requires {
		var err error
		if strings.HasSuffix(req, ".toml") {
			// Schema validation
			err = runValidate(cmd, req, []string{".env"})
		} else {
			// Example check
			err = runCheck(cmd, ".env", req)
		}

		if err != nil {
			if exitErr, ok := err.(*exitError); ok {
				if exitErr.code == 1 {
					anyFailure = true
					totalErrors++
				} else if exitErr.code == 2 {
					return err // propagate hard errors
				}
			} else {
				return err
			}
		}
	}

	// Print summary
	fmt.Printf("envdiff: %d errors, %d warnings\n", totalErrors, totalWarnings)

	if anyFailure {
		return &exitError{code: 1}
	}
	return nil
}
