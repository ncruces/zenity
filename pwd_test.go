package zenity_test

import (
	"context"
	"errors"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
)

func ExamplePassword() {
	zenity.Password(zenity.Title("Type your password"))
	// Output:
}

func ExamplePassword_username() {
	zenity.Password(
		zenity.Title("Type your username and password"),
		zenity.Username())
	// Output:
}

func TestPassword_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, _, err := zenity.Password(zenity.Context(ctx))
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestPassword_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := zenity.Password(zenity.Context(ctx))
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestPassword_username(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := zenity.Password(zenity.Context(ctx), zenity.Username())
	if err, skip := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		if !errors.Is(err, zenity.ErrUnsupported) {
			t.Error("was not unsupported:", err)
		}
	} else {
		if !errors.Is(err, context.Canceled) {
			t.Error("was not canceled:", err)
		}
	}
}
