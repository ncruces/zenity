package zenity

import (
	"bytes"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func entry(text string, opts options) (string, bool, error) {
	var data zenutil.Dialog
	data.Text = text
	data.Operation = "displayDialog"
	data.Options.Title = opts.title
	data.Options.Answer = &opts.entryText
	data.Options.Hidden = opts.hideText
	data.Options.Timeout = zenutil.Timeout
	data.Options.Icon = opts.icon.String()
	data.SetButtons(getButtons(true, true, opts))

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	out = bytes.TrimSuffix(out, []byte{'\n'})
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if opts.extraButton != nil &&
			*opts.extraButton == string(out) {
			return "", false, ErrExtraButton
		}
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(out), true, nil
}

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	pass, ok, err := entry("Password:", opts)
	return "", pass, ok, err
}
