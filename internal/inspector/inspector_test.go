package inspector_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/inspector"
	"envoy-cli/internal/masker"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DATABASE_URL", Value: "https://db.example.com"},
		{Key: "DEBUG", Value: "true"},
		{Key: "PORT", Value: "8080"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "EMPTY_VAR", Value: ""},
	}
}

func TestInspect_Counts(t *testing.T) {
	r := inspector.Inspect(entries(), nil)
	if r.TotalKeys != 6 {
		t.Errorf("expected 6 total keys, got %d", r.TotalKeys)
	}
	if r.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", r.EmptyValues)
	}
	if r.URLValues != 1 {
		t.Errorf("expected 1 URL value, got %d", r.URLValues)
	}
	if r.BooleanValues != 1 {
		t.Errorf("expected 1 boolean value, got %d", r.BooleanValues)
	}
	if r.NumericValues != 1 {
		t.Errorf("expected 1 numeric value, got %d", r.NumericValues)
	}
}

func TestInspect_SensitiveKeys(t *testing.T) {
	m := masker.New()
	r := inspector.Inspect(entries(), m)
	if len(r.SensitiveKeys) == 0 {
		t.Fatal("expected at least one sensitive key")
	}
	found := false
	for _, k := range r.SensitiveKeys {
		if k == "SECRET_KEY" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected SECRET_KEY in sensitive keys, got %v", r.SensitiveKeys)
	}
}

func TestInspect_NoMasker(t *testing.T) {
	r := inspector.Inspect(entries(), nil)
	if len(r.SensitiveKeys) != 0 {
		t.Errorf("expected no sensitive keys without masker, got %v", r.SensitiveKeys)
	}
}

func TestInspect_KeysPopulated(t *testing.T) {
	r := inspector.Inspect(entries(), nil)
	if len(r.Keys) != 6 {
		t.Errorf("expected 6 keys in report, got %d", len(r.Keys))
	}
	if r.Keys[0] != "APP_NAME" {
		t.Errorf("expected first key APP_NAME, got %s", r.Keys[0])
	}
}

func TestInspect_EmptyEntries(t *testing.T) {
	r := inspector.Inspect([]envfile.Entry{}, nil)
	if r.TotalKeys != 0 {
		t.Errorf("expected 0 total keys for empty input, got %d", r.TotalKeys)
	}
}
