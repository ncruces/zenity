package zenity

import (
	"context"
	"image/color"
	"reflect"
	"testing"
)

func Test_applyOptions(t *testing.T) {
	tests := []struct {
		name string
		args Option
		want options
	}{
		// General options
		{name: "Title", args: Title("Title"), want: options{title: stringPtr("Title")}},
		{name: "Width", args: Width(100), want: options{width: 100}},
		{name: "Height", args: Height(100), want: options{height: 100}},
		{name: "OKLabel", args: OKLabel("OK"), want: options{okLabel: stringPtr("OK")}},
		{name: "CancelLabel", args: CancelLabel("Cancel"), want: options{cancelLabel: stringPtr("Cancel")}},
		{name: "ExtraButton", args: ExtraButton("Extra"), want: options{extraButton: stringPtr("Extra")}},
		{name: "Icon", args: Icon(ErrorIcon), want: options{icon: ErrorIcon}},
		{name: "DefaultCancel", args: DefaultCancel(), want: options{defaultCancel: true}},

		// Message options
		{name: "NoWrap", args: NoWrap(), want: options{noWrap: true}},
		{name: "Ellipsize", args: Ellipsize(), want: options{ellipsize: true}},

		// Entry options
		{name: "EntryText", args: EntryText("text"), want: options{entryText: "text"}},
		{name: "HideText", args: HideText(), want: options{hideText: true}},
		{name: "Username", args: Username(), want: options{username: true}},

		// List options
		{name: "DisallowEmpty", args: DisallowEmpty(), want: options{disallowEmpty: true}},
		{name: "DefaultItems", args: DefaultItems("a", "b"), want: options{defaultItems: []string{"a", "b"}}},

		// File selection options
		{name: "Directory", args: Directory(), want: options{directory: true}},
		{name: "ConfirmOverwrite", args: ConfirmOverwrite(), want: options{confirmOverwrite: true}},
		{name: "ConfirmCreate", args: ConfirmCreate(), want: options{confirmCreate: true}},
		{name: "ShowHidden", args: ShowHidden(), want: options{showHidden: true}},
		{name: "Filename", args: Filename("file.go"), want: options{filename: "file.go"}},
		{name: "FileFilters", args: FileFilter{"Go files", []string{"*.go"}}, want: options{
			fileFilters: FileFilters{{"Go files", []string{"*.go"}}},
		}},

		// Progress indication options
		{name: "MaxValue", args: MaxValue(100), want: options{maxValue: 100}},
		{name: "NoCancel", args: NoCancel(), want: options{noCancel: true}},
		{name: "TimeRemaining", args: TimeRemaining(), want: options{timeRemaining: true}},

		// Color selection options
		{name: "Color", args: Color(color.Black), want: options{color: color.Black}},
		{name: "ShowPalette", args: ShowPalette(), want: options{showPalette: true}},

		// Context for timeout
		{name: "Context", args: Context(context.TODO()), want: options{ctx: context.TODO()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyOptions([]Option{tt.args}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
