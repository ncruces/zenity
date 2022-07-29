//go:build !windows && !darwin

package zenity

// Attach returns an Option to set the parent window to attach to.
//
// Attach accepts:
//   - a window id (int) on Unix
//   - a window handle (~uintptr) on Windows
//   - an application name (string) or process id (int) on macOS
func Attach(id int) Option {
	return funcOption(func(o *options) { o.attach = id })
}
