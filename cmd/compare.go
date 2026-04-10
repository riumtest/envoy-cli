package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/comparator"
	"github.com/user/envoy-cli/internal/loader"
)

var compareFormat string

var compareCmd = &cobra.Command{
	Use:   "compare <file1> <file2>",
	Short: "Compare two .env files and report overlap and mismatches",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().StringVarP(&compareFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	left, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading %s: %w", args[0], err)
	}
	right, err := loader.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading %s: %w", args[1], err)
	}

	result := comparator.Compare(left, right)

	if compareFormat == "json" {
		return printCompareJSON(result)
	}
	return printCompareText(result, args[0], args[1])
}

func printCompareText(r comparator.Result, leftFile, rightFile string) error {
	fmt.Printf("Comparing %s ↔ %s\n", leftFile, rightFile)
	fmt.Printf("Overlap ratio: %.1f%%\n\n", comparator.OverlapRatio(r)*100)

	if len(r.SharedKeys) > 0 {
		sort.Strings(r.SharedKeys)
		fmt.Printf("✔ Shared (%d):\n", len(r.SharedKeys))
		for _, k := range r.SharedKeys {
			fmt.Printf("  %s\n", k)
		}
	}
	if len(r.MismatchedKeys) > 0 {
		fmt.Printf("~ Mismatched (%d):\n", len(r.MismatchedKeys))
		for _, m := range r.MismatchedKeys {
			fmt.Printf("  %s: %q → %q\n", m.Key, m.LeftValue, m.RightValue)
		}
	}
	if len(r.OnlyInLeft) > 0 {
		sort.Strings(r.OnlyInLeft)
		fmt.Printf("← Only in %s (%d):\n", leftFile, len(r.OnlyInLeft))
		for _, k := range r.OnlyInLeft {
			fmt.Printf("  %s\n", k)
		}
	}
	if len(r.OnlyInRight) > 0 {
		sort.Strings(r.OnlyInRight)
		fmt.Printf("→ Only in %s (%d):\n", rightFile, len(r.OnlyInRight))
		for _, k := range r.OnlyInRight {
			fmt.Printf("  %s\n", k)
		}
	}
	return nil
}

func printCompareJSON(r comparator.Result) error {
	type output struct {
		OverlapRatio   float64                      `json:"overlap_ratio"`
		SharedKeys     []string                     `json:"shared_keys"`
		MismatchedKeys []comparator.MismatchEntry   `json:"mismatched_keys"`
		OnlyInLeft     []string                     `json:"only_in_left"`
		OnlyInRight    []string                     `json:"only_in_right"`
	}
	out := output{
		OverlapRatio:   comparator.OverlapRatio(r),
		SharedKeys:     r.SharedKeys,
		MismatchedKeys: r.MismatchedKeys,
		OnlyInLeft:     r.OnlyInLeft,
		OnlyInRight:    r.OnlyInRight,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
