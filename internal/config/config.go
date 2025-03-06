package config

import (
	"errors"
	"fmt"
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

	Patterns []string `mapstructure:"pattern"`

	Types []string `mapstructure:"type"`

	Exclude []string

	Hidden bool

	DryRun bool `mapstructure:"dry-run"`

	Paths []string
}

func (c Config) Display() bool {
	return c.Show
}

func (c *Config) Validate(config any) error {
	validTypes := map[string]bool{
		"go":         true,
		"python":     true,
		"dockerfile": true,
	}

	if len(c.Types) == 0 {
		c.Types = []string{"go", "python", "dockerfile"}
	} else {
		for _, t := range c.Types {
			if !validTypes[t] {
				return fmt.Errorf("%w: invalid file type: %s", ErrUsage, t)
			}
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

	if len(c.Paths) == 0 {
		return fmt.Errorf("%w: at least one path must be provided", ErrUsage)
	}

	return nil
}
