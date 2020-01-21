tell application (path to frontmost application as text)
	activate
	{{if .Color -}}
		set c to choose color default color { {{index .Color 0}}, {{index .Color 1}}, {{index .Color 2}} }
	{{else -}}
		set c to choose color
	{{end}}
	"rgb(" & (item 1 of c) div 256 & "," & (item 2 of c) div 256 & "," & (item 3 of c) div 256 & ")"
end tell