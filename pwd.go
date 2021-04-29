package zenity

// Password displays the password dialog.
//
// Valid options: Title, OKLabel, CancelLabel, ExtraButton, Icon, Username.
func Password(options ...Option) (usr string, pw string, err error) {
	return password(applyOptions(options))
}

// Username returns an Option to display the username (Unix only).
func Username() Option {
	return funcOption(func(o *options) { o.username = true })
}
