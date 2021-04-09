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

// EntryText returns an Option to set the entry text.
func EntryText(text string) Option {
	return funcOption(func(o *options) { o.entryText = text })
}

// HideText returns an Option to hide the entry text.
func HideText() Option {
	return funcOption(func(o *options) { o.hideText = true })
}
