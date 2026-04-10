// Package planner generates a migration plan from two sets of env entries,
// describing the steps needed to transform source into target.
package planner

import "github.com/envoy-cli/envoy-cli/internal/envfile"

// ActionType describes the kind of change in a plan step.
type ActionType string

const (
	ActionAdd    ActionType = "add"
	ActionRemove ActionType = "remove"
	ActionUpdate ActionType = "update"
	ActionKeep   ActionType = "keep"
)

// Step represents a single migration action.
type Step struct {
	Key      string
	Action   ActionType
	OldValue string
	NewValue string
}

// Plan holds the full set of migration steps.
type Plan struct {
	Steps []Step
}

// AddCount returns the number of add steps.
func (p *Plan) AddCount() int { return p.count(ActionAdd) }

// RemoveCount returns the number of remove steps.
func (p *Plan) RemoveCount() int { return p.count(ActionRemove) }

// UpdateCount returns the number of update steps.
func (p *Plan) UpdateCount() int { return p.count(ActionUpdate) }

// KeepCount returns the number of keep steps.
func (p *Plan) KeepCount() int { return p.count(ActionKeep) }

func (p *Plan) count(a ActionType) int {
	n := 0
	for _, s := range p.Steps {
		if s.Action == a {
			n++
		}
	}
	return n
}

// Build computes the migration plan from src entries to dst entries.
func Build(src, dst []envfile.Entry) Plan {
	srcMap := toMap(src)
	dstMap := toMap(dst)

	var steps []Step

	// Check existing src keys.
	for _, e := range src {
		if dstVal, ok := dstMap[e.Key]; ok {
			if dstVal == e.Value {
				steps = append(steps, Step{Key: e.Key, Action: ActionKeep, OldValue: e.Value, NewValue: e.Value})
			} else {
				steps = append(steps, Step{Key: e.Key, Action: ActionUpdate, OldValue: e.Value, NewValue: dstVal})
			}
		} else {
			steps = append(steps, Step{Key: e.Key, Action: ActionRemove, OldValue: e.Value})
		}
	}

	// Find keys only in dst.
	for _, e := range dst {
		if _, ok := srcMap[e.Key]; !ok {
			steps = append(steps, Step{Key: e.Key, Action: ActionAdd, NewValue: e.Value})
		}
	}

	return Plan{Steps: steps}
}

func toMap(entries []envfile.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
