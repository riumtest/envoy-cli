package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/loader"
	"envoy-cli/internal/validator"
)

var (
	validateWarnOnly bool
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a .env file for common issues",
	Long: `Validate checks a .env file for issues such as duplicate keys,
empty keys, empty values, and invalid formatting.

Exits with a non-zero status if any errors are found (unless --warn-only is set).`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&validateWarnOnly, "warn-only", false, "Print issues but exit with code 0")
}

func runValidate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if !loader.Exists(filePath) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	entries, err := loader.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	issues := validator.Validate(entries)

	if len(issues) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✔ %s is valid — no issues found\n", filePath)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Issues found in %s:\n", filePath)
	for _, issue := range issues {
		severity := "ERROR"
		if issue.Level == validator.LevelWarn {
			severity = "WARN"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] line %d: %s\n", severity, issue.Line, issue.Message)
	}

	hasErrors := false
	for _, issue := range issues {
		if issue.Level == validator.LevelError {
			hasErrors = true
			break
		}
	}

	if hasErrors && !validateWarnOnly {
		os.Exit(1)
	}

	return nil
}
