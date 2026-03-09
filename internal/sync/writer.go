package sync

import (
	"fmt"
	"os"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/env"
)

func WriteEnvFile(path string, file *env.EnvFile) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	for _, entry := range file.Entries {
		line := formatEntry(entry)
		fmt.Fprintln(f, line)
	}
	return nil
}

func formatEntry(e env.EnvEntry) string {
	prefix := ""
	if e.IsExported {
		prefix = "export "
	}
	val := e.Value
	if strings.ContainsAny(val, " \t#\"'") {
		val = `"` + strings.ReplaceAll(val, `"`, `\"`) + `"`
	}
	line := fmt.Sprintf("%s%s=%s", prefix, e.Key, val)
	if e.Comment != "" {
		line += " # " + e.Comment
	}
	return line
}

func ApplyChanges(target *env.EnvFile, entries []env.DiffEntry) *env.EnvFile {
	entryMap := make(map[string]int)
	for i, e := range target.Entries {
		entryMap[e.Key] = i
	}

	result := &env.EnvFile{
		Path:     target.Path,
		Comments: target.Comments,
	}
	result.Entries = make([]env.EnvEntry, len(target.Entries))
	copy(result.Entries, target.Entries)

	for _, d := range entries {
		switch d.Type {
		case env.DiffAdded:
			if d.Right != nil {
				result.Entries = append(result.Entries, *d.Right)
			}
		case env.DiffRemoved:
			if idx, ok := entryMap[d.Key]; ok {
				result.Entries = append(result.Entries[:idx], result.Entries[idx+1:]...)
				entryMap = make(map[string]int)
				for i, e := range result.Entries {
					entryMap[e.Key] = i
				}
			}
		case env.DiffChanged:
			if idx, ok := entryMap[d.Key]; ok && d.Right != nil {
				result.Entries[idx].Value = d.Right.Value
			}
		}
	}
	return result
}
