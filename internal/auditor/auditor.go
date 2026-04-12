// Package auditor provides functionality for auditing .env entries
// and producing a structured report of potential issues and observations.
package auditor

import (
	"fmt"
	"strings"

	"github.com/yourusername/envoy-cli/internal/masker"
)

// Severity represents the importance level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding represents a single audit observation for an entry.
type Finding struct {
	Key      string   `json:"key"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
}

// Report holds the full result of an audit run.
type Report struct {
	Findings []Finding `json:"findings"`
	Total    int       `json:"total"`
	Errors   int       `json:"errors"`
	Warnings int       `json:"warnings"`
	Infos    int       `json:"infos"`
}

// Entry represents a key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// Audit inspects a slice of entries and returns a structured Report.
func Audit(entries []Entry, m *masker.Masker) Report {
	var findings []Finding

	for _, e := range entries {
		if e.Key == "" {
			findings = append(findings, Finding{
				Key:      "(empty)",
				Severity: SeverityError,
				Message:  "key must not be empty",
			})
			continue
		}

		if e.Value == "" {
			findings = append(findings, Finding{
				Key:      e.Key,
				Severity: SeverityWarning,
				Message:  "value is empty",
			})
		}

		if m != nil && m.IsSensitive(e.Key) && len(e.Value) < 8 && e.Value != "" {
			findings = append(findings, Finding{
				Key:      e.Key,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("sensitive key %q has a suspiciously short value", e.Key),
			})
		}

		if strings.Contains(e.Value, "TODO") || strings.Contains(e.Value, "FIXME") {
			findings = append(findings, Finding{
				Key:      e.Key,
				Severity: SeverityInfo,
				Message:  "value contains a TODO/FIXME marker",
			})
		}

		if strings.HasPrefix(e.Key, "_") {
			findings = append(findings, Finding{
				Key:      e.Key,
				Severity: SeverityInfo,
				Message:  "key starts with underscore, may be a private/internal variable",
			})
		}
	}

	report := Report{Findings: findings, Total: len(findings)}
	for _, f := range findings {
		switch f.Severity {
		case SeverityError:
			report.Errors++
		case SeverityWarning:
			report.Warnings++
		case SeverityInfo:
			report.Infos++
		}
	}
	return report
}
