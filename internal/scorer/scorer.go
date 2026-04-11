// Package scorer provides a quality scoring mechanism for .env files,
// evaluating entries based on key naming, value presence, and sensitivity exposure.
package scorer

import (
	"strings"

	"github.com/envoy-cli/internal/masker"
)

// Result holds the scoring outcome for a set of env entries.
type Result struct {
	Score       int            // 0–100
	MaxScore    int
	Penalties   []string
	Suggestions []string
}

// Entry represents a single key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Score evaluates a slice of entries and returns a Result.
func Score(entries []Entry, m *masker.Masker) Result {
	if len(entries) == 0 {
		return Result{Score: 100, MaxScore: 100}
	}

	total := len(entries) * 3 // 3 checks per entry
	penalties := []string{}
	suggestions := []string{}
	deductions := 0

	for _, e := range entries {
		// Check 1: key is uppercase snake_case
		if e.Key != strings.ToUpper(e.Key) {
			deductions++
			penalties = append(penalties, "key '"+e.Key+"' is not uppercase")
			suggestions = append(suggestions, "rename '"+e.Key+"' to '"+strings.ToUpper(e.Key)+"'")
		}

		// Check 2: value is not empty
		if strings.TrimSpace(e.Value) == "" {
			deductions++
			penalties = append(penalties, "key '"+e.Key+"' has an empty value")
			suggestions = append(suggestions, "provide a value or remove '"+e.Key+"'")
		}

		// Check 3: sensitive keys should not have plaintext-looking short values
		if m != nil && m.IsSensitive(e.Key) {
			if len(e.Value) > 0 && len(e.Value) < 8 {
				deductions++
				penalties = append(penalties, "sensitive key '"+e.Key+"' has a suspiciously short value")
				suggestions = append(suggestions, "ensure '"+e.Key+"' contains a strong secret")
			}
		}
	}

	score := 100
	if total > 0 {
		score = 100 - (deductions*100)/total
	}
	if score < 0 {
		score = 0
	}

	return Result{
		Score:       score,
		MaxScore:    100,
		Penalties:   penalties,
		Suggestions: suggestions,
	}
}
