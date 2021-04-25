package zenity

// Progress displays the progress indication dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, MaxValue, Pulsate, NoCancel, TimeRemaining.
func Progress(options ...Option) (ProgressDialog, error) {
	return progress(applyOptions(options))
}

type ProgressDialog interface {
	Text(string) error
	Value(int) error
	Close() error
}

// MaxValue returns an Option to set the maximum value (macOS only).
// The default value is 100.
func MaxValue(value int) Option {
	return funcOption(func(o *options) { o.maxValue = value })
}

// Pulsate returns an Option to pulsate the progress bar.
func Pulsate() Option {
	return funcOption(func(o *options) { o.maxValue = -1 })
}

// NoCancel returns an Option to hide the Cancel button (Unix only).
func NoCancel() Option {
	return funcOption(func(o *options) { o.noCancel = true })
}

// TimeRemaining returns an Option to estimate when progress will reach 100% (Unix only).
func TimeRemaining() Option {
	return funcOption(func(o *options) { o.timeRemaining = true })
}
