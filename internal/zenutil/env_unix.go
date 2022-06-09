//go:build !windows && !darwin

package zenutil

// These are internal.
var (
	Separator = "\x1e"
	LineBreak = "\n"
)
