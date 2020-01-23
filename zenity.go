// Package zenity provides cross-platform access to simple dialogs that interact
// graphically with the user.
//
// It is inspired by, and closely follows the API of, the zenity program, which
// it uses to provide the functionality on various Unixes. See:
//
// https://help.gnome.org/users/zenity/
//
// This package does not require cgo, and it does not impose any threading or
// initialization requirements.
package zenity

import "image/color"

type constError string

func (e constError) Error() string { return string(e) }

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

// An Option is an argument passed to dialog functions to customize their
// behavior.
type Option func(*options)

func optsParse(options []Option) (res options) {
	for _, o := range options {
		o(&res)
	}
	return
}

// Title returns an Option to set the dialog title.
func Title(title string) Option {
	return func(o *options) { o.title = title }
}
