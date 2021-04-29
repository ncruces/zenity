// +build !windows,!darwin,!js

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

// Run is internal.
func Run(ctx context.Context, args []string) ([]byte, error) {
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
func RunProgress(ctx context.Context, max int, args []string) (*progressDialog, error) {
	if Command && path != "" {
		if Timeout > 0 {
			args = append(args, "--timeout", strconv.Itoa(Timeout))
		}
		syscall.Exec(path, append([]string{tool}, args...), os.Environ())
	}

	cmd := exec.Command(tool, args...)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	dlg := &progressDialog{
		done:    make(chan struct{}),
		lines:   make(chan string),
		percent: true,
		max:     max,
	}
	go func() {
		err := cmd.Wait()
		select {
		case _, ok := <-dlg.lines:
			if !ok {
				err = nil
			}
		default:
		}
		if cerr := ctx.Err(); cerr != nil {
			err = cerr
		}
		dlg.err = err
		close(dlg.done)
	}()
	go func() {
		defer cmd.Process.Signal(syscall.SIGTERM)
		for {
			var line string
			select {
			case s, ok := <-dlg.lines:
				if !ok {
					return
				}
				line = s
			case <-ctx.Done():
				return
			}
			if _, err := pipe.Write([]byte(line + "\n")); err != nil {
				return
			}
		}
	}()
	return dlg, nil
}
