package templater_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/templater"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRender_NoPlaceholders(t *testing.T) {
	tmpl := entries("HOST", "localhost", "PORT", "5432")
	subs := map[string]string{}
	r := templater.Render(tmpl, subs)
	if len(r.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", r.Missing)
	}
	if r.Entries[0].Value != "localhost" || r.Entries[1].Value != "5432" {
		t.Fatalf("unexpected values: %+v", r.Entries)
	}
}

func TestRender_AllSubstituted(t *testing.T) {
	tmpl := entries("DB_URL", "postgres://{{DB_USER}}:{{DB_PASS}}@{{DB_HOST}}/mydb")
	subs := map[string]string{
		"DB_USER": "admin",
		"DB_PASS": "secret",
		"DB_HOST": "db.example.com",
	}
	r := templater.Render(tmpl, subs)
	if len(r.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", r.Missing)
	}
	want := "postgres://admin:secret@db.example.com/mydb"
	if r.Entries[0].Value != want {
		t.Errorf("got %q, want %q", r.Entries[0].Value, want)
	}
}

func TestRender_MissingPlaceholders(t *testing.T) {
	tmpl := entries("API_URL", "https://{{API_HOST}}/{{API_VERSION}}")
	subs := map[string]string{"API_HOST": "api.example.com"}
	r := templater.Render(tmpl, subs)
	if len(r.Missing) != 1 || r.Missing[0] != "API_VERSION" {
		t.Fatalf("expected [API_VERSION] missing, got %v", r.Missing)
	}
	// Unresolved placeholder should remain in place.
	if r.Entries[0].Value != "https://api.example.com/{{API_VERSION}}" {
		t.Errorf("unexpected value: %q", r.Entries[0].Value)
	}
}

func TestRender_PreservesComment(t *testing.T) {
	tmpl := []envfile.Entry{{Key: "FOO", Value: "{{BAR}}", Comment: "# a comment"}}
	subs := map[string]string{"BAR": "baz"}
	r := templater.Render(tmpl, subs)
	if r.Entries[0].Comment != "# a comment" {
		t.Errorf("comment not preserved: %q", r.Entries[0].Comment)
	}
}

func TestBuildSubsFromEntries(t *testing.T) {
	env := entries("HOST", "localhost", "PORT", "8080")
	subs := templater.BuildSubsFromEntries(env)
	if subs["HOST"] != "localhost" || subs["PORT"] != "8080" {
		t.Errorf("unexpected subs: %v", subs)
	}
}

func TestSummary_NoMissing(t *testing.T) {
	r := templater.Result{Entries: make([]envfile.Entry, 3), Missing: nil}
	got := templater.Summary(r)
	if got != "rendered 3 entries, no missing placeholders" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithMissing(t *testing.T) {
	r := templater.Result{
		Entries: make([]envfile.Entry, 2),
		Missing: []string{"SECRET_KEY"},
	}
	got := templater.Summary(r)
	expected := "rendered 2 entries, 1 missing placeholder(s): SECRET_KEY"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
