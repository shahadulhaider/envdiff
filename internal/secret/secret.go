package secret

import (
	"math"
	"regexp"
	"strings"

	"github.com/shahadulhaider/envdiff/internal/env"
)

var (
	awsAccessKeyRe = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
	privateKeyRe   = regexp.MustCompile(`-----BEGIN (RSA |EC |)PRIVATE KEY-----`)
	dbURLRe        = regexp.MustCompile(`://[^:]+:[^@]+@`)
)

var secretKeyNames = []string{
	"secret", "password", "passwd", "token", "api_key", "apikey",
	"private", "credential", "auth", "access_key", "secret_key",
}

// IsSecret returns true if the key/value combination looks like a secret.
func IsSecret(key, value string) bool {
	lowerKey := strings.ToLower(key)
	for _, name := range secretKeyNames {
		if strings.Contains(lowerKey, name) {
			return true
		}
	}

	if awsAccessKeyRe.MatchString(value) {
		return true
	}
	if privateKeyRe.MatchString(value) {
		return true
	}
	if dbURLRe.MatchString(value) {
		return true
	}
	if len(value) > 20 && shannonEntropy(value) > 4.5 {
		return true
	}
	return false
}

// DetectSecrets returns entries whose values are likely secrets.
func DetectSecrets(entries []env.EnvEntry) []env.EnvEntry {
	var secrets []env.EnvEntry
	for _, e := range entries {
		if IsSecret(e.Key, e.Value) {
			secrets = append(secrets, e)
		}
	}
	return secrets
}

// MaskValue returns a fixed-length mask string.
func MaskValue(value string) string {
	return "****"
}

// shannonEntropy calculates the Shannon entropy of a string.
func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	freq := make(map[rune]float64)
	for _, c := range s {
		freq[c]++
	}
	n := float64(len(s))
	entropy := 0.0
	for _, count := range freq {
		p := count / n
		entropy -= p * math.Log2(p)
	}
	return entropy
}
