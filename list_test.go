package zenity_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ncruces/zenity"
)

func ExampleList() {
	zenity.List(
		"Select items from the list below:",
		[]string{"apples", "oranges", "bananas", "strawberries"},
		zenity.Title("Select items from the list"),
		zenity.DisallowEmpty(),
	)
	// Output:
}

func ExampleListItems() {
	zenity.ListItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
	// Output:
}

func ExampleListMultiple() {
	zenity.ListMultiple(
		"Select items from the list below:",
		[]string{"apples", "oranges", "bananas", "strawberries"},
		zenity.Title("Select items from the list"),
		zenity.DefaultItems("apples", "bananas"),
	)
	// Output:
}

func ExampleListMultipleItems() {
	zenity.ListMultipleItems(
		"Select items from the list below:",
		"apples", "oranges", "bananas", "strawberries")
	// Output:
}

func TestListTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)

	_, _, err := zenity.List("", nil, zenity.Context(ctx))
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}

	cancel()
}

func TestListCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := zenity.List("", nil, zenity.Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}
