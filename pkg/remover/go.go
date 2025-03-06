package remover

import (
	"fmt"
	"strings"
)

// GoRemover implements Remover for Go source files.
type GoRemover struct{}

// Process removes comments from Go code.
func (r *GoRemover) Process(lines []string) ([]string, []string) {
	var result []string
	issues := []string{}

	// Handle multi-line comments
	inMultiLineComment := false
	multiLineStart := 0

	// Process each line
	for i, line := range lines {
		// Skip empty lines but add them to result
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		// Handle multi-line comments
		if inMultiLineComment {
			if strings.Contains(line, "*/") {
				// End of multi-line comment found
				endIndex := strings.Index(line, "*/") + 2
				remaining := strings.TrimSpace(line[endIndex:])
				if remaining != "" {
					result = append(result, remaining)
				}
				inMultiLineComment = false
				issues = append(issues, fmt.Sprintf("Multi-line comment from line %d to %d", multiLineStart+1, i+1))
			}
			// Skip this line if we're inside a comment
			continue
		}

		// Handle single-line comments - ONLY those that START with // (after trimming whitespace)
		if strings.HasPrefix(trimmed, "//") {
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))
			continue // Skip the line entirely
		}

		// Check for start of multi-line comment - ONLY at the beginning of the line
		if strings.HasPrefix(trimmed, "/*") {
			// Check if there's also an end of multi-line comment on the same line
			if strings.Contains(trimmed, "*/") {
				endIndex := strings.Index(trimmed, "*/") + 2
				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}
				issues = append(issues, fmt.Sprintf("Multi-line comment on line %d", i+1))
			} else {
				// Start of multi-line comment
				inMultiLineComment = true
				multiLineStart = i
			}
		} else {
			// No comment, keep the line
			result = append(result, line)
		}
	}

	return result, issues
}
