package linter

import (
	"testing"
)

func mkEntries() []Entry {
	return []Entry{
		{Line: 1, Key: "APP_ENV", Value: "production"},
		{Line: 2, Key: "DB_HOST", Value: "localhost"},
	}
}

func TestLint_Clean(t *testing.T) {
	res := Lint(mkEntries())
	if len(res.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(res.Issues))
	}
	if res.HasErrors() {
		t.Fatal("expected HasErrors to be false")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := []Entry{{Line: 1, Key: "app_env", Value: "dev"}}
	res := Lint(entries)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if res.Issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", res.Issues[0].Severity)
	}
}

func TestLint_MixedCaseKey(t *testing.T) {
	entries := []Entry{{Line: 3, Key: "AppEnv", Value: "staging"}}
	res := Lint(entries)
	if len(res.Issues) == 0 {
		t.Fatal("expected at least one issue for mixed-case key")
	}
}

func TestLint_KeyStartsWithDigit(t *testing.T) {
	entries := []Entry{{Line: 5, Key: "1BAD_KEY", Value: "value"}}
	res := Lint(entries)

	var hasError bool
	for _, iss := range res.Issues {
		if iss.Severity == "error" {
			hasError = true
		}
	}
	if !hasError {
		t.Fatal("expected an error-severity issue for key starting with digit")
	}
	if !res.HasErrors() {
		t.Fatal("expected HasErrors to return true")
	}
}

func TestLint_ValueWithWhitespace(t *testing.T) {
	entries := []Entry{{Line: 7, Key: "CLEAN_KEY", Value: " spaced "}}
	res := Lint(entries)
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue for whitespace value, got %d", len(res.Issues))
	}
	if res.Issues[0].Severity != "warn" {
		t.Errorf("expected warn, got %s", res.Issues[0].Severity)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	entries := []Entry{
		{Line: 1, Key: "good_key", Value: " bad value "},
		{Line: 2, Key: "GOOD_KEY", Value: "clean"},
		{Line: 3, Key: "3INVALID", Value: "x"},
	}
	res := Lint(entries)
	// good_key → warn (case) + warn (whitespace) = 2; 3INVALID → warn (case) + error (digit) = 2
	if len(res.Issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d", len(res.Issues))
	}
}
