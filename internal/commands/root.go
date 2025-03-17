package commands

import (
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/internal/processor"
)

func NewRootCommand(cfg *config.Config, version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "tcisd [flags] command [flags] [pattern ...]",
		Short: "Strip comments from code files",
		Long: heredoc.Doc(`
		tcisd is a tool for stripping comments from code files.
		It can verify if files have comments (lint mode) or remove them (format mode).
		Patterns defaults to ['**'].`),
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Root().Flags()); err != nil {
				return fmt.Errorf("binding flags: %w", err)
			}

			if err := viper.Unmarshal(cfg); err != nil {
				return fmt.Errorf("unmarshalling config for %q: %w", cmd.Name(), err)
			}

			cfg.Paths = args
			if len(args) == 0 {
				cfg.Paths = []string{"**"}
			}

			switch cmd.Name() {
			case "lint":
				cfg.Mode = config.LintMode
			case "format":
				cfg.Mode = config.FormatMode
			}

			if cfg.Show {
				pretty.PrintYAML(cfg)

				return nil
			}

			if err := cfg.Validate(); err != nil {
				return err
			}

			proc := processor.New(
				cfg,
			)

			if err := proc.Process(); err != nil {
				return err
			}

			hasIssues := proc.Summary()

			if hasIssues {
				log.Println(color.RedString("Comments found in files"))

				return errors.New("comments found")
			}

			log.Println(color.GreenString("No comments found in files"))

			return nil
		},
		RunE: cobraext.UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	root.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	root.Flags().StringArrayP("types", "t", []string{"go", "python", "dockerfile"}, "File types to process (go, python, dockerfile)")
	root.Flags().StringArrayP("exclude", "e", []string{"**/.git/**"}, "Patterns to exclude")
	root.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")
	root.Flags().IntP("parallel", "j", runtime.NumCPU(), "Number of concurrent jobs (default: number of CPUs)")

	root.AddCommand(NewLintCommand(cfg), NewFormatCommand(cfg))

	root.CompletionOptions.DisableDefaultCmd = true

	return root
}
