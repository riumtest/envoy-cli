package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envoy-cli/internal/grouper"
	"github.com/envoy-cli/internal/loader"
	"github.com/spf13/cobra"
)

var (
	groupDelimiter string
	groupMinSize   int
	groupFormat    string
)

var groupCmd = &cobra.Command{
	Use:   "group <file>",
	Short: "Group env keys by common prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runGroup,
}

func init() {
	groupCmd.Flags().StringVarP(&groupDelimiter, "delimiter", "d", "_", "Key delimiter for prefix detection")
	groupCmd.Flags().IntVarP(&groupMinSize, "min-size", "m", 1, "Minimum entries per group")
	groupCmd.Flags().StringVarP(&groupFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	opts := grouper.Options{
		Delimiter:  groupDelimiter,
		MinSize:    groupMinSize,
		SortGroups: true,
	}

	groups := grouper.GroupByPrefix(entries, opts)

	switch groupFormat {
	case "json":
		type jsonGroup struct {
			Name    string            `json:"name"`
			Count   int               `json:"count"`
			Entries map[string]string `json:"entries"`
		}
		out := make([]jsonGroup, 0, len(groups))
		for _, g := range groups {
			em := make(map[string]string, len(g.Entries))
			for _, e := range g.Entries {
				em[e.Key] = e.Value
			}
			out = append(out, jsonGroup{Name: g.Name, Count: len(g.Entries), Entries: em})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	default:
		if len(groups) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No groups found.")
			return nil
		}
		for _, g := range groups {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] (%d keys)\n", g.Name, len(g.Entries))
			for _, e := range g.Entries {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", e.Key, e.Value)
			}
		}
	}
	return nil
}
