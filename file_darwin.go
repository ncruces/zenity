package zenity

import (
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
		data.Options.Type = opts.fileFilters.types()
	}

	out, err := zenutil.Run(opts.ctx, "file", data)
	return strResult(opts, out, err)
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
		data.Options.Type = opts.fileFilters.types()
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
	return strResult(opts, out, err)
}
