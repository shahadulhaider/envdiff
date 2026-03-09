package main

import (
	"errors"
	"os"
)

// handleExitError checks if err is an exitError and exits with its code.
// Returns true if it was an exit error.
func handleExitError(err error) bool {
	var exitErr *exitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.code)
		return true
	}
	return false
}
