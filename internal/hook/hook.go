package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	beginMarker = "# BEGIN envdiff"
	endMarker   = "# END envdiff"
	hookScript  = `# BEGIN envdiff
if [ -f .env.example ]; then
  envdiff check || exit 1
fi
if [ -f .env.schema.toml ]; then
  envdiff validate --schema .env.schema.toml .env || exit 1
fi
# END envdiff`
)

// hookPath returns the path to the pre-commit hook file.
func hookPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
}

// Install installs the envdiff pre-commit hook.
// If the hook file already exists, the envdiff section is appended.
// If it doesn't exist, a new hook file is created.
func Install(repoRoot string) error {
	path := hookPath(repoRoot)

	// Ensure hooks directory exists
	hooksDir := filepath.Dir(path)
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	var existing string
	data, err := os.ReadFile(path)
	if err == nil {
		existing = string(data)
		// Check if already installed
		if strings.Contains(existing, beginMarker) {
			return fmt.Errorf("envdiff hook already installed")
		}
	}

	var content string
	if existing == "" {
		content = "#!/bin/sh\n\n" + hookScript + "\n"
	} else {
		// Append to existing hook
		if !strings.HasSuffix(existing, "\n") {
			existing += "\n"
		}
		content = existing + "\n" + hookScript + "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	return nil
}

// Uninstall removes the envdiff section from the pre-commit hook.
func Uninstall(repoRoot string) error {
	path := hookPath(repoRoot)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no pre-commit hook found")
		}
		return fmt.Errorf("read hook: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, beginMarker) {
		return fmt.Errorf("envdiff hook not installed")
	}

	// Remove the envdiff section
	lines := strings.Split(content, "\n")
	var result []string
	inSection := false
	for _, line := range lines {
		if line == beginMarker {
			inSection = true
			continue
		}
		if line == endMarker {
			inSection = false
			continue
		}
		if !inSection {
			result = append(result, line)
		}
	}

	// Clean up extra blank lines
	newContent := strings.TrimRight(strings.Join(result, "\n"), "\n") + "\n"

	if err := os.WriteFile(path, []byte(newContent), 0755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	return nil
}

// Status returns true if the envdiff hook is installed.
func Status(repoRoot string) (bool, error) {
	path := hookPath(repoRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read hook: %w", err)
	}
	return strings.Contains(string(data), beginMarker), nil
}
