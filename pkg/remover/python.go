package remover

import (
	"fmt"
	"strings"
)

// PythonRemover implements Remover for Python source files.
type PythonRemover struct{}

// Process removes comments from Python code.
func (r *PythonRemover) Process(lines []string) ([]string, []string) {
	var result []string
	issues := []string{}

	// Handle multi-line docstrings
	inMultiLineDocstring := false
	docstringStart := 0
	docstringType := ""

	// Process each line
	for i, line := range lines {
		// Skip empty lines but add them to result
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		// Handle multi-line docstrings
		if inMultiLineDocstring {
			if strings.Contains(line, docstringType) {
				// End of docstring found
				endIndex := strings.Index(line, docstringType) + len(docstringType)
				remaining := strings.TrimSpace(line[endIndex:])
				if remaining != "" {
					result = append(result, remaining)
				}
				inMultiLineDocstring = false
				issues = append(issues, fmt.Sprintf("Docstring from line %d to %d", docstringStart+1, i+1))
			}
			// Skip this line if we're inside a docstring
			continue
		}

		// Handle lines that start with # (after trimming whitespace)
		if strings.HasPrefix(trimmed, "#") {
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))
			continue // Skip the line entirely
		}

		// Check for triple-quoted docstrings at the beginning of the line
		if strings.HasPrefix(trimmed, "\"\"\"") && !inMultiLineDocstring {
			// Check if there's also an end on the same line
			if strings.Count(line, "\"\"\"") >= 2 {
				// Both start and end on the same line
				afterStartIndex := strings.Index(trimmed, "\"\"\"") + 3
				endIndex := afterStartIndex + strings.Index(trimmed[afterStartIndex:], "\"\"\"") + 3
				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}
				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				// Start of docstring
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "\"\"\""
			}
		} else if strings.HasPrefix(trimmed, "'''") && !inMultiLineDocstring {
			// Check if there's also an end on the same line
			if strings.Count(line, "'''") >= 2 {
				// Both start and end on the same line
				afterStartIndex := strings.Index(trimmed, "'''") + 3
				endIndex := afterStartIndex + strings.Index(trimmed[afterStartIndex:], "'''") + 3
				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}
				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				// Start of docstring
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "'''"
			}
		} else {
			// No comment, keep the line
			result = append(result, line)
		}
	}

	return result, issues
}
