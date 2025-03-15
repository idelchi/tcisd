package remover

import (
	"fmt"
	"strings"
)

type GoRemover struct{}

func (r *GoRemover) Process(lines []string) ([]string, []string) {
	var result []string

	issues := []string{}

	inMultiLineComment := false
	multiLineStart := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)

			continue
		}

		if inMultiLineComment {
			if strings.Contains(line, "*/") {
				endIndex := strings.Index(line, "*/") + 2

				remaining := strings.TrimSpace(line[endIndex:])
				if remaining != "" {
					result = append(result, remaining)
				}

				inMultiLineComment = false

				issues = append(issues, fmt.Sprintf("Multi-line comment from line %d to %d:\n\t\t%q", multiLineStart+1, i+1, line[:endIndex]))
			}

			continue
		}

		if strings.HasPrefix(trimmed, "//") {
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d: %q", i+1, trimmed))

			continue // Skip the line entirely
		}

		if strings.HasPrefix(trimmed, "/*") {
			if strings.Contains(trimmed, "*/") {
				endIndex := strings.Index(trimmed, "*/") + 2

				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}

				issues = append(issues, fmt.Sprintf("Multi-line comment on line %d:\n\t\t%q", i+1, trimmed[:endIndex]))
			} else {
				inMultiLineComment = true
				multiLineStart = i
			}
		} else {
			result = append(result, line)
		}
	}

	return result, issues
}
