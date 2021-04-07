// +build !windows,!darwin

package zenity

func list(text string, items []string, opts options) (string, bool, error) {
	return "", false, nil
}

func listMultiple(text string, items []string, opts options) ([]string, error) {
	return nil, nil
}
