// Package highlighter provides syntax highlighting for .env file output,
// annotating keys, values, comments, and sensitive entries with ANSI colour codes.
package highlighter

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/masker"
)

// ANSI colour codes.
const (
	colorReset  = "\033[0m"
	colorKey     = "\033[36m"  // cyan
	colorValue   = "\033[32m"  // green
	colorComment = "\033[90m"  // dark grey
	colorSensitive = "\033[33m" // yellow
	colorEmpty   = "\033[31m"  // red
)

// Options controls highlighting behaviour.
type Options struct {
	// MaskSensitive replaces sensitive values with asterisks before colouring.
	MaskSensitive bool
	// NoColor disables ANSI codes entirely (plain output).
	NoColor bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaskSensitive: false,
		NoColor:       false,
	}
}

// Highlight returns a coloured string representation of the given entries.
// Each entry is rendered as KEY=VALUE, with comments shown above where present.
func Highlight(entries []envfile.Entry, opts Options, m *masker.Masker) []string {
	lines := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.Key == "" {
			continue
		}

		value := e.Value
		sensitive := m != nil && m.IsSensitive(e.Key)

		if sensitive && opts.MaskSensitive && m != nil {
			value = m.Mask(e.Key, value)
		}

		if opts.NoColor {
			lines = append(lines, fmt.Sprintf("%s=%s", e.Key, value))
			continue
		}

		keyColor := colorKey
		if sensitive {
			keyColor = colorSensitive
		}

		valColor := colorValue
		if strings.TrimSpace(value) == "" {
			valColor = colorEmpty
		}

		lines = append(lines, fmt.Sprintf("%s%s%s=%s%s%s",
			keyColor, e.Key, colorReset,
			valColor, value, colorReset,
		))
	}

	return lines
}

// HighlightComment returns a coloured comment line.
func HighlightComment(line string, noColor bool) string {
	if noColor {
		return line
	}
	return colorComment + line + colorReset
}
