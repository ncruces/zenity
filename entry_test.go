package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleEntry() {
	zenity.Entry("Enter new text:",
		zenity.Title("Add a new entry"))
	// Output:
}

func TestEntryTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

	_, err := zenity.Entry("", zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}

	cancel()
}

func TestEntryCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.Entry("", zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
