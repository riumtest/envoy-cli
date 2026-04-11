package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/internal/aliaser"
	"github.com/envoy-cli/internal/exporter"
	"github.com/envoy-cli/internal/loader"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias <file>",
	Short: "Create key aliases within a .env file",
	Long: `Alias copies the value of an existing key into one or more new keys.

Example:
  envoy-cli alias .env --rule DB_HOST=DATABASE_HOST --rule DB_HOST=HOST`,
	Args: cobra.ExactArgs(1),
	RunE: runAlias,
}

func init() {
	aliasCmd.Flags().StringArrayP("rule", "r", nil, "Alias rule in FROM=TO format (repeatable)")
	aliasCmd.Flags().Bool("overwrite", false, "Overwrite existing keys on conflict")
	aliasCmd.Flags().Bool("skip-missing", false, "Skip rules whose source key is missing")
	aliasCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv, export, json")
	_ = aliasCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(aliasCmd)
}

func runAlias(cmd *cobra.Command, args []string) error {
	path := args[0]
	entries, err := loader.Load(path)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	rawRules, _ := cmd.Flags().GetStringArray("rule")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	skipMissing, _ := cmd.Flags().GetBool("skip-missing")
	format, _ := cmd.Flags().GetString("format")

	var rules []aliaser.Rule
	for _, r := range rawRules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("alias: invalid rule %q, expected FROM=TO", r)
		}
		rules = append(rules, aliaser.Rule{From: parts[0], To: parts[1]})
	}

	opts := aliaser.DefaultOptions()
	opts.Overwrite = overwrite
	opts.SkipMissing = skipMissing

	result, err := aliaser.Alias(entries, rules, opts)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	if len(result.Conflicts) > 0 {
		for _, c := range result.Conflicts {
			fmt.Fprintf(os.Stderr, "warning: alias conflict skipped: %s -> %s\n", c.From, c.To)
		}
	}

	if format == "json" {
		type row struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		rows := make([]row, len(result.Entries))
		for i, e := range result.Entries {
			rows[i] = row{Key: e.Key, Value: e.Value}
		}
		return json.NewEncoder(os.Stdout).Encode(rows)
	}

	return exporter.Write(os.Stdout, result.Entries, exporter.Options{
		Format: format,
	})
}
