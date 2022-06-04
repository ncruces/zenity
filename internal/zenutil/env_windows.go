package zenutil

import "strconv"

// These are internal.
var (
	Separator string
	LineBreak = "\r\n"
)

func ParseWindowId(id string) uintptr {
	hwnd, _ := strconv.ParseUint(id, 0, 64)
	return uintptr(hwnd)
}
