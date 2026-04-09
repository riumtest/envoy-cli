package sorter_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/sorter"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "ZEBRA", Value: "last"},
		{Key: "ALPHA", Value: "first"},
		{Key: "MANGO", Value: "middle"},
	}
}

func TestSort_AscendingByKey(t *testing.T) {
	result := sorter.Sort(entries(), sorter.DefaultOptions())
	keys := []string{result[0].Key, result[1].Key, result[2].Key}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSort_DescendingByKey(t *testing.T) {
	opts := sorter.Options{Order: sorter.Descending, ByValue: false}
	result := sorter.Sort(entries(), opts)
	if result[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %q", result[0].Key)
	}
}

func TestSort_ByValue(t *testing.T) {
	opts := sorter.Options{Order: sorter.Ascending, ByValue: true}
	result := sorter.Sort(entries(), opts)
	if result[0].Value != "first" {
		t.Errorf("expected 'first' first, got %q", result[0].Value)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	orig := entries()
	firstKey := orig[0].Key
	sorter.SortByKey(orig)
	if orig[0].Key != firstKey {
		t.Error("original slice was mutated")
	}
}

func TestSortByKey_Convenience(t *testing.T) {
	result := sorter.SortByKey(entries())
	if result[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA, got %q", result[0].Key)
	}
}

func TestSort_EmptySlice(t *testing.T) {
	result := sorter.SortByKey([]envfile.Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty slice, got len %d", len(result))
	}
}
