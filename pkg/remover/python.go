package remover

import (
	"fmt"
	"strings"
)

type PythonRemover struct{}

func (r *PythonRemover) Process(lines []string) ([]string, []string) {
	var result []string

	issues := []string{}

	inMultiLineDocstring := false
	docstringStart := 0
	docstringType := ""

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)

			continue
		}

		if inMultiLineDocstring {
			if strings.Contains(line, docstringType) {
				endIndex := strings.Index(line, docstringType) + len(docstringType)

				remaining := strings.TrimSpace(line[endIndex:])
				if remaining != "" {
					result = append(result, remaining)
				}

				inMultiLineDocstring = false

				issues = append(issues, fmt.Sprintf("Docstring from line %d to %d", docstringStart+1, i+1))
			}

			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			issues = append(issues, fmt.Sprintf("Single-line comment on line %d", i+1))

			continue // Skip the line entirely
		}

		if strings.HasPrefix(trimmed, "\"\"\"") && !inMultiLineDocstring {
			if strings.Count(line, "\"\"\"") >= 2 {
				afterStartIndex := strings.Index(trimmed, "\"\"\"") + 3
				endIndex := afterStartIndex + strings.Index(trimmed[afterStartIndex:], "\"\"\"") + 3

				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}

				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "\"\"\""
			}
		} else if strings.HasPrefix(trimmed, "'''") && !inMultiLineDocstring {
			if strings.Count(line, "'''") >= 2 {
				afterStartIndex := strings.Index(trimmed, "'''") + 3
				endIndex := afterStartIndex + strings.Index(trimmed[afterStartIndex:], "'''") + 3

				afterComment := strings.TrimSpace(trimmed[endIndex:])
				if afterComment != "" {
					result = append(result, afterComment)
				}

				issues = append(issues, fmt.Sprintf("Docstring on line %d", i+1))
			} else {
				inMultiLineDocstring = true
				docstringStart = i
				docstringType = "'''"
			}
		} else {
			result = append(result, line)
		}
	}

	return result, issues
}
