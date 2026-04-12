package freezer_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/freezer"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestFreeze_CapturesEntries(t *testing.T) {
	e := entries("HOST", "localhost", "PORT", "8080")
	frozen := freezer.Freeze(e)
	if len(frozen) != 2 {
		t.Fatalf("expected 2 frozen entries, got %d", len(frozen))
	}
	if frozen[0].Key != "HOST" || frozen[0].Value != "localhost" {
		t.Errorf("unexpected frozen entry: %+v", frozen[0])
	}
}

func TestCheck_NoViolations(t *testing.T) {
	e := entries("HOST", "localhost", "PORT", "8080")
	frozen := freezer.Freeze(e)
	violations, err := freezer.Check(frozen, e, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_MutatedValue(t *testing.T) {
	e := entries("HOST", "localhost", "PORT", "8080")
	frozen := freezer.Freeze(e)
	current := entries("HOST", "remotehost", "PORT", "8080")
	violations, _ := freezer.Check(frozen, current, false)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Kind != "mutated" || violations[0].Key != "HOST" {
		t.Errorf("unexpected violation: %+v", violations[0])
	}
}

func TestCheck_DeletedKey(t *testing.T) {
	e := entries("HOST", "localhost", "PORT", "8080")
	frozen := freezer.Freeze(e)
	current := entries("HOST", "localhost")
	violations, _ := freezer.Check(frozen, current, false)
	if len(violations) != 1 || violations[0].Kind != "deleted" {
		t.Errorf("expected deleted violation, got %+v", violations)
	}
}

func TestCheck_AddedKey(t *testing.T) {
	e := entries("HOST", "localhost")
	frozen := freezer.Freeze(e)
	current := entries("HOST", "localhost", "PORT", "9000")
	violations, _ := freezer.Check(frozen, current, false)
	if len(violations) != 1 || violations[0].Kind != "added" {
		t.Errorf("expected added violation, got %+v", violations)
	}
}

func TestCheck_FailOnViolation_ReturnsError(t *testing.T) {
	e := entries("HOST", "localhost")
	frozen := freezer.Freeze(e)
	current := entries("HOST", "changed")
	_, err := freezer.Check(frozen, current, true)
	if err == nil {
		t.Error("expected error due to violation, got nil")
	}
}
