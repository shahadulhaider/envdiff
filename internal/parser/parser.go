package parser

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/env"
)

// Parse parses .env content from a reader.
func Parse(r io.Reader) (*env.EnvFile, error) {
	ef := &env.EnvFile{}
	seen := make(map[string]int)

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		line := raw

		if lineNum == 1 {
			line = strings.TrimPrefix(line, "\xef\xbb\xbf")
		}

		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			ef.Comments = append(ef.Comments, trimmed)
			continue
		}

		isExported := false
		if strings.HasPrefix(trimmed, "export ") {
			isExported = true
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "export "))
		}

		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx < 0 {
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		rawVal := trimmed[eqIdx+1:]
		value, comment := parseValue(rawVal)

		entry := env.EnvEntry{
			Key:        key,
			Value:      value,
			LineNum:    lineNum,
			Comment:    comment,
			IsExported: isExported,
			Raw:        raw,
		}

		if idx, exists := seen[key]; exists {
			ef.Duplicates = append(ef.Duplicates, ef.Entries[idx])
			ef.Entries[idx] = entry
		} else {
			seen[key] = len(ef.Entries)
			ef.Entries = append(ef.Entries, entry)
		}
	}

	return ef, scanner.Err()
}

func parseValue(raw string) (value, comment string) {
	if len(raw) == 0 {
		return "", ""
	}
	if raw[0] == '"' {
		end := strings.LastIndex(raw, "\"")
		if end > 0 {
			v := strings.ReplaceAll(raw[1:end], `\"`, `"`)
			return v, ""
		}
		return raw[1:], ""
	}
	if raw[0] == '\'' {
		end := strings.LastIndex(raw, "'")
		if end > 0 {
			return raw[1:end], ""
		}
		return raw[1:], ""
	}
	if idx := strings.Index(raw, " # "); idx >= 0 {
		return strings.TrimSpace(raw[:idx]), strings.TrimSpace(raw[idx+3:])
	}
	return strings.TrimSpace(raw), ""
}

// ParseFile parses a .env file by path.
func ParseFile(path string) (*env.EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ef, err := Parse(f)
	if err != nil {
		return nil, err
	}
	ef.Path = path
	return ef, nil
}
