package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/envoy-cli/internal/differ"
	"github.com/envoy-cli/envoy-cli/internal/loader"
	"github.com/envoy-cli/envoy-cli/internal/planner"
	"github.com/spf13/cobra"
)

var planFormat string

func init() {
	planCmd := &cobra.Command{
		Use:   "plan <base> <target>",
		Short: "Show an execution plan derived from diffing two .env files",
		Args:  cobra.ExactArgs(2),
		RunE:  runPlan,
	}
	planCmd.Flags().StringVarP(&planFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(planCmd)
}

func runPlan(cmd *cobra.Command, args []string) error {
	base, target, err := loader.LoadPair(args[0], args[1])
	if err != nil {
		return fmt.Errorf("loading files: %w", err)
	}

	result := differ.Compare(base, target)
	plan := planner.Build(result)

	switch planFormat {
	case "json":
		return printPlanJSON(plan)
	default:
		return printPlanText(plan)
	}
}

func printPlanText(plan planner.Plan) error {
	if !plan.HasChanges() {
		fmt.Println("No changes. Environment is up to date.")
		return nil
	}
	for _, a := range plan.Actions {
		if a.Kind == planner.ActionNoop {
			continue
		}
		fmt.Println(a.String())
	}
	return nil
}

func printPlanJSON(plan planner.Plan) error {
	type jsonAction struct {
		Kind     string `json:"kind"`
		Key      string `json:"key"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}
	var out []jsonAction
	for _, a := range plan.Actions {
		if a.Kind == planner.ActionNoop {
			continue
		}
		out = append(out, jsonAction{
			Kind:     string(a.Kind),
			Key:      a.Key,
			OldValue: a.OldValue,
			NewValue: a.NewValue,
		})
	}
	if out == nil {
		out = []jsonAction{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
