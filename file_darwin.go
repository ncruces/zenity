package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/cmd"
	"github.com/ncruces/zenity/internal/osa"
)

func SelectFile(options ...Option) (string, error) {
	opts := optsParse(options)
	dir, _ := splitDirAndName(opts.filename)
	out, err := osa.Run("file", osa.File{
		Operation: "chooseFile",
		Prompt:    opts.title,
		Type:      appleFilters(opts.filters),
		Location:  dir,
	})
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
	dir, _ := splitDirAndName(opts.filename)
	out, err := osa.Run("file", osa.File{
		Operation: "chooseFile",
		Multiple:  true,
		Prompt:    opts.title,
		Separator: cmd.Separator,
		Type:      appleFilters(opts.filters),
		Location:  dir,
	})
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	if len(out) == 0 {
		return nil, nil
	}
	return strings.Split(string(out), cmd.Separator), nil
}

func SelectFileSave(options ...Option) (string, error) {
	opts := optsParse(options)
	dir, name := splitDirAndName(opts.filename)
	out, err := osa.Run("file", osa.File{
		Operation: "chooseFileName",
		Prompt:    opts.title,
		Location:  dir,
		Name:      name,
	})
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
	dir, _ := splitDirAndName(opts.filename)
	out, err := osa.Run("file", osa.File{
		Operation: "chooseFolder",
		Prompt:    opts.title,
		Location:  dir,
	})
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

func appleFilters(filters []FileFilter) []string {
	var filter []string
	for _, f := range filters {
		for _, p := range f.Patterns {
			star := strings.LastIndexByte(p, '*')
			if star >= 0 {
				dot := strings.LastIndexByte(p, '.')
				if star > dot {
					return nil
				}
				filter = append(filter, p[dot+1:])
			} else {
				filter = append(filter, p)
			}
		}
	}
	return filter
}
