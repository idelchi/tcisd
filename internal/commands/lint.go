package commands

import (
	"fmt"
	"log"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/internal/processor"
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
			// Set up processor with appropriate number of workers
			proc := processor.New(
				cfg,
				min(runtime.NumCPU(), len(cfg.Paths)),
				cfg.Types,
			)

			// Process the files
			if err := proc.Process(); err != nil {
				return err
			}

			// Print summary
			hasIssues := proc.Summary()

			if hasIssues {
				log.Println(color.RedString("Comments found in files"))
				return fmt.Errorf("comments found")
			}

			log.Println(color.GreenString("No comments found in files"))
			return nil
		},
	}

	cmd.Flags().StringArrayP("pattern", "p", []string{"**/*.go"}, "File pattern to match (doublestar format)")
	cmd.Flags().StringArrayP("type", "t", []string{"go"}, "File types to process (go, python)")
	cmd.Flags().StringArrayP("exclude", "e", nil, "Patterns to exclude")
	cmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")

	return cmd
}
