// +build !windows,!darwin

package zenity

import (
	"fmt"
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
	if opts.width > 0 {
		args = append(args, "--width", fmt.Sprintf("%d", opts.width))
	}
	if opts.height > 0 {
		args = append(args, "--height", fmt.Sprintf("%d", opts.height))
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
		args = append(args, "--window-icon=error", "--icon-name=dialog-error")
	case WarningIcon:
		args = append(args, "--window-icon=warning", "--icon-name=dialog-warning")
	case InfoIcon:
		args = append(args, "--window-icon=info", "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--window-icon=question", "--icon-name=dialog-question")
	}

	out, err := zenutil.Run(opts.ctx, args)

	if err == nil {
		return true, nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		switch err.ExitCode() {
		case 1:
			if len(out) > 0 && string(out[:len(out)-1]) == opts.extraButton {
				return false, ErrExtraButton
			}

			return false, ErrCancelOrClosed
		}
	}

	return false, err
}
