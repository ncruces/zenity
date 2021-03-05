package zenity

import (
	"os"
	"path/filepath"
)

// SelectFile displays the file selection dialog.
//
// Returns an empty string on cancel.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
func SelectFile(options ...Option) (string, error) {
	return selectFile(applyOptions(options))
}

// SelectFileMutiple displays the multiple file selection dialog.
//
// Returns a nil slice on cancel.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
func SelectFileMutiple(options ...Option) ([]string, error) {
	return selectFileMutiple(applyOptions(options))
}

// SelectFileSave displays the save file selection dialog.
//
// Returns an empty string on cancel.
//
// Valid options: Title, Filename, ConfirmOverwrite, ConfirmCreate, ShowHidden,
// FileFilter(s).
func SelectFileSave(options ...Option) (string, error) {
	return selectFileSave(applyOptions(options))
}

// Directory returns an Option to activate directory-only selection.
func Directory() Option {
	return funcOption(func(o *options) { o.directory = true })
}

// ConfirmOverwrite returns an Option to confirm file selection if filename
// already exists.
func ConfirmOverwrite() Option {
	return funcOption(func(o *options) { o.confirmOverwrite = true })
}

// ConfirmCreate returns an Option to confirm file selection if filename does
// not yet exist (Windows only).
func ConfirmCreate() Option {
	return funcOption(func(o *options) { o.confirmCreate = true })
}

// ShowHidden returns an Option to show hidden files (Windows and macOS only).
func ShowHidden() Option {
	return funcOption(func(o *options) { o.showHidden = true })
}

// Filename returns an Option to set the filename.
//
// You can specify a file name, a directory path, or both.
// Specifying a file name, makes it the default selected file.
// Specifying a directory path, makes it the default dialog location.
func Filename(filename string) Option {
	return funcOption(func(o *options) { o.filename = filename })
}

// FileFilter is an Option that sets a filename filter.
//
// macOS hides filename filters from the user,
// and only supports filtering by extension (or "type").
type FileFilter struct {
	Name     string   // display string that describes the filter (optional)
	Patterns []string // filter patterns for the display string
}

func (f FileFilter) apply(o *options) {
	o.fileFilters = append(o.fileFilters, f)
}

// FileFilters is an Option that sets multiple filename filters.
type FileFilters []FileFilter

func (f FileFilters) apply(o *options) {
	o.fileFilters = append(o.fileFilters, f...)
}

func splitDirAndName(path string) (dir, name string) {
	if path != "" {
		path = filepath.Clean(path)
		fi, err := os.Stat(path)
		if err == nil && fi.IsDir() {
			return path, ""
		}
	}
	return filepath.Split(path)
}
