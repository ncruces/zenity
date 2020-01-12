// +build !windows,!darwin

package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/cmd"
	"github.com/ncruces/zenity/internal/zen"
)

// Display file selection dialog.
//
// Returns an empty string on cancel.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
func SelectFile(options ...Option) (string, error) {
	opts := optsParse(options)

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
	args = append(args, initFilters(opts.filters)...)

	out, err := zen.Run(args)
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

// Display multiple file selection dialog.
//
// Returns a nil slice on cancel.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
func SelectFileMutiple(options ...Option) ([]string, error) {
	opts := optsParse(options)

	args := []string{"--file-selection", "--multiple", "--separator", cmd.Separator}
	if opts.directory {
		args = append(args, "--directory")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, initFilters(opts.filters)...)

	out, err := zen.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return strings.Split(string(out), cmd.Separator), nil
}

// Display save file selection dialog.
//
// Returns an empty string on cancel.
//
// Valid options: Title, Filename, ConfirmOverwrite, ConfirmCreate, ShowHidden, FileFilter(s).
func SelectFileSave(options ...Option) (string, error) {
	opts := optsParse(options)

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
	if opts.overwrite {
		args = append(args, "--confirm-overwrite")
	}
	args = append(args, initFilters(opts.filters)...)

	out, err := zen.Run(args)
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
