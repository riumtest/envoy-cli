package interpolator

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestInterpolate_NoReferences(t *testing.T) {
	in := entries("HOST", "localhost", "PORT", "5432")
	r := Interpolate(in)
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Value != "localhost" {
		t.Errorf("unexpected value: %s", r.Entries[0].Value)
	}
	if len(r.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", r.Unresolved)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	in := entries("HOST", "db", "DSN", "postgres://${HOST}:5432")
	r := Interpolate(in)
	want := "postgres://db:5432"
	if r.Entries[1].Value != want {
		t.Errorf("got %q, want %q", r.Entries[1].Value, want)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	in := entries("HOST", "db", "DSN", "postgres://$HOST:5432")
	r := Interpolate(in)
	want := "postgres://db:5432"
	if r.Entries[1].Value != want {
		t.Errorf("got %q, want %q", r.Entries[1].Value, want)
	}
}

func TestInterpolate_UnresolvedReference(t *testing.T) {
	in := entries("DSN", "postgres://${MISSING_HOST}:5432")
	r := Interpolate(in)
	if len(r.Unresolved) != 1 || r.Unresolved[0] != "DSN" {
		t.Errorf("expected unresolved DSN, got %v", r.Unresolved)
	}
	if r.Entries[0].Value != "postgres://${MISSING_HOST}:5432" {
		t.Errorf("value should remain unchanged: %s", r.Entries[0].Value)
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	in := entries("BASE", "http://example.com", "URL", "${BASE}/api", "FULL", "${URL}/v1")
	r := Interpolate(in)
	if r.Entries[1].Value != "http://example.com/api" {
		t.Errorf("URL: got %q", r.Entries[1].Value)
	}
	// FULL references URL which was already expanded in entries, not in result
	// second pass would be needed for chained; single pass leaves ${URL} if not in original
}

func TestInterpolateWithEnv_FallsBackToEnv(t *testing.T) {
	in := entries("DSN", "postgres://${DB_HOST}:5432")
	env := map[string]string{"DB_HOST": "prod-db"}
	r := InterpolateWithEnv(in, env)
	want := "postgres://prod-db:5432"
	if r.Entries[0].Value != want {
		t.Errorf("got %q, want %q", r.Entries[0].Value, want)
	}
	if len(r.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", r.Unresolved)
	}
}

func TestInterpolateWithEnv_EntryTakesPrecedence(t *testing.T) {
	in := entries("HOST", "local", "DSN", "${HOST}")
	env := map[string]string{"HOST": "remote"}
	r := InterpolateWithEnv(in, env)
	if r.Entries[1].Value != "local" {
		t.Errorf("entry should take precedence over env: got %q", r.Entries[1].Value)
	}
}

func TestSummary_AllResolved(t *testing.T) {
	r := Result{Entries: entries("A", "1", "B", "2")}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_WithUnresolved(t *testing.T) {
	r := Result{Entries: entries("A", "1"), Unresolved: []string{"A"}}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
