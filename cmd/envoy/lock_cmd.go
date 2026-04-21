package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
)

func newLockCmd(cfgDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage locks on environment sets",
	}
	cmd.AddCommand(newLockAcquireCmd(cfgDir))
	cmd.AddCommand(newLockReleaseCmd(cfgDir))
	cmd.AddCommand(newLockListCmd(cfgDir))
	return cmd
}

func newLockAcquireCmd(cfgDir string) *cobra.Command {
	var owner string
	cmd := &cobra.Command{
		Use:   "acquire <set>",
		Short: "Acquire a lock on an environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			if owner == "" {
				owner = os.Getenv("USER")
			}
			if owner == "" {
				owner = "unknown"
			}
			if err := config.AddLock(cfgDir, setName, owner); err != nil {
				return fmt.Errorf("acquire lock: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "locked %q as %q\n", setName, owner)
			return nil
		},
	}
	cmd.Flags().StringVar(&owner, "owner", "", "lock owner (defaults to $USER)")
	return cmd
}

func newLockReleaseCmd(cfgDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "release <set>",
		Short: "Release a lock on an environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			if err := config.RemoveLock(cfgDir, setName); err != nil {
				return fmt.Errorf("release lock: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "unlocked %q\n", setName)
			return nil
		},
	}
}

func newLockListCmd(cfgDir string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all active locks",
		RunE: func(cmd *cobra.Command, args []string) error {
			records, err := config.LoadLocks(cfgDir)
			if err != nil {
				return fmt.Errorf("load locks: %w", err)
			}
			if len(records) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no active locks")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SET\tOWNER\tLOCKED AT")
			for name, r := range records {
				fmt.Fprintf(w, "%s\t%s\t%s\n", name, r.LockedBy, r.LockedAt.Format(time.RFC3339))
			}
			return w.Flush()
		},
	}
}
