package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/archiver"
	"github.com/user/envoy-cli/internal/loader"
)

func init() {
	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive and restore .env file snapshots",
	}

	saveCmd := &cobra.Command{
		Use:   "save [file]",
		Short: "Save a .env file to the archive",
		Args:  cobra.ExactArgs(1),
		RunE:  runArchiveSave,
	}
	saveCmd.Flags().String("label", "", "archive label (default: timestamp)")
	saveCmd.Flags().String("dir", ".envoy-archives", "archive directory")

	loadCmd := &cobra.Command{
		Use:   "load [label]",
		Short: "Load an archived snapshot by label",
		Args:  cobra.ExactArgs(1),
		RunE:  runArchiveLoad,
	}
	loadCmd.Flags().String("dir", ".envoy-archives", "archive directory")
	loadCmd.Flags().String("format", "dotenv", "output format: dotenv, json")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved archives",
		RunE:  runArchiveList,
	}
	listCmd.Flags().String("dir", ".envoy-archives", "archive directory")

	archiveCmd.AddCommand(saveCmd, loadCmd, listCmd)
	rootCmd.AddCommand(archiveCmd)
}

func runArchiveSave(cmd *cobra.Command, args []string) error {
	entries, err := loader.Load(args[0])
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	label, _ := cmd.Flags().GetString("label")
	dir, _ := cmd.Flags().GetString("dir")
	opts := archiver.Options{Dir: dir, MaxKeep: 20}
	path, err := archiver.Save(entries, label, opts)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "archived to %s\n", path)
	return nil
}

func runArchiveLoad(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	format, _ := cmd.Flags().GetString("format")
	opts := archiver.Options{Dir: dir, MaxKeep: 20}
	a, err := archiver.Load(args[0], opts)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if format == "json" {
		return json.NewEncoder(out).Encode(a.Entries)
	}
	for _, e := range a.Entries {
		fmt.Fprintf(out, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}

func runArchiveList(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	opts := archiver.Options{Dir: dir, MaxKeep: 20}
	labels, err := archiver.List(opts)
	if err != nil {
		return err
	}
	if len(labels) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no archives found")
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(labels, "\n"))
	_ = os.Stdout
	return nil
}
