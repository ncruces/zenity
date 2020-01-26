package zenity_test

import "github.com/ncruces/zenity"

func ExampleNotify() {
	zenity.Notify("There are system updates necessary!",
		zenity.Title("Warning"),
		zenity.Icon(zenity.InfoIcon))
	// Output:
}
