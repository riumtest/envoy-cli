package differ_test

import (
	"testing"

	"envoy-cli/internal/differ"
	"envoy-cli/internal/envfile"
)

func mkEntries(pairs ...string) []envfile.Entry {
	var entries []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, envfile.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCompare_NoChanges(t *testing.T) {
	base := mkEntries("HOST", "localhost", "PORT", "5432")
	head := mkEntries("HOST", "localhost", "PORT", "5432")
	res := differ.Compare(base, head)
	if res.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(res.Diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(res.Diffs))
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	base := mkEntries("HOST", "localhost")
	head := mkEntries("HOST", "localhost", "PORT", "5432")
	res := differ.Compare(base, head)
	if !res.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, d := range res.Diffs {
		if d.Key == "PORT" && d.Change == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected PORT to be marked as added")
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	base := mkEntries("HOST", "localhost", "PORT", "5432")
	head := mkEntries("HOST", "localhost")
	res := differ.Compare(base, head)
	if !res.HasChanges() {
		t.Fatal("expected changes")
	}
	for _, d := range res.Diffs {
		if d.Key == "PORT" && d.Change != differ.Removed {
			t.Errorf("expected PORT removed, got %s", d.Change)
		}
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	base := mkEntries("HOST", "localhost")
	head := mkEntries("HOST", "production.db")
	res := differ.Compare(base, head)
	if !res.HasChanges() {
		t.Fatal("expected changes")
	}
	for _, d := range res.Diffs {
		if d.Key == "HOST" {
			if d.Change != differ.Changed {
				t.Errorf("expected HOST changed, got %s", d.Change)
			}
			if d.OldValue != "localhost" || d.NewValue != "production.db" {
				t.Errorf("unexpected values: old=%s new=%s", d.OldValue, d.NewValue)
			}
		}
	}
}

func TestCompare_ResultSortedByKey(t *testing.T) {
	base := mkEntries("Z_KEY", "1", "A_KEY", "2")
	head := mkEntries("Z_KEY", "1", "A_KEY", "2")
	res := differ.Compare(base, head)
	if len(res.Diffs) < 2 {
		t.Fatal("expected at least 2 diffs")
	}
	if res.Diffs[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %s", res.Diffs[0].Key)
	}
	if res.Diffs[1].Key != "Z_KEY" {
		t.Errorf("expected Z_KEY second, got %s", res.Diffs[1].Key)
	}
}
