package config

import (
	"errors"
	"fmt"
)

// Mode represents the operation mode of the application.
type Mode string

const (
	// LintMode checks files for comments without modifying them.
	LintMode Mode = "lint"
	// FormatMode strips comments from files.
	FormatMode Mode = "format"
)

// ErrUsage indicates an error in command-line usage or configuration.
var ErrUsage = errors.New("usage error")

// Config holds the application's configuration parameters.
type Config struct {
	// Show enables output display
	Show bool

	// Mode is the operation mode (lint or format)
	Mode Mode

	// Patterns is a list of file patterns to match
	Patterns []string `mapstructure:"pattern"`

	// Types is a list of file types to process
	Types []string `mapstructure:"type"`

	// Exclude is a list of patterns to exclude
	Exclude []string

	// Hidden indicates whether to include hidden files
	Hidden bool

	// DryRun indicates whether to show changes without modifying files
	DryRun bool `mapstructure:"dry-run"`

	// Paths is a list of paths to process
	Paths []string
}

// Display returns the value of the Show field.
func (c Config) Display() bool {
	return c.Show
}

// Validate performs configuration validation.
// It returns a wrapped ErrUsage if any validation rules are violated.
func (c Config) Validate(config any) error {
	if len(c.Patterns) == 0 {
		return fmt.Errorf("%w: at least one pattern must be provided", ErrUsage)
	}

	if len(c.Types) == 0 {
		return fmt.Errorf("%w: at least one file type must be provided", ErrUsage)
	}

	// Validate file types
	validTypes := map[string]bool{
		"go":     true,
		"bash":   true,
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
