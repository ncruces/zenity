// +build !windows,!darwin

package zenutil

import (
	"context"
	"os"
	"os/exec"
	"strconv"
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

func Run(ctx context.Context, args []string) ([]byte, error) {
	if Command && path != "" {
		if Timeout > 0 {
			args = append(args, "--timeout", strconv.Itoa(Timeout))
		}
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}

	if ctx != nil {
		return exec.CommandContext(ctx, tool, args...).Output()
	}
	return exec.Command(tool, args...).Output()
}
