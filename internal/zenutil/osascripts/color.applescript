tell application (path to frontmost application as text)
	activate
	set c to choose color default color { {{index . 0}},{{index . 1}},{{index . 2}} }
	"rgb(" & (item 1 of c) div 256 & "," & (item 2 of c) div 256 & "," & (item 3 of c) div 256 & ")"
end tell