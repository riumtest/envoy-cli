package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/loader"
	"envoy-cli/internal/templater"
)

var (
	templateEnvFile  string
	templateOutput   string
	templateStrict   bool
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template [template-file]",
		Short: "Render a template file using values from a .env file",
		Long: `Render a template file by substituting placeholders of the form {{KEY}}
with values loaded from a .env file.

Example:
  envoy-cli template app.conf.tmpl --env .env
  envoy-cli template nginx.conf.tmpl --env production.env --output nginx.conf
  envoy-cli template deploy.yaml.tmpl --env .env --strict`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplate,
	}

	templateCmd.Flags().StringVarP(&templateEnvFile, "env", "e", ".env", "Path to the .env file to use for substitution")
	templateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Write rendered output to file instead of stdout")
	templateCmd.Flags().BoolVar(&templateStrict, "strict", false, "Fail if any placeholder cannot be resolved")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	tmplFile := args[0]

	// Read the template file
	tmplBytes, err := os.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("reading template file %q: %w", tmplFile, err)
	}

	// Load the env file
	if !loader.Exists(templateEnvFile) {
		return fmt.Errorf("env file not found: %s", templateEnvFile)
	}

	entries, err := loader.Load(templateEnvFile)
	if err != nil {
		return fmt.Errorf("loading env file %q: %w", templateEnvFile, err)
	}

	// Build substitution map from entries
	subs := templater.BuildSubsFromEntries(entries)

	// Render the template
	result, err := templater.Render(string(tmplBytes), subs)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	// Enforce strict mode: fail on unresolved placeholders
	if templateStrict && len(result.Missing) > 0 {
		summary := templater.Summary(result)
		return fmt.Errorf(
			"strict mode: %d placeholder(s) could not be resolved: %s",
			len(result.Missing),
			strings.Join(summary.Missing, ", "),
		)
	}

	// Write output
	if templateOutput != "" {
		if err := os.WriteFile(templateOutput, []byte(result.Output), 0644); err != nil {
			return fmt.Errorf("writing output file %q: %w", templateOutput, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Rendered template written to %s\n", templateOutput)

		// Print summary of substitutions
		summary := templater.Summary(result)
		fmt.Fprintf(cmd.OutOrStdout(), "  Substituted : %d\n", summary.Substituted)
		if len(summary.Missing) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  Unresolved  : %s\n", strings.Join(summary.Missing, ", "))
		}
	} else {
		fmt.Fprint(cmd.OutOrStdout(), result.Output)
	}

	return nil
}
