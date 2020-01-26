package zenity

// Notify displays a notification.
func Notify(text string, options ...Option) error {
	return notify(text, options)
}
