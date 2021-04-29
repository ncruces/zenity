// +build windows darwin

package zenity

func password(opts options) (string, string, error) {
	opts.hideText = true
	str, err := entry("Password:", opts)
	return "", str, err
}
