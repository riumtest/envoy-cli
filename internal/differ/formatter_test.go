package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/envfile"
	"github.com/envoy-cli/envoy-cli/internal/masker"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	target := []envfile.Entry{{Key: "FOO", Value: "bar"}}
	result := differ.Compare(base, target)
	var buf bytes.Buffer
	f := differ.NewFormatter("text", nil)
	if err := f.Format(&buf, result); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}

func TestTextFormatter_WithChanges(t *testing.T) {
	base := []envfile.Entry{{Key: "FOO", Value: "old"}, {Key: "GONE", Value: "x"}}
	target := []envfile.Entry{{Key: "FOO", Value: "new"}, {Key: "BAR", Value: "added"}}
	result := differ.Compare(base, target)
	var buf bytes.Buffer
	f := differ.NewFormatter("text", nil)
	if err := f.Format(&buf, result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ BAR") {
		t.Error("expected added BAR")
	}
	if !strings.Contains(out, "- GONE") {
		t.Error("expected removed GONE")
	}
	if !strings.Contains(out, "~ FOO") {
		t.Error("expected changed FOO")
	}
}

func TestTextFormatter_WithMasking(t *testing.T) {
	base := []envfile.Entry{{Key: "SECRET_KEY", Value: "oldpassword"}}
	target := []envfile.Entry{{Key: "SECRET_KEY", Value: "newpassword"}}
	result := differ.Compare(base, target)
	m := masker.New()
	var buf bytes.Buffer
	f := differ.NewFormatter("text", m)
	if err := f.Format(&buf, result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "oldpassword") || strings.Contains(out, "newpassword") {
		t.Error("expected masked values")
	}
}

func TestMaskValue(t *testing.T) {
	m := masker.New()
	base := []envfile.Entry{{Key: "DB_PASSWORD", Value: "secret"}}
	target := []envfile.Entry{{Key: "DB_PASSWORD", Value: "newsecret"}}
	result := differ.Compare(base, target)
	var buf bytes.Buffer
	f := differ.NewFormatter("text", m)
	if err := f.Format(&buf, result); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "secret") {
		t.Error("sensitive value should be masked")
	}
}

func TestJSONFormatter(t *testing.T) {
	base := []envfile.Entry{{Key: "A", Value: "1"}}
	target := []envfile.Entry{{Key: "A", Value: "2"}, {Key: "B", Value: "3"}}
	result := differ.Compare(base, target)
	var buf bytes.Buffer
	f := differ.NewFormatter("json", nil)
	if err := f.Format(&buf, result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"key\"") {
		t.Error("expected JSON with key field")
	}
	if !strings.Contains(out, "changed") && !strings.Contains(out, "added") {
		t.Error("expected change kinds in JSON")
	}
}
