package matcher

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Logger is an interface for logging formatted messages.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Matcher is a file matcher that compiles a list of files matching a given pattern.
type Matcher struct {
	// Exclude is a list of patterns that are used to exclude files.
	Exclude []string
	// Logger is a logger for debug messages.
	Logger Logger
	// files is the list of files that are added to the matcher.
	files []string
	// Additional file exclusion functions
	extraExcludes map[string]func(string) bool
}

// ListFiles returns all files found by the Matcher.
func (m *Matcher) ListFiles() []string {
	if m.files == nil {
		return []string{}
	}

	return m.files
}

// New creates a Matcher with default settings.
func New(hidden bool, exclude []string, logger Logger) Matcher {
	matcher := Matcher{
		Exclude: exclude,
		Logger:  logger,
	}

	// Get the name of the executable itself
	if exe, err := os.Executable(); err == nil {
		// Exclude the executable itself
		matcher.Exclude = append(matcher.Exclude, exe)
	}

	matcher.Exclude = append(matcher.Exclude,
		// Exclude all kinds of executables
		"**/*.exe",

		// Exclude some known directories
		"**/.git/**",
		"**/node_modules/**",
		"**/vendor/**",
		"**/.task/**",
		"**/.cache/**",
	)

	if !hidden {
		// Exclude hidden folders & files if hidden is false
		matcher.Exclude = append(matcher.Exclude, "**/.*", "**/.*/**/*")
	}

	matcher.extraExcludes = map[string]func(string) bool{
		"binary file": IsBinary,
	}

	return matcher
}

// Match finds all files matching the given pattern and applies the exclusion options.
func (m *Matcher) Match(pattern string) (err error) {
	// Get all files that match the pattern
	var matches []string

	if matches, err = doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly()); err != nil {
		return fmt.Errorf("matching pattern %q: %w", pattern, err)
	}

outer:
	for _, match := range matches {
		// Convert to absolute path
		match, _ = filepath.Abs(match)
		match = filepath.ToSlash(match)

		switch {
		// Skip files that are already found
		case contains(match, m.files):
			m.Logger.Printf("Skipped %q: already in matches", match)
		// If the file is explicitly included, include it immediately
		case IsExplicitlyIncluded(pattern):
			m.Logger.Printf("Including %q: explicitly included", match)
			m.files = append(m.files, match)
		case IsExcluded(match, m.Exclude) != "":
			m.Logger.Printf("Skipped %q: matches exclude pattern %q", match, IsExcluded(match, m.Exclude))
		default:
			for name, fn := range m.extraExcludes {
				if fn(match) {
					m.Logger.Printf("Skipped %q: %s", match, name)
					continue outer
				}
			}
			// Append the match to the matches slice
			m.files = append(m.files, match)
		}
	}

	return nil
}

// IsBinary checks if a file is binary.
func IsBinary(file string) bool {
	// Simple binary detection by reading the first 512 bytes
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil {
		return false
	}

	// Check for null bytes in the first 512 bytes
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

// IsExplicitlyIncluded returns true if the given pattern is considered to be explicitly included.
func IsExplicitlyIncluded(pattern string) bool {
	return !strings.Contains(pattern, "*")
}

// IsExcluded returns the exclude pattern that the given file matches, or an empty string if none.
func IsExcluded(file string, excludes []string) string {
	for _, pattern := range excludes {
		if matched, _ := doublestar.Match(pattern, file); matched {
			return pattern
		}
	}
	return ""
}

// contains checks if a string is in a slice of strings.
func contains(str string, slice []string) bool {
	return slices.Contains(slice, str)
}
