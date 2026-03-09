package output

import (
	"fmt"
	"io"

	"github.com/shahadulhaider/envdiff/internal/env"
)

// GitHubFormatter formats diff results as GitHub Actions annotations.
type GitHubFormatter struct {
	opts Options
}

func (f *GitHubFormatter) Format(result *env.DiffResult, w io.Writer) error {
	file := result.Right
	if file == "" {
		file = "unknown"
	}

	for _, e := range result.Added() {
		line := 0
		if e.Right != nil {
			line = e.Right.LineNum
		}
		fmt.Fprintf(w, "::notice file=%s,line=%d::%s is missing from %s\n", file, line, e.Key, result.Left)
	}

	for _, e := range result.Removed() {
		line := 0
		if e.Left != nil {
			line = e.Left.LineNum
		}
		fmt.Fprintf(w, "::error file=%s,line=%d::%s was removed from %s\n", result.Left, line, e.Key, result.Right)
	}

	for _, e := range result.Changed() {
		line := 0
		if e.Right != nil {
			line = e.Right.LineNum
		}
		val := ""
		if !f.opts.Mask && !f.opts.NoValues && e.Right != nil {
			val = fmt.Sprintf(" (was: %s, now: %s)", e.Left.Value, e.Right.Value)
		}
		fmt.Fprintf(w, "::warning file=%s,line=%d::%s value changed%s\n", file, line, e.Key, val)
	}

	return nil
}
