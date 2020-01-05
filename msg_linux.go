package zenity

import "os/exec"

func Error(text string, options ...Option) (bool, error) {
	return message("--error", text, options)
}

func Info(text string, options ...Option) (bool, error) {
	return message("--info", text, options)
}

func Question(text string, options ...Option) (bool, error) {
	return message("--question", text, options)
}

func Warning(text string, options ...Option) (bool, error) {
	return message("--warning", text, options)
}

func message(arg, text string, options []Option) (bool, error) {
	opts := optsParse(options)

	args := []string{arg, "--text", text, "--no-markup"}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.ok != "" {
		args = append(args, "--ok-label", opts.ok)
	}
	if opts.cancel != "" {
		args = append(args, "--cancel-label", opts.cancel)
	}
	if opts.extra != "" {
		args = append(args, "--extra-button", opts.extra)
	}
	if opts.nowrap {
		args = append(args, "--no-wrap")
	}
	if opts.ellipsize {
		args = append(args, "--ellipsize")
	}
	if opts.defcancel {
		args = append(args, "--default-cancel")
	}
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--icon-name=dialog-error")
	case InfoIcon:
		args = append(args, "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--icon-name=dialog-question")
	case WarningIcon:
		args = append(args, "--icon-name=dialog-warning")
	}

	cmd := exec.Command("zenity", args...)
	_, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
