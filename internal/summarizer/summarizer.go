// Package summarizer provides functionality for generating a high-level
// summary report of a parsed .env file, including key counts, value types,
// and sensitive key detection.
package summarizer

import (
	"regexp"
	"strconv"

	"github.com/user/envoy-cli/internal/envfile"
)

// Report holds the summary statistics for a set of env entries.
type Report struct {
	TotalKeys     int      `json:"total_keys"`
	EmptyValues   int      `json:"empty_values"`
	SensitiveKeys []string `json:"sensitive_keys"`
	NumericValues int      `json:"numeric_values"`
	BooleanValues int      `json:"boolean_values"`
	URLValues     int      `json:"url_values"`
	UniqueKeys    int      `json:"unique_keys"`
	DuplicateKeys int      `json:"duplicate_keys"`
}

var sensitivePattern = regexp.MustCompile(`(?i)(secret|password|passwd|token|api_?key|private|credential|auth)`)
var urlPattern = regexp.MustCompile(`(?i)^https?://`)

// Summarize analyzes a slice of env entries and returns a Report.
func Summarize(entries []envfile.Entry) Report {
	seen := make(map[string]struct{})
	report := Report{}

	for _, e := range entries {
		report.TotalKeys++

		if _, dup := seen[e.Key]; !dup {
			seen[e.Key] = struct{}{}
			report.UniqueKeys++
		} else {
			report.DuplicateKeys++
		}

		if e.Value == "" {
			report.EmptyValues++
		}

		if sensitivePattern.MatchString(e.Key) {
			report.SensitiveKeys = append(report.SensitiveKeys, e.Key)
		}

		if _, err := strconv.ParseFloat(e.Value, 64); err == nil {
			report.NumericValues++
		} else if _, err := strconv.ParseBool(e.Value); err == nil {
			report.BooleanValues++
		} else if urlPattern.MatchString(e.Value) {
			report.URLValues++
		}
	}

	if report.SensitiveKeys == nil {
		report.SensitiveKeys = []string{}
	}

	return report
}
