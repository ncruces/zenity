var app = Application.currentApplication()
app.includeStandardAdditions = true
app.activate()

var opts = {}

{{if .Title -}}
	opts.withTitle = {{json .Title}}
{{end -}}
{{if .Subtitle -}}
	opts.subtitle = {{json .Subtitle}}
{{end -}}

void app.displayNotification({{json .Text}}, opts)