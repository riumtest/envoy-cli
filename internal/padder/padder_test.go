package padder_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/padder"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "A", Value: "hi"},
		{Key: "B", Value: "hello"},
		{Key: "C", Value: "hey"},
	}
}

func TestPad_LeftAlign(t *testing.T) {
	opts := padder.DefaultOptions()
	out := padder.Pad(entries(), opts)
	for _, e := range out {
		if len(e.Value) != 5 {
			t.Errorf("expected width 5 for key %s, got %d", e.Key, len(e.Value))
		}
	}
	if out[0].Value != "hi   " {
		t.Errorf("unexpected value: %q", out[0].Value)
	}
}

func TestPad_RightAlign(t *testing.T) {
	opts := padder.DefaultOptions()
	opts.Align = padder.AlignRight
	out := padder.Pad(entries(), opts)
	if out[0].Value != "   hi" {
		t.Errorf("unexpected value: %q", out[0].Value)
	}
}

func TestPad_MinWidth(t *testing.T) {
	opts := padder.DefaultOptions()
	opts.MinWidth = 10
	out := padder.Pad(entries(), opts)
	for _, e := range out {
		if len(e.Value) != 10 {
			t.Errorf("expected width 10 for key %s, got %d", e.Key, len(e.Value))
		}
	}
}

func TestPad_MaxWidth(t *testing.T) {
	opts := padder.DefaultOptions()
	opts.MaxWidth = 3
	out := padder.Pad(entries(), opts)
	for _, e := range out {
		if len(e.Value) != 3 {
			t.Errorf("expected width 3 for key %s, got %d", e.Key, len(e.Value))
		}
	}
	if out[1].Value != "hel" {
		t.Errorf("expected truncated value, got %q", out[1].Value)
	}
}

func TestPad_DoesNotMutateOriginal(t *testing.T) {
	in := entries()
	origVal := in[0].Value
	padder.Pad(in, padder.DefaultOptions())
	if in[0].Value != origVal {
		t.Error("original entries were mutated")
	}
}

func TestPad_EmptyEntries(t *testing.T) {
	out := padder.Pad([]envfile.Entry{}, padder.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}
