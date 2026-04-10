package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/internal/loader"
	"github.com/envoy-cli/internal/renamer"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <file>",
	Short: "Rename one or more keys in an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runRename,
}

var renameRules []string
var renameOutput string
var renameFormat string

func init() {
	renameCmd.Flags().StringArrayVarP(&renameRules, "rule", "r", nil, "Rename rule in FROM=TO format (repeatable)")
	renameCmd.Flags().StringVarP(&renameOutput, "output", "o", "", "Write result to file (default: stdout)")
	renameCmd.Flags().StringVar(&renameFormat, "format", "dotenv", "Output format: dotenv|json")
	_ = renameCmd.MarkFlagRequired("rule")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	var rules []renamer.Rule
	for _, raw := range renameRules {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected FROM=TO", raw)
		}
		rules = append(rules, renamer.Rule{From: parts[0], To: parts[1]})
	}

	updated, res, err := renamer.Rename(entries, rules)
	if err != nil {
		return err
	}

	for _, s := range res.Skipped {
		fmt.Fprintf(os.Stderr, "warn: key %q not found, skipping\n", s.From)
	}
	for _, c := range res.Conflict {
		fmt.Fprintf(os.Stderr, "warn: target key %q already exists, skipping rename from %q\n", c.To, c.From)
	}

	w := os.Stdout
	if renameOutput != "" {
		f, err := os.Create(renameOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if renameFormat == "json" {
		m := make(map[string]string, len(updated))
		for _, e := range updated {
			m[e.Key] = e.Value
		}
		return json.NewEncoder(w).Encode(m)
	}

	for _, e := range updated {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
