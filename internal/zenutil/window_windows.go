package zenutil

import (
	"strconv"

	"golang.org/x/sys/windows"
)

// ParseWindowId is internal.
func ParseWindowId(id string) windows.HWND {
	hwnd, _ := strconv.ParseUint(id, 0, 64)
	return windows.HWND(uintptr(hwnd))
}

// GetParentWindowId is internal.
func GetParentWindowId(pid int) int {
	return 0
}
