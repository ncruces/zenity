// +build !windows,!darwin

package zenity

import (
	"os/exec"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, options []Option) (bool, error) {
	opts := applyOptions(options)

	var args []string
	switch kind {
	case questionKind:
		args = append(args, "--question")
	case infoKind:
		args = append(args, "--info")
	case warningKind:
		args = append(args, "--warning")
	case errorKind:
		args = append(args, "--error")
	}
	if text != "" {
		args = append(args, "--text", text, "--no-markup")
	}
	if opts.title != "" {
		args = append(args, "--title", opts.title)
	}
	if opts.okLabel != "" {
		args = append(args, "--ok-label", opts.okLabel)
	}
	if opts.cancelLabel != "" {
		args = append(args, "--cancel-label", opts.cancelLabel)
	}
	if opts.extraButton != "" {
		args = append(args, "--extra-button", opts.extraButton)
	}
	if opts.noWrap {
		args = append(args, "--no-wrap")
	}
	if opts.ellipsize {
		args = append(args, "--ellipsize")
	}
	if opts.defaultCancel {
		args = append(args, "--default-cancel")
	}
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--icon-name=dialog-error")
	case WarningIcon:
		args = append(args, "--icon-name=dialog-warning")
	case InfoIcon:
		args = append(args, "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--icon-name=dialog-question")
	}

	out, err := zenutil.Run(args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 255 {
		if len(out) > 0 && string(out[:len(out)-1]) == opts.extraButton {
			return false, ErrExtraButton
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
