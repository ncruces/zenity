//go:build !windows && !darwin

package zenutil

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
)

var (
	tool, path string
	pathOnce   sync.Once
)

func initPath() {
	for _, tool = range [3]string{"qarma", "zenity", "matedialog"} {
		path, _ = exec.LookPath(tool)
		if path != "" {
			return
		}
	}
	tool = "zenity"
}

// IsAvailable is internal.
func IsAvailable() bool {
	pathOnce.Do(initPath)
	return path != ""
}

// Run is internal.
func Run(ctx context.Context, args []string) ([]byte, error) {
	pathOnce.Do(initPath)
	if Command && path != "" {
		if Timeout > 0 {
			args = append(args, "--timeout", strconv.Itoa(Timeout))
		}
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}

	if ctx != nil {
		out, err := exec.CommandContext(ctx, tool, args...).Output()
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return out, err
	}
	return exec.Command(tool, args...).Output()
}

// RunProgress is internal.
func RunProgress(ctx context.Context, max int, close bool, extra *string, args []string) (*progressDialog, error) {
	pathOnce.Do(initPath)
	if Command && path != "" {
		if Timeout > 0 {
			args = append(args, "--timeout", strconv.Itoa(Timeout))
		}
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}
	if ctx == nil {
		ctx = context.Background()
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	var out *bytes.Buffer
	if extra != nil {
		out = &bytes.Buffer{}
		cmd.Stdout = out
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	dlg := &progressDialog{
		ctx:     ctx,
		cmd:     cmd,
		max:     max,
		percent: true,
		close:   close,
		lines:   make(chan string),
		done:    make(chan struct{}),
	}
	go dlg.pipe(pipe)
	go dlg.wait(extra, out)
	return dlg, nil
}
