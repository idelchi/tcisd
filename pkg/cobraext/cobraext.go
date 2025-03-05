package cobraext

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/pkg/pretty"
)

// ErrExitGracefully is an error that signals the program to exit gracefully.
var ErrExitGracefully = errors.New("exit")

// Validator is an interface for types that can show and validate their configuration.
type Validator interface {
	Validate() error
	Display() bool
}

// NewDefaultRootCommand creates a root command with default settings.
// It sets up integration with viper, with environment variable and flag binding.
func NewDefaultRootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}

// MustBindPFlags binds a command's flags to viper.
// It panics if it encounters an error.
func MustBindPFlags(cmd *cobra.Command, cfg interface{}) {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(fmt.Sprintf("binding flags: %v", err))
	}
}

// Validate unmarshals the configuration and performs validation checks.
// If cfg.Show is true, prints the configuration and exits.
func Validate(cfg Validator, validations ...interface{}) error {
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	if cfg.Display() {
		pretty.PrintJSONMasked(cfg)
		return ErrExitGracefully
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}
