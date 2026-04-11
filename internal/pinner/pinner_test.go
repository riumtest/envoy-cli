package pinner_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/pinner"
)

func entries(kv ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, envfile.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return out
}

func TestPin_NoDrift(t *testing.T) {
	pinned := entries("HOST", "localhost", "PORT", "8080")
	current := entries("HOST", "localhost", "PORT", "8080")

	r := pinner.Pin(pinned, current)

	if len(r.Drifted) != 0 || len(r.Missing) != 0 || len(r.New) != 0 {
		t.Fatalf("expected clean result, got %s", r.Summary())
	}
}

func TestPin_DriftedValue(t *testing.T) {
	pinned := entries("HOST", "localhost", "PORT", "8080")
	current := entries("HOST", "prod.example.com", "PORT", "8080")

	r := pinner.Pin(pinned, current)

	if len(r.Drifted) != 1 {
		t.Fatalf("expected 1 drifted, got %d", len(r.Drifted))
	}
	if r.Drifted[0].Key != "HOST" {
		t.Errorf("expected HOST to drift, got %s", r.Drifted[0].Key)
	}
	if r.Drifted[0].Pinned != "localhost" {
		t.Errorf("unexpected pinned value: %s", r.Drifted[0].Pinned)
	}
	if r.Drifted[0].Current != "prod.example.com" {
		t.Errorf("unexpected current value: %s", r.Drifted[0].Current)
	}
}

func TestPin_MissingKey(t *testing.T) {
	pinned := entries("HOST", "localhost", "PORT", "8080")
	current := entries("HOST", "localhost")

	r := pinner.Pin(pinned, current)

	if len(r.Missing) != 1 || r.Missing[0] != "PORT" {
		t.Errorf("expected PORT missing, got %v", r.Missing)
	}
}

func TestPin_NewKey(t *testing.T) {
	pinned := entries("HOST", "localhost")
	current := entries("HOST", "localhost", "DEBUG", "true")

	r := pinner.Pin(pinned, current)

	if len(r.New) != 1 || r.New[0] != "DEBUG" {
		t.Errorf("expected DEBUG as new key, got %v", r.New)
	}
}

func TestPin_Summary(t *testing.T) {
	pinned := entries("A", "1", "B", "2", "C", "3")
	current := entries("A", "changed", "D", "4")

	r := pinner.Pin(pinned, current)

	summary := r.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	// 1 drifted (A), 2 missing (B,C), 1 new (D)
	expected := "1 drifted, 2 missing, 1 new"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
