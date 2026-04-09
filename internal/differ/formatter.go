package differ

import (
	"fmt"
	"strings"
)

// Formatter defines the interface for formatting diff results
type Formatter interface {
	Format(*Result) string
}

// TextFormatter formats differences as plain text
type TextFormatter struct {
	MaskSecrets bool
}

// Format returns a plain text representation of differences
func (f *TextFormatter) Format(result *Result) string {
	if !result.HasChanges() {
		return "No differences found\n"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d difference(s):\n\n", len(result.Differences)))

	for _, diff := range result.Differences {
		switch diff.Type {
		case DiffTypeAdded:
			value := diff.NewValue
			if f.MaskSecrets {
				value = maskValue(value)
			}
			builder.WriteString(fmt.Sprintf("+ %s=%s\n", diff.Key, value))
		case DiffTypeRemoved:
			value := diff.OldValue
			if f.MaskSecrets {
				value = maskValue(value)
			}
			builder.WriteString(fmt.Sprintf("- %s=%s\n", diff.Key, value))
		case DiffTypeChanged:
			oldValue := diff.OldValue
			newValue := diff.NewValue
			if f.MaskSecrets {
				oldValue = maskValue(oldValue)
				newValue = maskValue(newValue)
			}
			builder.WriteString(fmt.Sprintf("~ %s\n", diff.Key))
			builder.WriteString(fmt.Sprintf("  - %s\n", oldValue))
			builder.WriteString(fmt.Sprintf("  + %s\n", newValue))
		}
	}

	builder.WriteString(fmt.Sprintf("\n%s\n", result.Summary()))
	return builder.String()
}

// maskValue masks a value for secret protection
func maskValue(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) <= 4 {
		return "****"
	}
	// Show first 2 and last 2 characters
	return value[:2] + "****" + value[len(value)-2:]
}

// JSONFormatter formats differences as JSON (placeholder for future implementation)
type JSONFormatter struct {
	MaskSecrets bool
}

// Format returns a JSON representation of differences
func (f *JSONFormatter) Format(result *Result) string {
	// Simplified JSON output for now
	var builder strings.Builder
	builder.WriteString("{\n")
	builder.WriteString(fmt.Sprintf("  \"total\": %d,\n", len(result.Differences)))
	builder.WriteString(fmt.Sprintf("  \"summary\": \"%s\"\n", result.Summary()))
	builder.WriteString("}\n")
	return builder.String()
}
