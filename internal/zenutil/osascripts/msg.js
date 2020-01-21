var app = Application.currentApplication()
app.includeStandardAdditions = true
app.activate()

var opts = {}

{{if .Message -}}
	opts.message = {{json .Message}}
{{end -}}
{{if .As -}}
	opts.as = {{json .As}}
{{end -}}
{{if .Title -}}
	opts.withTitle = {{json .Title}}
{{end -}}
{{if .Icon -}}
	opts.withIcon = {{json .Icon}}
{{end -}}
{{if .Buttons -}}
	opts.buttons = {{json .Buttons}}
{{end -}}
{{if .Default -}}
	opts.defaultButton = {{json .Default}}
{{end -}}
{{if .Cancel -}}
	opts.cancelButton = {{json .Cancel}}
{{end -}}

var res = app[{{json .Operation}}]({{json .Text}}, opts).buttonReturned
res === {{json .Extra}} ? res : void 0