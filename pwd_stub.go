// +build windows darwin

package zenity

func password(opts options) (string, string, bool, error) {
	opts.hideText = true
	str, ok, err := entry("Password:", opts)
	return "", str, ok, err
}
