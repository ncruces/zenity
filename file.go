package zenity

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// SelectFile displays the file selection dialog.
//
// Valid options: Title, WindowIcon, Attach, Modal, Directory, Filename,
// ShowHidden, FileFilter(s).
//
// May return: ErrCanceled.
func SelectFile(options ...Option) (string, error) {
	return selectFile(applyOptions(options))
}

// SelectFileMultiple displays the multiple file selection dialog.
//
// Valid options: Title, WindowIcon, Attach, Modal, Directory, Filename,
// ShowHidden, FileFilter(s).
//
// May return: ErrCanceled, ErrUnsupported.
func SelectFileMultiple(options ...Option) ([]string, error) {
	return selectFileMultiple(applyOptions(options))
}

// SelectFileSave displays the save file selection dialog.
//
// Valid options: Title, WindowIcon, Attach, Modal, Filename,
// ConfirmOverwrite, ConfirmCreate, ShowHidden, FileFilter(s).
//
// May return: ErrCanceled.
func SelectFileSave(options ...Option) (string, error) {
	return selectFileSave(applyOptions(options))
}

// Directory returns an Option to activate directory-only selection.
func Directory() Option {
	return funcOption(func(o *options) { o.directory = true })
}

// ConfirmOverwrite returns an Option to confirm file selection if the file
// already exists.
func ConfirmOverwrite() Option {
	return funcOption(func(o *options) { o.confirmOverwrite = true })
}

// ConfirmCreate returns an Option to confirm file selection if the file
// does not yet exist (Windows only).
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
// On Windows and macOS filtering is always case-insensitive.
//
// macOS hides filename filters from the user,
// and only supports filtering by extension
// (or "uniform type identifiers").
//
// Patterns may use the fnmatch syntax on all platforms:
// https://docs.python.org/3/library/fnmatch.html
type FileFilter struct {
	Name     string   // display string that describes the filter (optional)
	Patterns []string // filter patterns for the display string
	CaseFold bool     // if set patterns are matched case-insensitively
}

func (f FileFilter) apply(o *options) {
	o.fileFilters = append(o.fileFilters, f)
}

// FileFilters is an Option that sets multiple filename filters.
type FileFilters []FileFilter

func (f FileFilters) apply(o *options) {
	o.fileFilters = append(o.fileFilters, f...)
}

// Windows patterns need a name.
func (f FileFilters) name() {
	for i, filter := range f {
		if filter.Name == "" {
			f[i].Name = strings.Join(filter.Patterns, " ")
		}
	}
}

// Windows patterns are case-insensitive, don't support character classes or escaping.
//
// First we remove character classes, then escaping. Patterns with literal wildcards are invalid.
// The semicolon is a separator, so we replace it with the single character wildcard.
func (f FileFilters) simplify() {
	for i := range f {
		var j = 0
		for _, pattern := range f[i].Patterns {
			var escape, invalid bool
			var buf strings.Builder
			for _, b := range []byte(removeClasses(pattern)) {
				if !escape && b == '\\' {
					escape = true
					continue
				}
				if escape && (b == '*' || b == '?') {
					invalid = true
					break
				}
				if b == ';' {
					b = '?'
				}
				buf.WriteByte(b)
				escape = false
			}
			if buf.Len() > 0 && !invalid {
				f[i].Patterns[j] = buf.String()
				j++
			}
		}
		if j != 0 {
			f[i].Patterns = f[i].Patterns[:j]
		} else {
			f[i].Patterns = nil
		}
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
			dot := strings.LastIndexByte(pattern, '.')
			if dot < 0 {
				continue
			}

			var escape bool
			var buf strings.Builder
			for _, b := range []byte(removeClasses(pattern[dot+1:])) {
				switch {
				case escape:
					escape = false
				case b == '\\':
					escape = true
					continue
				case b == '*' || b == '?':
					return nil
				}
				buf.WriteByte(b)
			}
			res = append(res, buf.String())
		}
	}
	if res == nil {
		return nil
	}
	// Workaround for macOS bug: first type cannot be a four letter extension, so prepend dot string.
	return append([]string{"."}, res...)
}

// Unix patterns are case-sensitive. Fold them if requested.
func (f FileFilters) casefold() {
	for i := range f {
		if !f[i].CaseFold {
			continue
		}
		for j, pattern := range f[i].Patterns {
			var class = -1
			var escape bool
			var buf strings.Builder
			for i, r := range pattern {
				switch {
				case escape:
					escape = false
				case r == '\\':
					escape = true
				case class < 0:
					if r == '[' {
						class = i
					}
				case class < i-1:
					if r == ']' {
						class = -1
					}
				}

				nr := unicode.SimpleFold(r)
				if r == nr {
					buf.WriteRune(r)
					continue
				}

				if class < 0 {
					buf.WriteByte('[')
				}
				buf.WriteRune(r)
				for r != nr {
					buf.WriteRune(nr)
					nr = unicode.SimpleFold(nr)
				}
				if class < 0 {
					buf.WriteByte(']')
				}
			}
			f[i].Patterns[j] = buf.String()
		}
	}
}

// Remove character classes from pattern, assuming case insensitivity.
// Classes of one character (case-insensitive) are replaced by the character.
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

// Find a character class in the pattern.
func findClass(pattern string) (start, end int) {
	start = -1
	var escape bool
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
		case start < i-1:
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

func splitDirAndName(path string) (dir, name string, err error) {
	if path == "" {
		return "", "", nil
	}
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return path, "", nil
	}
	dir, name = filepath.Split(path)
	if dir == "" {
		return "", name, nil
	}
	_, err = os.Stat(dir)
	return dir, name, err
}
