// Package padder provides utilities for padding or aligning env entry values
// to a consistent width, useful for human-readable output formatting.
package padder

import (
	"strings"

	"envoy-cli/internal/envfile"
)

// DefaultOptions returns sensible defaults for padding.
func DefaultOptions() Options {
	return Options{
		Align:     AlignLeft,
		PadChar:   ' ',
		MinWidth:  0,
		MaxWidth:  0,
	}
}

// Alignment direction.
type Alignment int

const (
	AlignLeft  Alignment = iota
	AlignRight
)

// Options controls padding behaviour.
type Options struct {
	Align    Alignment
	PadChar  rune
	MinWidth int
	MaxWidth int // 0 means no limit
}

// Pad returns a copy of entries with values padded to a consistent width.
// The target width is determined by the longest value (clamped to MaxWidth if set),
// but never less than MinWidth.
func Pad(entries []envfile.Entry, opts Options) []envfile.Entry {
	if len(entries) == 0 {
		return []envfile.Entry{}
	}

	width := opts.MinWidth
	for _, e := range entries {
		if len(e.Value) > width {
			width = len(e.Value)
		}
	}
	if opts.MaxWidth > 0 && width > opts.MaxWidth {
		width = opts.MaxWidth
	}

	out := make([]envfile.Entry, len(entries))
	for i, e := range entries {
		v := e.Value
		if len(v) < width {
			pad := strings.Repeat(string(opts.PadChar), width-len(v))
			if opts.Align == AlignRight {
				v = pad + v
			} else {
				v = v + pad
			}
		} else if opts.MaxWidth > 0 && len(v) > opts.MaxWidth {
			v = v[:opts.MaxWidth]
		}
		out[i] = envfile.Entry{Key: e.Key, Value: v}
	}
	return out
}
