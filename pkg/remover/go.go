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

		// Check for start of multi-line comment
		if strings.Contains(line, "/*") {
			startIndex := strings.Index(line, "/*")

			// Check if there's also an end of multi-line comment on the same line
			if strings.Contains(line[startIndex:], "*/") {
				// Both start and end on the same line
				endIndex := startIndex + strings.Index(line[startIndex:], "*/") + 2
				beforeComment := line[:startIndex]
				afterComment := line[endIndex:]
				result[i] = strings.TrimSpace(beforeComment + afterComment)

				issues = append(issues, fmt.Sprintf("Multi-line comment on line %d", i+1))
			} else {
				// Start of multi-line comment
				result[i] = strings.TrimSpace(line[:startIndex])
				inMultiLineComment = true
				multiLineStart = i
			}
		} else {
			// Handle single-line comments - ONLY those that START with // (after trimming whitespace)
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				result[i] = ""
				issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))
			}
		}
	}

	return result, issues
}
