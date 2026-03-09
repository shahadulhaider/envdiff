package diff

import (
	"testing"

	"github.com/shahadulhaider/envdiff/internal/env"
)

func makeFile(path string, kvs ...string) *env.EnvFile {
	f := &env.EnvFile{Path: path}
	for i := 0; i+1 < len(kvs); i += 2 {
		f.Entries = append(f.Entries, env.EnvEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return f
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name        string
		left        *env.EnvFile
		right       *env.EnvFile
		wantHasDiff bool
		wantAdded   []string
		wantRemoved []string
		wantChanged []string
	}{
		{
			name:        "identical files",
			left:        makeFile("l.env", "A", "1", "B", "2"),
			right:       makeFile("r.env", "A", "1", "B", "2"),
			wantHasDiff: false,
		},
		{
			name:        "all added",
			left:        makeFile("l.env"),
			right:       makeFile("r.env", "A", "1", "B", "2"),
			wantHasDiff: true,
			wantAdded:   []string{"A", "B"},
		},
		{
			name:        "all removed",
			left:        makeFile("l.env", "A", "1", "B", "2"),
			right:       makeFile("r.env"),
			wantHasDiff: true,
			wantRemoved: []string{"A", "B"},
		},
		{
			name:        "mixed changes",
			left:        makeFile("l.env", "A", "1", "B", "2", "C", "3"),
			right:       makeFile("r.env", "B", "2", "C", "99", "D", "4"),
			wantHasDiff: true,
			wantRemoved: []string{"A"},
			wantChanged: []string{"C"},
			wantAdded:   []string{"D"},
		},
		{
			name:        "empty value vs missing key",
			left:        makeFile("l.env", "KEY", ""),
			right:       makeFile("r.env"),
			wantHasDiff: true,
			wantRemoved: []string{"KEY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Diff(tt.left, tt.right)

			if result.HasDiffs() != tt.wantHasDiff {
				t.Errorf("HasDiffs() = %v, want %v", result.HasDiffs(), tt.wantHasDiff)
			}

			checkKeys := func(entries []env.DiffEntry, want []string, label string) {
				t.Helper()
				got := make([]string, len(entries))
				for i, e := range entries {
					got[i] = e.Key
				}
				if len(got) != len(want) {
					t.Errorf("%s: got keys %v, want %v", label, got, want)
					return
				}
				for i, k := range want {
					if got[i] != k {
						t.Errorf("%s[%d]: got %q, want %q", label, i, got[i], k)
					}
				}
			}

			if tt.wantAdded != nil {
				checkKeys(result.Added(), tt.wantAdded, "Added")
			}
			if tt.wantRemoved != nil {
				checkKeys(result.Removed(), tt.wantRemoved, "Removed")
			}
			if tt.wantChanged != nil {
				checkKeys(result.Changed(), tt.wantChanged, "Changed")
			}
		})
	}
}

func TestMultiDiff(t *testing.T) {
	dev := makeFile("dev.env", "A", "1", "B", "2", "C", "3")
	stg := makeFile("stg.env", "A", "1", "B", "99")
	prd := makeFile("prd.env", "A", "1", "C", "3", "D", "4")

	result := MultiDiff([]*env.EnvFile{dev, stg, prd})

	if len(result.Files) != 3 {
		t.Errorf("Files count = %d, want 3", len(result.Files))
	}

	// All keys present: A, B, C, D
	keySet := make(map[string]bool)
	for _, k := range result.Keys {
		keySet[k] = true
	}
	for _, k := range []string{"A", "B", "C", "D"} {
		if !keySet[k] {
			t.Errorf("key %q missing from result.Keys", k)
		}
	}

	// B is missing from prd
	if result.Matrix["B"]["prd.env"] != nil {
		t.Errorf("B in prd.env should be nil (missing)")
	}

	// D is missing from dev and stg
	if result.Matrix["D"]["dev.env"] != nil {
		t.Errorf("D in dev.env should be nil (missing)")
	}
	if result.Matrix["D"]["stg.env"] != nil {
		t.Errorf("D in stg.env should be nil (missing)")
	}
	// D is present in prd
	if result.Matrix["D"]["prd.env"] == nil || *result.Matrix["D"]["prd.env"] != "4" {
		t.Errorf("D in prd.env should be '4'")
	}
}
