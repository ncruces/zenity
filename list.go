package zenity

// List displays the list dialog.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// WindowIcon, Attach, Modal, RadioList, DefaultItems, DisallowEmpty.
//
// May return: ErrCanceled, ErrExtraButton, ErrUnsupported.
func List(text string, items []string, options ...Option) (string, error) {
	return list(text, items, applyOptions(options))
}

// ListItems displays the list dialog.
//
// May return: ErrCanceled, ErrUnsupported.
func ListItems(text string, items ...string) (string, error) {
	return List(text, items)
}

// ListMultiple displays the list dialog, allowing multiple items to be selected.
//
// Valid options: Title, Width, Height, OKLabel, CancelLabel, ExtraButton,
// WindowIcon, Attach, Modal, CheckList, DefaultItems, DisallowEmpty.
//
// May return: ErrCanceled, ErrExtraButton, ErrUnsupported.
func ListMultiple(text string, items []string, options ...Option) ([]string, error) {
	return listMultiple(text, items, applyOptions(options))
}

// ListMultipleItems displays the list dialog, allowing multiple items to be selected.
//
// May return: ErrCanceled, ErrUnsupported.
func ListMultipleItems(text string, items ...string) ([]string, error) {
	return ListMultiple(text, items, CheckList())
}

// CheckList returns an Option to show check boxes (Unix only).
func CheckList() Option {
	return funcOption(func(o *options) { o.listKind = checkListKind })
}

// RadioList returns an Option to show radio boxes (Unix only).
func RadioList() Option {
	return funcOption(func(o *options) { o.listKind = radioListKind })
}

type listKind int

const (
	basicListKind listKind = iota
	checkListKind
	radioListKind
)

// MidSearch returns an Option to change list search to find text in the middle,
// not on the beginning (Unix only).
func MidSearch() Option {
	return funcOption(func(o *options) { o.midSearch = true })
}

// DefaultItems returns an Option to set the items to initially select (Windows and macOS only).
func DefaultItems(items ...string) Option {
	return funcOption(func(o *options) { o.defaultItems = items })
}

// DisallowEmpty returns an Option to not allow zero items to be selected (Windows and macOS only).
func DisallowEmpty() Option {
	return funcOption(func(o *options) { o.disallowEmpty = true })
}
