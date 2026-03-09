package diff

import (
	"github.com/shahadulhaider/envdiff/internal/env"
)

// Diff compares two env files and returns the differences.
// Order: removed first, then changed, then added (matches diff(1) convention).
func Diff(left, right *env.EnvFile) *env.DiffResult {
	result := &env.DiffResult{
		Left:  left.Path,
		Right: right.Path,
	}

	// Build lookup maps
	leftMap := make(map[string]env.EnvEntry)
	for _, e := range left.Entries {
		leftMap[e.Key] = e
	}
	rightMap := make(map[string]env.EnvEntry)
	for _, e := range right.Entries {
		rightMap[e.Key] = e
	}

	// Removed: in left but not right
	for _, e := range left.Entries {
		if _, exists := rightMap[e.Key]; !exists {
			le := e
			result.Entries = append(result.Entries, env.DiffEntry{
				Key:   e.Key,
				Type:  env.DiffRemoved,
				Left:  &le,
				Right: nil,
			})
		}
	}

	// Changed: in both but values differ
	for _, e := range left.Entries {
		if re, exists := rightMap[e.Key]; exists {
			if e.Value != re.Value {
				le := e
				re2 := re
				result.Entries = append(result.Entries, env.DiffEntry{
					Key:   e.Key,
					Type:  env.DiffChanged,
					Left:  &le,
					Right: &re2,
				})
			}
		}
	}

	// Added: in right but not left
	for _, e := range right.Entries {
		if _, exists := leftMap[e.Key]; !exists {
			re := e
			result.Entries = append(result.Entries, env.DiffEntry{
				Key:   e.Key,
				Type:  env.DiffAdded,
				Left:  nil,
				Right: &re,
			})
		}
	}

	return result
}

// MultiDiff compares N env files and returns a matrix result.
func MultiDiff(files []*env.EnvFile) *env.MultiDiffResult {
	result := &env.MultiDiffResult{
		Matrix: make(map[string]map[string]*string),
	}

	// Collect file paths and all keys (union, preserving first-file order)
	for _, f := range files {
		result.Files = append(result.Files, f.Path)
	}

	seen := make(map[string]bool)
	for _, f := range files {
		for _, e := range f.Entries {
			if !seen[e.Key] {
				seen[e.Key] = true
				result.Keys = append(result.Keys, e.Key)
			}
		}
	}

	// Build matrix
	for _, key := range result.Keys {
		result.Matrix[key] = make(map[string]*string)
		for _, f := range files {
			if entry, ok := f.Get(key); ok {
				val := entry.Value
				result.Matrix[key][f.Path] = &val
			} else {
				result.Matrix[key][f.Path] = nil
			}
		}
	}

	return result
}
