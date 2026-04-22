// Package linker resolves cross-file references between .env files,
// allowing one file to declare a key whose value is sourced from another.
package linker

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
)

// Rule describes a single link: the target key in the destination file
// is resolved from the source key in the source file.
type Rule struct {
	FromFile string
	FromKey  string
	ToKey    string
}

// Result holds the outcome of applying a single Rule.
type Result struct {
	Rule    Rule
	Value   string
	OK      bool
	Message string
}

// Link applies the given rules against the provided source entries map,
// returning the resolved entries and a result report per rule.
//
// sourceFiles maps a file path to its parsed entries. Entries in dst
// that already have the ToKey set will be overwritten.
func Link(dst []envfile.Entry, sourceFiles map[string][]envfile.Entry, rules []Rule) ([]envfile.Entry, []Result, error) {
	if len(rules) == 0 {
		return dst, nil, nil
	}

	// Build lookup maps per file.
	lookup := make(map[string]map[string]string, len(sourceFiles))
	for path, entries := range sourceFiles {
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		lookup[path] = m
	}

	// Build a mutable index of dst entries.
	dstIndex := make(map[string]int, len(dst))
	out := make([]envfile.Entry, len(dst))
	copy(out, dst)
	for i, e := range out {
		dstIndex[e.Key] = i
	}

	results := make([]Result, 0, len(rules))

	for _, rule := range rules {
		fileMap, ok := lookup[rule.FromFile]
		if !ok {
			results = append(results, Result{
				Rule:    rule,
				Message: fmt.Sprintf("source file %q not found", rule.FromFile),
			})
			continue
		}
		val, exists := fileMap[strings.TrimSpace(rule.FromKey)]
		if !exists {
			results = append(results, Result{
				Rule:    rule,
				Message: fmt.Sprintf("key %q not found in %q", rule.FromKey, rule.FromFile),
			})
			continue
		}

		if idx, found := dstIndex[rule.ToKey]; found {
			out[idx].Value = val
		} else {
			out = append(out, envfile.Entry{Key: rule.ToKey, Value: val})
			dstIndex[rule.ToKey] = len(out) - 1
		}

		results = append(results, Result{
			Rule:    rule,
			Value:   val,
			OK:      true,
			Message: "ok",
		})
	}

	return out, results, nil
}
