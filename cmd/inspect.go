package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/inspector"
/internal/masker"
 inspectFormat string

func init() {
	inspectCmd := &cobra.Command{
		Use:   "inspect <file>",
		Short: "Inspect an .env file and display a summary report",
		Args:  cobra.ExactArgs(1),
		RunE:  runInspect,
	}
	inspectCmd.Flags().StringVarP(&inspectFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	m := masker.New()
	report := inspector.Inspect(entries, m)

	switch strings.ToLower(inspectFormat) {
	case "json":
		return printInspectJSON(cmd, report)
	default:
		printInspectText(cmd, report)
		return nil
	}
}

func printInspectText(cmd *cobra.Command, r inspector.Report) {
	fmt.Fprintf(cmd.OutOrStdout(), "Total keys    : %d\n", r.TotalKeys)
	fmt.Fprintf(cmd.OutOrStdout(), "Empty values  : %d\n", r.EmptyValues)
	fmt.Fprintf(cmd.OutOrStdout(), "Sensitive keys: %s\n", joinOrNone(r.SensitiveKeys))
	fmt.Fprintf(cmd.OutOrStdout(), "URL keys      : %s\n", joinOrNone(r.URLKeys))
	fmt.Fprintf(cmd.OutOrStdout(), "Boolean keys  : %s\n", joinOrNone(r.BooleanKeys))
	fmt.Fprintf(cmd.OutOrStdout(), "Numeric keys  : %s\n", joinOrNone(r.NumericKeys))
}

func printInspectJSON(cmd *cobra.Command, r inspector.Report) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func joinOrNone(ss []string) string {
	if len(ss) == 0 {
		return "(none)"
	}
	return strings.Join(ss, ", ")
}
