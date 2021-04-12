package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleCalendar() {
	zenity.Calendar("Select a date from below:",
		zenity.DefaultDate(2006, time.January, 1))
	// Output:
}

func TestCalendarTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

	_, err := zenity.Calendar("", zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}

	cancel()
}

func TestCalendarCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Calendar("", zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
