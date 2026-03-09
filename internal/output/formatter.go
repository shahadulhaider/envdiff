package output

import (
	"github.com/shahadulhaider/envdiff/internal/env"
)

// Options configures formatter behavior.
type Options struct {
	Mask     bool
	NoValues bool
	Color    bool
}

// NewFormatter creates a Formatter for the given format type.
func NewFormatter(format env.FormatType, opts Options) env.Formatter {
	switch format {
	case env.FormatJSON:
		return &JSONFormatter{opts: opts}
	case env.FormatGitHub:
		return &GitHubFormatter{opts: opts}
	default:
		return &TableFormatter{opts: opts}
	}
}
