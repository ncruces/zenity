//go:build windows

// Package win is internal. DO NOT USE.
package win

import "golang.org/x/sys/windows"

type (
	Handle  = windows.Handle
	HWND    = windows.HWND
	Pointer = windows.Pointer
	GUID    = windows.GUID
)

// https://docs.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-systemtime
type SYSTEMTIME struct {
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	Day          uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Milliseconds uint16
}

func RtlGetNtVersionNumbers() (majorVersion, minorVersion, buildNumber uint32) {
	return windows.RtlGetNtVersionNumbers()
}
