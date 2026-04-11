// Package stripper removes comments and blank lines from parsed env entries
// or from raw .env file content, producing a clean minimal output.
package stripper

import (
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls which elements are stripped.
type Options struct {
	StripComments  bool
	StripBlanks    bool
	TrimWhitespace bool
}

// DefaultOptions returns an Options with all stripping enabled.
func DefaultOptions() Options {
	return Options{
		StripComments:  true,
		StripBlanks:    true,
		TrimWhitespace: true,
	}
}

// Strip removes comments and/or blank lines from a slice of entries.
// It returns a new slice; the original is not mutated.
func Strip(entries []envfile.Entry, opts Options) []envfile.Entry {
	result := make([]envfile.Entry, 0, len(entries))
	for _, e := range entries {
		if opts.StripComments && isComment(e) {
			continue
		}
		if opts.StripBlanks && isBlank(e) {
			continue
		}
		if opts.TrimWhitespace {
			e.Key = strings.TrimSpace(e.Key)
			e.Value = strings.TrimSpace(e.Value)
		}
		result = append(result, e)
	}
	return result
}

// StripRaw removes comment lines and blank lines from raw .env file text.
// Lines beginning with '#' (after trimming) are treated as comments.
func StripRaw(content string, opts Options) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if opts.StripComments && strings.HasPrefix(trimmed, "#") {
			continue
		}
		if opts.StripBlanks && trimmed == "" {
			continue
		}
		if opts.TrimWhitespace {
			out = append(out, trimmed)
		} else {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func isComment(e envfile.Entry) bool {
	return strings.HasPrefix(strings.TrimSpace(e.Key), "#")
}

func isBlank(e envfile.Entry) bool {
	return strings.TrimSpace(e.Key) == "" && strings.TrimSpace(e.Value) == ""
}
