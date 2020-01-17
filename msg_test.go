package zenity_test

import "github.com/ncruces/zenity"

func ExampleError() {
	zenity.Error("An error has occurred.",
		zenity.Title("Error"),
		zenity.Icon(zenity.ErrorIcon))
	// Output:
}

func ExampleInfo() {
	zenity.Info("All updates are complete.",
		zenity.Title("Information"),
		zenity.Icon(zenity.InfoIcon))
	// Output:
}

func ExampleWarning() {
	zenity.Warning("Are you sure you want to proceed?",
		zenity.Title("Warning"),
		zenity.Icon(zenity.WarningIcon))
	// Output:
}

func ExampleQuestion() {
	zenity.Question("Are you sure you want to proceed?",
		zenity.Title("Question"),
		zenity.Icon(zenity.QuestionIcon))
	// Output:
}
