// Package blanker provides functionality to blank (zero-out) the values
// of selected environment entries while preserving their keys.
package blanker

import "github.com/user/envoy-cli/internal/envfile"

// Options controls which entries are blanked.
type Options struct {
	// Keys is an explicit list of keys to blank. If empty, all entries are blanked.
	Keys []string
	// SensitiveOnly blanks only entries whose keys match common sensitive patterns.
	SensitiveOnly bool
	// Patterns is a list of substring patterns; keys containing any pattern are blanked.
	Patterns []string
}

// DefaultOptions returns an Options value that blanks all entries.
func DefaultOptions() Options {
	return Options{}
}

// Blank returns a new slice of entries with selected values replaced by an
// empty string. The original slice is never mutated.
func Blank(entries []envfile.Entry, opts Options) []envfile.Entry {
	allowlist := buildAllowlist(opts.Keys)

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		copy := e
		if shouldBlank(e.Key, allowlist, opts) {
			copy.Value = ""
		}
		out[i] = copy
	}
	return out
}

// buildAllowlist converts a slice of keys into a fast lookup map.
func buildAllowlist(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

func shouldBlank(key string, allowlist map[string]struct{}, opts Options) bool {
	// Explicit key list takes priority.
	if len(allowlist) > 0 {
		_, ok := allowlist[key]
		return ok
	}
	// Pattern matching.
	if len(opts.Patterns) > 0 {
		for _, p := range opts.Patterns {
			if containsFold(key, p) {
				return true
			}
		}
		return false
	}
	// Sensitive-only heuristic.
	if opts.SensitiveOnly {
		return isSensitive(key)
	}
	// Default: blank everything.
	return true
}

var sensitiveSubstrings = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "CERT", "KEY",
}

func isSensitive(key string) bool {
	upper := toUpper(key)
	for _, s := range sensitiveSubstrings {
		if containsFold(upper, s) {
			return true
		}
	}
	return false
}

func containsFold(s, substr string) bool {
	return len(s) >= len(substr) && indexFold(s, substr) >= 0
}

func indexFold(s, substr string) int {
	sU := toUpper(s)
	subU := toUpper(substr)
	for i := 0; i <= len(sU)-len(subU); i++ {
		if sU[i:i+len(subU)] == subU {
			return i
		}
	}
	return -1
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
