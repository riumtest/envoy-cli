package inspector_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/inspector"
	"github.com/user/envoy-cli/internal/masker"
)

func entries(pairs ...string) []inspector.Entry {
	var out []inspector.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, inspector.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestInspect_Counts(t *testing.T) {
	m := masker.New()
	ents := entries(
		"DATABASE_URL", "postgres://localhost/db",
		"API_KEY", "abc123secret",
		"DEBUG", "true",
		"PORT", "8080",
		"EMPTY_VAR", "",
	)
	r := inspector.Inspect(ents, m)

	if r.Total != 5 {
		t.Errorf("expected Total=5, got %d", r.Total)
	}
	if r.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", r.Empty)
	}
	if r.URLs != 1 {
		t.Errorf("expected URLs=1, got %d", r.URLs)
	}
	if r.Booleans != 1 {
		t.Errorf("expected Booleans=1, got %d", r.Booleans)
	}
	if r.Numeric != 1 {
		t.Errorf("expected Numeric=1, got %d", r.Numeric)
	}
}

func TestInspect_SensitiveKeys(t *testing.T) {
	m := masker.New()
	ents := entries(
		"SECRET_KEY", "topsecret",
		"PASSWORD", "hunter2",
		"HOST", "localhost",
	)
	r := inspector.Inspect(ents, m)
	if r.Sensitive != 2 {
		t.Errorf("expected Sensitive=2, got %d", r.Sensitive)
	}
}

func TestInspect_NoMasker(t *testing.T) {
	ents := entries("KEY", "value")
	r := inspector.Inspect(ents, nil)
	if r.Sensitive != 0 {
		t.Errorf("expected Sensitive=0 with nil masker, got %d", r.Sensitive)
	}
}

func TestInspect_KeysPopulated(t *testing.T) {
	ents := entries("A", "1", "B", "2", "C", "3")
	r := inspector.Inspect(ents, nil)
	if len(r.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(r.Keys))
	}
}

func TestInspect_EmptyInput(t *testing.T) {
	r := inspector.Inspect(nil, nil)
	if r.Total != 0 {
		t.Errorf("expected Total=0, got %d", r.Total)
	}
}
