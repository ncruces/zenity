//go:build windows || darwin

package zenity

func password(opts options) (string, string, error) {
	if opts.username {
		return "", "", ErrUnsupported
	}
	opts.hideText = true
	str, err := entry("Password:", opts)
	return "", str, err
}
