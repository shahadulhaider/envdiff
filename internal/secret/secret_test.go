package secret

import (
	"testing"

	"github.com/shahadulhaider/envdiff/internal/env"
)

func TestIsSecret(t *testing.T) {
	tests := []struct {
		key   string
		value string
		want  bool
	}{
		// Key name heuristics
		{"DB_PASSWORD", "anything", true},
		{"API_KEY", "anything", true},
		{"SECRET_TOKEN", "anything", true},
		{"ACCESS_KEY", "anything", true},
		{"PRIVATE_KEY", "anything", true},
		// AWS key pattern
		{"ANY_KEY", "AKIAIOSFODNN7EXAMPLE", true},
		// Database URL with password
		{"DATABASE_URL", "postgres://user:secret@localhost/db", true},
		// Normal values (not secrets)
		{"PORT", "5432", false},
		{"HOST", "localhost", false},
		{"DEBUG", "true", false},
		{"APP_NAME", "myapp", false},
		{"LOG_LEVEL", "info", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			got := IsSecret(tt.key, tt.value)
			if got != tt.want {
				t.Errorf("IsSecret(%q, %q) = %v, want %v", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

func TestDetectSecrets(t *testing.T) {
	entries := []env.EnvEntry{
		{Key: "PORT", Value: "5432"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "anything"},
	}

	secrets := DetectSecrets(entries)
	if len(secrets) != 2 {
		t.Errorf("DetectSecrets() returned %d secrets, want 2", len(secrets))
	}
}

func TestMaskValue(t *testing.T) {
	tests := []string{"secret123", "a", "very long secret value here", ""}
	for _, val := range tests {
		got := MaskValue(val)
		if got != "****" {
			t.Errorf("MaskValue(%q) = %q, want %q", val, got, "****")
		}
	}
}
