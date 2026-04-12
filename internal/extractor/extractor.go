// Package extractor provides utilities for extracting subsets of env entries
// based on key patterns, value patterns, or custom predicates.
package extractor

import (
	"regexp"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls which entries are extracted.
type Options struct {
	// KeyPattern is an optional regex applied to keys.
	KeyPattern string
	// ValuePattern is an optional regex applied to values.
	ValuePattern string
	// Keys is an explicit allowlist of keys to extract.
	Keys []string
	// CaseSensitive controls regex matching (default false).
	CaseSensitive bool
}

// Result holds the extracted entries and metadata.
type Result struct {
	Entries   []envfile.Entry
	Extracted int
	Skipped   int
}

// Extract returns entries that match the given options.
// If no options are set all entries are returned unchanged.
func Extract(entries []envfile.Entry, opts Options) (Result, error) {
	var keyRe, valRe *regexp.Regexp

	if opts.KeyPattern != "" {
		pattern := opts.KeyPattern
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		var err error
		keyRe, err = regexp.Compile(pattern)
		if err != nil {
			return Result{}, err
		}
	}

	if opts.ValuePattern != "" {
		pattern := opts.ValuePattern
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		var err error
		valRe, err = regexp.Compile(pattern)
		if err != nil {
			return Result{}, err
		}
	}

	allowlist := buildAllowlist(opts.Keys)

	var result Result
	for _, e := range entries {
		if !matches(e, keyRe, valRe, allowlist) {
			result.Skipped++
			continue
		}
		result.Entries = append(result.Entries, e)
		result.Extracted++
	}
	return result, nil
}

func matches(e envfile.Entry, keyRe, valRe *regexp.Regexp, allowlist map[string]struct{}) bool {
	if len(allowlist) > 0 {
		if _, ok := allowlist[strings.ToUpper(e.Key)]; !ok {
			return false
		}
	}
	if keyRe != nil && !keyRe.MatchString(e.Key) {
		return false
	}
	if valRe != nil && !valRe.MatchString(e.Value) {
		return false
	}
	return true
}

func buildAllowlist(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToUpper(k)] = struct{}{}
	}
	return m
}
