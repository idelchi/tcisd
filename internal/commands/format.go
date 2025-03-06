package commands

import (
	"errors"
	"log"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/internal/processor"
)

func NewFormatCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format [flags] [path ...]",
		Short: "Strip comments from files",
		Long:  "Strip comments from files, modifying them in-place.",
		Args:  cobra.ArbitraryArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("at least one path must be provided")
			}

			cfg.Paths = args
			cfg.Mode = config.FormatMode

			return cobraext.Validate(cfg, cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			proc := processor.New(
				cfg,
				min(runtime.NumCPU(), len(cfg.Paths)),
				cfg.Types,
			)

			if err := proc.Process(); err != nil {
				return err
			}

			hasIssues := proc.Summary()

			if !hasIssues {
				log.Println(color.GreenString("No comments found in files"))
			}

			return nil
		},
	}

	cmd.Flags().StringArrayP("pattern", "p", nil, "File pattern to match (doublestar format)")
	cmd.Flags().StringArrayP("type", "t", nil, "File types to process (go, python, dockerfile)")
	cmd.Flags().StringArrayP("exclude", "e", nil, "Patterns to exclude")
	cmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")
	cmd.Flags().BoolP("dry-run", "d", false, "Show what would be changed without modifying files")

	return cmd
}
