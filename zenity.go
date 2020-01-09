package zenity

// Errors

type constError string

func (e constError) Error() string { return string(e) }

// Message errors

const ErrExtraButton = constError("Extra button pressed.")

// Options

type options struct {
	// General options
	title string

	// File selection options
	filename  string
	overwrite bool
	filters   []FileFilter

	// Message options
	icon      MessageIcon
	ok        string
	cancel    string
	extra     string
	nowrap    bool
	ellipsize bool
	defcancel bool
}

type Option func(*options)

func optsParse(options []Option) (res options) {
	for _, o := range options {
		o(&res)
	}
	return
}

// General options

func Title(title string) Option {
	return func(o *options) {
		o.title = title
	}
}

// File selection options

func Filename(filename string) Option {
	return func(o *options) {
		o.filename = filename
	}
}

func ConfirmOverwrite(o *options) {
	o.overwrite = true
}

type FileFilter struct {
	Name     string
	Patterns []string
}

type FileFilters []FileFilter

func (f FileFilters) New() Option {
	return func(o *options) {
		o.filters = f
	}
}

// Message options

type MessageIcon int

const (
	ErrorIcon MessageIcon = iota + 1
	InfoIcon
	QuestionIcon
	WarningIcon
)

func Icon(icon MessageIcon) Option {
	return func(o *options) {
		o.icon = icon
	}
}

func OKLabel(ok string) Option {
	return func(o *options) {
		o.ok = ok
	}
}

func CancelLabel(cancel string) Option {
	return func(o *options) {
		o.cancel = cancel
	}
}

func ExtraButton(extra string) Option {
	return func(o *options) {
		o.extra = extra
	}
}

func NoWrap(o *options) {
	o.nowrap = true
}

func Ellipsize(o *options) {
	o.ellipsize = true
}

func DefaultCancel(o *options) {
	o.defcancel = true
}
