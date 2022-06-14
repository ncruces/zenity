//go:build windows || darwin || dev

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/zencmd"
	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(opts ...zenity.Option) error {
	if !listen {
		if text == unspecified {
			return nil
		}
		return zenity.Notify(text, opts...)
	}

	zenutil.Command = false
	ico := zenity.NoIcon
	for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
		line := scanner.Text()
		var cmd, msg string
		if n := strings.IndexByte(line, ':'); n >= 0 {
			cmd = strings.TrimSpace(line[:n])
			msg = strings.TrimSpace(zencmd.Unescape(line[n+1:]))
		} else {
			fmt.Fprint(os.Stderr, "Could not parse command from stdin")
		}
		switch cmd {
		case "icon":
			switch msg {
			case "error", "dialog-error":
				ico = zenity.ErrorIcon
			case "info", "dialog-information":
				ico = zenity.InfoIcon
			case "question", "dialog-question":
				ico = zenity.QuestionIcon
			case "important", "warning", "dialog-warning":
				ico = zenity.WarningIcon
			case "dialog-password":
				ico = zenity.PasswordIcon
			default:
				ico = zenity.NoIcon
			}
		case "message", "tooltip":
			opts := []zenity.Option{ico}
			if n := strings.IndexByte(msg, '\n'); n >= 0 {
				opts = append(opts, zenity.Title(msg[:n]))
				msg = msg[n+1:]
			}
			if err := zenity.Notify(msg, opts...); err != nil {
				return err
			}
		case "visible", "hints":
			// ignored
		default:
			fmt.Fprintf(os.Stderr, "Unknown command %q", cmd)
		}
	}
	return nil
}
