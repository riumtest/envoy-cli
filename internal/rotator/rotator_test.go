package rotator_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/rotator"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "DB_PASSWORD", Value: "old-pass"},
		{Key: "API_KEY", Value: "old-key"},
		{Key: "APP_NAME", Value: "myapp"},
	}
}

func fixedGenerator(val string) func(string) (string, error) {
	return func(_ string) (string, error) { return val, nil }
}

func TestRotate_AllKeys(t *testing.T) {
	opts := rotator.DefaultOptions()
	opts.Generator = fixedGenerator("newvalue")

	out, results, err := rotator.Rotate(entries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, e := range out {
		if e.Value != "newvalue" {
			t.Errorf("key %s: expected newvalue, got %s", e.Key, e.Value)
		}
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	opts := rotator.DefaultOptions()
	opts.Keys = []string{"DB_PASSWORD"}
	opts.Generator = fixedGenerator("rotated")

	out, results, err := rotator.Rotate(entries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD rotated, got %s", results[0].Key)
	}
	for _, e := range out {
		if e.Key == "API_KEY" && e.Value != "old-key" {
			t.Errorf("API_KEY should not have been rotated")
		}
	}
}

func TestRotate_DryRun(t *testing.T) {
	opts := rotator.DefaultOptions()
	opts.DryRun = true
	opts.Generator = fixedGenerator("newvalue")

	out, results, err := rotator.Rotate(entries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results even in dry-run")
	}
	for _, e := range out {
		if e.Value == "newvalue" {
			t.Errorf("dry-run should not modify values, but %s was changed", e.Key)
		}
	}
}

func TestRotate_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	opts := rotator.DefaultOptions()
	opts.Generator = fixedGenerator("x")

	_, _, err := rotator.Rotate(orig, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orig[0].Value != "old-pass" {
		t.Error("original slice was mutated")
	}
}

func TestRotate_DefaultGenerator_ProducesHex(t *testing.T) {
	opts := rotator.DefaultOptions()
	opts.Keys = []string{"API_KEY"}

	_, results, err := rotator.Rotate(entries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	val := results[0].NewValue
	if len(val) != 64 {
		t.Errorf("expected 64-char hex string, got len %d: %s", len(val), val)
	}
	if strings.ContainsAny(val, "ghijklmnopqrstuvwxyzGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("value does not look like hex: %s", val)
	}
}
