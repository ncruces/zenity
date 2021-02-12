// +build !windows,!darwin

package zenity

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func getChoices(output string) (indices []uint) {
	out := strings.Split(strings.TrimSpace(output), "|")

	indices = make([]uint, 0, len(out))

	for i := range out {
		idx, err := strconv.ParseUint(out[i], 10, 32)
		if err == nil {
			indices = append(indices, uint(idx))
		}
	}

	return
}

func list(text string, choices []string, options []Option) ([]uint, error) {
	opts := applyOptions(options)
	args := []string{"--list", "--hide-header", "--column=id", "--column=", "--hide-column=1", "--print-column=1"}

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
		args = append(args, "--window-icon=error", "--icon-name=dialog-error")
	case WarningIcon:
		args = append(args, "--window-icon=warning", "--icon-name=dialog-warning")
	case InfoIcon:
		args = append(args, "--window-icon=info", "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--window-icon=question", "--icon-name=dialog-question")
	}

	// List options
	if opts.multipleSelection {
		args = append(args, "--multiple")
		args = append(args, "--separator=|") // Dont use zenutil.Separator, because we return indices.
	}

	// Add choices with ids.
	for i := range choices {
		args = append(args, fmt.Sprintf("%d", i), choices[i])
	}

	out, err := zenutil.Run(opts.ctx, args)
	if err == nil {
		return getChoices(string(out)), nil
	}

	if err, ok := err.(*exec.ExitError); ok {
		switch err.ExitCode() {
		case 1:
			return nil, ErrCancelOrClosed
		}
	}

	return nil, err
}
