package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/shahadulhaider/envdiff/internal/diff"
	"github.com/shahadulhaider/envdiff/internal/env"
	"github.com/shahadulhaider/envdiff/internal/parser"
	syncp "github.com/shahadulhaider/envdiff/internal/sync"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "sync <source> <target>",
		Short:        "Interactively sync keys from source to target",
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(args[0], args[1])
		},
	}
}

func runSync(source, target string) error {
	fi, err := os.Stdout.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice) == 0 {
		fmt.Fprintf(os.Stderr, "error: sync requires interactive terminal\n")
		return &exitError{code: 2}
	}

	sourceFile, err := parser.ParseFile(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", source, err)
		return &exitError{code: 2}
	}

	tf, err := parser.ParseFile(target)
	if err != nil {
		if os.IsNotExist(err) {
			tf = &env.EnvFile{Path: target}
		} else {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", target, err)
			return &exitError{code: 2}
		}
	}

	result := diff.Diff(sourceFile, tf)
	if !result.HasDiffs() {
		fmt.Println("No differences found.")
		return nil
	}

	model := syncp.NewModel(source, target, result.Entries)
	p := tea.NewProgram(model)
	finalModelI, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running TUI: %v\n", err)
		return &exitError{code: 2}
	}

	finalModel, ok := finalModelI.(syncp.Model)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unexpected model type\n")
		return &exitError{code: 2}
	}

	if finalModel.Quit {
		fmt.Println("Cancelled.")
		return nil
	}

	if finalModel.Applied {
		selected := finalModel.SelectedEntries()
		if len(selected) == 0 {
			fmt.Println("No changes selected.")
			return nil
		}
		updated := syncp.ApplyChanges(tf, selected)
		if err := syncp.WriteEnvFile(target, updated); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", target, err)
			return &exitError{code: 2}
		}
		fmt.Printf("Applied %d changes to %s\n", len(selected), target)
	}

	return nil
}
