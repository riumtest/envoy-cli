package validator_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/validator"
)

func entries(tuples ...interface{}) []validator.Entry {
	var out []validator.Entry
	for i := 0; i+2 < len(tuples); i += 3 {
		out = append(out, validator.Entry{
			Line:  tuples[i].(int),
			Key:   tuples[i+1].(string),
			Value: tuples[i+2].(string),
		})
	}
	return out
}

func TestValidate_Clean(t *testing.T) {
	result := validator.Validate(entries(1, "APP_ENV", "production", 2, "PORT", "8080"))
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	result := validator.Validate(entries(1, "APP_ENV", "dev", 2, "APP_ENV", "prod"))
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != validator.SeverityWarning {
		t.Errorf("expected warning severity, got %s", result.Issues[0].Severity)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	result := validator.Validate(entries(1, "SECRET_KEY", ""))
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != validator.SeverityWarning {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	result := validator.Validate(entries(1, "", "somevalue"))
	if !result.HasErrors() {
		t.Error("expected error for empty key")
	}
}

func TestValidate_KeyWithWhitespace(t *testing.T) {
	result := validator.Validate(entries(1, "MY KEY", "val"))
	if !result.HasErrors() {
		t.Error("expected error for key with whitespace")
	}
}

func TestValidate_HasErrors_False(t *testing.T) {
	result := validator.Validate(entries(1, "GOOD_KEY", ""))
	if result.HasErrors() {
		t.Error("expected no errors, only warnings")
	}
}

func TestIssue_String(t *testing.T) {
	issue := validator.Issue{
		Line:     3,
		Key:      "FOO",
		Message:  "value is empty",
		Severity: validator.SeverityWarning,
	}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
