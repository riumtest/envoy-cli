package sampler_test

import (
	"testing"

	"github.com/envoy-cli/internal/envfile"
	"github.com/envoy-cli/internal/sampler"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "API_TOKEN", Value: "tok123"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestSample_ReturnsN(t *testing.T) {
	opts := sampler.DefaultOptions()
	opts.N = 3
	got := sampler.Sample(entries(), opts)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestSample_NGreaterThanLen(t *testing.T) {
	opts := sampler.DefaultOptions()
	opts.N = 100
	got := sampler.Sample(entries(), opts)
	if len(got) != len(entries()) {
		t.Fatalf("expected all %d entries, got %d", len(entries()), len(got))
	}
}

func TestSample_Deterministic(t *testing.T) {
	opts := sampler.DefaultOptions()
	opts.N = 4
	opts.Seed = 99
	a := sampler.Sample(entries(), opts)
	b := sampler.Sample(entries(), opts)
	if len(a) != len(b) {
		t.Fatal("lengths differ")
	}
	for i := range a {
		if a[i].Key != b[i].Key {
			t.Fatalf("index %d: %s != %s", i, a[i].Key, b[i].Key)
		}
	}
}

func TestSample_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	copy := entries()
	opts := sampler.DefaultOptions()
	opts.N = 3
	sampler.Sample(orig, opts)
	for i := range orig {
		if orig[i].Key != copy[i].Key {
			t.Fatalf("original mutated at index %d", i)
		}
	}
}

func TestSample_SensitiveOnly(t *testing.T) {
	opts := sampler.DefaultOptions()
	opts.N = 10
	opts.SensitiveOnly = true
	got := sampler.Sample(entries(), opts)
	for _, e := range got {
		if e.Key != "DB_PASSWORD" && e.Key != "API_TOKEN" && e.Key != "SECRET_KEY" {
			t.Fatalf("non-sensitive key included: %s", e.Key)
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 sensitive entries, got %d", len(got))
	}
}

func TestSample_ZeroN_ReturnsAll(t *testing.T) {
	opts := sampler.DefaultOptions()
	opts.N = 0
	got := sampler.Sample(entries(), opts)
	if len(got) != len(entries()) {
		t.Fatalf("expected all entries for N=0, got %d", len(got))
	}
}
