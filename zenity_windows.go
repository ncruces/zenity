package zenity

import (
	"reflect"

	"github.com/ncruces/zenity/internal/win"
)

func isAvailable() bool { return true }

func attach(id any) Option {
	if v := reflect.ValueOf(id); v.Kind() == reflect.Uintptr {
		id = win.HWND(uintptr(v.Uint()))
	} else {
		panic("interface conversion: expected uintptr")
	}
	return funcOption(func(o *options) { o.attach = id })
}
