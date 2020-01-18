package zenity

import "image/color"

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
	directory bool
	overwrite bool
	create    bool
	hidden    bool
	filters   []FileFilter

	// Color selection options
	color   color.Color
	palette bool

	// Message options
	icon      MessageIcon
	ok        string
	cancel    string
	extra     string
	nowrap    bool
	ellipsize bool
	defcancel bool
}

// Options are arguments passed to dialog functions to customize their behavior.
type Option func(*options)

func optsParse(options []Option) (res options) {
	for _, o := range options {
		o(&res)
	}
	return
}

// General options

// Option to set the dialog title.
func Title(title string) Option {
	return func(o *options) { o.title = title }
}

// File selection options

// Option to set the filename.
//
// You can specify a file name, a directory path, or both.
// Specifying a file name, makes it the default selected file.
// Specifying a directory path, makes it the default dialog location.
func Filename(filename string) Option {
	return func(o *options) { o.filename = filename }
}

// Option to activate directory-only selection.
func Directory() Option {
	return func(o *options) { o.directory = true }
}

// Option to confirm file selection if filename already exists.
func ConfirmOverwrite() Option {
	return func(o *options) { o.overwrite = true }
}

// Option to confirm file selection if filename does not yet exist (Windows only).
func ConfirmCreate() Option {
	return func(o *options) { o.create = true }
}

// Option to show hidden files (Windows and macOS only).
func ShowHidden() Option {
	return func(o *options) { o.hidden = true }
}

// FileFilter encapsulates a filename filter.
//
// macOS hides filename filters from the user,
// and only supports filtering by extension (or "type").
type FileFilter struct {
	Name     string   // display string that describes the filter (optional)
	Patterns []string // filter patterns for the display string
}

// Build option to set a filename filter.
func (f FileFilter) Build() Option {
	return func(o *options) { o.filters = append(o.filters, f) }
}

// FileFilters is a list of filename filters.
type FileFilters []FileFilter

// Build option to set filename filters.
func (f FileFilters) Build() Option {
	return func(o *options) { o.filters = append(o.filters, f...) }
}

// Color selection options

func Color(c color.Color) Option {
	return func(o *options) { o.color = c }
}

func ShowPalette() Option {
	return func(o *options) { o.palette = true }
}

// Message options

// MessageIcon is the enumeration for message dialog icons.
type MessageIcon int

const (
	ErrorIcon MessageIcon = iota + 1
	WarningIcon
	InfoIcon
	QuestionIcon
)

// Option to set the dialog icon.
func Icon(icon MessageIcon) Option {
	return func(o *options) { o.icon = icon }
}

// Option to set the label of the OK button.
func OKLabel(ok string) Option {
	return func(o *options) { o.ok = ok }
}

// Option to set the label of the Cancel button.
func CancelLabel(cancel string) Option {
	return func(o *options) { o.cancel = cancel }
}

// Option to add an extra button.
func ExtraButton(extra string) Option {
	return func(o *options) { o.extra = extra }
}

// Option to disable enable text wrapping.
func NoWrap() Option {
	return func(o *options) { o.nowrap = true }
}

// Option to enable ellipsizing in the dialog text.
func Ellipsize() Option {
	return func(o *options) { o.ellipsize = true }
}

// Option to give Cancel button focus by default.
func DefaultCancel() Option {
	return func(o *options) { o.defcancel = true }
}
