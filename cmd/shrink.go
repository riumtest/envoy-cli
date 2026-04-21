package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/loader"
	"github.com/your-org/envoy-cli/internal/shrinker"
)

var shrinkCmd = &cobra.Command{
	Use:   "shrink <file>",
	Short: "Remove entries whose values fall outside a length range",
	Args:  cobra.ExactArgs(1),
	RunE:  runShrink,
}

func init() {
	shrinkCmd.Flags().Int("min", 0, "Minimum value length (inclusive); shorter values are dropped")
	shrinkCmd.Flags().Int("max", 0, "Maximum value length (inclusive); longer values are dropped (0 = unlimited)")
	shrinkCmd.Flags().StringSlice("keep", nil, "Keys to always retain regardless of value length")
	shrinkCmd.Flags().String("format", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(shrinkCmd)
}

func runShrink(cmd *cobra.Command, args []string) error {
	minLen, _ := cmd.Flags().GetInt("min")
	maxLen, _ := cmd.Flags().GetInt("max")
	keepKeys, _ := cmd.Flags().GetStringSlice("keep")
	format, _ := cmd.Flags().GetString("format")

	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	opts := shrinker.DefaultOptions()
	opts.MinLen = minLen
	opts.MaxLen = maxLen
	opts.KeepKeys = keepKeys

	result := shrinker.Shrink(entries, opts)

	switch format {
	case "json":
		type row struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		rows := make([]row, len(result))
		for i, e := range result {
			rows[i] = row{Key: e.Key, Value: e.Value}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	default:
		for _, e := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
