//go:build windows || darwin

package zenity

import "fmt"

func password(opts options) (string, string, error) {
	if opts.username {
		return "", "", fmt.Errorf("%w: username", ErrUnsupported)
	}
	opts.hideText = true
	str, err := entry("Password:", opts)
	return "", str, err
}
