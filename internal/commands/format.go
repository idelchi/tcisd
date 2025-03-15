package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/tcisd/internal/config"
)

func NewFormatCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format [flags] [path ...]",
		Short: "Strip comments from files",
		Long:  "Strip comments from files, modifying them in-place.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	return cmd
}
