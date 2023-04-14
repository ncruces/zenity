//go:build !windows && !darwin

package zenity

import (
	"fmt"

	"github.com/ncruces/zenity/internal/zenutil"
)

func forms(text string, opts options) ([]string, error) {
	args := []string{"--forms", "--text", quoteMarkup(text)}
	args = appendGeneral(args, opts)

	// password fields
	for _, name := range opts.passwords {
		args = append(args, "--add-password", quoteMarkup(name))
	}

	// calendar fields
	for _, name := range opts.calendars {
		args = append(args, "--add-calendar", quoteMarkup(name))
	}

	// entry fields
	for _, name := range opts.entries {
		args = append(args, "--add-entry", quoteMarkup(name))
	}

	out, err := zenutil.Run(opts.ctx, args)
	fmt.Println(err)
	fmt.Println(string(out))
	return formsResult(opts, out, err)
}
