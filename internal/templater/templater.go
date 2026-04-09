// Package templater provides functionality for rendering .env files
// from a template with placeholder substitution.
package templater

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// placeholderRe matches {{VAR_NAME}} style placeholders.
var placeholderRe = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)

// Result holds the output of a template render operation.
type Result struct {
	// Entries are the rendered key-value pairs.
	Entries []envfile.Entry
	// Missing contains placeholder names that had no substitution value.
	Missing []string
}

// Render takes a slice of template entries (values may contain {{PLACEHOLDER}}
// tokens) and a substitution map, and returns a Result with all placeholders
// replaced. Keys that reference missing substitutions are collected in
// Result.Missing rather than causing an error, allowing callers to decide
// how to handle incomplete renders.
func Render(template []envfile.Entry, subs map[string]string) Result {
	missingSet := map[string]struct{}{}
	result := make([]envfile.Entry, 0, len(template))

	for _, entry := range template {
		rendered, missing := substitute(entry.Value, subs)
		for _, m := range missing {
			missingSet[m] = struct{}{}
		}
		result = append(result, envfile.Entry{
			Key:     entry.Key,
			Value:   rendered,
			Comment: entry.Comment,
		})
	}

	missing := make([]string, 0, len(missingSet))
	for k := range missingSet {
		missing = append(missing, k)
	}

	return Result{Entries: result, Missing: missing}
}

// substitute replaces all {{PLACEHOLDER}} tokens in value using subs.
// It returns the substituted string and a slice of any placeholder names
// that were not found in subs.
func substitute(value string, subs map[string]string) (string, []string) {
	var missing []string
	result := placeholderRe.ReplaceAllStringFunc(value, func(match string) string {
		name := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		if v, ok := subs[name]; ok {
			return v
		}
		missing = append(missing, name)
		return match
	})
	return result, missing
}

// BuildSubsFromEntries converts a slice of entries into a substitution map
// keyed by entry Key, making it easy to use one .env file as the substitution
// source for another.
func BuildSubsFromEntries(entries []envfile.Entry) map[string]string {
	subs := make(map[string]string, len(entries))
	for _, e := range entries {
		subs[e.Key] = e.Value
	}
	return subs
}

// Summary returns a human-readable summary of the render result.
func Summary(r Result) string {
	if len(r.Missing) == 0 {
		return fmt.Sprintf("rendered %d entries, no missing placeholders", len(r.Entries))
	}
	return fmt.Sprintf("rendered %d entries, %d missing placeholder(s): %s",
		len(r.Entries), len(r.Missing), strings.Join(r.Missing, ", "))
}
