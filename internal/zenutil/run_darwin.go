package zenutil

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Run(script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	script = buf.String()
	lang := "AppleScript"
	if strings.HasPrefix(script, "var app") {
		lang = "JavaScript"
	}

	if Command {
		path, err := exec.LookPath("osascript")
		if err == nil {
			os.Stderr.Close()
			syscall.Exec(path, []string{"osascript", "-l", lang, "-e", script}, nil)
		}
	}

	cmd := exec.Command("osascript", "-l", lang)
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

type Color struct {
	Color []uint32
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
}
