package chainer_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/envoy-cli/envoy/internal/chainer"
	"github.com/envoy-cli/envoy/internal/envfile"
)

func entries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "debug", Value: "true"},
		{Key: "PORT", Value: "  8080  "},
	}
}

func TestChain_NoSteps(t *testing.T) {
	in := entries()
	out, results, err := chainer.Chain(in, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if len(out) != len(in) {
		t.Errorf("expected %d entries, got %d", len(in), len(out))
	}
}

func TestChain_SingleStep(t *testing.T) {
	upperStep := chainer.Step{
		Name: "uppercase",
		Fn: func(es []envfile.Entry) ([]envfile.Entry, error) {
			out := make([]envfile.Entry, len(es))
			for i, e := range es {
				e.Key = strings.ToUpper(e.Key)
				out[i] = e
			}
			return out, nil
		},
	}

	out, results, err := chainer.Chain(entries(), []chainer.Step{upperStep})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	for _, e := range out {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("key %q was not uppercased", e.Key)
		}
	}
}

func TestChain_MultipleSteps_Ordered(t *testing.T) {
	var order []string

	mkStep := func(name string) chainer.Step {
		return chainer.Step{
			Name: name,
			Fn: func(es []envfile.Entry) ([]envfile.Entry, error) {
				order = append(order, name)
				return es, nil
			},
		}
	}

	steps := []chainer.Step{mkStep("a"), mkStep("b"), mkStep("c")}
	_, _, err := chainer.Chain(entries(), steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Join(order, ",") != "a,b,c" {
		t.Errorf("unexpected execution order: %v", order)
	}
}

func TestChain_StopsOnError(t *testing.T) {
	ran := map[string]bool{}
	expectedErr := errors.New("boom")

	steps := []chainer.Step{
		{Name: "first", Fn: func(es []envfile.Entry) ([]envfile.Entry, error) { ran["first"] = true; return es, nil }},
		{Name: "fail", Fn: func(es []envfile.Entry) ([]envfile.Entry, error) { return nil, expectedErr }},
		{Name: "third", Fn: func(es []envfile.Entry) ([]envfile.Entry, error) { ran["third"] = true; return es, nil }},
	}

	_, results, err := chainer.Chain(entries(), steps)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("unexpected error: %v", err)
	}
	if !ran["first"] {
		t.Error("expected first step to run")
	}
	if ran["third"] {
		t.Error("third step should not have run after failure")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results (including failing step), got %d", len(results))
	}
}

func TestChain_NilStepFnReturnsError(t *testing.T) {
	steps := []chainer.Step{{Name: "broken", Fn: nil}}
	_, _, err := chainer.Chain(entries(), steps)
	if err == nil {
		t.Fatal("expected error for nil step function")
	}
}
