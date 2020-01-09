# Zenity dialogs for Golang, Windows and macOS

[![GoDoc](https://godoc.org/github.com/ncruces/zenity?status.svg)](https://godoc.org/github.com/ncruces/zenity)

This repo includes both a cross-platform Go package providing [Zenity](https://help.gnome.org/users/zenity/)-like dialogs,
as well as a *“port”* of the `zenity` command to both Windows and macOS based on that library.

**This is a work in progress.**

Lots of things are missing.
For now, these are the only implemented dialogs:
* message (error, info, question, warning); and
* file selection.

Behavior on Windows, macOS and other UNIXes might differ sliglty.
Some of that is intended (reflecting platform differences),
other bits are unfortunate limitations,
others still open to be fixed.

## Why?

There are a bunch of other dialog packages for Go.
Why reinvent this particular wheel?

#### Requirements:

* no `cgo` (see [benefits](https://dave.cheney.net/2016/01/18/cgo-is-not-go), mostly cross-compilation)
* no main loop (or other threading requirements)
* no initialization
* on Windows:
  * Explorer shell not required (works in Server Core)
  * no other dependencies
  * Unicode support
* on macOS:
  * only dependency is `osascript`, JXA
* on other UNIXes:
  * wraps either one of `matedialog`, `qarma`, `zenity` (in that order of preference)
  * no command line support
