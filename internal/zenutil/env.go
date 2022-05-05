// Package zenutil is internal. DO NOT USE.
package zenutil

import "time"

// These are internal.
const (
	ErrCanceled    = stringErr("dialog canceled")
	ErrExtraButton = stringErr("extra button pressed")
	ErrUnsupported = stringErr("unsupported option")
)

// These are internal.
var (
	Command bool
	Timeout int

	DateFormat = "%Y-%m-%d"
	DateUTS35  = func() (string, error) { return "yyyy-MM-dd", nil }
	DateParse  = func(s string) (time.Time, error) { return time.Parse("2006-01-02", s) }
)

type stringErr string

func (e stringErr) Error() string { return string(e) }
