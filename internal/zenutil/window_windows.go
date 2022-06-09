package zenutil

import "strconv"

// ParseWindowId is internal.
func ParseWindowId(id string) uintptr {
	hwnd, _ := strconv.ParseUint(id, 0, 64)
	return uintptr(hwnd)
}

// GetParentWindowId is internal.
func GetParentWindowId() int {
	return 0
}
