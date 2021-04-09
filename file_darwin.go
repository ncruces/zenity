package zenity

import (
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(opts options) (string, error) {
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
	str, _, err := strResult(opts, out, err)
	return str, err
}

func selectFileMutiple(opts options) ([]string, error) {
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
	return lstResult(opts, out, err)
}

func selectFileSave(opts options) (string, error) {
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
	str, _, err := strResult(opts, out, err)
	return str, err
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
