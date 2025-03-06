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

func (c Config) Validate(config any) error {
	if len(c.Patterns) == 0 {
		return fmt.Errorf("%w: at least one pattern must be provided", ErrUsage)
	}

	if len(c.Types) == 0 {
		return fmt.Errorf("%w: at least one file type must be provided", ErrUsage)
	}

	validTypes := map[string]bool{
		"go":     true,
		"python": true,
	}

	for _, t := range c.Types {
		if !validTypes[t] {
			return fmt.Errorf("%w: invalid file type: %s", ErrUsage, t)
		}
	}

	if len(c.Paths) == 0 {
		return fmt.Errorf("%w: at least one path must be provided", ErrUsage)
	}

	return nil
}
