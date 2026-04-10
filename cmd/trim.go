package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/exporter"
	"github.com/user/envoy-cli/internal/loader"
	"github.com/user/envoy-cli/internal/trimmer"
)

var (
	trimRemoveEmpty  bool
	trimDeduplicate  bool
	trimNoTrimValues bool
	trimOutput       string
	trimFormat       string
)

var trimCmd = &cobra.Command{
	Use:   "trim <file>",
	Short: "Trim whitespace, empty values, and duplicate keys from a .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrim,
}

func init() {
	trimCmd.Flags().BoolVar(&trimRemoveEmpty, "remove-empty", false, "Remove entries with empty values")
	trimCmd.Flags().BoolVar(&trimDeduplicate, "deduplicate", false, "Keep only the last occurrence of duplicate keys")
	trimCmd.Flags().BoolVar(&trimNoTrimValues, "no-trim-values", false, "Disable trimming of whitespace from values")
	trimCmd.Flags().StringVarP(&trimOutput, "output", "o", "", "Write result to file instead of stdout")
	trimCmd.Flags().StringVar(&trimFormat, "format", "dotenv", "Output format: dotenv, export, json")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	opts := trimmer.Options{
		TrimValues:      !trimNoTrimValues,
		RemoveEmpty:     trimRemoveEmpty,
		DeduplicateKeys: trimDeduplicate,
	}
	trimmed := trimmer.Trim(entries, opts)

	w := os.Stdout
	if trimOutput != "" {
		f, err := os.Create(trimOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if trimFormat == "json" {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(trimmed)
	}

	return exporter.Write(w, trimmed, nil, exporter.Options{
		Format: trimFormat,
	})
}
