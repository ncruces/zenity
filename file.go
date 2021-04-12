package zenity

import (
	"os"
	"path/filepath"
	"strings"
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

// macOS filters by case insensitive literal extension.
// Extract all literal extensions from all patterns.
// If those contain wildcards, or classes with more than one character, accept anything.
func (f FileFilters) darwin() []string {
	var res []string
	for _, filter := range f {
		for _, pattern := range filter.Patterns {
			ext := pattern[strings.LastIndexByte(pattern, '.')+1:]

			var escape bool
			var buf strings.Builder
			for _, r := range removeClasses(ext) {
				switch {
				case escape:
					escape = false
				case r == '\\':
					escape = true
					continue
				case r == '*' || r == '?':
					return nil
				}
				buf.WriteRune(r)
			}
			res = append(res, buf.String())
		}
	}
	return res
}

func removeClasses(pattern string) string {
	var res strings.Builder
	for {
		i, j := findClass(pattern)
		if i < 0 {
			res.WriteString(pattern)
			return res.String()
		}
		res.WriteString(pattern[:i])

		var char string
		var escape, many bool
		for _, r := range pattern[i+1 : j-1] {
			if escape {
				escape = false
			} else if r == '\\' {
				escape = true
				continue
			}
			if char == "" {
				char = string(r)
			} else if !strings.EqualFold(char, string(r)) {
				many = true
				break
			}
		}
		if many {
			res.WriteByte('?')
		} else {
			res.WriteByte('\\')
			res.WriteString(char)
		}
		pattern = pattern[j:]
	}
}

func findClass(pattern string) (start, end int) {
	start = -1
	escape := false
	for i, b := range []byte(pattern) {
		switch {
		case escape:
			escape = false
		case b == '\\':
			escape = true
		case start < 0:
			if b == '[' {
				start = i
			}
		case 0 <= start && start < i-1:
			if b == ']' {
				return start, i + 1
			}
		}
	}
	return -1, -1
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
