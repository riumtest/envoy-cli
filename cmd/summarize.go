package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/"
	"github.)

var string

varse:   "<file>",
	Short: "Display a summary report of an .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSummarize,
}

func init() {
	summarizeCmd.Flags().StringVarP(&summarizeFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(summarizeCmd)
}

func runSummarize(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	report := summarizer.Summarize(entries)

	switch strings.ToLower(summarizeFormat) {
	case "json":
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	default:
		fmt.Fprintf(cmd.OutOrStdout(), "Total keys:      %d\n", report.TotalKeys)
		fmt.Fprintf(cmd.OutOrStdout(), "Unique keys:     %d\n", report.UniqueKeys)
		fmt.Fprintf(cmd.OutOrStdout(), "Empty values:    %d\n", report.EmptyValues)
		fmt.Fprintf(cmd.OutOrStdout(), "Numeric values:  %d\n", report.NumericValues)
		fmt.Fprintf(cmd.OutOrStdout(), "Boolean values:  %d\n", report.BooleanValues)
		fmt.Fprintf(cmd.OutOrStdout(), "URL values:      %d\n", report.URLValues)
		if len(report.SensitiveKeys) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Sensitive keys:  none\n")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Sensitive keys:  %s\n", strings.Join(report.SensitiveKeys, ", "))
		}
	}
	return nil
}
