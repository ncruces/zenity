package zenity

import (
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

	switch opts.icon {
	case ErrorIcon:
		data.Options.Icon = "stop"
	case WarningIcon:
		data.Options.Icon = "caution"
	case InfoIcon, QuestionIcon:
		data.Options.Icon = "note"
	}

	if opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil {
		if opts.okLabel == nil {
			opts.okLabel = stringPtr("OK")
		}
		if opts.cancelLabel == nil {
			opts.cancelLabel = stringPtr("Cancel")
		}
		if opts.extraButton == nil {
			data.Options.Buttons = []string{*opts.cancelLabel, *opts.okLabel}
			data.Options.Default = 2
			data.Options.Cancel = 1
		} else {
			data.Options.Buttons = []string{*opts.extraButton, *opts.cancelLabel, *opts.okLabel}
			data.Options.Default = 3
			data.Options.Cancel = 2
		}
		data.Extra = opts.extraButton
	}

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if len(out) > 0 && opts.extraButton != nil &&
			string(out[:len(out)-1]) == *opts.extraButton {
			return "", false, ErrExtraButton
		}
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(out[:len(out)-1]), true, nil
}

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	pass, ok, err := entry("Type your password", opts)
	return "", pass, ok, err
}
