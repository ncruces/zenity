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

	"github.com/ncruces/zenity/internal/zenutil"
)

func stringPtr(s string) *string { return &s }

// ErrCanceled is returned when the cancel button is pressed,
// or window functions are used to close the dialog.
const ErrCanceled = zenutil.ErrCanceled

// ErrExtraButton is returned when the extra button is pressed.
const ErrExtraButton = zenutil.ErrExtraButton

// ErrUnsupported is returned when a combination of options is not supported.
const ErrUnsupported = zenutil.ErrUnsupported

type options struct {
	// General options
	title         *string
	width         uint
	height        uint
	okLabel       *string
	cancelLabel   *string
	extraButton   *string
	icon          DialogIcon
	defaultCancel bool

	// Message options
	noWrap    bool
	ellipsize bool

	// Entry options
	entryText string
	hideText  bool
	username  bool

	// List options
	disallowEmpty bool
	defaultItems  []string

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

// ExtraButton returns an Option to add an extra button.
func ExtraButton(extra string) Option {
	return funcOption(func(o *options) { o.extraButton = &extra })
}

// DialogIcon is the enumeration for dialog icons.
type DialogIcon int

func (i DialogIcon) apply(o *options) { o.icon = i }

// The stock dialog icons.
const (
	ErrorIcon DialogIcon = iota + 1
	WarningIcon
	InfoIcon
	QuestionIcon
	PasswordIcon
	NoIcon
)

// Icon returns an Option to set the dialog icon.
func Icon(icon DialogIcon) Option { return icon }

// Context returns an Option to set a Context that can dismiss the dialog.
//
// Dialogs dismissed by the Context return Context.Err.
func Context(ctx context.Context) Option {
	return funcOption(func(o *options) { o.ctx = ctx })
}
