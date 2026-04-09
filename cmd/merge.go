package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/envfile"
	"github.com/user/envoy-cli/internal/exporter"
	"github.com/user/envoy-cli/internal/merger"
)

var (
	mergeStrategy string
	mergeOutput   string
	mergeFormat   string
)

var mergeCmd = &cobra.Command{
	Use:   "merge <file1> <file2> [fileN...]",
	Short: "Merge multiple .env files into one",
	Long: `Merge two or more .env files into a single output.

Conflict resolution strategies:
  first  – keep the value from the first file that defines the key (default)
  last   – overwrite with the value from the last file
  error  – abort if any key is defined more than once`,
	Args: cobra.MinimumNArgs(2),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "first",
		"Conflict strategy: first | last | error")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "",
		"Write merged output to file (default: stdout)")
	mergeCmd.Flags().StringVarP(&mergeFormat, "format", "f", "dotenv",
		"Output format: dotenv | export | json")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	files := make(map[string][]envfile.Entry, len(args))
	for _, path := range args {
		entries, err := envfile.Parse(path)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		files[path] = entries
	}

	var strategy merger.Strategy
	switch strings.ToLower(mergeStrategy) {
	case "first":
		strategy = merger.StrategyFirst
	case "last":
		strategy = merger.StrategyLast
	case "error":
		strategy = merger.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q; choose first, last, or error", mergeStrategy)
	}

	result, err := merger.Merge(files, args, strategy)
	if err != nil {
		return err
	}

	if len(result.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "warning: %d conflict(s) resolved using strategy %q\n",
			len(result.Conflicts), mergeStrategy)
		for _, c := range result.Conflicts {
			fmt.Fprintf(os.Stderr, "  key %q: %s\n", c.Key, strings.Join(c.Sources, " vs "))
		}
	}

	out := os.Stdout
	if mergeOutput != "" {
		f, err := os.Create(mergeOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	return exporter.Write(out, result.Entries, mergeFormat, false, nil)
}
