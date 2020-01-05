package zenity

import (
	"os/exec"
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

func message(dialog int, text string, options []Option) (bool, error) {
	opts := optsParse(options)

	data := osaMsg{
		Text:   text,
		Title:  opts.title,
		Dialog: opts.icon != 0 || dialog == 2,
	}

	if data.Dialog {
		switch opts.icon {
		case ErrorIcon:
			data.Icon = "stop"
		case InfoIcon, QuestionIcon:
			data.Icon = "note"
		case WarningIcon:
			data.Icon = "caution"
		}
	} else {
		switch dialog {
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

	if dialog != 2 {
		opts.cancel = ""
		if data.Dialog {
			opts.ok = "OK"
		}
	}
	if opts.ok != "" || opts.cancel != "" || opts.extra != "" || true {
		if opts.ok == "" {
			opts.ok = "OK"
		}
		if opts.cancel == "" {
			opts.cancel = "Cancel"
		}
		if dialog == 2 {
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
	}
	if opts.defcancel {
		if data.Cancel != 0 {
			data.Default = data.Cancel
		}
		if data.Dialog && data.Buttons == nil {
			data.Default = 1
		}
	}

	_, err := osaRun("msg", data)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

type osaMsg struct {
	Dialog  bool
	Text    string
	Message string
	As      string
	Title   string
	Icon    string
	Buttons []string
	Cancel  int
	Default int
}
