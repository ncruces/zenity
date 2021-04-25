package zenity

import (
	"strconv"

	"github.com/ncruces/zenity/internal/zenutil"
)

func progress(opts options) (ProgressDialog, error) {
	var env []string
	if opts.title != nil {
		env = append(env, "description="+*opts.title)
	}
	if opts.maxValue == 0 {
		opts.maxValue = 100
	}
	if opts.maxValue >= 0 {
		env = append(env, "total="+strconv.Itoa(opts.maxValue))
	}
	return zenutil.RunProgress(opts.ctx, opts.maxValue, env)
}
