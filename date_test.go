package zenity_test

import (
	"context"
	"errors"
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

func TestCalendarTimeout(t *testing.T) {
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

func TestCalendarCancel(t *testing.T) {
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
