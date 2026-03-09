package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/env"
	"github.com/shahadulhaider/envdiff/internal/parser"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate .env.example or .env.schema.toml from .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, _ := cmd.Flags().GetString("source")
			genSchema, _ := cmd.Flags().GetBool("schema")
			force, _ := cmd.Flags().GetBool("force")
			return runInit(source, genSchema, force)
		},
		SilenceUsage: true,
	}
	cmd.Flags().String("source", ".env", "source .env file to read")
	cmd.Flags().Bool("schema", false, "generate .env.schema.toml instead of .env.example")
	cmd.Flags().Bool("force", false, "overwrite existing output file")
	return cmd
}

func runInit(source string, genSchema bool, force bool) error {
	ef, err := parser.ParseFile(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", source, err)
		return &exitError{code: 2}
	}

	if genSchema {
		return writeSchema(ef, force)
	}
	return writeExample(ef, force)
}

func writeExample(ef *env.EnvFile, force bool) error {
	outPath := ".env.example"
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			fmt.Fprintf(os.Stderr, "%s already exists; use --force to overwrite\n", outPath)
			return &exitError{code: 1}
		}
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", outPath, err)
		return &exitError{code: 2}
	}
	defer f.Close()

	for _, entry := range ef.Entries {
		if entry.Comment != "" {
			fmt.Fprintf(f, "%s=  # %s\n", entry.Key, entry.Comment)
		} else {
			fmt.Fprintf(f, "%s=\n", entry.Key)
		}
	}

	fmt.Printf("wrote %s\n", outPath)
	return nil
}

func writeSchema(ef *env.EnvFile, force bool) error {
	outPath := ".env.schema.toml"
	if !force {
		if _, err := os.Stat(outPath); err == nil {
			fmt.Fprintf(os.Stderr, "%s already exists; use --force to overwrite\n", outPath)
			return &exitError{code: 1}
		}
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating %s: %v\n", outPath, err)
		return &exitError{code: 2}
	}
	defer f.Close()

	fmt.Fprintf(f, "allow_extra = true\n\n")

	for _, entry := range ef.Entries {
		typ := inferType(entry.Value)
		fmt.Fprintf(f, "[vars.%s]\n", entry.Key)
		fmt.Fprintf(f, "required = true\n")
		fmt.Fprintf(f, "type = %q\n\n", typ)
	}

	fmt.Printf("wrote %s\n", outPath)
	return nil
}

func inferType(val string) string {
	if val == "" {
		return "string"
	}
	lower := strings.ToLower(val)
	switch lower {
	case "true", "false", "yes", "no", "1", "0":
		return "bool"
	}
	allDigits := true
	for _, c := range val {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return "number"
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "mysql://") ||
		strings.HasPrefix(lower, "redis://") {
		return "url"
	}
	return "string"
}
