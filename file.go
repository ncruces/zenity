package zenity

import (
	"os"
	"path/filepath"
	"strings"
)

// SelectFile displays the file selection dialog.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
//
// May return: ErrCanceled.
func SelectFile(options ...Option) (string, error) {
	return selectFile(applyOptions(options))
}

// SelectFileMutiple displays the multiple file selection dialog.
//
// Valid options: Title, Directory, Filename, ShowHidden, FileFilter(s).
//
// May return: ErrCanceled, ErrUnsupported.
func SelectFileMutiple(options ...Option) ([]string, error) {
	return selectFileMutiple(applyOptions(options))
}

// SelectFileSave displays the save file selection dialog.
//
// Valid options: Title, Filename, ConfirmOverwrite, ConfirmCreate, ShowHidden,
// FileFilter(s).
//
// May return: ErrCanceled.
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
// and only supports filtering by extension
// (or "uniform type identifiers").
//
// Patterns may use the GTK syntax on all platforms:
// https://developer.gnome.org/pygtk/stable/class-gtkfilefilter.html#method-gtkfilefilter--add-pattern
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

// Windows' patterns need a name.
func (f FileFilters) name() {
	for i, filter := range f {
		if filter.Name == "" {
			f[i].Name = strings.Join(filter.Patterns, " ")
		}
	}
}

// Windows' patterns are case insensitive, don't support character classes or escaping.
//
// First we remove character classes, then escaping. Patterns with literal wildcards are invalid.
// The semicolon is a separator, so we replace it with the single character wildcard.
// Empty and invalid filters/patterns are ignored.
func (f FileFilters) simplify() {
	var i = 0
	for _, filter := range f {
		var j = 0
		for _, pattern := range filter.Patterns {
			var escape, invalid bool
			var buf strings.Builder
			for _, r := range removeClasses(pattern) {
				if !escape && r == '\\' {
					escape = true
					continue
				}
				if escape && (r == '*' || r == '?') {
					invalid = true
					break
				}
				if r == ';' {
					r = '?'
				}
				buf.WriteRune(r)
				escape = false
			}
			if buf.Len() > 0 && !invalid {
				filter.Patterns[j] = buf.String()
				j++
			}
		}
		if j > 0 {
			filter.Patterns = filter.Patterns[:j]
			f[i] = filter
			i++
		}
	}
	for ; i < len(f); i++ {
		f[i] = FileFilter{}
	}
}

// macOS types may be specified as extension strings without the leading period,
// or as uniform type identifiers:
// https://developer.apple.com/library/archive/documentation/LanguagesUtilities/Conceptual/MacAutomationScriptingGuide/PromptforaFileorFolder.html
//
// First check for uniform type identifiers.
// Then we extract the extension from each pattern, remove character classes, then escaping.
// If an extension contains a wildcard, any type is accepted.
func (f FileFilters) types() []string {
	var res []string
	for _, filter := range f {
		for _, pattern := range filter.Patterns {
			if isUniformTypeIdentifier(pattern) {
				res = append(res, pattern)
				continue
			}

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
			if buf.Len() > 0 {
				res = append(res, buf.String())
			}
		}
	}
	if res == nil {
		return nil
	}
	// Workaround for macOS bug: first type cannot be a four letter extension, so prepend empty string.
	return append([]string{""}, res...)
}

// Remove character classes from pattern, assuming case insensitivity.
// Classes of one character (case insensitive) are replaced by the character.
// Others are replaced by the single character wildcard.
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

// Uniform type identifiers use the reverse-DNS format:
// https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/understanding_utis/understand_utis_conc/understand_utis_conc.html
func isUniformTypeIdentifier(pattern string) bool {
	labels := strings.Split(pattern, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if len := len(label); len == 0 || label[0] == '-' || label[len-1] == '-' {
			return false
		}
		for _, r := range label {
			switch {
			case r == '-' || r > '\x7f' ||
				'a' <= r && r <= 'z' ||
				'A' <= r && r <= 'Z' ||
				'0' <= r && r <= '9':
				continue
			default:
				return false
			}
		}
	}

	return true
}

func splitDirAndName(path string) (dir, name string) {
	if path != "" {
		fi, err := os.Stat(path)
		if err == nil && fi.IsDir() {
			return path, ""
		}
	}
	return filepath.Split(path)
}
