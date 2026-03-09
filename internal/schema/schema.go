package schema

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shahadulhaider/envdiff/internal/env"
)

// tomlSchema is the raw TOML structure.
type tomlSchema struct {
	AllowExtra bool                     `toml:"allow_extra"`
	Vars       map[string]tomlSchemaVar `toml:"vars"`
}

type tomlSchemaVar struct {
	Required bool     `toml:"required"`
	Type     string   `toml:"type"`
	Pattern  string   `toml:"pattern"`
	Default  string   `toml:"default"`
	Enum     []string `toml:"enum"`
}

// LoadSchema parses a TOML schema file.
func LoadSchema(path string) (*env.SchemaConfig, error) {
	var raw tomlSchema
	raw.AllowExtra = true // default

	if _, err := toml.DecodeFile(path, &raw); err != nil {
		return nil, fmt.Errorf("parse schema %s: %w", path, err)
	}

	cfg := &env.SchemaConfig{
		AllowExtra: raw.AllowExtra,
		Rules:      make(map[string]env.SchemaRule),
	}

	for key, v := range raw.Vars {
		cfg.Rules[key] = env.SchemaRule{
			Required: v.Required,
			Type:     v.Type,
			Pattern:  v.Pattern,
			Default:  v.Default,
			Enum:     v.Enum,
		}
	}

	return cfg, nil
}

// Validate validates an env file against a schema.
func Validate(envFile *env.EnvFile, schema *env.SchemaConfig) *env.ValidationResult {
	result := &env.ValidationResult{}

	// Build key set from env file
	keySet := make(map[string]string)
	for _, e := range envFile.Entries {
		keySet[e.Key] = e.Value
	}

	// Check each rule
	for key, rule := range schema.Rules {
		val, exists := keySet[key]

		if rule.Required {
			if !exists || val == "" {
				result.Errors = append(result.Errors, env.ValidationError{
					Key:     key,
					Message: fmt.Sprintf("required key %q is missing or empty", key),
				})
				continue
			}
		}

		if !exists {
			continue
		}

		// Type validation
		if rule.Type != "" {
			if err := validateType(key, val, rule.Type, rule.Enum); err != nil {
				result.Errors = append(result.Errors, env.ValidationError{
					Key:     key,
					Message: err.Error(),
				})
				continue
			}
		}

		// Pattern validation
		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				result.Warnings = append(result.Warnings, env.ValidationError{
					Key:     key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
			} else if !re.MatchString(val) {
				result.Errors = append(result.Errors, env.ValidationError{
					Key:     key,
					Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern),
				})
			}
		}
	}

	// Check for extra keys
	if !schema.AllowExtra {
		for key := range keySet {
			if _, defined := schema.Rules[key]; !defined {
				result.Errors = append(result.Errors, env.ValidationError{
					Key:     key,
					Message: fmt.Sprintf("unexpected key %q not defined in schema", key),
				})
			}
		}
	}

	return result
}

func validateType(key, val, typ string, enum []string) error {
	switch strings.ToLower(typ) {
	case "string":
		// Any non-empty string (required check handles empty)
		return nil
	case "number", "int", "integer":
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return fmt.Errorf("key %q: expected number, got %q", key, val)
		}
	case "bool", "boolean":
		lower := strings.ToLower(val)
		valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
		if !valid[lower] {
			return fmt.Errorf("key %q: expected bool (true/false/1/0/yes/no), got %q", key, val)
		}
	case "url":
		u, err := url.Parse(val)
		if err != nil || u.Scheme == "" {
			return fmt.Errorf("key %q: expected URL with scheme, got %q", key, val)
		}
	case "email":
		if !strings.Contains(val, "@") || !strings.Contains(val, ".") {
			return fmt.Errorf("key %q: expected email address, got %q", key, val)
		}
	case "enum":
		for _, e := range enum {
			if val == e {
				return nil
			}
		}
		return fmt.Errorf("key %q: value %q not in enum %v", key, val, enum)
	}
	return nil
}
