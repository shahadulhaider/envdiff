package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/hook"
	"github.com/spf13/cobra"
)

func newHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage pre-commit hooks",
		Long:  "Install, uninstall, or check status of envdiff pre-commit hooks.",
	}

	cmd.AddCommand(
		newHookInstallCmd(),
		newHookUninstallCmd(),
		newHookStatusCmd(),
	)

	return cmd
}

func gitRepoRoot() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not a git repository")
	}
	return strings.TrimSpace(out.String()), nil
}

func newHookInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "install",
		Short:        "Install envdiff pre-commit hook",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := gitRepoRoot()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 2}
			}
			if err := hook.Install(root); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 1}
			}
			fmt.Println("envdiff hook installed")
			return nil
		},
	}
}

func newHookUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "uninstall",
		Short:        "Remove envdiff pre-commit hook",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := gitRepoRoot()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 2}
			}
			if err := hook.Uninstall(root); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 1}
			}
			fmt.Println("envdiff hook removed")
			return nil
		},
	}
}

func newHookStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "status",
		Short:        "Show envdiff hook installation status",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := gitRepoRoot()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 2}
			}
			installed, err := hook.Status(root)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return &exitError{code: 2}
			}
			if installed {
				fmt.Println("envdiff hook: installed")
			} else {
				fmt.Println("envdiff hook: not installed")
			}
			return nil
		},
	}
}
