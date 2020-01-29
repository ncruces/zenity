package zenity_test

import (
	"context"
	"errors"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleSelectColor() {
	zenity.SelectColor(
		zenity.Color(color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0x80}))
	// Output:
}

func ExampleSelectColor_palette() {
	zenity.SelectColor(
		zenity.ShowPalette(),
		zenity.Color(color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff}))
	// Output:
}

func TestSelectColorTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

	_, err := zenity.SelectColor(zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}

	cancel()
}

func TestSelectColorCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.SelectColor(zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
