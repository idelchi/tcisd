package logic

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"

	"github.com/idelchi/tcisd/internal/commands"
	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/internal/processor"
	"github.com/idelchi/tcisd/pkg/cobraext"
)

// Run is the main function of the application.
func Run(version string) int {
	// Create the Config instance
	cfg := &config.Config{}
	root := commands.NewRootCommand(cfg, version)

	// Execute the root command
	switch err := root.Execute(); {
	case errors.Is(err, cobraext.ErrExitGracefully):
		return 0
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Set up processor with appropriate number of workers
	proc := processor.New(
		cfg,
		min(runtime.NumCPU(), len(cfg.Paths)),
		cfg.Types,
	)

	// Process the files
	if err := proc.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", color.RedString("Error"), err)
		return 1
	}

	// Print summary
	hasIssues := proc.Summary()

	if hasIssues && cfg.Mode == config.LintMode {
		log.Println(color.RedString("Comments found in files"))
		return 1
	}

	if !hasIssues {
		log.Println(color.GreenString("No comments found in files"))
	}

	return 0
}
