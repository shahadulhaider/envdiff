package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/shahadulhaider/envdiff/internal/env"
)

var (
	addedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	removedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	changedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
)

// TableFormatter formats diff results as a table.
type TableFormatter struct {
	opts Options
}

func (f *TableFormatter) Format(result *env.DiffResult, w io.Writer) error {
	entries := result.Entries

	for _, e := range entries {
		var line string
		switch e.Type {
		case env.DiffAdded:
			val := f.formatValue(e.Right)
			line = fmt.Sprintf("+ %s=%s", e.Key, val)
			if f.opts.Color {
				line = addedStyle.Render(line)
			}
		case env.DiffRemoved:
			val := f.formatValue(e.Left)
			line = fmt.Sprintf("- %s=%s", e.Key, val)
			if f.opts.Color {
				line = removedStyle.Render(line)
			}
		case env.DiffChanged:
			leftVal := f.formatValue(e.Left)
			rightVal := f.formatValue(e.Right)
			if f.opts.NoValues {
				line = fmt.Sprintf("~ %s", e.Key)
			} else {
				line = fmt.Sprintf("~ %s=%s -> %s", e.Key, leftVal, rightVal)
			}
			if f.opts.Color {
				line = changedStyle.Render(line)
			}
		}
		fmt.Fprintln(w, line)
	}

	added := len(result.Added())
	removed := len(result.Removed())
	changed := len(result.Changed())
	fmt.Fprintf(w, "\n%d added, %d removed, %d changed\n", added, removed, changed)
	return nil
}

func (f *TableFormatter) formatValue(e *env.EnvEntry) string {
	if e == nil {
		return ""
	}
	if f.opts.Mask {
		return "****"
	}
	if f.opts.NoValues {
		return ""
	}
	return e.Value
}
