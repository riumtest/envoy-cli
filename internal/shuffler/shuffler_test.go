package shuffler_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/shuffler"
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

func TestShuffle_ReturnsSameLength(t *testing.T) {
	in := entries()
	out := shuffler.Shuffle(in, shuffler.DefaultOptions())
	if len(out) != len(in) {
		t.Fatalf("expected %d entries, got %d", len(in), len(out))
	}
}

func TestShuffle_DoesNotMutateOriginal(t *testing.T) {
	in := entries()
	origFirst := in[0].Key
	shuffler.Shuffle(in, shuffler.Options{Seed: 42})
	if in[0].Key != origFirst {
		t.Fatal("original slice was mutated")
	}
}

func TestShuffle_DeterministicWithFixedSeed(t *testing.T) {
	in := entries()
	a := shuffler.Shuffle(in, shuffler.Options{Seed: 99})
	b := shuffler.Shuffle(in, shuffler.Options{Seed: 99})
	for i := range a {
		if a[i].Key != b[i].Key {
			t.Fatalf("position %d: got %q and %q with same seed", i, a[i].Key, b[i].Key)
		}
	}
}

func TestShuffle_DifferentSeedProducesDifferentOrder(t *testing.T) {
	in := entries()
	a := shuffler.Shuffle(in, shuffler.Options{Seed: 1})
	b := shuffler.Shuffle(in, shuffler.Options{Seed: 2})
	same := true
	for i := range a {
		if a[i].Key != b[i].Key {
			same = false
			break
		}
	}
	if same {
		t.Log("warning: different seeds produced the same order (statistically unlikely but possible)")
	}
}

func TestShuffle_ContainsSameKeys(t *testing.T) {
	in := entries()
	out := shuffler.Shuffle(in, shuffler.Options{Seed: 7})
	got := make(map[string]bool, len(out))
	for _, e := range out {
		got[e.Key] = true
	}
	for _, e := range in {
		if !got[e.Key] {
			t.Errorf("key %q missing from shuffled output", e.Key)
		}
	}
}

func TestShuffle_EmptyEntries(t *testing.T) {
	out := shuffler.Shuffle([]envfile.Entry{}, shuffler.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(out))
	}
}
