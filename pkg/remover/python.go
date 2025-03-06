package remover

import (
	"fmt"
	"strings"
)

// PythonRemover implements Remover for Python source files.
type PythonRemover struct{}

// Process removes comments from Python code.
func (r *PythonRemover) Process(lines []string) ([]string, []string) {
	result := make([]string, len(lines))
	copy(result, lines)

	issues := []string{}

	// Handle multi-line docstrings
	inMultiLineDocstring := false
	docstringStart := 0
	docstringType := ""

	// Process each line
	for i, line := range result {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Handle multi-line docstrings
		if inMultiLineDocstring {
			if strings.Contains(line, docstringType) {
				// End of docstring found
				endIndex := strings.Index(line, docstringType) + len(docstringType)
				result[i] = strings.TrimSpace(line[endIndex:])
				inMultiLineDocstring = false

				issues = append(issues, fmt.Sprintf("Docstring from line %d to %d", docstringStart+1, i+1))
			} else {
				// Still in docstring
				result[i] = ""
			}
			continue
		}

		// Handle lines that start with # (after trimming whitespace)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			result[i] = ""
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))
			continue
		}

		// Check for start of triple-quoted docstrings
		if strings.Contains(line, "\"\"\"") && !inMultiLineDocstring {
			startIndex := strings.Index(line, "\"\"\"")

			// Check if there's also an end on the same line
			if strings.Count(line, "\"\"\"") >= 2 {
				// Both start and end on the same line
				endIndex := startIndex + 3 + strings.Index(line[startIndex+3:], "\"\"\"") + 3
				beforeComment := line[:startIndex]
				afterComment := line[endIndex:]
				result[i] = strings.TrimSpace(beforeComment + afterComment)

				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				// Start of docstring
				result[i] = strings.TrimSpace(line[:startIndex])
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "\"\"\""
			}
		} else if strings.Contains(line, "'''") && !inMultiLineDocstring {
			startIndex := strings.Index(line, "'''")

			// Check if there's also an end on the same line
			if strings.Count(line, "'''") >= 2 {
				// Both start and end on the same line
				endIndex := startIndex + 3 + strings.Index(line[startIndex+3:], "'''") + 3
				beforeComment := line[:startIndex]
				afterComment := line[endIndex:]
				result[i] = strings.TrimSpace(beforeComment + afterComment)

				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				// Start of docstring
				result[i] = strings.TrimSpace(line[:startIndex])
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "'''"
			}
		}
	}

	return result, issues
}
