// +build !windows,!darwin,!js

// Package zenutil is internal. DO NOT USE.
package zenutil

// These are internal.
const (
	LineBreak = "\n"
)

// These are internal.
var (
	Command   bool
	Timeout   int
	Separator = "\x1e"
)
