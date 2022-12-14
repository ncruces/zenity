//go:build !windows && !darwin

package zenity

import "github.com/ncruces/zenity/internal/zenutil"

func isAvailable() bool { return zenutil.IsAvailable() }

func attach(id any) Option {
	return funcOption(func(o *options) { o.attach = id.(int) })
}
