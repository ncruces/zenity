// +build !windows,!darwin

package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(options ...Option) (string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection"}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

func selectFileMutiple(options ...Option) ([]string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection", "--multiple", "--separator", zenutil.Separator}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return strings.Split(string(out), zenutil.Separator), nil
}

func selectFileSave(options ...Option) (string, error) {
	opts := applyOptions(options)

	args := []string{"--file-selection", "--save"}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	if opts.confirmOverwrite {
		args = append(args, "--confirm-overwrite")
	}
	args = append(args, initFilters(opts.fileFilters)...)

	out, err := zenutil.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return string(out), nil
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
