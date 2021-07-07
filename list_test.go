package zenity_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ncruces/zenity"
	"go.uber.org/goleak"
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

func TestList_timeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	_, err := zenity.List("", nil, zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !os.IsTimeout(err) {
		t.Error("did not timeout:", err)
	}
}

func TestList_cancel(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := zenity.List("", nil, zenity.Context(ctx))
	if skip, err := skip(err); skip {
		t.Skip("skipping:", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Error("was not canceled:", err)
	}
}

func TestList_script(t *testing.T) {
	items := []string{"apples", "oranges", "bananas", "strawberries"}
	tests := []struct {
		name string
		call string
		opts []zenity.Option
		want string
		err  error
	}{
		{name: "Cancel", call: "cancel", want: "", err: zenity.ErrCanceled},
		{name: "Apples", call: "select apples", want: "apples", err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := zenity.List(fmt.Sprintf("Please, %s.", tt.call), items, tt.opts...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if text != tt.want || err != tt.err {
				t.Errorf("List() = %q, %v; want %q, %v", text, err, tt.want, tt.err)
			}
		})
	}
}

func TestListMultiple_script(t *testing.T) {
	items := []string{"apples", "oranges", "bananas", "strawberries"}
	tests := []struct {
		name string
		call string
		opts []zenity.Option
		want []string
		err  error
	}{
		{name: "Cancel", call: "cancel", want: nil, err: zenity.ErrCanceled},
		{name: "Nothing", call: "select nothing", want: []string{}, err: nil},
		{name: "Apples", call: "select apples", want: []string{"apples"}, err: nil},
		{name: "Apples & Oranges", call: "select apples and oranges",
			want: []string{"apples", "oranges"}, err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zenity.ListMultiple(fmt.Sprintf("Please, %s.", tt.call), items, tt.opts...)
			if skip, err := skip(err); skip {
				t.Skip("skipping:", err)
			}
			if !reflect.DeepEqual(got, tt.want) || err != tt.err {
				t.Errorf("ListMultiple() = %q, %v; want %v, %v", got, err, tt.want, tt.err)
			}
		})
	}
}
