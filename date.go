package zenity

import (
	"time"
)

// Calendar displays the calendar dialog.
//
// Returns zero on cancel.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, Date.
func Calendar(text string, options ...Option) (time.Time, error) {
	return calendar(text, applyOptions(options))
}

// DefaultDate returns an Option to set the date.
func DefaultDate(year int, month time.Month, day int) Option {
	return funcOption(func(o *options) {
		o.year, o.month, o.day = year, int(month), day
	})
}
