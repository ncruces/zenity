// +build !windows,!darwin

package zenity

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(options []Option) (string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection"}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", fmt.Sprint(opts.width))
	}
	if opts.height > 0 {
		args = append(args, "--height", fmt.Sprint(opts.height))
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)

	if err == nil {
		if len(out) > 0 {
			out = out[:len(out)-1]
		}
		return string(out), nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		switch err.ExitCode() {
		case 1:
			return "", ErrCancelOrClosed
		}
	}

	return "", err
}

func selectFileMutiple(options []Option) ([]string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection", "--multiple", "--separator", zenutil.Separator}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", fmt.Sprint(opts.width))
	}
	if opts.height > 0 {
		args = append(args, "--height", fmt.Sprint(opts.height))
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)

	if err == nil {
		if len(out) > 0 {
			out = out[:len(out)-1]
		}
		return strings.Split(string(out), zenutil.Separator), nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		switch err.ExitCode() {
		case 1:
			return nil, ErrCancelOrClosed
		}
	}

	return nil, err
}

func selectFileSave(options []Option) (string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection", "--save"}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", fmt.Sprint(opts.width))
	}
	if opts.height > 0 {
		args = append(args, "--height", fmt.Sprint(opts.height))
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	if opts.confirmOverwrite {
		args = append(args, "--confirm-overwrite")
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)

	if err == nil {
		if len(out) > 0 {
			out = out[:len(out)-1]
		}
		return string(out), nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		switch err.ExitCode() {
		case 1:
			return "", ErrCancelOrClosed
		}
	}

	return "", err
}

func initFilters(filters []FileFilter) []string {
	var res []string
	for _, f := range filters {
		var buf strings.Builder
		buf.WriteString("--file-filter=")
		if f.Name != "" {
			buf.WriteString(f.Name)
			buf.WriteRune('|')
		}
		for _, p := range f.Patterns {
			buf.WriteString(p)
			buf.WriteRune(' ')
		}
		res = append(res, buf.String())
	}
	return res
}
