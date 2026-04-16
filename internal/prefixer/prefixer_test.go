package prefixer_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/prefixer"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "NAME", Value: "mydb"},
	}
}

func TestAdd_AppendsPrefix(t *testing.T) {
	opts := prefixer.DefaultOptions()
	result := prefixer.Add(entries(), "DB", opts)

	expected := []string{"DB_HOST", "DB_PORT", "DB_NAME"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("[%d] got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestAdd_PreservesValues(t *testing.T) {
	opts := prefixer.DefaultOptions()
	result := prefixer.Add(entries(), "DB", opts)

	if result[0].Value != "localhost" {
		t.Errorf("value mutated: got %q", result[0].Value)
	}
}

func TestAdd_EmptyPrefixIsNoop(t *testing.T) {
	opts := prefixer.DefaultOptions()
	src := entries()
	result := prefixer.Add(src, "", opts)

	for i, e := range result {
		if e.Key != src[i].Key {
			t.Errorf("key changed with empty prefix: %q", e.Key)
		}
	}
}

func TestAdd_DoesNotMutateOriginal(t *testing.T) {
	opts := prefixer.DefaultOptions()
	src := entries()
	origKey := src[0].Key
	prefixer.Add(src, "APP", opts)

	if src[0].Key != origKey {
		t.Error("original slice was mutated")
	}
}

func TestAdd_PreservesLength(t *testing.T) {
	opts := prefixer.DefaultOptions()
	src := entries()
	result := prefixer.Add(src, "DB", opts)

	if len(result) != len(src) {
		t.Errorf("length changed: got %d, want %d", len(result), len(src))
	}
}

func TestRemove_StripsPrefix(t *testing.T) {
	opts := prefixer.DefaultOptions()
	withPrefix := prefixer.Add(entries(), "DB", opts)
	result := prefixer.Remove(withPrefix, "DB", opts)

	expected := []string{"HOST", "PORT", "NAME"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("[%d] got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestRemove_LeavesUnmatchedKeysAlone(t *testing.T) {
	opts := prefixer.DefaultOptions()
	src := []envfile.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	result := prefixer.Remove(src, "DB", opts)

	if result[0].Key != "HOST" {
		t.Errorf("expected DB_HOST stripped to HOST, got %q", result[0].Key)
	}
	if result[1].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME unchanged, got %q", result[1].Key)
	}
}

func TestRemove_EmptyPrefixIsNoop(t *testing.T) {
	opts := prefixer.DefaultOptions()
	src := entries()
	result := prefixer.Remove(src, "", opts)

	for i, e := range result {
		if e.Key != src[i].Key {
			t.Errorf("key changed with empty prefix: %q", e.Key)
		}
	}
}
