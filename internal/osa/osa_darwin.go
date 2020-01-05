package osa

import (
	"os/exec"
	"strings"
)

//go:generate go run scripts/generate.go scripts/

func Run(script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := scripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	res := buf.String()
	res = res[len("<script>") : len(res)-len("</script>")]
	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = strings.NewReader(res)
	return cmd.Output()
}

type File struct {
	Operation string
	Prompt    string
	Location  string
	Type      []string
	Multiple  bool
}

type Msg struct {
	Dialog  bool
	Text    string
	Message string
	As      string
	Title   string
	Icon    string
	Buttons []string
	Cancel  int
	Default int
}
