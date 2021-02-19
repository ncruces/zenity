package zenutil

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Run is internal.
func Run(ctx context.Context, script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	script = buf.String()
	if Command {
		path, err := exec.LookPath("osascript")
		if err == nil {
			os.Stderr.Close()
			syscall.Exec(path, []string{"osascript", "-l", "JavaScript", "-e", script}, nil)
		}
	}

	if ctx != nil {
		cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript")
		cmd.Stdin = strings.NewReader(script)
		out, err := cmd.Output()
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return out, err
	}
	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = strings.NewReader(script)
	return cmd.Output()
}

// File is internal.
type File struct {
	Operation string
	Separator string
	Options   FileOptions
}

// FileOptions is internal.
type FileOptions struct {
	Prompt     string   `json:"withPrompt,omitempty"`
	Type       []string `json:"ofType,omitempty"`
	Name       string   `json:"defaultName,omitempty"`
	Location   string   `json:"defaultLocation,omitempty"`
	Multiple   bool     `json:"multipleSelectionsAllowed,omitempty"`
	Invisibles bool     `json:"invisibles,omitempty"`
}

// Msg is internal.
type Msg struct {
	Operation string
	Text      string
	Extra     string
	Options   MsgOptions
}

// MsgOptions is internal.
type MsgOptions struct {
	Message string   `json:"message,omitempty"`
	As      string   `json:"as,omitempty"`
	Title   string   `json:"withTitle,omitempty"`
	Icon    string   `json:"withIcon,omitempty"`
	Buttons []string `json:"buttons,omitempty"`
	Cancel  int      `json:"cancelButton,omitempty"`
	Default int      `json:"defaultButton,omitempty"`
	Timeout int      `json:"givingUpAfter,omitempty"`
}

// Notify is internal.
type Notify struct {
	Text    string
	Options NotifyOptions
}

// NotifyOptions is internal.
type NotifyOptions struct {
	Title    string `json:"withTitle,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}
