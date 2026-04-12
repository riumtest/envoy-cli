// Package tagger provides functionality for tagging .env entries with
// arbitrary labels, enabling grouping, filtering, and annotation workflows.
package tagger

import (
	"strings"

	"github.com/envoy-cli/envoy/internal/envfile"
)

// TaggedEntry wraps an envfile.Entry with an associated set of tags.
type TaggedEntry struct {
	envfile.Entry
	Tags []string
}

// Options controls how tagging is applied.
type Options struct {
	// Rules maps a tag name to a key prefix that triggers it.
	Rules map[string]string
	// CaseSensitive controls whether prefix matching is case-sensitive.
	CaseSensitive bool
}

// DefaultOptions returns sensible defaults for the tagger.
func DefaultOptions() Options {
	return Options{
		Rules:         map[string]string{},
		CaseSensitive: false,
	}
}

// Tag applies tag rules to a slice of entries and returns TaggedEntry values.
// Each entry may receive zero or more tags based on matching rules.
func Tag(entries []envfile.Entry, opts Options) []TaggedEntry {
	result := make([]TaggedEntry, 0, len(entries))
	for _, e := range entries {
		tagged := TaggedEntry{Entry: e}
		for tag, prefix := range opts.Rules {
			key := e.Key
			p := prefix
			if !opts.CaseSensitive {
				key = strings.ToLower(key)
				p = strings.ToLower(p)
			}
			if strings.HasPrefix(key, p) {
				tagged.Tags = append(tagged.Tags, tag)
			}
		}
		result = append(result, tagged)
	}
	return result
}

// FilterByTag returns only the TaggedEntry values that carry the given tag.
func FilterByTag(tagged []TaggedEntry, tag string) []TaggedEntry {
	var out []TaggedEntry
	for _, t := range tagged {
		for _, tg := range t.Tags {
			if tg == tag {
				out = append(out, t)
				break
			}
		}
	}
	return out
}

// ToEntries strips tag metadata and returns plain entries.
func ToEntries(tagged []TaggedEntry) []envfile.Entry {
	out := make([]envfile.Entry, len(tagged))
	for i, t := range tagged {
		out[i] = t.Entry
	}
	return out
}
