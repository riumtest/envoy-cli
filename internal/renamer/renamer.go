package renamer

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/envfile"
)

// Rule describes a single rename operation.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of a rename operation.
type Result struct {
	Renamed  []Rule
	Skipped  []Rule // rules where From key was not found
	Conflict []Rule // rules where To key already exists
}

// Rename applies the given rules to entries, returning updated entries and a Result.
// If a target key already exists the rule is recorded as a conflict and skipped.
func Rename(entries []envfile.Entry, rules []Rule) ([]envfile.Entry, Result, error) {
	if len(rules) == 0 {
		return entries, Result{}, nil
	}

	// build index for fast lookup
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	out := make([]envfile.Entry, len(entries))
	copy(out, entries)

	var res Result

	for _, rule := range rules {
		if strings.TrimSpace(rule.From) == "" || strings.TrimSpace(rule.To) == "" {
			return nil, Result{}, fmt.Errorf("rename rule has empty From or To: %+v", rule)
		}

		fromIdx, fromExists := index[rule.From]
		if !fromExists {
			res.Skipped = append(res.Skipped, rule)
			continue
		}

		if _, toExists := index[rule.To]; toExists {
			res.Conflict = append(res.Conflict, rule)
			continue
		}

		out[fromIdx].Key = rule.To
		delete(index, rule.From)
		index[rule.To] = fromIdx
		res.Renamed = append(res.Renamed, rule)
	}

	return out, res, nil
}
