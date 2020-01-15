# Zenity dialogs for Golang, Windows and macOS

[![GoDoc](https://godoc.org/github.com/ncruces/zenity?status.svg)](https://godoc.org/github.com/ncruces/zenity)

This repo includes both a cross-platform Go package providing [Zenity](https://help.gnome.org/users/zenity/)-like dialogs,
as well as a *“port”* of the `zenity` command to both Windows and macOS based on that library.

**This is a work in progress.**

Lots of things are missing.
For now, these are the only implemented dialogs:
* [message](https://github.com/ncruces/zenity/wiki/Message-dialog) (error, info, question, warning)
* [file selection](https://github.com/ncruces/zenity/wiki/File-Selection-dialog)

Behavior on Windows, macOS and other Unixes might differ slightly.
Some of that is intended (reflecting platform differences),
other bits are unfortunate limitations,
others still are open to be fixed.

## Why?

There are a bunch of other dialog packages for Go.\
Why reinvent this particular wheel?

#### Requirements:

* no `cgo` (see [benefits](https://dave.cheney.net/2016/01/18/cgo-is-not-go), mostly cross-compilation)
* no main loop (or other threading requirements)
* no initialization
* on Windows:
  * no additional dependencies
    * Explorer shell not required
    * works in Server Core
  * Unicode support
* on macOS:
  * only dependency is `osascript` (with [JXA](https://developer.apple.com/library/archive/releasenotes/InterapplicationCommunication/RN-JavaScriptForAutomation/Articles/Introduction.html))\
    JavaScript is easier to template (with `html/template`)
* on other Unixes:
  * wraps either one of `qarma`, `zenity`, `matedialog`,\
    in that order of preference
