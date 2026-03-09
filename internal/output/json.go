package output

import (
	"encoding/json"
	"io"

	"github.com/shahadulhaider/envdiff/internal/env"
)

// JSONFormatter formats diff results as JSON.
type JSONFormatter struct {
	opts Options
}

type jsonDiffEntry struct {
	Key        string `json:"key"`
	LeftValue  string `json:"left_value,omitempty"`
	RightValue string `json:"right_value,omitempty"`
}

type jsonOutput struct {
	Left    string          `json:"left"`
	Right   string          `json:"right"`
	Added   []jsonDiffEntry `json:"added"`
	Removed []jsonDiffEntry `json:"removed"`
	Changed []jsonDiffEntry `json:"changed"`
}

func (f *JSONFormatter) Format(result *env.DiffResult, w io.Writer) error {
	out := jsonOutput{
		Left:    result.Left,
		Right:   result.Right,
		Added:   make([]jsonDiffEntry, 0),
		Removed: make([]jsonDiffEntry, 0),
		Changed: make([]jsonDiffEntry, 0),
	}

	for _, e := range result.Added() {
		val := ""
		if e.Right != nil {
			val = e.Right.Value
			if f.opts.Mask {
				val = "****"
			}
		}
		out.Added = append(out.Added, jsonDiffEntry{Key: e.Key, RightValue: val})
	}

	for _, e := range result.Removed() {
		val := ""
		if e.Left != nil {
			val = e.Left.Value
			if f.opts.Mask {
				val = "****"
			}
		}
		out.Removed = append(out.Removed, jsonDiffEntry{Key: e.Key, LeftValue: val})
	}

	for _, e := range result.Changed() {
		leftVal, rightVal := "", ""
		if e.Left != nil {
			leftVal = e.Left.Value
			if f.opts.Mask {
				leftVal = "****"
			}
		}
		if e.Right != nil {
			rightVal = e.Right.Value
			if f.opts.Mask {
				rightVal = "****"
			}
		}
		out.Changed = append(out.Changed, jsonDiffEntry{Key: e.Key, LeftValue: leftVal, RightValue: rightVal})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
