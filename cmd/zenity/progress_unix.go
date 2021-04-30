// +build !windows,!js

package main

import (
	"os"
	"syscall"
)

func killParent() {
	syscall.Kill(os.Getppid(), syscall.SIGHUP)
}
