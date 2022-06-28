//go:build darwin || dev

package zencmd

import (
	"os"
	"syscall"
)

// KillParent is internal.
func KillParent() {
	syscall.Kill(os.Getppid(), syscall.SIGHUP)
}
