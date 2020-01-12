package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/cmd"
	"github.com/ncruces/zenity/internal/osa"
)

func SelectFile(options ...Option) (string, error) {
	opts := optsParse(options)

	data := osa.File{
		Prompt:     opts.title,
		Invisibles: opts.hidden,
	}
	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Type = initFilters(opts.filters)
	}
	data.Location, _ = splitDirAndName(opts.filename)

	out, err := osa.Run("file", data)
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

	data := osa.File{
		Prompt:     opts.title,
		Invisibles: opts.hidden,
		Multiple:   true,
		Separator:  cmd.Separator,
	}
	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Type = initFilters(opts.filters)
	}
	data.Location, _ = splitDirAndName(opts.filename)

	out, err := osa.Run("file", data)
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

	data := osa.File{
		Prompt: opts.title,
	}
	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFileName"
		data.Type = initFilters(opts.filters)
	}
	data.Location, data.Name = splitDirAndName(opts.filename)

	out, err := osa.Run("file", data)
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

func initFilters(filters []FileFilter) []string {
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
