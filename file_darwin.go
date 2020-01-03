package zenity

import (
	"bytes"
	"html/template"
	"io"
	"os/exec"
	"strings"
)

func SelectFile(options ...FileOption) (string, error) {
	opts := fileoptsParse(options)

	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = scriptExpand(scriptData{
		Operation: "chooseFile",
		Prompt:    opts.title,
		Location:  opts.filename,
		Type:      appleFilters(opts.filters),
	})
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

func SelectFileMutiple(options ...FileOption) ([]string, error) {
	opts := fileoptsParse(options)

	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = scriptExpand(scriptData{
		Operation: "chooseFile",
		Multiple:  true,
		Prompt:    opts.title,
		Location:  opts.filename,
		Type:      appleFilters(opts.filters),
	})
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	if len(out) == 0 {
		return nil, nil
	}
	return strings.Split(string(out), "\x00"), nil
}

func SelectFileSave(options ...FileOption) (string, error) {
	opts := fileoptsParse(options)

	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = scriptExpand(scriptData{
		Operation: "chooseFileName",
		Prompt:    opts.title,
		Location:  opts.filename,
	})
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

func SelectDirectory(options ...FileOption) (string, error) {
	opts := fileoptsParse(options)

	cmd := exec.Command("osascript", "-l", "JavaScript")
	cmd.Stdin = scriptExpand(scriptData{
		Operation: "chooseFolder",
		Prompt:    opts.title,
		Location:  opts.filename,
	})
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

func appleFilters(filters []FileFilter) []string {
	var filter []string
	for _, f := range filters {
		for _, e := range f.Exts {
			filter = append(filter, strings.TrimPrefix(e, "."))
		}
	}
	return filter
}

type scriptData struct {
	Operation string
	Prompt    string
	Location  string
	Type      []string
	Multiple  bool
}

func scriptExpand(data scriptData) io.Reader {
	var buf bytes.Buffer

	err := script.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	var slice = buf.Bytes()
	return bytes.NewReader(slice[len("<script>") : len(slice)-len("</script>")])
}

var script = template.Must(template.New("").Parse(`<script>
var app = Application.currentApplication();
app.includeStandardAdditions = true;
app.activate();

var opts = {};
opts.withPrompt = {{.Prompt}};
opts.multipleSelectionsAllowed = {{.Multiple}};

{{if .Location}}
	opts.defaultLocation = {{.Location}};
{{end}}
{{if .Type}}
	opts.ofType = {{.Type}};
{{end}}

var res;
try {
	res = app[{{.Operation}}](opts);
} catch (e) {
	if (e.errorNumber !== -128) throw e;
}
if (Array.isArray(res)) {
	res.join('\0');
} else if (res != null) {
	res.toString();
} else {
	void 0;
}
</script>`))
