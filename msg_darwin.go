package zenity

import (
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := applyOptions(options)
	data := zenutil.Msg{
		Text:    text,
		Timeout: zenutil.Timeout,
	}
	dialog := kind == questionKind || opts.icon != 0

	if dialog {
		data.Operation = "displayDialog"
		data.Title = opts.title

		switch opts.icon {
		case ErrorIcon:
			data.Icon = "stop"
		case WarningIcon:
			data.Icon = "caution"
		case InfoIcon, QuestionIcon:
			data.Icon = "note"
		}
	} else {
		data.Operation = "displayAlert"
		if opts.title != "" {
			data.Message = text
			data.Text = opts.title
		}

		switch kind {
		case infoKind:
			data.As = "informational"
		case warningKind:
			data.As = "warning"
		case errorKind:
			data.As = "critical"
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
				data.Buttons = []string{opts.cancelLabel, opts.okLabel}
				data.Default = 2
				data.Cancel = 1
			} else {
				data.Buttons = []string{opts.extraButton, opts.cancelLabel, opts.okLabel}
				data.Default = 3
				data.Cancel = 2
			}
		} else {
			if opts.extraButton == "" {
				data.Buttons = []string{opts.okLabel}
				data.Default = 1
			} else {
				data.Buttons = []string{opts.extraButton, opts.okLabel}
				data.Default = 2
			}
		}
		data.Extra = opts.extraButton
	}
	if opts.defaultCancel {
		if data.Cancel != 0 {
			data.Default = data.Cancel
		}
		if dialog && data.Buttons == nil {
			data.Default = 1
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
