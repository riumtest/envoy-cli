package stenciler_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/stenciler"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestStencil_ClearsValues(t *testing.T) {
	in := entries("DB_HOST", "localhost", "DB_PORT", "5432")
	out := stenciler.Stencil(in, stenciler.Options{TypedPlaceholders: false})
	for _, e := range out {
		if e.Value != "" {
			t.Errorf("expected empty value for %s, got %q", e.Key, e.Value)
		}
	}
}

func TestStencil_TypedPlaceholders(t *testing.T) {
	in := entries(
		"APP_DEBUG", "true",
		"APP_PORT", "8080",
		"APP_URL", "https://example.com",
		"APP_NAME", "myapp",
	)
	out := stenciler.Stencil(in, stenciler.DefaultOptions())
	want := map[string]string{
		"APP_DEBUG": "<bool>",
		"APP_PORT":  "<number>",
		"APP_URL":   "<url>",
		"APP_NAME":  "<string>",
	}
	for _, e := range out {
		if got, ok := want[e.Key]; ok {
			if e.Value != got {
				t.Errorf("key %s: want %q, got %q", e.Key, got, e.Value)
			}
		}
	}
}

func TestStencil_PreservesComments(t *testing.T) {
	in := []envfile.Entry{
		{Key: "", Value: "# database config"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	out := stenciler.Stencil(in, stenciler.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "# database config" {
		t.Errorf("expected comment to be preserved")
	}
}

func TestStencil_StripsComments_WhenDisabled(t *testing.T) {
	in := []envfile.Entry{
		{Key: "", Value: "# database config"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	opts := stenciler.Options{TypedPlaceholders: false, PreserveComments: false}
	out := stenciler.Stencil(in, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry (comment stripped), got %d", len(out))
	}
}

func TestRender_ProducesEnvFormat(t *testing.T) {
	in := entries("FOO", "<string>", "BAR", "<number>")
	got := stenciler.Render(in)
	if !strings.Contains(got, "FOO=<string>") {
		t.Errorf("expected FOO=<string> in output, got:\n%s", got)
	}
	if !strings.Contains(got, "BAR=<number>") {
		t.Errorf("expected BAR=<number> in output, got:\n%s", got)
	}
}

func TestRender_EmptyEntries(t *testing.T) {
	got := stenciler.Render([]envfile.Entry{})
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
