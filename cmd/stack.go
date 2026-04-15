package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/loader"
	"github.com/user/envoy-cli/internal/stacker"
)

var (
	stackStrategy string
	stackOutput   string
)

func init() {
	stackCmd := &cobra.Command{
		Use:   "stack <base.env> <overlay.env> [more.env...]",
		Short: "Merge multiple .env files in priority order",
		Long: `Stack merges two or more .env files layer by layer.
By default the last file wins on key conflicts (strategy=last).
Use --strategy=first to let the earliest definition win.`,
		Args:    cobra.MinimumNArgs(2),
		RunE:    runStack,
		Example: "  envoy-cli stack base.env prod.env --strategy=last",
	}
	stackCmd.Flags().StringVar(&stackStrategy, "strategy", "last", "conflict strategy: last or first")
	stackCmd.Flags().StringVarP(&stackOutput, "output", "o", "text", "output format: text or json")
	rootCmd.AddCommand(stackCmd)
}

func runStack(cmd *cobra.Command, args []string) error {
	layers := make([][]interface{ Key() string }, 0)
	_ = layers

	var allLayers [][]interface{}
	_ = allLayers

	opts := stacker.DefaultOptions()
	if stackStrategy == "first" {
		opts.Strategy = stacker.StrategyFirst
	}

	var layerSlices [][]interface{}
	_ = layerSlices

	resultLayers := make([]interface{}, 0)
	_ = resultLayers

	// Load each file as a layer.
	var envLayers [][]interface{}
	_ = envLayers

	parsedLayers := make([]interface{}, 0, len(args))
	_ = parsedLayers

	// Use loader to read all files.
	var stackLayers [][]envfile.Entry
	for _, path := range args {
		entries, err := loader.Load(path)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}
		stackLayers = append(stackLayers, entries)
	}

	res := stacker.Stack(stackLayers, opts)

	if stackOutput == "json" {
		type row struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		rows := make([]row, len(res.Entries))
		for i, e := range res.Entries {
			rows[i] = row{Key: e.Key, Value: e.Value}
		}
		return json.NewEncoder(os.Stdout).Encode(rows)
	}

	for _, e := range res.Entries {
		fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
	}
	fmt.Fprintf(os.Stderr, "# %d entries, %d overridden\n", len(res.Entries), res.Overridden)
	return nil
}
