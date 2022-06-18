//go:build windows

package win

import "golang.org/x/sys/windows"

func GetCurrentThreadId() (id uint32)     { return windows.GetCurrentThreadId() }
func GetSystemDirectory() (string, error) { return windows.GetSystemDirectory() }

//sys GetConsoleWindow() (ret HWND) = kernel32.GetConsoleWindow
//sys GetModuleHandle(moduleName *uint16) (ret Handle, err error) = kernel32.GetModuleHandleW
