package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shahadulhaider/envdiff/internal/env"
)

func makeEnvFile(kvs ...string) *env.EnvFile {
	f := &env.EnvFile{}
	for i := 0; i+1 < len(kvs); i += 2 {
		f.Entries = append(f.Entries, env.EnvEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return f
}

func writeSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeSchemaFile: %v", err)
	}
	return path
}

func TestLoadSchema(t *testing.T) {
	content := `
allow_extra = true

[vars.DB_HOST]
required = true
type = "string"

[vars.DB_PORT]
required = true
type = "number"
`
	path := writeSchemaFile(t, content)
	cfg, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("LoadSchema() error = %v", err)
	}
	if !cfg.AllowExtra {
		t.Error("expected AllowExtra = true")
	}
	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cfg.Rules))
	}
	if !cfg.Rules["DB_HOST"].Required {
		t.Error("DB_HOST should be required")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		schema     *env.SchemaConfig
		envFile    *env.EnvFile
		wantValid  bool
		wantErrKey string
	}{
		{
			name: "valid env",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"HOST": {Required: true, Type: "string"},
					"PORT": {Required: true, Type: "number"},
				},
			},
			envFile:   makeEnvFile("HOST", "localhost", "PORT", "5432"),
			wantValid: true,
		},
		{
			name: "missing required key",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"API_KEY": {Required: true, Type: "string"},
				},
			},
			envFile:    makeEnvFile("HOST", "localhost"),
			wantValid:  false,
			wantErrKey: "API_KEY",
		},
		{
			name: "type mismatch number",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"PORT": {Type: "number"},
				},
			},
			envFile:    makeEnvFile("PORT", "not_a_number"),
			wantValid:  false,
			wantErrKey: "PORT",
		},
		{
			name: "bool validation",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"DEBUG": {Type: "bool"},
				},
			},
			envFile:   makeEnvFile("DEBUG", "true"),
			wantValid: true,
		},
		{
			name: "enum validation pass",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"LEVEL": {Type: "enum", Enum: []string{"debug", "info", "warn", "error"}},
				},
			},
			envFile:   makeEnvFile("LEVEL", "info"),
			wantValid: true,
		},
		{
			name: "enum validation fail",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"LEVEL": {Type: "enum", Enum: []string{"debug", "info", "warn", "error"}},
				},
			},
			envFile:    makeEnvFile("LEVEL", "invalid"),
			wantValid:  false,
			wantErrKey: "LEVEL",
		},
		{
			name: "extra key not allowed",
			schema: &env.SchemaConfig{
				AllowExtra: false,
				Rules: map[string]env.SchemaRule{
					"HOST": {Type: "string"},
				},
			},
			envFile:    makeEnvFile("HOST", "localhost", "EXTRA", "value"),
			wantValid:  false,
			wantErrKey: "EXTRA",
		},
		{
			name: "pattern match fail",
			schema: &env.SchemaConfig{
				AllowExtra: true,
				Rules: map[string]env.SchemaRule{
					"API_KEY": {Pattern: `^[A-Z0-9]{32}$`},
				},
			},
			envFile:    makeEnvFile("API_KEY", "tooshort"),
			wantValid:  false,
			wantErrKey: "API_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.envFile, tt.schema)
			if result.IsValid() != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v; errors: %v", result.IsValid(), tt.wantValid, result.Errors)
			}
			if tt.wantErrKey != "" {
				found := false
				for _, e := range result.Errors {
					if e.Key == tt.wantErrKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error for key %q, got errors: %v", tt.wantErrKey, result.Errors)
				}
			}
		})
	}
}
