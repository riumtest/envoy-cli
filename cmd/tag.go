package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/loader"
	"github.com/envoy-cli/envoy/internal/tagger"
	"github.com/spf13/cobra"
)

var (
	tagRules  []string
	tagFilter string
	tagFormat string
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag <file>",
		Short: "Tag .env entries by key prefix rules",
		Args:  cobra.ExactArgs(1),
		RunE:  runTag,
	}
	tagCmd.Flags().StringArrayVarP(&tagRules, "rule", "r", nil, "Tag rules in tag=PREFIX format (e.g. db=DB_)")
	tagCmd.Flags().StringVarP(&tagFilter, "filter", "f", "", "Only output entries with this tag")
	tagCmd.Flags().StringVar(&tagFormat, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	opts := tagger.DefaultOptions()
	opts.Rules = make(map[string]string)
	for _, rule := range tagRules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected tag=PREFIX", rule)
		}
		opts.Rules[parts[0]] = parts[1]
	}

	tagged := tagger.Tag(entries, opts)
	if tagFilter != "" {
		tagged = tagger.FilterByTag(tagged, tagFilter)
	}

	if tagFormat == "json" {
		return json.NewEncoder(os.Stdout).Encode(tagged)
	}

	for _, t := range tagged {
		label := "(none)"
		if len(t.Tags) > 0 {
			label = strings.Join(t.Tags, ",")
		}
		fmt.Fprintf(os.Stdout, "%-30s = %-30s [%s]\n", t.Key, t.Value, label)
	}
	return nil
}
