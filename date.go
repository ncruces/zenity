package zenity

import "time"

// Calendar displays the calendar dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// WindowIcon, Attach, Modal, DefaultDate.
//
// May return: ErrCanceled, ErrExtraButton.
func Calendar(text string, options ...Option) (time.Time, error) {
	return calendar(text, applyOptions(options))
}

// DefaultDate returns an Option to set the date.
func DefaultDate(year int, month time.Month, day int) Option {
	return funcOption(func(o *options) {
		o.time = ptr(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
	})
}
