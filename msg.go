package zenity

// Question displays the question dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, NoWrap, Ellipsize, DefaultCancel.
//
// May return: ErrCanceled, ErrExtraButton.
func Question(text string, options ...Option) error {
	return message(questionKind, text, applyOptions(options))
}

// Info displays the info dialog.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
//
// May return: ErrCanceled, ErrExtraButton.
func Info(text string, options ...Option) error {
	return message(infoKind, text, applyOptions(options))
}

// Warning displays the warning dialog.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
//
// May return: ErrCanceled, ErrExtraButton.
func Warning(text string, options ...Option) error {
	return message(warningKind, text, applyOptions(options))
}

// Error displays the error dialog.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
//
// May return: ErrCanceled, ErrExtraButton.
func Error(text string, options ...Option) error {
	return message(errorKind, text, applyOptions(options))
}

type messageKind int

const (
	questionKind messageKind = iota
	infoKind
	warningKind
	errorKind
)

// NoWrap returns an Option to disable enable text wrapping (Unix only).
func NoWrap() Option {
	return funcOption(func(o *options) { o.noWrap = true })
}

// Ellipsize returns an Option to enable ellipsizing in the dialog text (Unix only).
func Ellipsize() Option {
	return funcOption(func(o *options) { o.ellipsize = true })
}

// DefaultCancel returns an Option to give the Cancel button focus by default.
func DefaultCancel() Option {
	return funcOption(func(o *options) { o.defaultCancel = true })
}
