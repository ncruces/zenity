package zenity

// List displays the list dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, DefaultItems, DisallowEmpty.
//
// May return: ErrCanceled, ErrExtraButton, ErrUnsupported.
func List(text string, items []string, options ...Option) (string, error) {
	return list(text, items, applyOptions(options))
}

// ListItems displays the list dialog.
//
// May return: ErrCanceled, ErrExtraButton.
func ListItems(text string, items ...string) (string, error) {
	return List(text, items)
}

// ListMultiple displays the list dialog, allowing multiple items to be selected.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// Icon, DefaultItems, DisallowEmpty.
//
// May return: ErrCanceled, ErrExtraButton, ErrUnsupported.
func ListMultiple(text string, items []string, options ...Option) ([]string, error) {
	return listMultiple(text, items, applyOptions(options))
}

// ListMultipleItems displays the list dialog, allowing multiple items to be selected.
//
// May return: ErrCanceled, ErrExtraButton.
func ListMultipleItems(text string, items ...string) ([]string, error) {
	return ListMultiple(text, items)
}

// DefaultItems returns an Option to set the items to initially select (macOS only).
func DefaultItems(items ...string) Option {
	return funcOption(func(o *options) { o.defaultItems = items })
}

// DisallowEmpty returns an Option to not allow zero items to be selected (macOS only).
func DisallowEmpty() Option {
	return funcOption(func(o *options) { o.disallowEmpty = true })
}
