package zenity

// ErrExtraButton is returned by dialog functions when the extra button is
// pressed.
const ErrExtraButton = stringErr("Extra button pressed")

// Question displays the question dialog.
//
// Returns true on OK, false on Cancel, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, NoWrap, Ellipsize, DefaultCancel.
func Question(text string, options ...Option) (bool, error) {
	return message(questionKind, text, applyOptions(options))
}

// Info displays the info dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
func Info(text string, options ...Option) (bool, error) {
	return message(infoKind, text, applyOptions(options))
}

// Warning displays the warning dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
func Warning(text string, options ...Option) (bool, error) {
	return message(warningKind, text, applyOptions(options))
}

// Error displays the error dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, ExtraButton, Icon,
// NoWrap, Ellipsize.
func Error(text string, options ...Option) (bool, error) {
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
