package zenutil

import (
	"fmt"
	"image/color"
	"strings"

	"golang.org/x/image/colornames"
)

// ParseColor is internal.
func ParseColor(s string) color.Color {
	if len(s) == 4 || len(s) == 5 {
		c := color.NRGBA{A: 0xf}
		n, _ := fmt.Sscanf(s, "#%1x%1x%1x%1x", &c.R, &c.G, &c.B, &c.A)
		c.R, c.G, c.B, c.A = c.R*0x11, c.G*0x11, c.B*0x11, c.A*0x11
		if n >= 3 {
			return c
		}
	}

	if len(s) == 7 || len(s) == 9 {
		c := color.NRGBA{A: 0xff}
		n, _ := fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
		if n >= 3 {
			return c
		}
	}

	if len(s) >= 10 {
		c := color.NRGBA{A: 0xff}
		if _, err := fmt.Sscanf(s, "rgb(%d,%d,%d)", &c.R, &c.G, &c.B); err == nil {
			return c
		}

		var a float32
		if _, err := fmt.Sscanf(s, "rgba(%d,%d,%d,%f)", &c.R, &c.G, &c.B, &a); err == nil {
			c.A = uint8(255*a + 0.5)
			return c
		}
	}

	c, ok := colornames.Map[strings.ToLower(s)]
	if ok {
		return c
	}
	return nil
}

// UnparseColor is internal.
func UnparseColor(c color.Color) string {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	if n.A == 255 {
		return fmt.Sprintf("rgb(%d,%d,%d)", n.R, n.G, n.B)
	} else {
		return fmt.Sprintf("rgba(%d,%d,%d,%f)", n.R, n.G, n.B, float32(n.A)/255)
	}
}
