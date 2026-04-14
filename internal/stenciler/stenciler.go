// Package stenciler generates a redacted .env template from a set of
// entries, replacing all values with empty strings or typed placeholders
// so the file can be safely committed as a reference template.
package stenciler

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Options controls how the stencil is generated.
type Options struct {
	// TypedPlaceholders replaces values with a hint of the inferred type
	// e.g. "<string>", "<number>", "<bool>", "<url>".
	TypedPlaceholders bool
	// PreserveComments keeps comment lines from the original entries.
	PreserveComments bool
}

// DefaultOptions returns sensible defaults for stencil generation.
func DefaultOptions() Options {
	return Options{
		TypedPlaceholders: true,
		PreserveComments:  true,
	}
}

// Stencil converts a slice of entries into a template where all values
// are replaced with empty strings or typed placeholders.
func Stencil(entries []envfile.Entry, opts Options) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(entries))
	for _, e := range entries {
		if e.Key == "" {
			if opts.PreserveComments {
				out = append(out, e)
			}
			continue
		}
		placeholder := ""
		if opts.TypedPlaceholders {
			placeholder = inferPlaceholder(e.Value)
		}
		out = append(out, envfile.Entry{
			Key:   e.Key,
			Value: placeholder,
		})
	}
	return out
}

// Render serialises stencilled entries to a .env-style string.
func Render(entries []envfile.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Key == "" {
			// comment or blank line stored in Value by the parser
			fmt.Fprintf(&sb, "%s\n", e.Value)
			continue
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

// inferPlaceholder returns a typed placeholder string based on the value.
func inferPlaceholder(v string) string {
	if v == "" {
		return "<string>"
	}
	lower := strings.ToLower(v)
	if lower == "true" || lower == "false" {
		return "<bool>"
	}
	if isNumeric(v) {
		return "<number>"
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return "<url>"
	}
	return "<string>"
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
