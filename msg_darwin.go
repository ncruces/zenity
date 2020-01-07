package zenity

import (
	"os/exec"

	"github.com/ncruces/zenity/internal/osa"
)

func Error(text string, options ...Option) (bool, error) {
	return message(0, text, options)
}

func Info(text string, options ...Option) (bool, error) {
	return message(1, text, options)
}

func Question(text string, options ...Option) (bool, error) {
	return message(2, text, options)
}

func Warning(text string, options ...Option) (bool, error) {
	return message(3, text, options)
}

func message(typ int, text string, options []Option) (bool, error) {
	opts := optsParse(options)

	dialog := opts.icon != 0 || typ == 2
	var op string
	if dialog {
		op = "displayDialog"
	} else {
		op = "displayAlert"
	}

	data := osa.Msg{
		Operation: op,
		Text:      text,
		Title:     opts.title,
	}

	if dialog {
		switch opts.icon {
		case ErrorIcon:
			data.Icon = "stop"
		case InfoIcon, QuestionIcon:
			data.Icon = "note"
		case WarningIcon:
			data.Icon = "caution"
		}
	} else {
		switch typ {
		case 0:
			data.As = "critical"
		case 1:
			data.As = "informational"
		case 3:
			data.As = "warning"
		}

		if opts.title != "" {
			data.Text = opts.title
			data.Message = text
		}
	}

	if typ != 2 {
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
		if typ == 2 {
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

	out, err := osa.Run("msg", data)
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
