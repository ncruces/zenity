package zenity

import (
	"fmt"

	"github.com/ncruces/zenity/internal/zenutil"
)

func progress(opts options) (ProgressDialog, error) {
	if opts.extraButton != nil {
		return nil, fmt.Errorf("%w: extra button", ErrUnsupported)
	}

	var data zenutil.Progress
	data.Description = opts.title
	if opts.maxValue == 0 {
		opts.maxValue = 100
	}
	if opts.maxValue >= 0 {
		data.Total = &opts.maxValue
	}
	if i, ok := opts.windowIcon.(string); ok {
		data.WindowIcon = i
	}

	return zenutil.RunProgress(opts.ctx, opts.maxValue, opts.autoClose, data)
}
