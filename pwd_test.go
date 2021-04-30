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

func ExamplePassword() {
	zenity.Password(zenity.Title("Type your password"))
	// Output:
}

func TestPassword_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, _, err := zenity.Password(zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestPassword_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := zenity.Password(zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
