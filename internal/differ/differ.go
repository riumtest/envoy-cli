package differ

import (
	"fmt"
	"sort"
)

// DiffType represents the type of difference between env files
type DiffType string

const (
	DiffTypeAdded   DiffType = "added"
	DiffTypeRemoved DiffType = "removed"
	DiffTypeChanged DiffType = "changed"
)

// Difference represents a single difference between two env files
type Difference struct {
	Key      string
	Type     DiffType
	OldValue string
	NewValue string
}

// Result contains all differences between two env files
type Result struct {
	Differences []Difference
}

// HasChanges returns true if there are any differences
func (r *Result) HasChanges() bool {
	return len(r.Differences) > 0
}

// Compare compares two environment variable maps and returns differences
func Compare(source, target map[string]string) *Result {
	result := &Result{
		Differences: []Difference{},
	}

	// Find added and changed keys
	for key, targetValue := range target {
		sourceValue, exists := source[key]
		if !exists {
			result.Differences = append(result.Differences, Difference{
				Key:      key,
				Type:     DiffTypeAdded,
				NewValue: targetValue,
			})
		} else if sourceValue != targetValue {
			result.Differences = append(result.Differences, Difference{
				Key:      key,
				Type:     DiffTypeChanged,
				OldValue: sourceValue,
				NewValue: targetValue,
			})
		}
	}

	// Find removed keys
	for key, sourceValue := range source {
		if _, exists := target[key]; !exists {
			result.Differences = append(result.Differences, Difference{
				Key:      key,
				Type:     DiffTypeRemoved,
				OldValue: sourceValue,
			})
		}
	}

	// Sort differences by key for consistent output
	sort.Slice(result.Differences, func(i, j int) bool {
		return result.Differences[i].Key < result.Differences[j].Key
	})

	return result
}

// Summary returns a human-readable summary of the differences
func (r *Result) Summary() string {
	if !r.HasChanges() {
		return "No differences found"
	}

	added, removed, changed := 0, 0, 0
	for _, diff := range r.Differences {
		switch diff.Type {
		case DiffTypeAdded:
			added++
		case DiffTypeRemoved:
			removed++
		case DiffTypeChanged:
			changed++
		}
	}

	return fmt.Sprintf("%d added, %d removed, %d changed", added, removed, changed)
}
