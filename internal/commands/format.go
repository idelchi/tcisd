package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/config"
)

// NewFormatCommand creates the format subcommand.
// It handles stripping comments from files.
func NewFormatCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format [flags] [path ...]",
		Short: "Strip comments from files",
		Long:  "Strip comments from files, modifying them in-place.",
		Args:  cobra.ArbitraryArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one path must be provided")
			}

			cfg.Paths = args
			cfg.Mode = config.FormatMode

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
	cmd.Flags().BoolP("dry-run", "d", false, "Show what would be changed without modifying files")

	return cmd
}
