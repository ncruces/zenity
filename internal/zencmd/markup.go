package zencmd

import (
	"encoding/xml"
	"io"
	"strings"
)

// StripMarkup is internal.
func StripMarkup(s string) string {
	// Strips XML markup described in:
	// https://docs.gtk.org/Pango/pango_markup.html

	dec := xml.NewDecoder(strings.NewReader(s))
	var res strings.Builder

	for {
		t, err := dec.Token()
		if err == io.EOF {
			return res.String()
		}
		if err != nil {
			return s
		}
		if t, ok := t.(xml.CharData); ok {
			res.Write(t)
		}
	}
}
