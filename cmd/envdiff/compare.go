package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/diff"
	"github.com/shahadulhaider/envdiff/internal/env"
	"github.com/shahadulhaider/envdiff/internal/parser"
	"github.com/spf13/cobra"
)

func newCompareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "compare <file1> <file2> [file3...]",
		Short: "Compare multiple .env files in a matrix view",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(cmd, args)
		},
		SilenceUsage: true,
	}
}

func runCompare(cmd *cobra.Command, paths []string) error {
	var files []*env.EnvFile
	for _, p := range paths {
		f, err := parser.ParseFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", p, err)
			return &exitError{code: 2}
		}
		files = append(files, f)
	}

	result := diff.MultiDiff(files)

	ignorePattern, _ := cmd.Flags().GetString("ignore")
	mask, _ := cmd.Flags().GetBool("mask")
	formatStr, _ := cmd.Flags().GetString("format")

	if formatStr == "json" {
		return printCompareJSON(result, mask)
	}

	return printCompareTable(result, ignorePattern, mask)
}

func printCompareJSON(result *env.MultiDiffResult, mask bool) error {
	type jsonMatrix struct {
		Keys   []string                      `json:"keys"`
		Files  []string                      `json:"files"`
		Matrix map[string]map[string]*string `json:"matrix"`
	}

	out := jsonMatrix{
		Keys:   result.Keys,
		Files:  result.Files,
		Matrix: make(map[string]map[string]*string),
	}

	for key, fileMap := range result.Matrix {
		out.Matrix[key] = make(map[string]*string)
		for file, val := range fileMap {
			if mask && val != nil {
				masked := "****"
				out.Matrix[key][file] = &masked
			} else {
				out.Matrix[key][file] = val
			}
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func printCompareTable(result *env.MultiDiffResult, ignorePattern string, mask bool) error {
	headers := make([]string, len(result.Files))
	for i, f := range result.Files {
		parts := strings.Split(f, "/")
		headers[i] = parts[len(parts)-1]
	}

	keyWidth := 10
	for _, k := range result.Keys {
		if len(k) > keyWidth {
			keyWidth = len(k)
		}
	}
	colWidth := 10
	for _, h := range headers {
		if len(h) > colWidth {
			colWidth = len(h)
		}
	}
	if colWidth > 20 {
		colWidth = 20
	}

	groups := groupByPrefix(result.Keys)

	header := fmt.Sprintf("%-*s", keyWidth, "KEY")
	for _, h := range headers {
		header += fmt.Sprintf("  %-*s", colWidth, truncate(h, colWidth))
	}
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	for _, group := range groups {
		if group.prefix != "" && len(group.keys) > 1 {
			fmt.Printf("\n[%s]\n", group.prefix)
		}
		for _, key := range group.keys {
			if ignorePattern != "" {
				if matched, _ := filepath.Match(ignorePattern, key); matched {
					continue
				}
			}

			row := fmt.Sprintf("%-*s", keyWidth, key)
			for _, file := range result.Files {
				val := "<missing>"
				if v := result.Matrix[key][file]; v != nil {
					if mask {
						val = "****"
					} else {
						val = *v
					}
				}
				row += fmt.Sprintf("  %-*s", colWidth, truncate(val, colWidth))
			}
			fmt.Println(row)
		}
	}
	return nil
}

type keyGroup struct {
	prefix string
	keys   []string
}

func groupByPrefix(keys []string) []keyGroup {
	prefixCount := make(map[string]int)
	for _, k := range keys {
		parts := strings.SplitN(k, "_", 2)
		if len(parts) > 1 {
			prefixCount[parts[0]]++
		}
	}

	var groups []keyGroup
	seen := make(map[string]bool)
	var ungrouped []string

	for _, k := range keys {
		parts := strings.SplitN(k, "_", 2)
		if len(parts) > 1 && prefixCount[parts[0]] >= 2 {
			prefix := parts[0]
			if !seen[prefix] {
				seen[prefix] = true
				var grpKeys []string
				for _, k2 := range keys {
					if strings.HasPrefix(k2, prefix+"_") {
						grpKeys = append(grpKeys, k2)
					}
				}
				groups = append(groups, keyGroup{prefix: prefix, keys: grpKeys})
			}
		} else {
			ungrouped = append(ungrouped, k)
		}
	}

	if len(ungrouped) > 0 {
		result := []keyGroup{{prefix: "", keys: ungrouped}}
		result = append(result, groups...)
		sort.Slice(result[1:], func(i, j int) bool {
			return result[i+1].prefix < result[j+1].prefix
		})
		return result
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].prefix < groups[j].prefix
	})
	return groups
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max > 3 {
		return s[:max-3] + "..."
	}
	return s[:max]
}
