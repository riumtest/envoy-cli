package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/extractor"
	"github.com/user/envoy-cli/internal/loader"
)

func init() {
	extractCmd := &cobra.Command{
		Use:   "extract <file>",
		Short: "Extract a subset of keys from an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runExtract,
	}

	extractCmd.Flags().StringP("key-pattern", "k", "", "Regex pattern to match keys")
	extractCmd.Flags().StringP("value-pattern", "v", "", "Regex pattern to match values")
	extractCmd.Flags().StringSliceP("keys", "K", nil, "Explicit list of keys to extract")
	extractCmd.Flags().Bool("case-sensitive", false, "Use case-sensitive regex matching")
	extractCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv, export, json")

	rootCmd.AddCommand(extractCmd)
}

func runExtract(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	entries, err := loader.Load(filePath)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	keyPattern, _ := cmd.Flags().GetString("key-pattern")
	valuePattern, _ := cmd.Flags().GetString("value-pattern")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
	format, _ := cmd.Flags().GetString("format")

	res, err := extractor.Extract(entries, extractor.Options{
		KeyPattern:    keyPattern,
		ValuePattern:  valuePattern,
		Keys:          keys,
		CaseSensitive: caseSensitive,
	})
	if err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	switch strings.ToLower(format) {
	case "json":
		m := make(map[string]string, len(res.Entries))
		for _, e := range res.Entries {
			m[e.Key] = e.Value
		}
		return json.NewEncoder(os.Stdout).Encode(m)
	case "export":
		for _, e := range res.Entries {
			fmt.Fprintf(os.Stdout, "export %s=%s\n", e.Key, e.Value)
		}
	default:
		for _, e := range res.Entries {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
	}

	fmt.Fprintf(os.Stderr, "extracted %d keys, skipped %d\n", res.Extracted, res.Skipped)
	return nil
}
