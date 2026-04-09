// Package validator provides utilities for validating .env file contents,
// including checks for duplicate keys, empty values, and invalid formats.
package validator

import (
	"fmt"
	"strings"
)

// Issue represents a single validation problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
	Severity Severity
}

// Severity indicates how serious a validation issue is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

func (i Issue) String() string {
	return fmt.Sprintf("[%s] line %d: %s — %s", i.Severity, i.Line, i.Key, i.Message)
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has error severity.
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Validate checks the provided key-value pairs (with line metadata) for common issues.
// entries is a slice of [lineNumber, key, value] tuples.
func Validate(entries []Entry) *Result {
	result := &Result{}
	seen := make(map[string]int)

	for _, e := range entries {
		key := e.Key
		lineno := e.Line

		if strings.TrimSpace(key) == "" {
			result.Issues = append(result.Issues, Issue{
				Line:     lineno,
				Key:      "(empty)",
				Message:  "key must not be empty",
				Severity: SeverityError,
			})
			continue
		}

		if prev, exists := seen[key]; exists {
			result.Issues = append(result.Issues, Issue{
				Line:     lineno,
				Key:      key,
				Message:  fmt.Sprintf("duplicate key (first defined on line %d)", prev),
				Severity: SeverityWarning,
			})
		} else {
			seen[key] = lineno
		}

		if strings.TrimSpace(e.Value) == "" {
			result.Issues = append(result.Issues, Issue{
				Line:     lineno,
				Key:      key,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}

		if strings.ContainsAny(key, " \t") {
			result.Issues = append(result.Issues, Issue{
				Line:     lineno,
				Key:      key,
				Message:  "key contains whitespace",
				Severity: SeverityError,
			})
		}
	}

	return result
}

// Entry represents a parsed line with positional metadata.
type Entry struct {
	Line  int
	Key   string
	Value string
}
