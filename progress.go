package zenity

// Progress displays the progress indication dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, WindowIcon, Attach, Modal, MaxValue, Pulsate, NoCancel, TimeRemaining.
//
// May return: ErrUnsupported.
func Progress(options ...Option) (ProgressDialog, error) {
	return progress(applyOptions(options))
}

// ProgressDialog allows you to interact with the progress indication dialog.
type ProgressDialog interface {
	// Text sets the dialog text.
	Text(string) error

	// Value sets how much of the task has been completed.
	Value(int) error

	// MaxValue gets how much work the task requires in total.
	MaxValue() int

	// Complete marks the task completed.
	Complete() error

	// Close closes the dialog.
	Close() error

	// Done returns a channel that is closed when the dialog is closed.
	Done() <-chan struct{}
}

// MaxValue returns an Option to set the maximum value.
// The default maximum value is 100.
func MaxValue(value int) Option {
	return funcOption(func(o *options) { o.maxValue = value })
}

// Pulsate returns an Option to pulsate the progress bar.
func Pulsate() Option {
	return funcOption(func(o *options) { o.maxValue = -1 })
}

// NoCancel returns an Option to hide the Cancel button (Windows and Unix only).
func NoCancel() Option {
	return funcOption(func(o *options) { o.noCancel = true })
}

// AutoClose returns an Option to dismiss the dialog when 100% has been reached.
func AutoClose() Option {
	return funcOption(func(o *options) { o.autoClose = true })
}

// TimeRemaining returns an Option to estimate when progress will reach 100% (Unix only).
func TimeRemaining() Option {
	return funcOption(func(o *options) { o.timeRemaining = true })
}
