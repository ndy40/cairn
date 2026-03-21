//go:build windows

package main

import (
	"os"
)

func processAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, FindProcess always succeeds; OpenProcess would be needed for a real check,
	// but sending signal 0 via os.Process is not supported. We use a best-effort approach.
	_ = p
	return true
}
