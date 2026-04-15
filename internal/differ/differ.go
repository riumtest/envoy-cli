package differ

import "github.com/envoy-cli/envoy-cli/internal/envfile"

// ChangeKind represents the type of change between two env files.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff output.
type Result struct {
	Changes []Change
}

// Summary returns a short human-readable summary of the diff.
func (r *Result) Summary() string {
	added, removed, changed := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	if added == 0 && removed == 0 && changed == 0 {
		return "no changes"
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed", added, removed, changed)
}

// Compare diffs two slices of env entries and returns a Result.
func Compare(base, target []envfile.Entry) Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var changes []Change

	for _, e := range base {
		if nv, ok := targetMap[e.Key]; ok {
			if nv != e.Value {
				changes = append(changes, Change{Key: e.Key, Kind: Changed, OldValue: e.Value, NewValue: nv})
			} else {
				changes = append(changes, Change{Key: e.Key, Kind: Unchanged, OldValue: e.Value, NewValue: nv})
			}
		} else {
			changes = append(changes, Change{Key: e.Key, Kind: Removed, OldValue: e.Value})
		}
	}

	for _, e := range target {
		if _, ok := baseMap[e.Key]; !ok {
			changes = append(changes, Change{Key: e.Key, Kind: Added, NewValue: e.Value})
		}
	}

	return Result{Changes: changes}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
