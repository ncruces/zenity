package zenity

import (
	"os/exec"
	"strings"
)

//go:generate go run osa_scripts/generate.go osa_scripts/

func osaRun(script string, data interface{}) ([]byte, error) {
	var buf strings.Builder

	err := osaScripts.ExecuteTemplate(&buf, script, data)
	if err != nil {
		return nil, err
	}

	var res = buf.String()
	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = strings.NewReader(res[len("<script>") : len(res)-len("</script>")])
	return cmd.Output()
}
