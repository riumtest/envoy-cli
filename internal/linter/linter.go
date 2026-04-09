// Package linter provides style and convention checks for .env files.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single linting problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// Result holds all issues found during a lint pass.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r *Result) HasErrors() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return true
		}
	}
	return false
}

// Entry is a parsed key/value pair with its source line number.
type Entry struct {
	Line  int
	Key   string
	Value string
}

// Lint runs all built-in lint rules against the provided entries and returns
// a Result containing any issues found.
func Lint(entries []Entry) Result {
	var issues []Issue

	for _, e := range entries {
		// Rule: key must be UPPER_SNAKE_CASE
		if !isUpperSnakeCase(e.Key) {
			issues = append(issues, Issue{
				Line:     e.Line,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q should be UPPER_SNAKE_CASE", e.Key),
				Severity: "warn",
			})
		}

		// Rule: key must not start with a digit
		if len(e.Key) > 0 && unicode.IsDigit(rune(e.Key[0])) {
			issues = append(issues, Issue{
				Line:     e.Line,
				Key:      e.Key,
				Message:  fmt.Sprintf("key %q must not start with a digit", e.Key),
				Severity: "error",
			})
		}

		// Rule: value should not contain unquoted leading/trailing whitespace
		if e.Value != strings.TrimSpace(e.Value) {
			issues = append(issues, Issue{
				Line:     e.Line,
				Key:      e.Key,
				Message:  fmt.Sprintf("value for %q has leading or trailing whitespace", e.Key),
				Severity: "warn",
			})
		}
	}

	return Result{Issues: issues}
}

// isUpperSnakeCase returns true when s consists only of uppercase letters,
// digits, and underscores, and is non-empty.
func isUpperSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
