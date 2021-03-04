package zenity

// Notify displays a notification.
//
// Valid options: Title, Icon.
func Notify(text string, options ...Option) error {
	return notify(text, applyOptions(options))
}
