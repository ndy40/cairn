package updater

import "os/exec"

// detach is a no-op on Windows, where the desktop app is not updated.
func detach(*exec.Cmd) {}
