package zenity

import (
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := applyOptions(options)

	var data zenutil.Msg
	data.Text = text
	data.Options.Timeout = zenutil.Timeout

	dialog := kind == questionKind || opts.icon != 0

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
		if opts.title != "" {
			data.Options.Message = text
			data.Text = opts.title
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

	if kind != questionKind {
		if dialog {
			opts.okLabel = "OK"
		}
		opts.cancelLabel = ""
	}
	if opts.okLabel != "" || opts.cancelLabel != "" || opts.extraButton != "" {
		if opts.okLabel == "" {
			opts.okLabel = "OK"
		}
		if opts.cancelLabel == "" {
			opts.cancelLabel = "Cancel"
		}
		if kind == questionKind {
			if opts.extraButton == "" {
				data.Options.Buttons = []string{opts.cancelLabel, opts.okLabel}
				data.Options.Default = 2
				data.Options.Cancel = 1
			} else {
				data.Options.Buttons = []string{opts.extraButton, opts.cancelLabel, opts.okLabel}
				data.Options.Default = 3
				data.Options.Cancel = 2
			}
		} else {
			if opts.extraButton == "" {
				data.Options.Buttons = []string{opts.okLabel}
				data.Options.Default = 1
			} else {
				data.Options.Buttons = []string{opts.extraButton, opts.okLabel}
				data.Options.Default = 2
			}
		}
		data.Extra = opts.extraButton
	}
	if opts.defaultCancel {
		if data.Options.Cancel != 0 {
			data.Options.Default = data.Options.Cancel
		}
		if dialog && data.Options.Buttons == nil {
			data.Options.Default = 1
		}
	}

	out, err := zenutil.Run(opts.ctx, "msg", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(out) > 0 && string(out[:len(out)-1]) == opts.extraButton {
		return false, ErrExtraButton
	}
	return true, err
}
