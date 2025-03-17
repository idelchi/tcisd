package matcher

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/bmatcuk/doublestar/v4"
)

type Globber struct {
	Exclude []string
	files   []string
}

func (m *Globber) List() []string {
	return m.files
}

func New(hidden bool, excludes []string) Globber {
	matcher := Globber{
		Exclude: excludes,
		files:   []string{},
	}

	if !hidden {
		matcher.Exclude = append(matcher.Exclude, "**/.*", "**/.*/**/*")
	}

	return matcher
}

func (m *Globber) Match(pattern string) (err error) {
	var matches []string

	if matches, err = doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly()); err != nil {
		return fmt.Errorf("matching pattern %q: %w", pattern, err)
	}

	for _, match := range matches {
		match, _ = filepath.Abs(match)
		match = filepath.ToSlash(match)

		if slices.Contains(m.files, match) || IsExcluded(match, m.Exclude) != "" {
			continue
		}

		m.files = append(m.files, match)
	}

	return nil
}

func IsExcluded(file string, excludes []string) (pattern string) {
	for _, pattern := range excludes {
		if matched, _ := doublestar.Match(pattern, file); matched {
			return pattern
		}
	}

	return
}
