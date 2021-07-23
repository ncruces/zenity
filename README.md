# Zenity dialogs for Golang, Windows and macOS

[![Go Reference](https://pkg.go.dev/badge/image)](https://pkg.go.dev/github.com/ncruces/zenity)
[![Go Report](https://goreportcard.com/badge/github.com/ncruces/zenity)](https://goreportcard.com/report/github.com/ncruces/zenity)
[![Go Cover](https://gocover.io/_badge/github.com/ncruces/zenity)](https://gocover.io/github.com/ncruces/zenity)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

This repo includes both a cross-platform Go package providing
[Zenity](https://help.gnome.org/users/zenity/stable/)-like dialogs
(simple dialogs that interact graphically with the user),
as well as a *“port”* of the `zenity` command to both Windows and macOS based on that library.

Implemented dialogs:
* [message](https://github.com/ncruces/zenity/wiki/Message-dialog) (error, info, question, warning)
* [text entry](https://github.com/ncruces/zenity/wiki/Text-entry-dialog)
* [list](https://github.com/ncruces/zenity/wiki/List-dialog) (simple)
* [password](https://github.com/ncruces/zenity/wiki/Password-dialog)
* [file selection](https://github.com/ncruces/zenity/wiki/File-selection-dialog)
* [color selection](https://github.com/ncruces/zenity/wiki/Color-selection-dialog)
* [progress](https://github.com/ncruces/zenity/wiki/Progress-dialog)
* [notification](https://github.com/ncruces/zenity/wiki/Notification)

Behavior on Windows, macOS and other Unixes might differ slightly.
Some of that is intended (reflecting platform differences),
other bits are unfortunate limitations.

## Installing

The Go package:

    go get github.com/ncruces/zenity

The `zenity` command on macOS/WSL using [Homebrew](https://brew.sh/):

    brew install ncruces/tap/zenity

Or download the [latest release](https://github.com/ncruces/zenity/releases/latest).

## Why?

There are a bunch of other dialog packages for Go.\
Why reinvent this particular wheel?

#### Benefits:

* no `cgo` (see [benefits](https://dave.cheney.net/2016/01/18/cgo-is-not-go), mostly cross-compilation)
* no main loop (or any other threading or initialization requirements)
* cancelation through [`context`](https://golang.org/pkg/context/)
* on Windows:
  * no additional dependencies
    * Explorer shell not required
    * works in Server Core
  * Unicode support
  * High DPI (no manifest required)
  * Visual Styles (no manifest required)
  * WSL/Cygwin/MSYS2 [support](https://github.com/ncruces/zenity/wiki/Zenity-for-WSL,-Cygwin,-MSYS2)
* on macOS:
  * only dependency is `osascript`
* on other Unixes:
  * wraps either one of `zenity`, `qarma`, `matedialog`
