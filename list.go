package zenity

// MultipleSelection returns an Option to enable multiple selection of the list dialog.
func MultipleSelection() Option {
	return funcOption(func(o *options) { o.multipleSelection = true })
}

func List(text string, choices []string, options ...Option) ([]uint, error) {
	return list(text, choices, options)
}
