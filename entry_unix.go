// +build !windows,!darwin

package zenity

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func entry(text string, opts options) (string, bool, error) {
	args := []string{"--entry", "--text", text, "--entry-text", opts.entryText}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", strconv.FormatUint(uint64(opts.width), 10))
	}
	if opts.height > 0 {
		args = append(args, "--height", strconv.FormatUint(uint64(opts.height), 10))
	}
	if opts.okLabel != nil {
		args = append(args, "--ok-label", *opts.okLabel)
	}
	if opts.cancelLabel != nil {
		args = append(args, "--cancel-label", *opts.cancelLabel)
	}
	if opts.extraButton != nil {
		args = append(args, "--extra-button", *opts.extraButton)
	}
	if opts.hideText {
		args = append(args, "--hide-text")
	}
	switch opts.icon {
	case ErrorIcon:
		args = append(args, "--window-icon=error")
	case WarningIcon:
		args = append(args, "--window-icon=warning")
	case InfoIcon:
		args = append(args, "--window-icon=info")
	case QuestionIcon:
		args = append(args, "--window-icon=question")
	}

	out, err := zenutil.Run(opts.ctx, args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if len(out) > 0 && opts.extraButton != nil &&
			string(out[:len(out)-1]) == *opts.extraButton {
			return "", false, ErrExtraButton
		}
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(out[:len(out)-1]), true, nil
}

func password(opts options) (string, string, bool, error) {
	args := []string{"--password"}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.okLabel != nil {
		args = append(args, "--ok-label", *opts.okLabel)
	}
	if opts.cancelLabel != nil {
		args = append(args, "--cancel-label", *opts.cancelLabel)
	}
	if opts.extraButton != nil {
		args = append(args, "--extra-button", *opts.extraButton)
	}
	if opts.username {
		args = append(args, "--username")
	}

	out, err := zenutil.Run(opts.ctx, args)
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		if len(out) > 0 && opts.extraButton != nil &&
			string(out[:len(out)-1]) == *opts.extraButton {
			return "", "", false, ErrExtraButton
		}
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	if opts.username {
		if split := strings.SplitN(string(out[:len(out)-1]), "|", 2); len(split) == 2 {
			return split[0], split[1], true, nil
		}
	}
	return "", string(out[:len(out)-1]), true, nil
}
