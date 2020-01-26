package zenity_test

import "github.com/ncruces/zenity"

func ExampleNotify() {
	zenity.Notify("An error has occurred.",
		zenity.Title("Error"),
		zenity.Icon(zenity.ErrorIcon))
	// Output:
}
