package remover

import (
	"fmt"
	"strings"
)

// BashRemover implements Remover for Bash scripts.
type BashRemover struct{}

// Process removes comments from Bash scripts.
func (r *BashRemover) Process(lines []string) ([]string, []string) {
	result := make([]string, len(lines))
	copy(result, lines)

	issues := []string{}

	for i, line := range result {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Preserve shebang
		if i == 0 && strings.HasPrefix(line, "#!") {
			continue
		}

		// Handle comments
		if strings.Contains(line, "#") {
			commentIndex := strings.Index(line, "#")

			// Make sure it's not inside a string
			inSingleQuotes := false
			inDoubleQuotes := false
			isComment := true

			for j := 0; j < commentIndex; j++ {
				char := line[j]
				if char == '\'' && !inDoubleQuotes {
					inSingleQuotes = !inSingleQuotes
				} else if char == '"' && !inSingleQuotes {
					inDoubleQuotes = !inDoubleQuotes
				}
			}

			// If we're inside quotes, it's not a comment
			if inSingleQuotes || inDoubleQuotes {
				isComment = false
			}

			if isComment {
				result[i] = strings.TrimSpace(line[:commentIndex])
				issues = append(issues, fmt.Sprintf("Comment on line %d", i+1))
			}
		}
	}

	return result, issues
}
