//go:build darwin || dev

package main

import (
	"os"
	"syscall"
)

func killParent() {
	syscall.Kill(os.Getppid(), syscall.SIGHUP)
}
