// Package planner builds an ordered execution plan from a diff result.
package planner

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/differ"
)

// ActionKind describes what action should be taken.
type ActionKind string

const (
	ActionSet    ActionKind = "set"
	ActionDelete ActionKind = "delete"
	ActionUpdate ActionKind = "update"
	ActionNoop   ActionKind = "noop"
)

// Action is a single planned operation.
type Action struct {
	Kind     ActionKind
	Key      string
	OldValue string
	NewValue string
}

// String returns a human-readable description of the action.
func (a Action) String() string {
	switch a.Kind {
	case ActionSet:
		return fmt.Sprintf("SET %s=%s", a.Key, a.NewValue)
	case ActionDelete:
		return fmt.Sprintf("DELETE %s", a.Key)
	case ActionUpdate:
		return fmt.Sprintf("UPDATE %s: %s -> %s", a.Key, a.OldValue, a.NewValue)
	default:
		return fmt.Sprintf("NOOP %s", a.Key)
	}
}

// Plan is an ordered list of actions.
type Plan struct {
	Actions []Action
}

// HasChanges reports whether the plan contains any non-noop actions.
func (p *Plan) HasChanges() bool {
	for _, a := range p.Actions {
		if a.Kind != ActionNoop {
			return true
		}
	}
	return false
}

// Build constructs a Plan from a differ.Result.
func Build(result differ.Result) Plan {
	actions := make([]Action, 0, len(result.Changes))
	for _, c := range result.Changes {
		var a Action
		switch c.Kind {
		case differ.Added:
			a = Action{Kind: ActionSet, Key: c.Key, NewValue: c.NewValue}
		case differ.Removed:
			a = Action{Kind: ActionDelete, Key: c.Key, OldValue: c.OldValue}
		case differ.Changed:
			a = Action{Kind: ActionUpdate, Key: c.Key, OldValue: c.OldValue, NewValue: c.NewValue}
		default:
			a = Action{Kind: ActionNoop, Key: c.Key}
		}
		actions = append(actions, a)
	}
	return Plan{Actions: actions}
}
