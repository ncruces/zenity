package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"github.com/ncruces/zenity/internal/zenutil"
	"go.uber.org/goleak"
)

func ExampleSelectColor() {
	zenity.SelectColor(
		zenity.Color(color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0x80}))
}

func ExampleSelectColor_palette() {
	zenity.SelectColor(
		zenity.ShowPalette(),
		zenity.Color(color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff}))
}

func TestSelectColor_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/5)
	defer cancel()

	_, err := zenity.SelectColor(zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestSelectColor_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.SelectColor(zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestSelectColor_script(t *testing.T) {
	tests := []struct {
		name string
		call string
		opts []zenity.Option
		want color.Color
		err  error
	}{
		{name: "Cancel", call: "cancel", err: zenity.ErrCanceled},
		{name: "Black", call: "choose black", want: color.Black},
		{name: "White", call: "choose white", want: color.White},
		{name: "Rebecca", call: "press OK", want: color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff},
			opts: []zenity.Option{zenity.ShowPalette(), zenity.Color(color.NRGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xff})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zenity.Info(fmt.Sprintf("In the color selection dialog, %s.", tt.call))
			got, err := zenity.SelectColor(tt.opts...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if !zenutil.ColorEquals(got, tt.want) || err != tt.err {
				t.Errorf("SelectColor() = %v, %v; want %v, %v", got, err, tt.want, tt.err)
			}
		})
	}
}
