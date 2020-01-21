var app = Application.currentApplication()
app.includeStandardAdditions = true
app.activate()

var opts = {}

{{if .Prompt -}}
	opts.withPrompt = {{json .Prompt}}
{{end -}}
{{if .Type -}}
	opts.ofType = {{json .Type}}
{{end -}}
{{if .Name -}}
	opts.defaultName = {{json .Name}}
{{end -}}
{{if .Location -}}
	opts.defaultLocation = {{json .Location}}
{{end -}}
{{if .Invisibles -}}
	opts.invisibles = {{json .Invisibles}}
{{end -}}
{{if .Multiple -}}
	opts.multipleSelectionsAllowed = {{json .Multiple}}
{{end -}}

var res = app[{{json .Operation}}](opts)
if (Array.isArray(res)) {
	res.join({{json .Separator}})
} else {
	res.toString()
}