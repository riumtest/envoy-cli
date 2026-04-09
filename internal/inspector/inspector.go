// Package inspector provides functionality to inspect and summarize
// the contents of a parsed .env file, including key statistics,
// sensitive key detection, and value type heuristics.
package inspector

import (
	"strings"

	"github.com/user/envoy-cli/internal/masker"
)

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// Report holds the inspection results for a set of env entries.
type Report struct {
	Total     int
	Empty     int
	Sensitive int
	URLs      int
	Booleans  int
	Numeric   int
	Keys      []string
}

// Inspect analyses a slice of entries and returns a Report.
func Inspect(entries []Entry, m *masker.Masker) Report {
	r := Report{}
	for _, e := range entries {
		r.Total++
		r.Keys = append(r.Keys, e.Key)

		if e.Value == "" {
			r.Empty++
			continue
		}

		if m != nil && m.IsSensitive(e.Key) {
			r.Sensitive++
		}

		if isURL(e.Value) {
			r.URLs++
		} else if isBoolean(e.Value) {
			r.Booleans++
		} else if isNumeric(e.Value) {
			r.Numeric++
		}
	}
	return r
}

func isURL(v string) bool {
	return strings.HasPrefix(v, "http://") ||
		strings.HasPrefix(v, "https://") ||
		strings.HasPrefix(v, "postgres://") ||
		strings.HasPrefix(v, "mysql://") ||
		strings.HasPrefix(v, "redis://")
}

func isBoolean(v string) bool {
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "1", "0":
		return true
	}
	return false
}

func isNumeric(v string) bool {
	if v == "" {
		return false
	}
	for _, c := range v {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
