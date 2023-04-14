package zenity

type formFieldKind int

const (
	FormFieldEntry formFieldKind = iota
	FormFieldPassword
	FormFieldCalendar
	FormFieldComboBox
	FormFieldList
)

type formFields struct {
	kind       formFieldKind
	name       string
	cols       []string
	values     []string
	showHeader bool
}

func Forms(text string, options ...Option) ([]string, error) {
	return forms(text, applyOptions(options))
}
