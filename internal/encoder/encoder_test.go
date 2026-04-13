package encoder

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
)

func entries(pairs ...string) []envfile.Entry {
	out := make([]envfile.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestEncode_Base64(t *testing.T) {
	results, err := Encode(entries("SECRET", "hello"), FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Encoded != "aGVsbG8=" {
		t.Errorf("expected base64 'aGVsbG8=', got %q", results[0].Encoded)
	}
	if results[0].Original != "hello" {
		t.Errorf("original should be preserved, got %q", results[0].Original)
	}
}

func TestEncode_Hex(t *testing.T) {
	results, err := Encode(entries("TOKEN", "abc"), FormatHex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Encoded != "616263" {
		t.Errorf("expected hex '616263', got %q", results[0].Encoded)
	}
}

func TestEncode_Escape(t *testing.T) {
	results, err := Encode(entries("MSG", "say \"hi\"\nbye"), FormatEscape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `say \"hi\"\nbye`
	if results[0].Encoded != want {
		t.Errorf("expected %q, got %q", want, results[0].Encoded)
	}
}

func TestEncode_EmptyValue(t *testing.T) {
	results, err := Encode(entries("EMPTY", ""), FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Encoded != "" {
		t.Errorf("empty value should remain empty, got %q", results[0].Encoded)
	}
}

func TestEncode_UnsupportedFormat(t *testing.T) {
	_, err := Encode(entries("KEY", "val"), Format("unknown"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestEncode_MultipleEntries(t *testing.T) {
	input := entries("A", "foo", "B", "bar", "C", "baz")
	results, err := Encode(input, FormatBase64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Key != input[i].Key {
			t.Errorf("result[%d]: expected key %q, got %q", i, input[i].Key, r.Key)
		}
		if r.Original != input[i].Value {
			t.Errorf("result[%d]: expected original %q, got %q", i, input[i].Value, r.Original)
		}
	}
}

func TestToEntries(t *testing.T) {
	results := []Result{
		{Key: "A", Original: "foo", Encoded: "Zm9v", Format: FormatBase64},
		{Key: "B", Original: "", Encoded: "", Format: FormatBase64},
	}
	out := ToEntries(results)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "A" || out[0].Value != "Zm9v" {
		t.Errorf("unexpected entry: %+v", out[0])
	}
	if out[1].Value != "" {
		t.Errorf("expected empty value for B, got %q", out[1].Value)
	}
}
