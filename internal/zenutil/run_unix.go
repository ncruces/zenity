// +build !windows,!darwin

package zenutil

import (
	"os"
	"os/exec"
	"syscall"
)

var tool, path string

func init() {
	for _, tool = range [3]string{"qarma", "zenity", "matedialog"} {
		path, _ = exec.LookPath(tool)
		if path != "" {
			return
		}
	}
	tool = "zenity"
}

func Run(args []string) ([]byte, error) {
	if Command && path != "" {
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}
	return exec.Command(tool, args...).Output()
}
