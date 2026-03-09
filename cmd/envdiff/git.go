package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/diff"
	gitpkg "github.com/shahadulhaider/envdiff/internal/git"
	"github.com/shahadulhaider/envdiff/internal/output"
	"github.com/shahadulhaider/envdiff/internal/parser"
	"github.com/spf13/cobra"
)

func newGitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git <file>",
		Short: "Diff .env file against a git ref",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, _ := cmd.Flags().GetString("ref")
			return runGitDiff(cmd, ref, args[0])
		},
		SilenceUsage: true,
	}
	cmd.Flags().String("ref", "HEAD", "git ref to compare against")
	return cmd
}

func runGitDiff(cmd *cobra.Command, ref, filePath string) error {
	if !gitpkg.IsGitRepo() {
		fmt.Fprintf(os.Stderr, "error: not a git repository\n")
		return &exitError{code: 2}
	}

	// Get file content at ref
	content, err := gitpkg.ShowFileAtRef(ref, filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s at ref %s: %v\n", filePath, ref, err)
		return &exitError{code: 2}
	}

	// Parse the ref version (empty string → empty env file)
	refFile, err := parser.Parse(strings.NewReader(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s at ref %s: %v\n", filePath, ref, err)
		return &exitError{code: 2}
	}
	refFile.Path = ref + ":" + filePath

	// Parse the current version
	currentFile, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", filePath, err)
		return &exitError{code: 2}
	}

	result := diff.Diff(refFile, currentFile)

	// Apply --ignore filter
	ignorePattern, _ := cmd.Flags().GetString("ignore")
	if ignorePattern != "" {
		result = filterDiffResult(result, ignorePattern)
	}

	if !result.HasDiffs() {
		return nil
	}

	mask, _ := cmd.Flags().GetBool("mask")
	noValues, _ := cmd.Flags().GetBool("no-values")
	colorMode, _ := cmd.Flags().GetString("color")
	formatStr, _ := cmd.Flags().GetString("format")

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

	return &exitError{code: 1}
}
