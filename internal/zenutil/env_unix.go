//go:build !windows && !darwin

// Package zenutil is internal. DO NOT USE.
package zenutil

// These are internal.
var (
	Command   bool
	Timeout   int
	LineBreak = "\n"
	Separator = "\x1e"
)
