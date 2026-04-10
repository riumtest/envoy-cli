package comparator_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/comparator"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_SharedKeys(t *testing.T) {
	left := entries("HOST", "localhost", "PORT", "8080")
	right := entries("HOST", "localhost", "PORT", "8080")
	r := comparator.Compare(left, right)
	if len(r.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(r.SharedKeys))
	}
	if len(r.MismatchedKeys) != 0 || len(r.OnlyInLeft) != 0 || len(r.OnlyInRight) != 0 {
		t.Error("expected no mismatches or unique keys")
	}
}

func TestCompare_MismatchedKeys(t *testing.T) {
	left := entries("PORT", "8080")
	right := entries("PORT", "9090")
	r := comparator.Compare(left, right)
	if len(r.MismatchedKeys) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(r.MismatchedKeys))
	}
	m := r.MismatchedKeys[0]
	if m.Key != "PORT" || m.LeftValue != "8080" || m.RightValue != "9090" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := entries("HOST", "localhost", "SECRET", "abc")
	right := entries("HOST", "localhost")
	r := comparator.Compare(left, right)
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0] != "SECRET" {
		t.Errorf("expected SECRET only in left, got %v", r.OnlyInLeft)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := entries("HOST", "localhost")
	right := entries("HOST", "localhost", "DEBUG", "true")
	r := comparator.Compare(left, right)
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "DEBUG" {
		t.Errorf("expected DEBUG only in right, got %v", r.OnlyInRight)
	}
}

func TestCompare_EmptyInputs(t *testing.T) {
	r := comparator.Compare(nil, nil)
	if len(r.SharedKeys) != 0 || len(r.MismatchedKeys) != 0 {
		t.Error("expected empty result for nil inputs")
	}
}

func TestOverlapRatio_Full(t *testing.T) {
	left := entries("A", "1", "B", "2")
	right := entries("A", "1", "B", "2")
	r := comparator.Compare(left, right)
	if ratio := comparator.OverlapRatio(r); ratio != 1.0 {
		t.Errorf("expected ratio 1.0, got %f", ratio)
	}
}

func TestOverlapRatio_Zero(t *testing.T) {
	r := comparator.Result{}
	if ratio := comparator.OverlapRatio(r); ratio != 0.0 {
		t.Errorf("expected ratio 0.0, got %f", ratio)
	}
}
