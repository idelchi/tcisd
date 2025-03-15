package config

import (
	"errors"
	"fmt"
	"slices"
)

type Mode string

const (
	LintMode   Mode = "lint"
	FormatMode Mode = "format"
)

var ErrUsage = errors.New("usage error")

type Config struct {
	Show bool

	Mode Mode

	Patterns []string

	Types []string

	Exclude []string

	Hidden bool

	Parallel int

	Paths []string
}

func (c Config) Display() bool {
	return c.Show
}

func (c *Config) Validate() error {
	validTypes := []string{
		"go",
		"python",
		"dockerfile",
	}

	for _, t := range c.Types {
		if !slices.Contains(validTypes, t) {
			return fmt.Errorf("%w: invalid file type: %s", ErrUsage, t)
		}
	}

	if len(c.Patterns) == 0 {
		for _, t := range c.Types {
			switch t {
			case "go":
				c.Patterns = append(c.Patterns, "**/*.go")
			case "python":
				c.Patterns = append(c.Patterns, "**/*.py")
			case "dockerfile":
				c.Patterns = append(c.Patterns, "**/Dockerfile")
				c.Patterns = append(c.Patterns, "**/Dockerfile.*")
			}
		}
	}

	if c.Parallel < 1 {
		return fmt.Errorf("%w: invalid number of parallel jobs: %d", ErrUsage, c.Parallel)
	}

	return nil
}
