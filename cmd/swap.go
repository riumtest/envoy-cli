package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/envoy-cli/internal/loader"
	"github.com/envoy-cli/envoy-cli/internal/swapper"
	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap <file>",
	Short: "Swap keys and values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSwap,
}

var swapKeepEmpty bool
var swapAllowDuplicates bool
var swapFormat string

func init() {
	swapCmd.Flags().BoolVar(&swapKeepEmpty, "keep-empty", false, "Include entries with empty values (key becomes empty)")
	swapCmd.Flags().BoolVar(&swapAllowDuplicates, "allow-duplicates", false, "Keep duplicate swapped keys")
	swapCmd.Flags().StringVar(&swapFormat, "format", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(swapCmd)
}

func runSwap(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	opts := swapper.Options{
		SkipEmpty:      !swapKeepEmpty,
		SkipDuplicates: !swapAllowDuplicates,
	}
	result := swapper.Swap(entries, opts)

	switch swapFormat {
	case "json":
		m := make(map[string]string, len(result))
		for _, e := range result {
			m[e.Key] = e.Value
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m)
	default:
		for _, e := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
