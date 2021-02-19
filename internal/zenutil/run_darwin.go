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

type File struct {
	Operation  string
	Prompt     string
	Name       string
	Location   string
	Separator  string
	Type       []string
	Invisibles bool
	Multiple   bool
}

type Msg struct {
	Operation string
	Text      string
	Message   string
	As        string
	Title     string
	Icon      string
	Extra     string
	Buttons   []string
	Cancel    int
	Default   int
	Timeout   int
}

type Notify struct {
	Text     string
	Title    string
	Subtitle string
}
