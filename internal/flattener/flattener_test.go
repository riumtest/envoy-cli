package flattener_test

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/flattener"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestFlatten_AddsPrefix(t *testing.T) {
	in := entries("HOST", "localhost", "PORT", "5432")
	out := flattener.Flatten(in, "DB", flattener.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" || out[1].Key != "DB_PORT" {
		t.Errorf("unexpected keys: %v", out)
	}
}

func TestFlatten_EmptyNamespace(t *testing.T) {
	in := entries("KEY", "val")
	out := flattener.Flatten(in, "", flattener.DefaultOptions())
	if out[0].Key != "KEY" {
		t.Errorf("expected KEY, got %s", out[0].Key)
	}
}

func TestFlatten_DeduplicatesAfterPrefix(t *testing.T) {
	in := entries("HOST", "a", "HOST", "b")
	out := flattener.Flatten(in, "APP", flattener.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry after dedup, got %d", len(out))
	}
	if out[0].Value != "a" {
		t.Errorf("expected first value 'a', got %s", out[0].Value)
	}
}

func TestFlatten_LowercaseOption(t *testing.T) {
	opts := flattener.Options{Separator: "_", Uppercase: false}
	in := entries("host", "localhost")
	out := flattener.Flatten(in, "db", opts)
	if out[0].Key != "db_host" {
		t.Errorf("expected db_host, got %s", out[0].Key)
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	opts := flattener.Options{Separator: ".", Uppercase: false}
	in := entries("host", "localhost")
	out := flattener.Flatten(in, "db", opts)
	if out[0].Key != "db.host" {
		t.Errorf("expected db.host, got %s", out[0].Key)
	}
}

func TestUnflatten_StripsPrefix(t *testing.T) {
	in := entries("DB_HOST", "localhost", "DB_PORT", "5432")
	out := flattener.Unflatten(in, "DB", flattener.DefaultOptions())
	if out[0].Key != "HOST" || out[1].Key != "PORT" {
		t.Errorf("unexpected keys after unflatten: %v", out)
	}
}

func TestUnflatten_LeavesNonMatchingKeysIntact(t *testing.T) {
	in := entries("DB_HOST", "localhost", "APP_ENV", "prod")
	out := flattener.Unflatten(in, "DB", flattener.DefaultOptions())
	if out[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", out[0].Key)
	}
	if out[1].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV unchanged, got %s", out[1].Key)
	}
}

func TestUnflatten_EmptyNamespace(t *testing.T) {
	in := entries("KEY", "val")
	out := flattener.Unflatten(in, "", flattener.DefaultOptions())
	if out[0].Key != "KEY" {
		t.Errorf("expected KEY unchanged, got %s", out[0].Key)
	}
}
