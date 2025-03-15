package remover

import (
	"fmt"
	"strings"
)

type DockerfileRemover struct{}

func (r *DockerfileRemover) Process(lines []string) ([]string, []string) {
	var result []string
	issues := []string{}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			issues = append(issues, fmt.Sprintf("Comment on line %d: %q", i+1, trimmed))
			continue // Skip the line entirely
		}

		result = append(result, line)
	}

	return result, issues
}
