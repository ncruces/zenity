package zenity

func Forms(text string, options ...Option) ([]string, error) {
	return forms(text, applyOptions(options))
}
