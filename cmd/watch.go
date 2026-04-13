package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envoy-cli/envoy/internal/watchdog"
	"github.com/spf13/cobra"
)

var (
	watchSnapshotPath string
	watchInterval     int
	watchMaxDrift     int
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <file>",
		Short: "Watch an .env file for drift against a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().StringVarP(&watchSnapshotPath, "snapshot", "s", "", "Path to snapshot file (required)")
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 5, "Poll interval in seconds")
	watchCmd.Flags().IntVar(&watchMaxDrift, "max-drift", 0, "Alert when drift exceeds this many keys")
	_ = watchCmd.MarkFlagRequired("snapshot")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	opts := watchdog.Options{
		PollInterval: time.Duration(watchInterval) * time.Second,
		MaxDrift:     watchMaxDrift,
	}

	alerts := make(chan watchdog.Alert, 16)
	stop, err := watchdog.Watch(filePath, watchSnapshotPath, opts, alerts)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}
	defer stop()

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (interval: %ds, max-drift: %d)...\n", filePath, watchInterval, watchMaxDrift)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			fmt.Fprintln(cmd.OutOrStdout(), "\nStopping watchdog.")
			return nil
		case alert := <-alerts:
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] Drift detected in %s (%d change(s)):\n",
				alert.At.Format(time.RFC3339), alert.File, len(alert.Changes))
			for _, c := range alert.Changes {
				switch c.Type {
				case "added":
					fmt.Fprintf(cmd.OutOrStdout(), "  + %s=%s\n", c.Key, c.NewValue)
				case "removed":
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", c.Key)
				case "changed":
					fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
				}
			}
		}
	}
}
