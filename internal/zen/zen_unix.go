// +build !windows,!darwin

package zen

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/ncruces/zenity/internal/cmd"
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
	if cmd.Command && path != "" {
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}
	return exec.Command(tool, args...).Output()
}
