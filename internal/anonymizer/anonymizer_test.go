package anonymizer_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/anonymizer"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "PORT", Value: "8080"},
		{Key: "EMPTY_VAL", Value: ""},
	}
}

func TestAnonymize_AllKeys(t *testing.T) {
	out := anonymizer.Anonymize(entries(), anonymizer.DefaultOptions())
	for _, e := range out {
		if e.Value == "" {
			continue // empty values stay empty
		}
		if !strings.HasPrefix(e.Value, "anon_") {
			t.Errorf("key %s: expected anon_ prefix, got %q", e.Key, e.Value)
		}
	}
}

func TestAnonymize_SensitiveOnly(t *testing.T) {
	opts := anonymizer.Options{Prefix: "anon", SensitiveOnly: true}
	out := anonymizer.Anonymize(entries(), opts)

	byKey := map[string]string{}
	for _, e := range out {
		byKey[e.Key] = e.Value
	}

	// sensitive keys should be replaced
	if byKey["DB_PASSWORD"] == "s3cr3t" {
		t.Error("DB_PASSWORD should have been anonymized")
	}
	if byKey["API_KEY"] == "abc123" {
		t.Error("API_KEY should have been anonymized")
	}
	// non-sensitive keys should be untouched
	if byKey["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", byKey["APP_NAME"])
	}
	if byKey["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", byKey["PORT"])
	}
}

func TestAnonymize_Deterministic(t *testing.T) {
	opts := anonymizer.DefaultOptions()
	a := anonymizer.Anonymize(entries(), opts)
	b := anonymizer.Anonymize(entries(), opts)
	for i := range a {
		if a[i].Value != b[i].Value {
			t.Errorf("token not deterministic for key %s", a[i].Key)
		}
	}
}

func TestAnonymize_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	anonymizer.Anonymize(orig, anonymizer.DefaultOptions())
	if orig[0].Value != "myapp" {
		t.Error("original entries were mutated")
	}
}

func TestAnonymize_CustomPrefix(t *testing.T) {
	opts := anonymizer.Options{Prefix: "REDACTED"}
	out := anonymizer.Anonymize(entries()[:1], opts)
	if !strings.HasPrefix(out[0].Value, "redacted_") {
		t.Errorf("expected redacted_ prefix, got %q", out[0].Value)
	}
}

func TestAnonymize_EmptyValueUnchanged(t *testing.T) {
	out := anonymizer.Anonymize(entries(), anonymizer.DefaultOptions())
	for _, e := range out {
		if e.Key == "EMPTY_VAL" && e.Value != "" {
			t.Errorf("empty value should remain empty, got %q", e.Value)
		}
	}
}
