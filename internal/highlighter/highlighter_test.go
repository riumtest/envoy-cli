package highlighter_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/highlighter"
	"github.com/user/envoy-cli/internal/masker"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "DEBUG", Value: ""},
	}
}

func TestHighlight_NoColor(t *testing.T) {
	opts := highlighter.DefaultOptions()
	opts.NoColor = true

	lines := highlighter.Highlight(entries(), opts, nil)

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "APP_NAME=envoy" {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestHighlight_WithColor_ContainsANSI(t *testing.T) {
	opts := highlighter.DefaultOptions()

	lines := highlighter.Highlight(entries(), opts, nil)

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// ANSI reset code should be present
	if !strings.Contains(lines[0], "\033[") {
		t.Errorf("expected ANSI codes in output, got: %s", lines[0])
	}
}

func TestHighlight_SensitiveKeyColored(t *testing.T) {
	m := masker.New()
	opts := highlighter.DefaultOptions()

	lines := highlighter.Highlight(entries(), opts, m)

	// DB_PASSWORD is sensitive — its line should use the sensitive colour (yellow \033[33m)
	if !strings.Contains(lines[1], "\033[33m") {
		t.Errorf("expected sensitive colour for DB_PASSWORD, got: %s", lines[1])
	}
}

func TestHighlight_MaskSensitive(t *testing.T) {
	m := masker.New()
	opts := highlighter.DefaultOptions()
	opts.MaskSensitive = true

	lines := highlighter.Highlight(entries(), opts, m)

	if strings.Contains(lines[1], "s3cr3t") {
		t.Errorf("expected value to be masked, got: %s", lines[1])
	}
}

func TestHighlight_EmptyValue_RedColored(t *testing.T) {
	opts := highlighter.DefaultOptions()

	lines := highlighter.Highlight(entries(), opts, nil)

	// DEBUG has empty value — should use red \033[31m
	if !strings.Contains(lines[2], "\033[31m") {
		t.Errorf("expected red colour for empty value, got: %s", lines[2])
	}
}

func TestHighlight_SkipsEmptyKey(t *testing.T) {
	e := []envfile.Entry{{Key: "", Value: "orphan"}, {Key: "X", Value: "1"}}
	opts := highlighter.DefaultOptions()
	opts.NoColor = true

	lines := highlighter.Highlight(e, opts, nil)

	if len(lines) != 1 {
		t.Errorf("expected 1 line (empty key skipped), got %d", len(lines))
	}
}

func TestHighlightComment_NoColor(t *testing.T) {
	out := highlighter.HighlightComment("# hello", true)
	if out != "# hello" {
		t.Errorf("unexpected: %s", out)
	}
}

func TestHighlightComment_WithColor(t *testing.T) {
	out := highlighter.HighlightComment("# hello", false)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes, got: %s", out)
	}
}
