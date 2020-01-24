package zenity

// ErrExtraButton is returned by dialog functions when the extra button is
// pressed.
const ErrExtraButton = constError("Extra button pressed.")

// Question displays the question dialog.
//
// Returns true on OK, false on Cancel, or ErrExtraButton.
//
// Valid options: Title, Icon, OKLabel, CancelLabel, ExtraButton, NoWrap,
// Ellipsize, DefaultCancel.
func Question(text string, options ...Option) (bool, error) {
	return message(questionKind, text, options)
}

// Info displays the info dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Icon, OKLabel, ExtraButton, NoWrap, Ellipsize.
func Info(text string, options ...Option) (bool, error) {
	return message(infoKind, text, options)
}

// Warning displays the warning dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Icon, OKLabel, ExtraButton, NoWrap, Ellipsize.
func Warning(text string, options ...Option) (bool, error) {
	return message(warningKind, text, options)
}

// Error displays the error dialog.
//
// Returns true on OK, false on dismiss, or ErrExtraButton.
//
// Valid options: Title, Icon, OKLabel, ExtraButton, NoWrap, Ellipsize.
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

// MessageIcon is the enumeration for message dialog icons.
type MessageIcon int

// Icons for
const (
	ErrorIcon MessageIcon = iota + 1
	WarningIcon
	InfoIcon
	QuestionIcon
)

// Icon returns an Option to set the dialog icon.
func Icon(icon MessageIcon) Option {
	return funcOption(func(o *options) { o.icon = icon })
}

// OKLabel returns an Option to set the label of the OK button.
func OKLabel(ok string) Option {
	return funcOption(func(o *options) { o.okLabel = ok })
}

// CancelLabel returns an Option to set the label of the Cancel button.
func CancelLabel(cancel string) Option {
	return funcOption(func(o *options) { o.cancelLabel = cancel })
}

// ExtraButton returns an Option to add an extra button.
func ExtraButton(extra string) Option {
	return funcOption(func(o *options) { o.extraButton = extra })
}

// NoWrap returns an Option to disable enable text wrapping.
func NoWrap() Option {
	return funcOption(func(o *options) { o.noWrap = true })
}

// Ellipsize returns an Option to enable ellipsizing in the dialog text.
func Ellipsize() Option {
	return funcOption(func(o *options) { o.ellipsize = true })
}

// DefaultCancel returns an Option to give Cancel button focus by default.
func DefaultCancel() Option {
	return funcOption(func(o *options) { o.defaultCancel = true })
}
