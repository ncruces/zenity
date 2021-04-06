package zenity

import (
	"bytes"
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, opts options) (bool, error) {
	var data zenutil.Dialog
	data.Text = text
	data.Options.Timeout = zenutil.Timeout

	// dialog is more flexible, alert prettier
	var dialog bool
	if opts.icon != 0 { // use if we want to show a specific icon
		dialog = true
	} else if kind == questionKind && opts.cancelLabel == nil { // use for questions with default buttons
		dialog = true
	}

	if dialog {
		data.Operation = "displayDialog"
		data.Options.Title = opts.title
		data.Options.Icon = opts.icon.String()
	} else {
		data.Operation = "displayAlert"
		data.Options.As = kind.String()
		if opts.title != nil {
			data.Text = *opts.title
			data.Options.Message = text
		}
	}

	data.SetButtons(getButtons(dialog, kind == questionKind, opts))

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if opts.extraButton != nil &&
			*opts.extraButton == string(bytes.TrimSuffix(out, []byte{'\n'})) {
			return false, ErrExtraButton
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
