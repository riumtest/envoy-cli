package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/loader"
	"github.com/yourusername/envoy-cli/internal/redactor"
)

var (
	redactPlaceholder string
	redactPatterns    []string
	redactFormat      string
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact <file>",
		Short: "Print env file with sensitive values masked",
		Args:  cobra.ExactArgs(1),
		RunE:  runRedact,
	}
	redactCmd.Flags().StringVar(&redactPlaceholder, "placeholder", "***", "Replacement string for sensitive values")
	redactCmd.Flags().StringSliceVar(&redactPatterns, "extra-patterns", nil, "Additional key patterns to treat as sensitive")
	redactCmd.Flags().StringVar(&redactFormat, "format", "dotenv", "Output format: dotenv or json")
	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	opts := redactor.Options{
		Placeholder:   redactPlaceholder,
		ExtraPatterns: redactPatterns,
	}
	redacted := redactor.Redact(entries, opts)

	switch strings.ToLower(redactFormat) {
	case "json":
		m := make(map[string]string, len(redacted))
		for _, e := range redacted {
			m[e.Key] = e.Value
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m)
	default:
		for _, e := range redacted {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
