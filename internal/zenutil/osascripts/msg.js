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
{{if .Cancel -}}
	opts.cancelButton = {{json .Cancel}}
{{end -}}
{{if .Default -}}
	opts.defaultButton = {{json .Default}}
{{end -}}
{{if .Timeout -}}
	opts.givingUpAfter = {{json .Timeout}}
{{end -}}

var res = app[{{json .Operation}}]({{json .Text}}, opts)
if (res.gaveUp) {
	ObjC.import("stdlib")
	$.exit(5)
}
if (res.buttonReturned === {{json .Extra}}) {
	res
} else {
	void 0
}