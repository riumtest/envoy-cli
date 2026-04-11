package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/loader"
	"github.com/user/envoy-cli/internal/splitter"
)

var splitCmd = &cobra.Command{
	Use:   "split <file>",
	Short: "Split an env file into groups by key prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runSplit,
}

func init() {
	splitCmd.Flags().StringSliceP("prefix", "p", []string{}, "Prefixes to split by (comma-separated)")
	splitCmd.Flags().Bool("strip", false, "Strip the matched prefix from keys in output")
	splitCmd.Flags().Bool("include-unmatched", true, "Include unmatched entries under empty group")
	splitCmd.Flags().StringP("format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(splitCmd)
}

func runSplit(cmd *cobra.Command, args []string) error {
	prefixes, _ := cmd.Flags().GetStringSlice("prefix")
	strip, _ := cmd.Flags().GetBool("strip")
	includeUnmatched, _ := cmd.Flags().GetBool("include-unmatched")
	format, _ := cmd.Flags().GetString("format")

	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	opts := splitter.DefaultOptions()
	opts.Prefixes = prefixes
	opts.StripPrefix = strip
	opts.IncludeUnmatched = includeUnmatched

	groups := splitter.Split(entries, opts)

	if format == "json" {
		out := map[string]map[string]string{}
		for group, es := range groups {
			label := group
			if label == "" {
				label = "(unmatched)"
			}
			out[label] = map[string]string{}
			for _, e := range es {
				out[label][e.Key] = e.Value
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	for _, prefix := range append(prefixes, "") {
		es, ok := groups[prefix]
		if !ok {
			continue
		}
		label := prefix
		if label == "" {
			label = "(unmatched)"
		}
		fmt.Printf("[%s]\n", strings.TrimSuffix(label, "_"))
		for _, e := range es {
			fmt.Printf("  %s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
