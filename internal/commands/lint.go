package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/tcisd/internal/config"
)

func NewLintCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [flags] [path ...]",
		Short: "Check files for comments",
		Long:  "Check files for comments without modifying them.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	return cmd
}
