package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shahadulhaider/envdiff/internal/diff"
	"github.com/shahadulhaider/envdiff/internal/env"
	"github.com/shahadulhaider/envdiff/internal/output"
	"github.com/shahadulhaider/envdiff/internal/parser"
	"github.com/shahadulhaider/envdiff/internal/secret"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <file1> <file2>",
		Short: "Compare two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(cmd, args[0], args[1])
		},
		SilenceUsage: true,
	}
}

func runDiff(cmd *cobra.Command, left, right string) error {
	leftFile, err := parser.ParseFile(left)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", left, err)
		return &exitError{code: 2}
	}

	rightFile, err := parser.ParseFile(right)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", right, err)
		return &exitError{code: 2}
	}

	result := diff.Diff(leftFile, rightFile)

	// Apply --ignore filter
	ignorePattern, _ := cmd.Flags().GetString("ignore")
	if ignorePattern != "" {
		result = filterDiffResult(result, ignorePattern)
	}

	if !result.HasDiffs() {
		return nil // exit 0
	}

	// Build formatter options
	mask, _ := cmd.Flags().GetBool("mask")
	noValues, _ := cmd.Flags().GetBool("no-values")
	colorMode, _ := cmd.Flags().GetString("color")
	formatStr, _ := cmd.Flags().GetString("format")

	// Apply masking to result if requested
	if mask {
		applyMask(result)
	}

	useColor := shouldUseColor(colorMode)
	opts := output.Options{Mask: false, NoValues: noValues, Color: useColor}

	fmtType := parseFormatType(formatStr)
	formatter := output.NewFormatter(fmtType, opts)

	if err := formatter.Format(result, os.Stdout); err != nil {
		return &exitError{code: 2}
	}

	return &exitError{code: 1} // diffs exist → exit 1
}

// filterDiffResult removes entries whose keys match the ignore glob pattern.
func filterDiffResult(result *env.DiffResult, pattern string) *env.DiffResult {
	filtered := &env.DiffResult{Left: result.Left, Right: result.Right}
	for _, e := range result.Entries {
		matched, _ := filepath.Match(pattern, e.Key)
		if !matched {
			filtered.Entries = append(filtered.Entries, e)
		}
	}
	return filtered
}

// applyMask replaces all values in a diff result with "****".
func applyMask(result *env.DiffResult) {
	for i := range result.Entries {
		if result.Entries[i].Left != nil {
			v := secret.MaskValue(result.Entries[i].Left.Value)
			e := *result.Entries[i].Left
			e.Value = v
			result.Entries[i].Left = &e
		}
		if result.Entries[i].Right != nil {
			v := secret.MaskValue(result.Entries[i].Right.Value)
			e := *result.Entries[i].Right
			e.Value = v
			result.Entries[i].Right = &e
		}
	}
}

func shouldUseColor(colorMode string) bool {
	// Never use color in CI environments (unless explicitly forced)
	if colorMode != "always" {
		if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" ||
			os.Getenv("GITLAB_CI") != "" || os.Getenv("JENKINS_URL") != "" {
			return false
		}
	}
	switch colorMode {
	case "always":
		return true
	case "never":
		return false
	default: // "auto"
		// Check if stdout is a TTY
		fi, err := os.Stdout.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
}

func parseFormatType(s string) env.FormatType {
	// Auto-detect GitHub Actions format when not explicitly set
	if s == "table" && os.Getenv("GITHUB_ACTIONS") != "" {
		return env.FormatGitHub
	}
	switch s {
	case "json":
		return env.FormatJSON
	case "github":
		return env.FormatGitHub
	default:
		return env.FormatTable
	}
}

// exitError is a cobra-compatible error that carries an exit code.
type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}
