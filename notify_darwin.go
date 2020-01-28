package zenity

import (
	"strings"

	"github.com/ncruces/zenity/internal/zenutil"
)

func notify(text string, options []Option) error {
	opts := applyOptions(options)
	data := zenutil.Notify{
		Text:  text,
		Title: opts.title,
	}
	if i := strings.IndexByte(text, '\n'); i >= 0 && i < len(text) {
		data.Subtitle = text[:i]
		data.Text = text[i+1:]
	}
	_, err := zenutil.Run(opts.ctx, "notify", data)
	if err != nil {
		return err
	}
	return nil
}
