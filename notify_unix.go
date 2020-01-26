// +build !windows,!darwin

package zenity

func notify(text string, options []Option) error {
	panic("not implemented")
}
