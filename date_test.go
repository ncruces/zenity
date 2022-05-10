package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

func ExampleCalendar() {
	zenity.Calendar("Select a date from below:",
		zenity.DefaultDate(2006, time.January, 1))
	// Output:
}

func TestCalendar_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/5)
	defer cancel()

	_, err := zenity.Calendar("", zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestCalendar_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Calendar("", zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestCalendar_script(t *testing.T) {
	tests := []struct {
		name string
		call string
		opts []zenity.Option
		want time.Time
		err  error
	}{
		{name: "Cancel", call: "cancel", err: zenity.ErrCanceled},
		{name: "Black", call: "choose today", want: time.Now()},
		{name: "Rebecca", call: "press OK", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
			opts: []zenity.Option{zenity.DefaultDate(2000, 1, 1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zenity.Calendar(fmt.Sprintf("Please, %s.", tt.call), tt.opts...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			got.Date()
			if !DateEquals(got, tt.want) || err != tt.err {
				t.Errorf("Calendar() = %v, %v; want %v, %v", got, err, tt.want, tt.err)
			}
		})
	}
}

func DateEquals(t1, t2 time.Time) bool {
	d1, m1, y1 := t1.Date()
	d2, m2, y2 := t2.Date()
	return d1 == d2 && m1 == m2 && y1 == y2
}
