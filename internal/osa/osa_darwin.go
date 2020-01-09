package osa

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ncruces/zenity/internal/cmd"
)

func Run(script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	res := buf.String()
	res = res[len("<script>") : len(res)-len("\n</script>")]

	if cmd.Command {
		path, err := exec.LookPath("osascript")
		if err == nil {
			os.Stderr.Close()
			syscall.Exec(path, []string{"osascript", "-l", "JavaScript", "-e", res}, nil)
		}
	}

	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = strings.NewReader(res)
	return cmd.Output()
}

type File struct {
	Operation string
	Prompt    string
	Name      string
	Location  string
	Separator string
	Type      []string
	Multiple  bool
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
