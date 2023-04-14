package zenity

type formFieldKind int

const (
	FormFieldEntry formFieldKind = iota
	FormFieldPassword
	FormFieldCalendar
	FormFieldComboBox
)

type formFields struct {
	kind   formFieldKind
	name   string
	values []string
}

func Forms(text string, options ...Option) ([]string, error) {
	return forms(text, applyOptions(options))
}
