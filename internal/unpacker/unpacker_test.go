package unpacker_test

import (
	"testing"

	"envoy-cli/internal/envfile"
	"envoy-cli/internal/unpacker"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestUnpack_Base64(t *testing.T) {
	// "hello" base64 → "aGVsbG8="
	in := entries("GREETING", "aGVsbG8=")
	out, err := unpacker.Unpack(in, unpacker.Options{Format: unpacker.FormatBase64})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "hello" {
		t.Errorf("expected hello, got %q", out[0].Value)
	}
}

func TestUnpack_Hex(t *testing.T) {
	// "world" hex → "776f726c64"
	in := entries("WORD", "776f726c64")
	out, err := unpacker.Unpack(in, unpacker.Options{Format: unpacker.FormatHex})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "world" {
		t.Errorf("expected world, got %q", out[0].Value)
	}
}

func TestUnpack_Escape(t *testing.T) {
	in := entries("MSG", `line1\nline2`)
	out, err := unpacker.Unpack(in, unpacker.Options{Format: unpacker.FormatEscape})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "line1\nline2" {
		t.Errorf("unexpected value: %q", out[0].Value)
	}
}

func TestUnpack_SpecificKeys(t *testing.T) {
	in := entries("A", "aGVsbG8=", "B", "aGVsbG8=")
	out, err := unpacker.Unpack(in, unpacker.Options{
		Format: unpacker.FormatBase64,
		Keys:   []string{"A"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "hello" {
		t.Errorf("A should be decoded, got %q", out[0].Value)
	}
	if out[1].Value != "aGVsbG8=" {
		t.Errorf("B should be unchanged, got %q", out[1].Value)
	}
}

func TestUnpack_InvalidBase64(t *testing.T) {
	in := entries("BAD", "!!!notbase64!!!")
	_, err := unpacker.Unpack(in, unpacker.Options{Format: unpacker.FormatBase64})
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestUnpack_UnknownFormat(t *testing.T) {
	in := entries("K", "v")
	_, err := unpacker.Unpack(in, unpacker.Options{Format: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestUnpack_DoesNotMutateOriginal(t *testing.T) {
	in := entries("X", "aGVsbG8=")
	orig := in[0].Value
	_, _ = unpacker.Unpack(in, unpacker.DefaultOptions())
	if in[0].Value != orig {
		t.Error("original slice was mutated")
	}
}
