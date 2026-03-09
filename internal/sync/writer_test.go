package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shahadulhaider/envdiff/internal/env"
)

func TestWriteEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env")

	file := &env.EnvFile{
		Entries: []env.EnvEntry{
			{Key: "HOST", Value: "localhost"},
			{Key: "PORT", Value: "5432"},
		},
	}

	if err := WriteEnvFile(path, file); err != nil {
		t.Fatalf("WriteEnvFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "HOST=localhost") {
		t.Errorf("expected HOST=localhost, got:\n%s", content)
	}
	if !strings.Contains(content, "PORT=5432") {
		t.Errorf("expected PORT=5432, got:\n%s", content)
	}
}

func TestApplyChanges(t *testing.T) {
	target := &env.EnvFile{
		Entries: []env.EnvEntry{
			{Key: "A", Value: "1"},
			{Key: "B", Value: "2"},
		},
	}

	rightC := env.EnvEntry{Key: "C", Value: "3"}
	rightB := env.EnvEntry{Key: "B", Value: "99"}

	entries := []env.DiffEntry{
		{Key: "C", Type: env.DiffAdded, Right: &rightC},
		{Key: "B", Type: env.DiffChanged, Right: &rightB},
	}

	result := ApplyChanges(target, entries)

	for _, e := range result.Entries {
		if e.Key == "B" && e.Value != "99" {
			t.Errorf("expected B=99, got B=%s", e.Value)
		}
	}

	foundC := false
	for _, e := range result.Entries {
		if e.Key == "C" {
			foundC = true
		}
	}
	if !foundC {
		t.Error("C not found in result after apply")
	}
}
