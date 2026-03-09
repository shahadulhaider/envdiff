package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantKeys     []string
		wantValues   map[string]string
		wantExported map[string]bool
		wantComment  map[string]string
		wantDups     int
		wantComments int
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "whitespace only",
			input: "   \n  \n",
		},
		{
			name:         "comment only",
			input:        "# this is a comment\n# another comment",
			wantComments: 2,
		},
		{
			name:       "simple key=value",
			input:      "KEY=value",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "value"},
		},
		{
			name:       "double quoted",
			input:      `KEY="hello world"`,
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "hello world"},
		},
		{
			name:       "single quoted",
			input:      "KEY='hello world'",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "hello world"},
		},
		{
			name:         "export prefix",
			input:        "export KEY=value",
			wantKeys:     []string{"KEY"},
			wantValues:   map[string]string{"KEY": "value"},
			wantExported: map[string]bool{"KEY": true},
		},
		{
			name:       "empty value",
			input:      "KEY=",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": ""},
		},
		{
			name:       "malformed no equals skipped",
			input:      "NOEQUALS\nKEY=value",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "value"},
		},
		{
			name:       "duplicate keys last wins",
			input:      "KEY=first\nKEY=second",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "second"},
			wantDups:   1,
		},
		{
			name:        "inline comment",
			input:       "KEY=value # this is comment",
			wantKeys:    []string{"KEY"},
			wantValues:  map[string]string{"KEY": "value"},
			wantComment: map[string]string{"KEY": "this is comment"},
		},
		{
			name:       "whitespace around equals",
			input:      "KEY = value",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "value"},
		},
		{
			name:       "bom handling",
			input:      "\xef\xbb\xbfKEY=value",
			wantKeys:   []string{"KEY"},
			wantValues: map[string]string{"KEY": "value"},
		},
		{
			name:       "crlf line endings",
			input:      "KEY=value\r\nOTHER=val2\r\n",
			wantKeys:   []string{"KEY", "OTHER"},
			wantValues: map[string]string{"KEY": "value", "OTHER": "val2"},
		},
		{
			name:       "key ordering preserved",
			input:      "Z=1\nA=2\nM=3",
			wantKeys:   []string{"Z", "A", "M"},
			wantValues: map[string]string{"Z": "1", "A": "2", "M": "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// Check key count
			if len(tt.wantKeys) != len(got.Entries) {
				t.Errorf("got %d entries, want %d; keys=%v", len(got.Entries), len(tt.wantKeys), got.Keys())
			}

			// Check key ordering and values
			for i, key := range tt.wantKeys {
				if i >= len(got.Entries) {
					break
				}
				if got.Entries[i].Key != key {
					t.Errorf("entry[%d].Key = %q, want %q", i, got.Entries[i].Key, key)
				}
				if tt.wantValues != nil {
					if got.Entries[i].Value != tt.wantValues[key] {
						t.Errorf("entry[%d].Value = %q, want %q", i, got.Entries[i].Value, tt.wantValues[key])
					}
				}
				if tt.wantExported != nil && tt.wantExported[key] {
					if !got.Entries[i].IsExported {
						t.Errorf("entry[%d].IsExported = false, want true", i)
					}
				}
				if tt.wantComment != nil && tt.wantComment[key] != "" {
					if got.Entries[i].Comment != tt.wantComment[key] {
						t.Errorf("entry[%d].Comment = %q, want %q", i, got.Entries[i].Comment, tt.wantComment[key])
					}
				}
			}

			// Check duplicates
			if len(got.Duplicates) != tt.wantDups {
				t.Errorf("got %d duplicates, want %d", len(got.Duplicates), tt.wantDups)
			}

			// Check comments
			if len(got.Comments) != tt.wantComments {
				t.Errorf("got %d comments, want %d", len(got.Comments), tt.wantComments)
			}
		})
	}
}
