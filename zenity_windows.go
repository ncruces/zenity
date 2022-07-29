package zenity

import (
	"reflect"

	"github.com/ncruces/zenity/internal/win"
)

// Attach returns an Option to set the parent window to attach to.
//
// Attach accepts:
//   - a window id (int) on Unix
//   - a window handle (~uintptr) on Windows
//   - an application name (string) or process id (int) on macOS
func Attach(id any) Option {
	if v := reflect.ValueOf(id); v.Kind() == reflect.Uintptr {
		id = win.HWND(uintptr(v.Uint()))
	} else {
		panic("interface conversion: expected uintptr")
	}
	return funcOption(func(o *options) { o.attach = id })
}
