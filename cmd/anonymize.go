package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/anonymizer"
	"github.com/user/envoy-cli/internal/loader"
)

var (
	anonPrefix        string
	anonSensitiveOnly bool
	anonExtraPatterns []string
	anonJSONOutput    bool
)

func init() {
	anonCmd := &cobra.Command{
		Use:   "anonymize <file>",
		Short: "Replace env values with deterministic opaque tokens",
		Args:  cobra.ExactArgs(1),
		RunE:  runAnonymize,
	}
	anonCmd.Flags().StringVar(&anonPrefix, "prefix", "anon", "token prefix")
	anonCmd.Flags().BoolVar(&anonSensitiveOnly, "sensitive-only", false, "only anonymize sensitive keys")
	anonCmd.Flags().StringSliceVar(&anonExtraPatterns, "pattern", nil, "extra sensitive key patterns")
	anonCmd.Flags().BoolVar(&anonJSONOutput, "json", false, "output as JSON")
	rootCmd.AddCommand(anonCmd)
}

func runAnonymize(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	opts := anonymizer.Options{
		Prefix:        anonPrefix,
		SensitiveOnly: anonSensitiveOnly,
		ExtraPatterns: anonExtraPatterns,
	}
	out := anonymizer.Anonymize(entries, opts)

	if anonJSONOutput {
		type kv struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var rows []kv
		for _, e := range out {
			rows = append(rows, kv{Key: e.Key, Value: e.Value})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}

	for _, e := range out {
		fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
