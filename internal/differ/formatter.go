package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envoy-cli/envoy-cli/internal/masker"
)

// Formatter writes a diff Result to a writer.
type Formatter interface {
	Format(w io.Writer, result Result) error
}

// NewFormatter returns a Formatter for the given format string ("text" or "json").
func NewFormatter(format string, m *masker.Masker) Formatter {
	switch strings.ToLower(format) {
	case "json":
		return &jsonFormatter{masker: m}
	default:
		return &textFormatter{masker: m}
	}
}

type textFormatter struct {
	masker *masker.Masker
}

func (f *textFormatter) Format(w io.Writer, result Result) error {
	changes := result.Changes
	sort.Slice(changes, func(i, j int) bool { return changes[i].Key < changes[j].Key })
	for _, c := range changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(w, "+ %s=%s\n", c.Key, maskValue(f.masker, c.Key, c.NewValue))
		case Removed:
			fmt.Fprintf(w, "- %s=%s\n", c.Key, maskValue(f.masker, c.Key, c.OldValue))
		case Changed:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", c.Key,
				maskValue(f.masker, c.Key, c.OldValue),
				maskValue(f.masker, c.Key, c.NewValue))
		}
	}
	return nil
}

type jsonFormatter struct {
	masker *masker.Masker
}

func (f *jsonFormatter) Format(w io.Writer, result Result) error {
	type jsonChange struct {
		Key      string `json:"key"`
		Kind     string `json:"kind"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}
	out := make([]jsonChange, 0, len(result.Changes))
	for _, c := range result.Changes {
		if c.Kind == Unchanged {
			continue
		}
		out = append(out, jsonChange{
			Key:      c.Key,
			Kind:     string(c.Kind),
			OldValue: maskValue(f.masker, c.Key, c.OldValue),
			NewValue: maskValue(f.masker, c.Key, c.NewValue),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func maskValue(m *masker.Masker, key, value string) string {
	if m != nil && m.IsSensitive(key) {
		return m.Mask(value)
	}
	return value
}
