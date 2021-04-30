package zenity

import (
	"github.com/ncruces/zenity/internal/zenutil"
)

func progress(opts options) (ProgressDialog, error) {
	if opts.extraButton != nil {
		return nil, ErrUnsupported
	}

	var data zenutil.Progress
	data.Description = opts.title
	if opts.maxValue == 0 {
		opts.maxValue = 100
	}
	if opts.maxValue >= 0 {
		data.Total = &opts.maxValue
	}

	return zenutil.RunProgress(opts.ctx, opts.maxValue, data)
}
