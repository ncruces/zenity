package zenity

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func list(text string, items []string, opts options) (string, bool, error) {
	var data zenutil.List
	data.Items = items
	data.Options.Prompt = &text
	data.Options.Title = opts.title
	data.Options.OK = opts.okLabel
	data.Options.Cancel = opts.cancelLabel
	data.Options.Default = opts.defaultItems
	data.Options.Empty = !opts.disallowEmpty

	out, err := zenutil.Run(opts.ctx, "list", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(bytes.TrimSuffix(out, []byte{'\n'})), true, nil
}

func listMultiple(text string, items []string, opts options) ([]string, error) {
	var data zenutil.List
	data.Items = items
	data.Options.Prompt = &text
	data.Options.Title = opts.title
	data.Options.OK = opts.okLabel
	data.Options.Cancel = opts.cancelLabel
	data.Options.Default = opts.defaultItems
	data.Options.Empty = !opts.disallowEmpty
	data.Options.Multiple = true
	data.Separator = zenutil.Separator

	out, err := zenutil.Run(opts.ctx, "list", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out = bytes.TrimSuffix(out, []byte{'\n'})
	if len(out) == 0 {
		return nil, nil
	}
	return strings.Split(string(out), zenutil.Separator), nil
}
