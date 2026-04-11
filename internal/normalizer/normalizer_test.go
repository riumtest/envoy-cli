package normalizer_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/normalizer"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	in := entries("db_host", "localhost", "api_key", "secret")
	opts := normalizer.DefaultOptions()
	out := normalizer.Normalize(in, opts)

	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
	if out[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", out[1].Key)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	in := entries("HOST", "  localhost  ")
	opts := normalizer.DefaultOptions()
	out := normalizer.Normalize(in, opts)

	if out[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", out[0].Value)
	}
}

func TestNormalize_QuoteValues(t *testing.T) {
	in := entries("HOST", "localhost")
	opts := normalizer.DefaultOptions()
	opts.QuoteValues = true
	out := normalizer.Normalize(in, opts)

	if out[0].Value != `"localhost"` {
		t.Errorf("expected quoted value, got %s", out[0].Value)
	}
}

func TestNormalize_QuoteValues_AlreadyQuoted(t *testing.T) {
	in := entries("HOST", `"localhost"`)
	opts := normalizer.DefaultOptions()
	opts.QuoteValues = true
	out := normalizer.Normalize(in, opts)

	if out[0].Value != `"localhost"` {
		t.Errorf("expected value unchanged, got %s", out[0].Value)
	}
}

func TestNormalize_StripExport(t *testing.T) {
	in := entries("export MY_VAR", "value")
	opts := normalizer.DefaultOptions()
	out := normalizer.Normalize(in, opts)

	if out[0].Key != "MY_VAR" {
		t.Errorf("expected MY_VAR, got %s", out[0].Key)
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	in := entries("db_host", "  val  ")
	opts := normalizer.DefaultOptions()
	normalizer.Normalize(in, opts)

	if in[0].Key != "db_host" {
		t.Error("original entry key was mutated")
	}
	if in[0].Value != "  val  " {
		t.Error("original entry value was mutated")
	}
}

func TestNormalize_EmptyEntries(t *testing.T) {
	out := normalizer.Normalize([]envfile.Entry{}, normalizer.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}
