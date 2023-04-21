//go:build !windows && !darwin

package zenity

import (
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func forms(text string, opts options) ([]string, error) {
	args := []string{"--forms", "--text", quoteMarkup(text)}
	args = appendGeneral(args, opts)
	args = appendButtons(args, opts)

	// fields
	for _, field := range opts.fields {
		switch field.kind {
		case FormFieldEntry:
			args = append(args, "--add-entry", quoteMarkup(field.name))
		case FormFieldPassword:
			args = append(args, "--add-password", quoteMarkup(field.name))
		case FormFieldCalendar:
			args = append(args, "--add-calendar", quoteMarkup(field.name))
		case FormFieldComboBox:
			args = append(args, "--add-combo", quoteMarkup(field.name))
			if len(field.values) > 0 {
				args = append(args, "--combo-values", quoteMarkup(strings.Join(field.values, "|")))
			}
		case FormFieldList:
			args = append(args, "--add-list", quoteMarkup(field.name))
			if field.showHeader {
				args = append(args, "--show-header")
			}
			if len(field.cols) > 0 {
				args = append(args, "--column-values", quoteMarkup(strings.Join(field.cols, "|")))
			}
			if len(field.values) > 0 {
				args = append(args, "--list-values", quoteMarkup(strings.Join(field.values, "|")))
			}
		}
	}

	out, err := zenutil.Run(opts.ctx, args)
	return formsResult(opts, out, err)
}
