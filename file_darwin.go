package zenity

import (
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(options []Option) (string, error) {
	opts := applyOptions(options)

	var data zenutil.File
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, _ = splitDirAndName(opts.filename)

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Options.Type = initFilters(opts.fileFilters)
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
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

func selectFileMutiple(options []Option) ([]string, error) {
	opts := applyOptions(options)

	var data zenutil.File
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, _ = splitDirAndName(opts.filename)
	data.Options.Multiple = true
	data.Separator = zenutil.Separator

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Options.Type = initFilters(opts.fileFilters)
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
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
	return strings.Split(string(out), zenutil.Separator), nil
}

func selectFileSave(options []Option) (string, error) {
	opts := applyOptions(options)

	var data zenutil.File
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, data.Options.Name = splitDirAndName(opts.filename)

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFileName"
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
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
