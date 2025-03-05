package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/pkg/cobraext"
)

// NewLintCommand creates the lint subcommand.
// It handles checking files for comments without modifying them.
func NewLintCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [flags] [path ...]",
		Short: "Check files for comments",
		Long:  "Check files for comments without modifying them.",
		Args:  cobra.ArbitraryArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one path must be provided")
			}

			cfg.Paths = args
			cfg.Mode = config.LintMode

			return cobraext.Validate(cfg, cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil // Processing happens in the parent logic
		},
	}

	cmd.Flags().StringArrayP("pattern", "p", []string{"**/*.go"}, "File pattern to match (doublestar format)")
	cmd.Flags().StringArrayP("type", "t", []string{"go"}, "File types to process (go, bash, python)")
	cmd.Flags().StringArrayP("exclude", "e", nil, "Patterns to exclude")
	cmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")

	// Bind flags to config
	cobraext.MustBindPFlags(cmd, cfg)

	return cmd
}
