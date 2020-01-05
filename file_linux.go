package zenity

import (
	"os/exec"
	"strings"
)

func SelectFile(options ...Option) (string, error) {
	opts := optsParse(options)

	args := []string{"--file-selection"}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, zenityFilters(opts.filters)...)
	cmd := exec.Command("zenity", args...)
	out, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
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

func SelectFileMutiple(options ...Option) ([]string, error) {
	opts := optsParse(options)

	args := []string{"--file-selection", "--multiple", "--separator=\x1e"}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	args = append(args, zenityFilters(opts.filters)...)
	cmd := exec.Command("zenity", args...)
	out, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return strings.Split(string(out), "\x1e"), nil
}

func SelectFileSave(options ...Option) (string, error) {
	opts := optsParse(options)

	args := []string{"--file-selection", "--save"}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	if opts.overwrite {
		args = append(args, "--confirm-overwrite")
	}
	args = append(args, zenityFilters(opts.filters)...)
	cmd := exec.Command("zenity", args...)
	out, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
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

func SelectDirectory(options ...Option) (string, error) {
	opts := optsParse(options)

	args := []string{"--file-selection", "--directory"}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.filename != "" {
		args = append(args, "--filename", opts.filename)
	}
	cmd := exec.Command("zenity", args...)
	out, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
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

func zenityFilters(filters []FileFilter) []string {
	var res []string
	for _, f := range filters {
		var buf strings.Builder
		buf.WriteString("--file-filter=")
		if f.Name != "" {
			buf.WriteString(f.Name)
			buf.WriteRune('|')
		}
		for _, e := range f.Exts {
			buf.WriteRune('*')
			buf.WriteString(e)
			buf.WriteRune(' ')
		}
		res = append(res, buf.String())
	}
	return res
}
