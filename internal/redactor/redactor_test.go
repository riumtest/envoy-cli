package redactor_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/envfile"
	"github.com/yourusername/envoy-cli/internal/redactor"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRedact_SensitiveKeysAreReplaced(t *testing.T) {
	in := entries("SECRET_KEY", "abc123", "APP_NAME", "myapp")
	out := redactor.Redact(in, redactor.Options{})

	if out[0].Value != "***" || !out[0].Redacted {
		t.Errorf("expected SECRET_KEY to be redacted")
	}
	if out[1].Value != "myapp" || out[1].Redacted {
		t.Errorf("expected APP_NAME to remain unredacted")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	in := entries("DB_PASSWORD", "hunter2")
	out := redactor.Redact(in, redactor.Options{Placeholder: "<hidden>"})

	if out[0].Value != "<hidden>" {
		t.Errorf("expected custom placeholder, got %q", out[0].Value)
	}
}

func TestRedact_ExtraPatterns(t *testing.T) {
	in := entries("STRIPE_TOKEN", "tok_live_xyz", "APP_PORT", "8080")
	out := redactor.Redact(in, redactor.Options{ExtraPatterns: []string{"token"}})

	if !out[0].Redacted {
		t.Errorf("expected STRIPE_TOKEN to be redacted via extra pattern")
	}
	if out[1].Redacted {
		t.Errorf("expected APP_PORT to remain unredacted")
	}
}

func TestRedact_EmptyEntries(t *testing.T) {
	out := redactor.Redact([]envfile.Entry{}, redactor.Options{})
	if len(out) != 0 {
		t.Errorf("expected empty output")
	}
}

func TestRedactString_ReplacesInlineValues(t *testing.T) {
	raw := "DB_PASSWORD=secret\nAPP_ENV=production"
	got := redactor.RedactString(raw, redactor.Options{})

	if got != "DB_PASSWORD=***\nAPP_ENV=production" {
		t.Errorf("unexpected output:\n%s", got)
	}
}

func TestRedactString_PreservesComments(t *testing.T) {
	raw := "# This is a comment\nAPI_KEY=mykey"
	got := redactor.RedactString(raw, redactor.Options{})

	if got != "# This is a comment\nAPI_KEY=***" {
		t.Errorf("unexpected output:\n%s", got)
	}
}
