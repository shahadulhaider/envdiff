package root

import (
	"github.com/spf13/cobra"
)

var (
	format   string
	mask     bool
	noValues bool
	ignore   string
	color    string
	ci       bool
)

func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "envdiff",
		Short:   "Compare, validate, and sync .env files across environments",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Root command does nothing by itself
			return cmd.Help()
		},
	}

	cmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	// Add persistent flags
	cmd.PersistentFlags().StringVar(&format, "format", "table", "output format (table|json|github)")
	cmd.PersistentFlags().BoolVar(&mask, "mask", false, "mask sensitive values in output")
	cmd.PersistentFlags().BoolVar(&noValues, "no-values", false, "hide values in output")
	cmd.PersistentFlags().StringVar(&ignore, "ignore", "", "glob pattern for keys to ignore (e.g. DEBUG_*)")
	cmd.PersistentFlags().StringVar(&color, "color", "auto", "color output (auto|always|never)")
	cmd.PersistentFlags().BoolVar(&ci, "ci", false, "CI mode (no colors, machine-readable output)")

	return cmd
}
