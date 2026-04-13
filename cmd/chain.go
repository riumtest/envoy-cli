package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/chainer"
	"github.com/envoy-cli/envoy/internal/envfile"
	"github.com/envoy-cli/envoy/internal/loader"
	"github.com/envoy-cli/envoy/internal/normalizer"
	"github.com/envoy-cli/envoy/internal/sanitizer"
	"github.com/envoy-cli/envoy/internal/trimmer"
	"github.com/spf13/cobra"
)

var chainCmd = &cobra.Command{
	Use:   "chain <file>",
	Short: "Apply a named sequence of built-in transforms to an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runChain,
}

var chainSteps []string
var chainFormat string

func init() {
	chainCmd.Flags().StringSliceVarP(&chainSteps, "steps", "s", []string{"sanitize", "trim", "normalize"}, "Ordered list of steps to apply (sanitize, trim, normalize)")
	chainCmd.Flags().StringVarP(&chainFormat, "format", "f", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(chainCmd)
}

func runChain(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	stepRegistry := map[string]chainer.StepFn{
		"sanitize": func(es []envfile.Entry) ([]envfile.Entry, error) {
			return sanitizer.Sanitize(es, sanitizer.DefaultOptions()), nil
		},
		"trim": func(es []envfile.Entry) ([]envfile.Entry, error) {
			return trimmer.Trim(es, trimmer.DefaultOptions()), nil
		},
		"normalize": func(es []envfile.Entry) ([]envfile.Entry, error) {
			return normalizer.Normalize(es, normalizer.DefaultOptions()), nil
		},
	}

	var steps []chainer.Step
	for _, name := range chainSteps {
		fn, ok := stepRegistry[strings.ToLower(name)]
		if !ok {
			return fmt.Errorf("unknown step %q; available: sanitize, trim, normalize", name)
		}
		steps = append(steps, chainer.Step{Name: name, Fn: fn})
	}

	out, results, err := chainer.Chain(entries, steps)
	if err != nil {
		return err
	}

	if chainFormat == "json" {
		type jsonResult struct {
			Step    string `json:"step"`
			Count   int    `json:"output_count"`
		}
		type output struct {
			Steps   []jsonResult       `json:"steps"`
			Entries []envfile.Entry    `json:"entries"`
		}
		var jr []jsonResult
		for _, r := range results {
			jr = append(jr, jsonResult{Step: r.Step, Count: len(r.Entries)})
		}
		return json.NewEncoder(os.Stdout).Encode(output{Steps: jr, Entries: out})
	}

	for _, e := range out {
		fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
