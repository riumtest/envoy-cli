package limiter_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/envfile"
	"github.com/your-org/envoy-cli/internal/limiter"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "GAMMA", Value: "3"},
		{Key: "DELTA", Value: "4"},
		{Key: "EPSILON", Value: "5"},
	}
}

func TestLimit_NoConstraints(t *testing.T) {
	result := limiter.Limit(entries(), limiter.DefaultOptions())
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestLimit_WithLimit(t *testing.T) {
	result := limiter.Limit(entries(), limiter.Options{Limit: 3})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Key != "ALPHA" || result[2].Key != "GAMMA" {
		t.Errorf("unexpected entries: %+v", result)
	}
}

func TestLimit_WithOffset(t *testing.T) {
	result := limiter.Limit(entries(), limiter.Options{Offset: 2})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Key != "GAMMA" {
		t.Errorf("expected GAMMA as first entry, got %s", result[0].Key)
	}
}

func TestLimit_WithLimitAndOffset(t *testing.T) {
	result := limiter.Limit(entries(), limiter.Options{Limit: 2, Offset: 1})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "BETA" || result[1].Key != "GAMMA" {
		t.Errorf("unexpected entries: %+v", result)
	}
}

func TestLimit_OffsetBeyondLength(t *testing.T) {
	result := limiter.Limit(entries(), limiter.Options{Offset: 100})
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(result))
	}
}

func TestLimit_EmptyEntries(t *testing.T) {
	result := limiter.Limit([]envfile.Entry{}, limiter.Options{Limit: 5})
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(result))
	}
}

func TestLimit_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	result := limiter.Limit(orig, limiter.Options{Limit: 2})
	result[0].Key = "MUTATED"
	if orig[0].Key == "MUTATED" {
		t.Error("Limit mutated the original slice")
	}
}
