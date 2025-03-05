package cobraext

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PipeOrArg reads from either the first argument or stdin.
// If an argument is provided, it is returned.
// If no argument is provided but stdin is piped, the stdin content is returned.
// If neither an argument nor stdin is provided, an empty string is returned.
func PipeOrArg(args []string) (string, error) {
	if len(args) > 0 {
		// Prioritize argument if it exists, regardless of stdin
		return args[0], nil
	}

	if IsPiped() {
		// No arg but stdin is piped
		arg, err := Read()
		if err != nil {
			return "", fmt.Errorf("reading from stdin: %w", err)
		}

		return arg, nil
	}

	return "", nil
}

// IsPiped checks if something has been piped to stdin.
func IsPiped() bool {
	fi, err := os.Stdin.Stat()
	return (fi.Mode()&os.ModeCharDevice) == 0 && err == nil
}

// Read returns stdin as a string, trimming the trailing newline.
func Read() (string, error) {
	bytes, err := io.ReadAll(os.Stdin)
	return strings.TrimSuffix(string(bytes), "\n"), err
}
