package zenity_test

import "github.com/ncruces/zenity"

const defaultPath = ""
const defaultName = ""

func ExampleSelectFile() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		}.Build())
	// Output:
}

func ExampleSelectFileMutiple() {
	zenity.SelectFileMutiple(
		zenity.Filename(defaultPath),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		}.Build())
	// Output:
}

func ExampleSelectFileSave() {
	zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename(defaultName),
		zenity.FileFilters{
			{"Go files", []string{"*.go"}},
			{"Web files", []string{"*.html", "*.js", "*.css"}},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}},
		}.Build())
	// Output:
}

func ExampleSelectFile_directory() {
	zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory())
	// Output:
}

func ExampleSelectFileMutiple_directory() {
	zenity.SelectFileMutiple(
		zenity.Filename(defaultPath),
		zenity.Directory())
	// Output:
}
