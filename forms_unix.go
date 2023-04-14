//go:build !windows && !darwin

package zenity

import (
	"fmt"

	"github.com/ncruces/zenity/internal/zenutil"
)

func forms(text string, opts options) ([]string, error) {
	args := []string{"--forms", "--text", quoteMarkup(text)}
	args = appendGeneral(args, opts)


	out, err := zenutil.Run(opts.ctx, args)
	fmt.Println(err)
	fmt.Println(string(out))
	return formsResult(opts, out, err)
}
