package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/linter"
	"envoy-cli/internal/loader"
)

var lintFormat string

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Check a .env file for style and convention issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVarP(&lintFormat, "format", "f", "text", "Output format: text or json")
}

func runLint(cmd *cobra.Command, args []string) error {
	path := args[0]

	envEntries, err := loader.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	var entries []linter.Entry
	for i, e := range envEntries {
		entries = append(entries, linter.Entry{
			Line:  i + 1,
			Key:   e.Key,
			Value: e.Value,
		})
	}

	result := linter.Lint(entries)

	switch lintFormat {
	case "json":
		return printLintJSON(result)
	default:
		return printLintText(result)
	}
}

func printLintText(result linter.Result) error {
	if len(result.Issues) == 0 {
		fmt.Println("No issues found.")
		return nil
	}
	for _, iss := range result.Issues {
		fmt.Printf("[%s] line %d (%s): %s\n", iss.Severity, iss.Line, iss.Key, iss.Message)
	}
	if result.HasErrors() {
		os.Exit(1)
	}
	return nil
}

func printLintJSON(result linter.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result.Issues); err != nil {
		return fmt.Errorf("json encode error: %w", err)
	}
	if result.HasErrors() {
		os.Exit(1)
	}
	return nil
}
