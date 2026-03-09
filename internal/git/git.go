package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// IsGitRepo returns true if the current directory is inside a git repository.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

// ShowFileAtRef returns the content of a file at a given git ref.
// Returns ("", nil) if the file does not exist at that ref.
func ShowFileAtRef(ref, path string) (string, error) {
	arg := ref + ":" + path
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd := exec.Command("git", "show", arg)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		// If the file doesn't exist at this ref, return empty (not an error)
		errStr := errBuf.String()
		if strings.Contains(errStr, "does not exist") ||
			strings.Contains(errStr, "exists on disk") ||
			strings.Contains(errStr, "Path") ||
			strings.Contains(errStr, "unknown revision") {
			return "", nil
		}
		return "", err
	}
	return out.String(), nil
}
