package zenity

type opts struct {
	title string
}

type Option func(*opts)

func (o *opts) Title(title string) {
	o.title = title
}

type fileopts struct {
	opts
	filename  string
	overwrite bool
	filters   []FileFilter
}

type FileOption func(*fileopts)

func Filename(filename string) FileOption {
	return func(o *fileopts) {
		o.filename = filename
	}
}

func ConfirmOverwrite(o *fileopts) {
	o.overwrite = true
}

type FileFilter struct {
	Name string
	Exts []string
}

type FileFilters []FileFilter

func (f FileFilters) New() FileOption {
	return func(o *fileopts) {
		o.filters = f
	}
}

func fileoptsParse(options []FileOption) (res fileopts) {
	for _, o := range options {
		o(&res)
	}
	return
}
