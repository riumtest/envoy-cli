package differ

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	source := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	target := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result := Compare(source, target)

	if result.HasChanges() {
		t.Errorf("Expected no changes, got %d differences", len(result.Differences))
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	source := map[string]string{
		"KEY1": "value1",
	}
	target := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result := Compare(source, target)

	if len(result.Differences) != 1 {
		t.Fatalf("Expected 1 difference, got %d", len(result.Differences))
	}

	diff := result.Differences[0]
	if diff.Type != DiffTypeAdded || diff.Key != "KEY2" || diff.NewValue != "value2" {
		t.Errorf("Expected added KEY2=value2, got %+v", diff)
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	source := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	target := map[string]string{
		"KEY1": "value1",
	}

	result := Compare(source, target)

	if len(result.Differences) != 1 {
		t.Fatalf("Expected 1 difference, got %d", len(result.Differences))
	}

	diff := result.Differences[0]
	if diff.Type != DiffTypeRemoved || diff.Key != "KEY2" || diff.OldValue != "value2" {
		t.Errorf("Expected removed KEY2, got %+v", diff)
	}
}

func TestCompare_ChangedKeys(t *testing.T) {
	source := map[string]string{
		"KEY1": "old_value",
	}
	target := map[string]string{
		"KEY1": "new_value",
	}

	result := Compare(source, target)

	if len(result.Differences) != 1 {
		t.Fatalf("Expected 1 difference, got %d", len(result.Differences))
	}

	diff := result.Differences[0]
	if diff.Type != DiffTypeChanged || diff.OldValue != "old_value" || diff.NewValue != "new_value" {
		t.Errorf("Expected changed KEY1, got %+v", diff)
	}
}

func TestSummary(t *testing.T) {
	result := &Result{
		Differences: []Difference{
			{Key: "KEY1", Type: DiffTypeAdded},
			{Key: "KEY2", Type: DiffTypeRemoved},
			{Key: "KEY3", Type: DiffTypeChanged},
		},
	}

	summary := result.Summary()
	expected := "1 added, 1 removed, 1 changed"

	if summary != expected {
		t.Errorf("Expected summary '%s', got '%s'", expected, summary)
	}
}
