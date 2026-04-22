package linker_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/linker"
)

func entries(kvs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestLink_ResolvesRule(t *testing.T) {
	src := map[string][]envfile.Entry{
		"prod.env": entries("DB_PASS", "s3cr3t"),
	}
	dst := entries("APP_NAME", "myapp")
	rules := []linker.Rule{
		{FromFile: "prod.env", FromKey: "DB_PASS", ToKey: "DB_PASS"},
	}

	out, results, err := linker.Link(dst, src, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].OK {
		t.Fatalf("expected 1 ok result, got %+v", results)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[1].Key != "DB_PASS" || out[1].Value != "s3cr3t" {
		t.Errorf("unexpected entry: %+v", out[1])
	}
}

func TestLink_OverwritesExistingKey(t *testing.T) {
	src := map[string][]envfile.Entry{
		"prod.env": entries("TOKEN", "new-token"),
	}
	dst := entries("TOKEN", "old-token")
	rules := []linker.Rule{
		{FromFile: "prod.env", FromKey: "TOKEN", ToKey: "TOKEN"},
	}

	out, _, err := linker.Link(dst, src, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "new-token" {
		t.Errorf("expected overwrite, got %q", out[0].Value)
	}
}

func TestLink_MissingSourceFile(t *testing.T) {
	rules := []linker.Rule{
		{FromFile: "ghost.env", FromKey: "X", ToKey: "X"},
	}
	out, results, err := linker.Link(nil, nil, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty dst, got %d", len(out))
	}
	if results[0].OK {
		t.Error("expected not-ok result for missing file")
	}
}

func TestLink_MissingKey(t *testing.T) {
	src := map[string][]envfile.Entry{
		"a.env": entries("FOO", "bar"),
	}
	rules := []linker.Rule{
		{FromFile: "a.env", FromKey: "MISSING", ToKey: "MISSING"},
	}
	_, results, err := linker.Link(nil, src, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].OK {
		t.Error("expected not-ok result for missing key")
	}
}

func TestLink_NoRules(t *testing.T) {
	dst := entries("A", "1", "B", "2")
	out, results, err := linker.Link(dst, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries unchanged, got %d", len(out))
	}
}
