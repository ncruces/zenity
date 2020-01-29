package zenity_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ncruces/zenity"
)

func ExampleNotify() {
	zenity.Notify("There are system updates necessary!",
		zenity.Title("Warning"),
		zenity.Icon(zenity.InfoIcon))
	// Output:
}

func TestNotifyCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := zenity.Notify("text", zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
