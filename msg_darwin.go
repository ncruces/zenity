package zenity

import (
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

		switch opts.icon {
		case ErrorIcon:
			data.Options.Icon = "stop"
		case WarningIcon:
			data.Options.Icon = "caution"
		case InfoIcon, QuestionIcon:
			data.Options.Icon = "note"
		}
	} else {
		data.Operation = "displayAlert"
		if opts.title != nil {
			data.Text = *opts.title
			data.Options.Message = text
		}

		switch kind {
		case infoKind:
			data.Options.As = "informational"
		case warningKind:
			data.Options.As = "warning"
		case errorKind:
			data.Options.As = "critical"
		}
	}

	if kind == questionKind {
		// alert defaults to a single button, we need two
		if opts.cancelLabel == nil && !dialog {
			opts.cancelLabel = stringPtr("Cancel")
		}
	} else {
		// dialog defaults to two buttons, we need one
		if opts.okLabel == nil && dialog {
			opts.okLabel = stringPtr("OK")
		}
		// only questions have cancel
		opts.cancelLabel = nil
	}

	if opts.okLabel != nil || opts.cancelLabel != nil || opts.extraButton != nil {
		if opts.okLabel == nil {
			opts.okLabel = stringPtr("OK")
		}
		if kind == questionKind {
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
		} else {
			if opts.extraButton == nil {
				data.Options.Buttons = []string{*opts.okLabel}
				data.Options.Default = 1
			} else {
				data.Options.Buttons = []string{*opts.extraButton, *opts.okLabel}
				data.Options.Default = 2
			}
		}
		data.Extra = opts.extraButton
	}

	if kind == questionKind && opts.defaultCancel {
		if data.Options.Cancel != 0 {
			data.Options.Default = data.Options.Cancel
		} else {
			data.Options.Default = 1
		}
	}

	out, err := zenutil.Run(opts.ctx, "dialog", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if len(out) > 0 && opts.extraButton != nil &&
			string(out[:len(out)-1]) == *opts.extraButton {
			return false, ErrExtraButton
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
