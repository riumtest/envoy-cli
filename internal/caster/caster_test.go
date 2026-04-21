package caster_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/caster"
	"github.com/user/envoy-cli/internal/envfile"
)

func entries(kvs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, envfile.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestCast_StringValue(t *testing.T) {
	res := caster.Cast(entries("HOST", "localhost"))
	if res[0].Kind != caster.KindString {
		t.Fatalf("expected string, got %s", res[0].Kind)
	}
	if res[0].Parsed != "localhost" {
		t.Fatalf("unexpected parsed value: %v", res[0].Parsed)
	}
}

func TestCast_BoolValue(t *testing.T) {
	for _, v := range []string{"true", "false", "True", "FALSE"} {
		res := caster.Cast(entries("FLAG", v))
		if res[0].Kind != caster.KindBool {
			t.Fatalf("value %q: expected bool, got %s", v, res[0].Kind)
		}
		if _, ok := res[0].Parsed.(bool); !ok {
			t.Fatalf("value %q: parsed is not bool", v)
		}
	}
}

func TestCast_IntValue(t *testing.T) {
	res := caster.Cast(entries("PORT", "8080"))
	if res[0].Kind != caster.KindInt {
		t.Fatalf("expected int, got %s", res[0].Kind)
	}
	if res[0].Parsed.(int64) != 8080 {
		t.Fatalf("unexpected parsed value: %v", res[0].Parsed)
	}
}

func TestCast_FloatValue(t *testing.T) {
	res := caster.Cast(entries("RATIO", "3.14"))
	if res[0].Kind != caster.KindFloat {
		t.Fatalf("expected float, got %s", res[0].Kind)
	}
	if res[0].Parsed.(float64) != 3.14 {
		t.Fatalf("unexpected parsed value: %v", res[0].Parsed)
	}
}

func TestCast_EmptyValue(t *testing.T) {
	res := caster.Cast(entries("EMPTY", ""))
	if res[0].Kind != caster.KindEmpty {
		t.Fatalf("expected empty, got %s", res[0].Kind)
	}
	if res[0].Parsed != nil {
		t.Fatalf("expected nil parsed for empty value")
	}
}

func TestCast_DoesNotMutateOriginal(t *testing.T) {
	orig := entries("NUM", "42")
	copy := orig[0]
	caster.Cast(orig)
	if orig[0].Key != copy.Key || orig[0].Value != copy.Value {
		t.Fatal("original entry was mutated")
	}
}

func TestCast_MultipleEntries(t *testing.T) {
	res := caster.Cast(entries("A", "hello", "B", "1", "C", "true", "D", ""))
	expected := []caster.Kind{caster.KindString, caster.KindInt, caster.KindBool, caster.KindEmpty}
	for i, r := range res {
		if r.Kind != expected[i] {
			t.Errorf("entry %d: expected %s, got %s", i, expected[i], r.Kind)
		}
	}
}
