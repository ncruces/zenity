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

func ExampleEntry() {
	zenity.Entry("Enter new text:",
		zenity.Title("Add a new entry"))
	// Output:
}

func TestEntry_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, err := zenity.Entry("", zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestEntry_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Entry("", zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
