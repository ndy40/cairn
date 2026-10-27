//go:build !windows

package updater

import (
	"os/exec"
	"syscall"
)

// detach starts cmd in its own session so it survives the terminal closing.
func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
