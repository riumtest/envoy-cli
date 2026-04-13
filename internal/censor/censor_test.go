package censor_test

import (
	"testing"

	"github.com/envoy-cli/internal/censor"
	"github.com/envoy-cli/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DATABASE_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "DEBUG", Value: "true"},
		{Key: "AUTH_TOKEN", Value: "tok-xyz"},
	}
}

func TestCensor_SensitiveKeysReplaced(t *testing.T) {
	out := censor.Censor(entries(), censor.Options{})
	for _, e := range out {
		switch e.Key {
		case "DATABASE_PASSWORD", "API_KEY", "AUTH_TOKEN":
			if e.Value != censor.DefaultPlaceholder {
				t.Errorf("expected %q to be censored, got %q", e.Key, e.Value)
			}
		case "APP_NAME", "DEBUG":
			if e.Value == censor.DefaultPlaceholder {
				t.Errorf("expected %q to be preserved, got censored", e.Key)
			}
		}
	}
}

func TestCensor_CustomPlaceholder(t *testing.T) {
	out := censor.Censor(entries(), censor.Options{Placeholder: "[REDACTED]"})
	for _, e := range out {
		if e.Key == "DATABASE_PASSWORD" && e.Value != "[REDACTED]" {
			t.Errorf("expected [REDACTED], got %q", e.Value)
		}
	}
}

func TestCensor_ExplicitKeys(t *testing.T) {
	out := censor.Censor(entries(), censor.Options{
		Keys: []string{"APP_NAME", "DEBUG"},
	})
	for _, e := range out {
		if (e.Key == "APP_NAME" || e.Key == "DEBUG") && e.Value != censor.DefaultPlaceholder {
			t.Errorf("expected %q to be censored via explicit list", e.Key)
		}
	}
}

func TestCensor_ExtraPatterns(t *testing.T) {
	out := censor.Censor(entries(), censor.Options{
		ExtraPatterns: []string{"debug"},
	})
	for _, e := range out {
		if e.Key == "DEBUG" && e.Value != censor.DefaultPlaceholder {
			t.Errorf("expected DEBUG to be censored via extra pattern")
		}
	}
}

func TestCensor_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	censor.Censor(orig, censor.Options{})
	if orig[1].Value != "s3cr3t" {
		t.Error("original entries were mutated")
	}
}

func TestCensor_EmptyEntries(t *testing.T) {
	out := censor.Censor([]envfile.Entry{}, censor.Options{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}
