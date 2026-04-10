// Package profiler analyzes .env files and produces a usage profile,
// reporting value types, key patterns, and overall file statistics.
package profiler

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Profile holds the analysis results for a set of env entries.
type Profile struct {
	TotalKeys    int
	EmptyValues  int
	BooleanKeys  []string
	NumericKeys  []string
	URLKeys      []string
	SecretKeys   []string
	PrefixGroups map[string]int // prefix (before first _) -> count
}

var secretPatterns = []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL"}

// Analyze builds a Profile from the provided entries.
func Analyze(entries []envfile.Entry) Profile {
	p := Profile{
		PrefixGroups: make(map[string]int),
	}

	for _, e := range entries {
		p.TotalKeys++

		if e.Value == "" {
			p.EmptyValues++
		}

		if isBoolean(e.Value) {
			p.BooleanKeys = append(p.BooleanKeys, e.Key)
		} else if isNumeric(e.Value) {
			p.NumericKeys = append(p.NumericKeys, e.Key)
		} else if isURL(e.Value) {
			p.URLKeys = append(p.URLKeys, e.Key)
		}

		if isSensitive(e.Key) {
			p.SecretKeys = append(p.SecretKeys, e.Key)
		}

		if idx := strings.Index(e.Key, "_"); idx > 0 {
			prefix := e.Key[:idx]
			p.PrefixGroups[prefix]++
		}
	}

	return p
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

func isURL(v string) bool {
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range secretPatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
