// Package scorer evaluates the quality of env entries and returns a score.
package scorer

import (
	"strings"
	"unicode"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/masker"
)

const (
	MaxScore = 100
)

type Result struct {
	Key   string
	Score int
	Notes []string
}

// Score evaluates each entry and returns a per-key result.
func Score(entries []envfile.Entry) []Result {
	m := masker.New()
	out := make([]Result, 0, len(entries))
	for _, e := range entries {
		out = append(out, scoreEntry(e, m))
	}
	return out
}

func scoreEntry(e envfile.Entry, m *masker.Masker) Result {
	s := MaxScore
	var notes []string

	if e.Key == "" {
		return Result{Key: e.Key, Score: 0, Notes: []string{"empty key"}}
	}

	if e.Key != strings.ToUpper(e.Key) {
		s -= 20
		notes = append(notes, "key is not uppercase")
	}

	if e.Value == "" {
		s -= 30
		notes = append(notes, "empty value")
	}

	if m.IsSensitive(e.Key) && len(e.Value) < 8 {
		s -= 25
		notes = append(notes, "sensitive key has short value")
	}

	if strings.ContainsAny(e.Key, " \t") {
		s -= 15
		notes = append(notes, "key contains whitespace")
	}

	if len(e.Key) > 0 && unicode.IsDigit(rune(e.Key[0])) {
		s -= 10
		notes = append(notes, "key starts with digit")
	}

	if s < 0 {
		s = 0
	}
	return Result{Key: e.Key, Score: s, Notes: notes}
}
