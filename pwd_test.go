package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExamplePassword() {
	zenity.Password(zenity.Title("Type your password"))
	// Output:
}

func TestPasswordTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

	_, _, err := zenity.Password(zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}

	cancel()
}

func TestPasswordCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := zenity.Password(zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
