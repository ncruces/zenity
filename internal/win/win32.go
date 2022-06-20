//go:build windows

// Package win is internal. DO NOT USE.
package win

import "golang.org/x/sys/windows"

type Handle = windows.Handle
type HWND = windows.HWND
type Pointer = windows.Pointer

func RtlGetNtVersionNumbers() (majorVersion, minorVersion, buildNumber uint32) {
	return windows.RtlGetNtVersionNumbers()
}
