package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/config"
)

func NewRootCommand(cfg *config.Config, version string) *cobra.Command {
	root := cobraext.NewDefaultRootCommand(version)

	root.Use = "tcisd [flags] command [flags]"
	root.Short = "Strip comments from code files"
	root.Long = "tcisd is a tool for stripping comments from code files.\n" +
		"It can verify if files have comments (lint mode) or remove them (format mode)."

	root.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	root.AddCommand(NewLintCommand(cfg), NewFormatCommand(cfg))

	root.CompletionOptions.DisableDefaultCmd = true

	return root
}
