package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shahadulhaider/envdiff/internal/env"
	"github.com/shahadulhaider/envdiff/internal/parser"
	"github.com/shahadulhaider/envdiff/internal/schema"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate --schema <schema.toml> <file> [file...]",
		Short: "Validate .env files against a TOML schema",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaPath, _ := cmd.Flags().GetString("schema")
			if schemaPath == "" {
				fmt.Fprintf(os.Stderr, "error: --schema flag is required\n")
				return &exitError{code: 2}
			}
			return runValidate(cmd, schemaPath, args)
		},
		SilenceUsage: true,
	}
	cmd.Flags().String("schema", "", "path to TOML schema file (required)")
	return cmd
}

type fileValidationResult struct {
	File     string
	Result   *env.ValidationResult
	ParseErr error
}

func runValidate(cmd *cobra.Command, schemaPath string, files []string) error {
	sc, err := schema.LoadSchema(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading schema %s: %v\n", schemaPath, err)
		return &exitError{code: 2}
	}

	formatStr, _ := cmd.Flags().GetString("format")
	mask, _ := cmd.Flags().GetBool("mask")

	var results []fileValidationResult
	anyFailure := false

	for _, f := range files {
		ef, err := parser.ParseFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", f, err)
			return &exitError{code: 2}
		}
		vr := schema.Validate(ef, sc)
		results = append(results, fileValidationResult{File: f, Result: vr})
		if !vr.IsValid() {
			anyFailure = true
		}
	}

	switch formatStr {
	case "json":
		printValidateJSON(results, mask)
	case "github":
		printValidateGitHub(results, mask)
	default:
		printValidateTable(results, mask)
	}

	if anyFailure {
		return &exitError{code: 1}
	}
	return nil
}

func printValidateTable(results []fileValidationResult, mask bool) {
	for _, r := range results {
		if len(results) > 1 {
			fmt.Printf("%s:\n", r.File)
		}
		if r.Result.IsValid() {
			fmt.Printf("  ok\n")
			continue
		}
		for _, e := range r.Result.Errors {
			fmt.Printf("  error: %s\n", e.Message)
		}
		for _, w := range r.Result.Warnings {
			fmt.Printf("  warning: %s\n", w.Message)
		}
	}
}

func printValidateGitHub(results []fileValidationResult, mask bool) {
	for _, r := range results {
		for _, e := range r.Result.Errors {
			fmt.Printf("::error file=%s::%s\n", r.File, e.Message)
		}
		for _, w := range r.Result.Warnings {
			fmt.Printf("::warning file=%s::%s\n", r.File, w.Message)
		}
	}
}

func printValidateJSON(results []fileValidationResult, mask bool) {
	type jsonError struct {
		Key     string `json:"key"`
		Message string `json:"message"`
	}
	type jsonResult struct {
		File     string      `json:"file"`
		Valid    bool        `json:"valid"`
		Errors   []jsonError `json:"errors"`
		Warnings []jsonError `json:"warnings"`
	}

	var out []jsonResult
	for _, r := range results {
		jr := jsonResult{
			File:     r.File,
			Valid:    r.Result.IsValid(),
			Errors:   []jsonError{},
			Warnings: []jsonError{},
		}
		for _, e := range r.Result.Errors {
			jr.Errors = append(jr.Errors, jsonError{Key: e.Key, Message: e.Message})
		}
		for _, w := range r.Result.Warnings {
			jr.Warnings = append(jr.Warnings, jsonError{Key: w.Key, Message: w.Message})
		}
		out = append(out, jr)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}
