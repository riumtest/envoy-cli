package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "A lightweight CLI for managing and diffing .env files",
	Long: `envoy-cli helps you compare .env files across environments
with support for secret masking and multiple output formats.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Addfmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/differ"
	"envoy-cli/internal/envfile"
	"envoy-cli/internal/masker"
)

var (
	maskSecrets bool
	outputFormat string
	extraPatterns []string
)

var diffCmd = &cobra.Command{
	Use:   "diff <base> <target>",
	Short: "Diff two .env files and display changes",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().BoolVarP(&maskSecrets, "mask", "m", false, "Mask sensitive values in output")
	diffCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text or json")
	diffCmd.Flags().StringArrayVarP(&extraPatterns, "pattern", "p", nil, "Additional key patterns to treat as sensitive")
}

func runDiff(cmd *cobra.Command, args []string) error {
	baseFile, targetFile := args[0], args[1]

	baseEnv, err := envfile.Parse(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file %q: %w", baseFile, err)
	}

	targetEnv, err := envfile.Parse(targetFile)
	if err != nil {
		return fmt.Errorf("reading target file %q: %w", targetFile, err)
	}

	var m *masker.Masker
	if maskSecrets {
		if len(extraPatterns) > 0 {
			m = masker.NewWithPatterns(extraPatterns)
		} else {
			m = masker.New()
		}
	}

	changes := differ.Compare(baseEnv, targetEnv)

	var formatter differ.Formatter
	switch outputFormat {
	case "json":
		formatter = &differ.JSONFormatter{Masker: m}
	default:
		formatter = &differ.TextFormatter{Masker: m}
	}

	output, err := formatter.Format(changes)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	fmt.Fprint(os.Stdout, output)
	return nil
}
