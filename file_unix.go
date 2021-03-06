// +build !windows,!darwin

package zenity

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(opts options) (string, error) {
	args := []string{"--file-selection"}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSuffix(out, []byte{'\n'})), nil
}

func selectFileMutiple(opts options) ([]string, error) {
	args := []string{"--file-selection", "--multiple", "--separator", zenutil.Separator}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out = bytes.TrimSuffix(out, []byte{'\n'})
	if len(out) == 0 {
		return nil, nil
	}
	return strings.Split(string(out), zenutil.Separator), nil
}

func selectFileSave(opts options) (string, error) {
	args := []string{"--file-selection", "--save"}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.confirmOverwrite {
		args = append(args, "--confirm-overwrite")
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(opts.ctx, args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSuffix(out, []byte{'\n'})), nil
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
