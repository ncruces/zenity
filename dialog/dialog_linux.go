package dialog

import (
	"os/exec"
	"strings"
)

func OpenFile(title, defaultPath string, filters []FileFilter) (string, error) {
	args := []string{"--file-selection"}
	if title != "" {
		args = append(args, "--title="+title)
	}
	if defaultPath != "" {
		args = append(args, "--filename="+defaultPath)
	}
	args = append(args, zenityFilters(filters)...)
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

func OpenFiles(title, defaultPath string, filters []FileFilter) ([]string, error) {
	args := []string{"--file-selection", "--multiple", "--separator=\x1e"}
	if title != "" {
		args = append(args, "--title="+title)
	}
	if defaultPath != "" {
		args = append(args, "--filename="+defaultPath)
	}
	args = append(args, zenityFilters(filters)...)
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

func SaveFile(title, defaultPath string, confirmOverwrite bool, filters []FileFilter) (string, error) {
	args := []string{"--file-selection", "--save"}
	if title != "" {
		args = append(args, "--title="+title)
	}
	if defaultPath != "" {
		args = append(args, "--filename="+defaultPath)
	}
	if confirmOverwrite {
		args = append(args, "--confirm-overwrite")
	}
	args = append(args, zenityFilters(filters)...)
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

func PickFolder(title, defaultPath string) (string, error) {
	args := []string{"--file-selection", "--directory"}
	if title != "" {
		args = append(args, "--title="+title)
	}
	if defaultPath != "" {
		args = append(args, "--filename="+defaultPath)
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
