package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/auditor"
	"github.com/yourusername/envoy-cli/internal/loader"
	"github.com/yourusername/envoy-cli/internal/masker"
)

var auditCmd = &cobra.Command{
	Use:   "audit <file>",
	Short: "Audit an .env file for issues and suspicious patterns",
	Args:  cobra.ExactArgs(1),
	RunE:  runAudit,
}

var auditFormat string

func init() {
	auditCmd.Flags().StringVarP(&auditFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	path := args[0]
	if !loader.Exists(path) {
		return fmt.Errorf("file not found: %s", path)
	}

	parsed, err := loader.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", path, err)
	}

	entries := make([]auditor.Entry, len(parsed))
	for i, e := range parsed {
		entries[i] = auditor.Entry{Key: e.Key, Value: e.Value}
	}

	m := masker.New()
	report := auditor.Audit(entries, m)

	if auditFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	if report.Total == 0 {
		fmt.Println("✔ No issues found.")
		return nil
	}

	for _, f := range report.Findings {
		var icon string
		switch f.Severity {
		case auditor.SeverityError:
			icon = "✖"
		case auditor.SeverityWarning:
			icon = "⚠"
		default:
			icon = "ℹ"
		}
		fmt.Printf("%s [%s] %s: %s\n", icon, f.Severity, f.Key, f.Message)
	}
	fmt.Printf("\nTotal: %d  Errors: %d  Warnings: %d  Infos: %d\n",
		report.Total, report.Errors, report.Warnings, report.Infos)
	return nil
}
