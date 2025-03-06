package parse

import (
	"errors"
	"fmt"

	"github.com/idelchi/gogen/pkg/cobraext"
	"github.com/idelchi/tcisd/internal/commands"
	"github.com/idelchi/tcisd/internal/config"
)

func Execute(version string) error {
	cfg := &config.Config{}
	root := commands.NewRootCommand(cfg, version)

	switch err := root.Execute(); {
	case errors.Is(err, cobraext.ErrExitGracefully):
		return nil
	case err != nil:
		return fmt.Errorf("executing command: %w", err)
	default:
		return nil
	}
}
