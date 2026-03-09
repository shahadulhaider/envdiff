package env

import "io"

// EnvEntry represents a single environment variable entry.
type EnvEntry struct {
	Key        string
	Value      string
	LineNum    int
	Comment    string
	IsExported bool
	Raw        string
}

// EnvFile represents a parsed .env file.
type EnvFile struct {
	Entries    []EnvEntry
	Path       string
	Comments   []string
	Duplicates []EnvEntry
}

// Keys returns all keys in the EnvFile.
func (f *EnvFile) Keys() []string {
	keys := make([]string, len(f.Entries))
	for i, entry := range f.Entries {
		keys[i] = entry.Key
	}
	return keys
}

// Get retrieves an EnvEntry by key.
func (f *EnvFile) Get(key string) (EnvEntry, bool) {
	for _, entry := range f.Entries {
		if entry.Key == key {
			return entry, true
		}
	}
	return EnvEntry{}, false
}

// Len returns the number of entries in the EnvFile.
func (f *EnvFile) Len() int {
	return len(f.Entries)
}

// DiffType represents the type of difference between two environment files.
type DiffType int

const (
	DiffAdded DiffType = iota
	DiffRemoved
	DiffChanged
)

// String returns the string representation of DiffType.
func (d DiffType) String() string {
	switch d {
	case DiffAdded:
		return "added"
	case DiffRemoved:
		return "removed"
	case DiffChanged:
		return "changed"
	default:
		return "unknown"
	}
}

// DiffEntry represents a single difference between two environment files.
type DiffEntry struct {
	Key   string
	Type  DiffType
	Left  *EnvEntry
	Right *EnvEntry
}

// DiffResult represents the differences between two environment files.
type DiffResult struct {
	Left    string
	Right   string
	Entries []DiffEntry
}

// HasDiffs returns true if there are any differences.
func (r *DiffResult) HasDiffs() bool {
	return len(r.Entries) > 0
}

// Added returns all added entries.
func (r *DiffResult) Added() []DiffEntry {
	var added []DiffEntry
	for _, entry := range r.Entries {
		if entry.Type == DiffAdded {
			added = append(added, entry)
		}
	}
	return added
}

// Removed returns all removed entries.
func (r *DiffResult) Removed() []DiffEntry {
	var removed []DiffEntry
	for _, entry := range r.Entries {
		if entry.Type == DiffRemoved {
			removed = append(removed, entry)
		}
	}
	return removed
}

// Changed returns all changed entries.
func (r *DiffResult) Changed() []DiffEntry {
	var changed []DiffEntry
	for _, entry := range r.Entries {
		if entry.Type == DiffChanged {
			changed = append(changed, entry)
		}
	}
	return changed
}

// MultiDiffResult represents differences across multiple environment files.
type MultiDiffResult struct {
	Files  []string
	Keys   []string
	Matrix map[string]map[string]*string
}

// FormatType represents the output format type.
type FormatType int

const (
	FormatTable FormatType = iota
	FormatJSON
	FormatGitHub
)

// String returns the string representation of FormatType.
func (f FormatType) String() string {
	switch f {
	case FormatTable:
		return "table"
	case FormatJSON:
		return "json"
	case FormatGitHub:
		return "github"
	default:
		return "unknown"
	}
}

// Formatter defines the interface for formatting diff results.
type Formatter interface {
	Format(result *DiffResult, w io.Writer) error
}

// SchemaRule defines validation rules for an environment variable.
type SchemaRule struct {
	Required bool
	Type     string
	Pattern  string
	Default  string
	Enum     []string
}

// SchemaConfig defines validation rules for an environment file.
type SchemaConfig struct {
	AllowExtra bool
	Rules      map[string]SchemaRule
}

// ValidationError represents a validation error.
type ValidationError struct {
	Key     string
	Message string
	Line    int
}

// ValidationResult represents the result of validation.
type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationError
}

// IsValid returns true if there are no errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}
