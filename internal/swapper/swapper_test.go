package swapper_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/swapper"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestSwap_Basic(t *testing.T) {
	in := entries("HOST", "localhost", "PORT", "8080")
	out := swapper.Swap(in, swapper.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "localhost" || out[0].Value != "HOST" {
		t.Errorf("unexpected first entry: %+v", out[0])
	}
	if out[1].Key != "8080" || out[1].Value != "PORT" {
		t.Errorf("unexpected second entry: %+v", out[1])
	}
}

func TestSwap_SkipEmpty(t *testing.T) {
	in := entries("HOST", "", "PORT", "8080")
	out := swapper.Swap(in, swapper.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "8080" {
		t.Errorf("unexpected key: %s", out[0].Key)
	}
}

func TestSwap_KeepEmpty(t *testing.T) {
	in := entries("HOST", "")
	opts := swapper.Options{SkipEmpty: false}
	out := swapper.Swap(in, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "" || out[0].Value != "HOST" {
		t.Errorf("unexpected entry: %+v", out[0])
	}
}

func TestSwap_SkipDuplicates(t *testing.T) {
	in := entries("A", "same", "B", "same")
	opts := swapper.Options{SkipEmpty: true, SkipDuplicates: true}
	out := swapper.Swap(in, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Value != "A" {
		t.Errorf("expected first winner A, got %s", out[0].Value)
	}
}

func TestSwap_DoesNotMutateOriginal(t *testing.T) {
	in := entries("KEY", "val")
	swapper.Swap(in, swapper.DefaultOptions())
	if in[0].Key != "KEY" || in[0].Value != "val" {
		t.Error("original entries were mutated")
	}
}
