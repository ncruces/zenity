//go:build !windows && !darwin

package zenutil

// These are internal.
var (
	Command    bool
	Timeout    int
	Separator  = "\x1e"
	LineBreak  = "\n"
	DateFormat = "%F"
)
