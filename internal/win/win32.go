//go:build windows

// Package win is internal. DO NOT USE.
package win

import "golang.org/x/sys/windows"

type Handle = windows.Handle
type HWND = windows.HWND
type Pointer = windows.Pointer

//sys RtlGetNtVersionNumbers(major *uint32, minor *uint32, build *uint32) = ntdll.RtlGetNtVersionNumbers
