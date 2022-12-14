package zenity

func isAvailable() bool { return true }

func attach(id any) Option {
	switch id.(type) {
	case int, string:
	default:
		panic("interface conversion: expected int or string")
	}
	return funcOption(func(o *options) { o.attach = id })
}
