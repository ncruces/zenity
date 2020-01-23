package zenity

import (
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := optsParse(options)
	data := zenutil.Msg{Text: text}
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
			opts.ok = "OK"
		}
		opts.cancel = ""
	}
	if opts.ok != "" || opts.cancel != "" || opts.extra != "" {
		if opts.ok == "" {
			opts.ok = "OK"
		}
		if opts.cancel == "" {
			opts.cancel = "Cancel"
		}
		if kind == questionKind {
			if opts.extra == "" {
				data.Buttons = []string{opts.cancel, opts.ok}
				data.Default = 2
				data.Cancel = 1
			} else {
				data.Buttons = []string{opts.extra, opts.cancel, opts.ok}
				data.Default = 3
				data.Cancel = 2
			}
		} else {
			if opts.extra == "" {
				data.Buttons = []string{opts.ok}
				data.Default = 1
			} else {
				data.Buttons = []string{opts.extra, opts.ok}
				data.Default = 2
			}
		}
		data.Extra = opts.extra
	}
	if opts.defcancel {
		if data.Cancel != 0 {
			data.Default = data.Cancel
		}
		if dialog && data.Buttons == nil {
			data.Default = 1
		}
	}

	out, err := zenutil.Run("msg", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(out) > 0 && string(out[:len(out)-1]) == opts.extra {
		return false, ErrExtraButton
	}
	return true, err
}
