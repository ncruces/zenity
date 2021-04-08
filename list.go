package zenity

// List displays the list dialog.
//
// Returns false on cancel, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, DefaultItems, DisallowEmpty.
func List(text string, items []string, options ...Option) (string, bool, error) {
	return list(text, items, applyOptions(options))
}

// ListItems displays the list dialog.
//
// Returns false on cancel, or ErrExtraButton.
func ListItems(text string, items ...string) (string, bool, error) {
	return List(text, items)
}

// ListMultiple displays the list dialog, allowing multiple items to be selected.
//
// Returns a nil slice on cancel, or ErrExtraButton.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, DefaultItems, DisallowEmpty.
func ListMultiple(text string, items []string, options ...Option) ([]string, error) {
	return listMultiple(text, items, applyOptions(options))
}

// ListMultiple displays the list dialog, allowing multiple items to be selected.
//
// Returns a nil slice on cancel, or ErrExtraButton.
func ListMultipleItems(text string, items ...string) ([]string, error) {
	return ListMultiple(text, items)
}

// DefaultItems returns an Option to set the items to initally select (Windows and macOS only).
func DefaultItems(items ...string) Option {
	return funcOption(func(o *options) { o.defaultItems = items })
}

// DisallowEmpty returns an Option to not allow zero items to be selected (Windows and macOS only).
func DisallowEmpty() Option {
	return funcOption(func(o *options) { o.disallowEmpty = true })
}
