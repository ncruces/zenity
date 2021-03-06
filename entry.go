package zenity

// Entry displays the text entry dialog.
//
// Returns false on cancel, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, EntryText, HideText.
func Entry(text string, options ...Option) (string, bool, error) {
	return entry(text, applyOptions(options))
}

// Password displays the password dialog.
//
// Returns false on cancel, or ErrExtraButton.
//
// Valid options: Title, OKLabel, CancelLabel, ExtraButton, Icon, Username.
func Password(options ...Option) (usr string, pw string, ok bool, err error) {
	return password(applyOptions(options))
}

// EntryText returns an Option to set the entry text.
func EntryText(text string) Option {
	return funcOption(func(o *options) { o.entryText = text })
}

// HideText returns an Option to hide the entry text.
func HideText() Option {
	return funcOption(func(o *options) { o.hideText = true })
}

// Username returns an Option to display the username (Unix only).
func Username() Option {
	return funcOption(func(o *options) { o.username = true })
}
