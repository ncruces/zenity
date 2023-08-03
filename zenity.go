// Package zenity provides cross-platform access to simple dialogs that interact
// graphically with the user.
//
// It is inspired by, and closely follows the API of, the zenity program, which
// it uses to provide the functionality on various Unixes. See:
//
// https://help.gnome.org/users/zenity/stable/
//
// This package does not require cgo, and it does not impose any threading or
// initialization requirements.
package zenity

import (
	"context"
	"image/color"
	"time"

	"github.com/ncruces/zenity/internal/zenutil"
)

func ptr[T any](v T) *T { return &v }

// ErrCanceled is returned when the cancel button is pressed,
// or window functions are used to close the dialog.
const ErrCanceled = zenutil.ErrCanceled

// ErrExtraButton is returned when the extra button is pressed.
const ErrExtraButton = zenutil.ErrExtraButton

// ErrUnsupported is returned when a combination of options is not supported.
const ErrUnsupported = zenutil.ErrUnsupported

// IsAvailable reports whether dependencies of the package are installed.
// It always returns true on Windows and macOS.
func IsAvailable() bool {
	return isAvailable()
}

type options struct {
	// General options
	title         *string
	width         uint
	height        uint
	okLabel       *string
	cancelLabel   *string
	extraButton   *string
	defaultCancel bool
	icon          any
	windowIcon    any
	attach        any
	modal         bool
	display       string
	class         string
	name          string

	// Message options
	noWrap    bool
	ellipsize bool

	// Entry options
	entryText string
	hideText  bool
	username  bool

	// List options
	listKind      listKind
	midSearch     bool
	disallowEmpty bool
	defaultItems  []string

	// Calendar options
	time *time.Time

	// File selection options
	directory        bool
	confirmOverwrite bool
	confirmCreate    bool
	showHidden       bool
	filename         string
	fileFilters      FileFilters

	// Color selection options
	color       color.Color
	showPalette bool

	// Progress indication options
	maxValue      int
	noCancel      bool
	autoClose     bool
	timeRemaining bool

	// Context for timeout
	ctx context.Context
}

// An Option is an argument passed to dialog functions to customize their
// behavior.
type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (f funcOption) apply(o *options) { f(o) }

func applyOptions(options []Option) (res options) {
	for _, o := range options {
		o.apply(&res)
	}
	return
}

// Title returns an Option to set the dialog title.
func Title(title string) Option {
	return funcOption(func(o *options) { o.title = &title })
}

// Width returns an Option to set the dialog width (Unix only).
func Width(width uint) Option {
	return funcOption(func(o *options) {
		o.width = width
	})
}

// Height returns an Option to set the dialog height (Unix only).
func Height(height uint) Option {
	return funcOption(func(o *options) {
		o.height = height
	})
}

// OKLabel returns an Option to set the label of the OK button.
func OKLabel(ok string) Option {
	return funcOption(func(o *options) { o.okLabel = &ok })
}

// CancelLabel returns an Option to set the label of the Cancel button.
func CancelLabel(cancel string) Option {
	return funcOption(func(o *options) { o.cancelLabel = &cancel })
}

// ExtraButton returns an Option to add one extra button.
func ExtraButton(extra string) Option {
	return funcOption(func(o *options) { o.extraButton = &extra })
}

// DefaultCancel returns an Option to give the Cancel button focus by default.
func DefaultCancel() Option {
	return funcOption(func(o *options) { o.defaultCancel = true })
}

// DialogIcon is an Option that sets the dialog icon.
type DialogIcon int

func (i DialogIcon) apply(o *options) {
	o.icon = i
}

// The stock dialog icons.
const (
	ErrorIcon DialogIcon = iota
	WarningIcon
	InfoIcon
	QuestionIcon
	PasswordIcon
	NoIcon
)

// Icon returns an Option to set the dialog icon.
//
// Icon accepts a DialogIcon, or a string.
// The string can be a GTK icon name (Unix), or a file path (Windows and macOS).
// Supported file formats depend on the plaftorm, but PNG should be cross-platform.
func Icon(icon any) Option {
	switch icon.(type) {
	case DialogIcon, string:
	default:
		panic("interface conversion: expected string or DialogIcon")
	}
	return funcOption(func(o *options) { o.icon = icon })
}

// WindowIcon returns an Option to set the window icon.
//
// WindowIcon accepts a DialogIcon, or a string file path.
// Supported file formats depend on the plaftorm, but PNG should be cross-platform.
func WindowIcon(icon any) Option {
	switch icon.(type) {
	case DialogIcon, string:
	default:
		panic("interface conversion: expected string or DialogIcon")
	}
	return funcOption(func(o *options) { o.windowIcon = icon })
}

// Attach returns an Option to set the parent window to attach to.
//
// Attach accepts:
//   - a window id (int) on Unix
//   - a window handle (~uintptr) on Windows
//   - an application name (string) or process id (int) on macOS
func Attach(id any) Option {
	return attach(id)
}

// Modal returns an Option to set the modal hint.
func Modal() Option {
	return funcOption(func(o *options) { o.modal = true })
}

// Display returns an Option to set the X display to use (Unix only).
func Display(display string) Option {
	return funcOption(func(o *options) { o.display = display })
}

// ClassHint returns an Option to set the program name and class
// as used by the window manager (Unix only).
func ClassHint(name, class string) Option {
	return funcOption(func(o *options) {
		if name != "" {
			o.name = name
		}
		if class != "" {
			o.class = class
		}
	})
}

// Context returns an Option to set a Context that can dismiss the dialog.
//
// Dialogs dismissed by ctx return ctx.Err().
func Context(ctx context.Context) Option {
	return funcOption(func(o *options) { o.ctx = ctx })
}
