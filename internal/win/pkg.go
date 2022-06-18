// Package win is internal. DO NOT USE.
package win

//go:generate -command mkwinsyscall go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go
//go:generate mkwinsyscall comctl32.go comdlg32.go ole32.go shell32.go
