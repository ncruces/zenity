package zenity

import (
	"path/filepath"

	"github.com/ncruces/zenity/internal/zenutil"
)

func selectFile(opts options) (name string, err error) {
	var data zenutil.File
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, _, err = splitDirAndName(opts.filename)
	if data.Options.Location != "" && err == nil {
		data.Options.Location, err = filepath.Abs(data.Options.Location)
	}
	if err != nil {
		return "", err
	}
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Options.Type = opts.fileFilters.types()
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
	return strResult(opts, out, err)
}

func selectFileMultiple(opts options) (list []string, err error) {
	var data zenutil.File
	data.Separator = zenutil.Separator
	data.Options.Multiple = true
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, _, err = splitDirAndName(opts.filename)
	if data.Options.Location != "" && err == nil {
		data.Options.Location, err = filepath.Abs(data.Options.Location)
	}
	if err != nil {
		return nil, err
	}
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFile"
		data.Options.Type = opts.fileFilters.types()
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
	return lstResult(opts, out, err)
}

func selectFileSave(opts options) (name string, err error) {
	var data zenutil.File
	data.Options.Prompt = opts.title
	data.Options.Invisibles = opts.showHidden
	data.Options.Location, data.Options.Name, err = splitDirAndName(opts.filename)
	if data.Options.Location != "" && err == nil {
		data.Options.Location, err = filepath.Abs(data.Options.Location)
	}
	if err != nil {
		return "", err
	}
	if opts.attach != nil {
		data.Application = opts.attach
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	if opts.directory {
		data.Operation = "chooseFolder"
	} else {
		data.Operation = "chooseFileName"
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
	return strResult(opts, out, err)
}
