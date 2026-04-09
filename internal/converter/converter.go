// Package converter provides utilities for converting .env file entries
// between different formats such as dotenv, export shell, JSON, and YAML.
package converter

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/envfile"
)

// Format represents the target output format for conversion.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
)

// ErrUnknownFormat is returned when an unsupported format is requested.
type ErrUnknownFormat struct {
	Format string
}

func (e ErrUnknownFormat) Error() string {
	return fmt.Sprintf("unknown format: %q", e.Format)
}

// Convert transforms a slice of envfile.Entry values into the specified format string.
func Convert(entries []envfile.Entry, format Format) (string, error) {
	switch format {
	case FormatDotEnv:
		return toDotEnv(entries), nil
	case FormatExport:
		return toExport(entries), nil
	case FormatJSON:
		return toJSON(entries), nil
	case FormatYAML:
		return toYAML(entries), nil
	default:
		return "", ErrUnknownFormat{Format: string(format)}
	}
}

func toDotEnv(entries []envfile.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func toExport(entries []envfile.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
	}
	return sb.String()
}

func toJSON(entries []envfile.Entry) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", e.Key, e.Value, comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func toYAML(entries []envfile.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s: %q\n", e.Key, e.Value)
	}
	return sb.String()
}
