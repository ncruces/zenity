package zenity

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func quoteAccelerators(text string) string {
	return strings.ReplaceAll(text, "&", "&&")
}

func quoteMnemonics(text string) string {
	return strings.ReplaceAll(text, "_", "__")
}

func quoteMarkup(text string) string {
	var res strings.Builder
	err := xml.EscapeText(&res, []byte(text))
	if err != nil {
		return text
	}
	return res.String()
}

func appendGeneral(args []string, opts options) []string {
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if id, ok := opts.attach.(int); ok {
		args = append(args, "--attach", strconv.Itoa(id))
	}
	if opts.modal {
		args = append(args, "--modal")
	}
	if opts.display != "" {
		args = append(args, "--display", opts.display)
	}
	if opts.class != "" {
		args = append(args, "--class", opts.class)
	}
	if opts.name != "" {
		args = append(args, "--name", opts.name)
	}
	return args
}

func appendButtons(args []string, opts options) []string {
	if opts.okLabel != nil {
		args = append(args, "--ok-label", *opts.okLabel)
	}
	if opts.cancelLabel != nil {
		args = append(args, "--cancel-label", *opts.cancelLabel)
	}
	if opts.extraButton != nil {
		args = append(args, "--extra-button", *opts.extraButton)
	}
	return args
}

func appendWidthHeight(args []string, opts options) []string {
	if opts.width > 0 {
		args = append(args, "--width", strconv.FormatUint(uint64(opts.width), 10))
	}
	if opts.height > 0 {
		args = append(args, "--height", strconv.FormatUint(uint64(opts.height), 10))
	}
	return args
}

func appendWindowIcon(args []string, opts options) []string {
	switch opts.windowIcon {
	case ErrorIcon:
		args = append(args, "--window-icon=error")
	case WarningIcon:
		args = append(args, "--window-icon=warning")
	case InfoIcon:
		args = append(args, "--window-icon=info")
	case QuestionIcon:
		args = append(args, "--window-icon=question")
	}
	if i, ok := opts.windowIcon.(string); ok {
		args = append(args, "--window-icon", i)
	}
	return args
}

func strResult(opts options, out []byte, err error) (string, error) {
	out = bytes.TrimSuffix(out, []byte{'\n'})
	if eerr, ok := err.(*exec.ExitError); ok {
		if eerr.ExitCode() == 1 {
			if opts.extraButton != nil && *opts.extraButton == string(out) {
				return "", ErrExtraButton
			}
			return "", ErrCanceled
		}
		return "", fmt.Errorf("%w: %s", eerr, eerr.Stderr)
	}
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func lstResult(opts options, out []byte, err error) ([]string, error) {
	str, err := strResult(opts, out, err)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return []string{}, nil
	}
	return strings.Split(str, zenutil.Separator), nil
}

func pwdResult(sep string, opts options, out []byte, err error) (string, string, error) {
	str, err := strResult(opts, out, err)
	if opts.username {
		usr, pwd, _ := strings.Cut(str, sep)
		return usr, pwd, err
	}
	return "", str, err
}
