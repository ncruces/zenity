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
	return selectFile(options...)
}

// SelectFileMutiple displays the multiple file selection dialog.
//
// Returns a nil slice on cancel.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
func SelectFileMutiple(options ...Option) ([]string, error) {
	return selectFileMutiple(options...)
}

// SelectFileSave displays the save file selection dialog.
//
// Returns an empty string on cancel.
//
// Valid options: Title, Filename, ConfirmOverwrite, ConfirmCreate, ShowHidden,
// FileFilter(s).
func SelectFileSave(options ...Option) (string, error) {
	return selectFileSave(options...)
}

// Filename returns an Option to set the filename.
//
// You can specify a file name, a directory path, or both.
// Specifying a file name, makes it the default selected file.
// Specifying a directory path, makes it the default dialog location.
func Filename(filename string) Option {
	return func(o *options) { o.filename = filename }
}

// Directory returns an Option to activate directory-only selection.
func Directory() Option {
	return func(o *options) { o.directory = true }
}

// ConfirmOverwrite returns an Option to confirm file selection if filename
// already exists.
func ConfirmOverwrite() Option {
	return func(o *options) { o.overwrite = true }
}

// ConfirmCreate returns an Option to confirm file selection if filename does
// not yet exist (Windows only).
func ConfirmCreate() Option {
	return func(o *options) { o.create = true }
}

// ShowHidden returns an Option to show hidden files (Windows and macOS only).
func ShowHidden() Option {
	return func(o *options) { o.hidden = true }
}

// FileFilter encapsulates a filename filter.
//
// macOS hides filename filters from the user,
// and only supports filtering by extension (or "type").
type FileFilter struct {
	Name     string   // display string that describes the filter (optional)
	Patterns []string // filter patterns for the display string
}

// Build returns an Option to set a filename filter.
func (f FileFilter) Build() Option {
	return func(o *options) { o.filters = append(o.filters, f) }
}

// FileFilters is a list of filename filters.
type FileFilters []FileFilter

// Build returns an Option to set filename filters.
func (f FileFilters) Build() Option {
	return func(o *options) { o.filters = append(o.filters, f...) }
}

func splitDirAndName(path string) (dir, name string) {
	path = filepath.Clean(path)
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return path, ""
	}
	return filepath.Split(path)
}
