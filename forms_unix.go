//go:build !windows && !darwin

package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func forms(text string, opts options) ([]string, error) {
	args := []string{"--forms", "--text", quoteMarkup(text)}
	args = appendGeneral(args, opts)

	// fields
	for _, field := range opts.fields {
		switch field.kind {
		case FormFieldEntry:
			args = append(args, "--add-entry", quoteMarkup(field.name))
		case FormFieldPassword:
			args = append(args, "--add-password", quoteMarkup(field.name))
		case FormFieldCalendar:
			args = append(args, "--add-calendar", quoteMarkup(field.name))
		}
	}

	out, err := zenutil.Run(opts.ctx, args)
	return formsResult(opts, out, err)
}
