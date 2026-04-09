package differ

import (
	"strings"
	"testing"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	formatter := &TextFormatter{MaskSecrets: false}
	result := &Result{Differences: []Difference{}}

	output := formatter.Format(result)

	if !strings.Contains(output, "No differences") {
		t.Errorf("Expected 'No differences' message, got: %s", output)
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	formatter := &TextFormatter{MaskSecrets: false}
	result := &Result{
		Differences: []Difference{
			{Key: "NEW_KEY", Type: DiffTypeAdded, NewValue: "new_value"},
			{Key: "OLD_KEY", Type: DiffTypeRemoved, OldValue: "old_value"},
			{Key: "CHANGED_KEY", Type: DiffTypeChanged, OldValue: "old", NewValue: "new"},
		},
	}

	output := formatter.Format(result)

	if !strings.Contains(output, "+ NEW_KEY=new_value") {
		t.Errorf("Expected added key in output")
	}
	if !strings.Contains(output, "- OLD_KEY=old_value") {
		t.Errorf("Expected removed key in output")
	}
	if !strings.Contains(output, "~ CHANGED_KEY") {
		t.Errorf("Expected changed key in output")
	}
}

func TestTextFormatter_WithMasking(t *testing.T) {
	formatter := &TextFormatter{MaskSecrets: true}
	result := &Result{
		Differences: []Difference{
			{Key: "SECRET_KEY", Type: DiffTypeAdded, NewValue: "supersecret123"},
		},
	}

	output := formatter.Format(result)

	if strings.Contains(output, "supersecret123") {
		t.Errorf("Secret should be masked, but found in output: %s", output)
	}
	if !strings.Contains(output, "su****23") {
		t.Errorf("Expected masked value in output, got: %s", output)
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"short", "****"},
		{"longer_secret", "lo****et"},
		{"supersecret123", "su****23"},
		{"", ""},
	}

	for _, tt := range tests {
		result := maskValue(tt.input)
		if result != tt.expected {
			t.Errorf("maskValue(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{MaskSecrets: false}
	result := &Result{
		Differences: []Difference{
			{Key: "KEY1", Type: DiffTypeAdded},
		},
	}

	output := formatter.Format(result)

	if !strings.Contains(output, "total") {
		t.Errorf("Expected JSON output with 'total' field")
	}
}
