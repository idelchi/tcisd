package remover

import (
	"fmt"
	"strings"
)

// GoRemover implements Remover for Go source files.
type GoRemover struct{}

// Process removes comments from Go code.
func (r *GoRemover) Process(lines []string) ([]string, []string) {
	result := make([]string, len(lines))
	copy(result, lines)

	issues := []string{}

	// Handle multi-line comments
	inMultiLineComment := false
	multiLineStart := 0

	// Process each line
	for i, line := range result {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Handle multi-line comments
		if inMultiLineComment {
			if strings.Contains(line, "*/") {
				// End of multi-line comment found
				endIndex := strings.Index(line, "*/") + 2
				result[i] = strings.TrimSpace(line[endIndex:])
				inMultiLineComment = false

				issues = append(issues, fmt.Sprintf("Multi-line comment from line %d to %d", multiLineStart+1, i+1))
			} else {
				// Still in multi-line comment
				result[i] = ""
			}
			continue
		}

		// Trim whitespace for checking prefixes
		trimmed := strings.TrimSpace(line)

		// Handle single-line comments - ONLY those that START with // (after trimming whitespace)
		if strings.HasPrefix(trimmed, "//") {
			result[i] = ""
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))
			continue
		}

		// Check for start of multi-line comment - ONLY at the beginning of the line
		if strings.HasPrefix(trimmed, "/*") {
			// Check if there's also an end of multi-line comment on the same line
			if strings.Contains(trimmed, "*/") {
				endIndex := strings.Index(trimmed, "*/") + 2
				afterComment := trimmed[endIndex:]
				result[i] = strings.TrimSpace(afterComment)

				issues = append(issues, fmt.Sprintf("Multi-line comment on line %d", i+1))
			} else {
				// Start of multi-line comment
				result[i] = ""
				inMultiLineComment = true
				multiLineStart = i
			}
		}
	}

	return result, issues
}
