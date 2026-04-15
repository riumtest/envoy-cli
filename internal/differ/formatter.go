package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/envoy-cli/internal/masker"
)

const defaultMask = "***"

// Formatter renders a diff Result to an output writer.
type Formatter interface {
	Format(w io.Writer, result *Result) error
}

// TextFormatter writes a human-readable diff.
type TextFormatter struct {
	Masker *masker.Masker
}

// JSONFormatter writes a JSON diff.
type JSONFormatter struct {
	Masker *masker.Masker
}

// NewFormatter returns a Formatter based on the format string ("text" or "json").
func NewFormatter(format string, m *masker.Masker) Formatter {
	switch strings.ToLower(format) {
	case "json":
		return &JSONFormatter{Masker: m}
	default:
		return &TextFormatter{Masker: m}
	}
}

// Format writes a text diff to w.
func (f *TextFormatter) Format(w io.Writer, result *Result) error {
	if len(result.Changes) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "+ %s=%s\n", c.Key, maskValue(f.Masker, c.Key, c.NewValue))
		case Removed:
			fmt.Fprintf(w, "- %s=%s\n", c.Key, maskValue(f.Masker, c.Key, c.OldValue))
		case Changed:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", c.Key,
				maskValue(f.Masker, c.Key, c.OldValue),
				maskValue(f.Masker, c.Key, c.NewValue))
		}
	}
	return nil
}

// Format writes a JSON diff to w.
func (f *JSONFormatter) Format(w io.Writer, result *Result) error {
	type jsonChange struct {
		Key      string `json:"key"`
		Type     string `json:"type"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}

	var out []jsonChange
	for _, c := range result.Changes {
		if c.Type == Unchanged {
			continue
		}
		out = append(out, jsonChange{
			Key:      c.Key,
			Type:     string(c.Type),
			OldValue: maskValue(f.Masker, c.Key, c.OldValue),
			NewValue: maskValue(f.Masker, c.Key, c.NewValue),
		})
	}
	if out == nil {
		out = []jsonChange{}
	}
	return json.NewEncoder(w).Encode(out)
}

func maskValue(m *masker.Masker, key, value string) string {
	if m != nil && m.IsSensitive(key) {
		return defaultMask
	}
	return value
}
