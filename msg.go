package zenity

// ErrExtraButton is returned by dialog functions when the extra button is
// pressed.
const ErrExtraButton = stringErr("Extra button pressed")

// Question displays the question dialog.
//
// Returns true on OK, false on Cancel, or ErrExtraButton.
//
// Valid options: Title, Width, Height, Icon, OKLabel, CancelLabel,
// ExtraButton, NoWrap, Ellipsize, DefaultCancel.
func Question(text string, options ...Option) (bool, error) {
	return message(questionKind, text, options)
}

// Info displays the info dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, Icon, OKLabel, ExtraButton,
// NoWrap, Ellipsize.
func Info(text string, options ...Option) (bool, error) {
	return message(infoKind, text, options)
}

// Warning displays the warning dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, Icon, OKLabel, ExtraButton,
// NoWrap, Ellipsize.
func Warning(text string, options ...Option) (bool, error) {
	return message(warningKind, text, options)
}

// Error displays the error dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Width, Height, Icon, OKLabel, ExtraButton,
// NoWrap, Ellipsize.
func Error(text string, options ...Option) (bool, error) {
	return message(errorKind, text, options)
}

type messageKind int

const (
	questionKind messageKind = iota
	infoKind
	warningKind
	errorKind
)

// OKLabel returns an Option to set the label of the OK button.
func OKLabel(ok string) Option {
	return funcOption(func(o *options) { o.okLabel = &ok })
}

// CancelLabel returns an Option to set the label of the Cancel button.
func CancelLabel(cancel string) Option {
	return funcOption(func(o *options) { o.cancelLabel = &cancel })
}

// ExtraButton returns an Option to add an extra button.
func ExtraButton(extra string) Option {
	return funcOption(func(o *options) { o.extraButton = &extra })
}

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
