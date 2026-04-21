package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
)

// newAuditCmd returns the root 'audit' command with subcommands.
func newAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "View audit logs for environment sets",
		Long:  "Display recorded audit entries showing who changed what and when.",
	}

	cmd.AddCommand(newAuditListCmd())
	cmd.AddCommand(newAuditSummaryCmd())
	cmd.AddCommand(newAuditClearCmd())

	return cmd
}

// newAuditListCmd returns a command that lists all audit entries for a set.
func newAuditListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list <set>",
		Short: "List audit log entries for an environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			log, err := config.LoadAudit()
			if err != nil {
				return fmt.Errorf("loading audit log: %w", err)
			}

			entries := config.AuditFor(log, setName)
			if entries == nil || len(entries) == 0 {
				fmt.Printf("No audit entries found for set %q.\n", setName)
				return nil
			}

			// Apply limit from the end of the slice.
			if limit > 0 && limit < len(entries) {
				entries = entries[len(entries)-limit:]
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tACTOR\tACTION\tDETAIL")
			for _, e := range entries {
				t := time.Unix(e.Timestamp, 0).Format(time.RFC3339)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t, e.Actor, e.Action, e.Detail)
			}
			return w.Flush()
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 0, "Show only the last N entries (0 = all)")
	return cmd
}

// newAuditSummaryCmd returns a command that prints a summary of audit activity.
func newAuditSummaryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "summary <set>",
		Short: "Print a summary of audit activity for an environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			log, err := config.LoadAudit()
			if err != nil {
				return fmt.Errorf("loading audit log: %w", err)
			}

			entries := config.AuditFor(log, setName)
			if entries == nil || len(entries) == 0 {
				fmt.Printf("No audit entries found for set %q.\n", setName)
				return nil
			}

			// Count actions.
			counts := make(map[string]int)
			for _, e := range entries {
				counts[e.Action]++
			}

			fmt.Printf("Audit summary for set %q (%d entries):\n", setName, len(entries))
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  ACTION\tCOUNT")
			for action, count := range counts {
				fmt.Fprintf(w, "  %s\t%d\n", action, count)
			}
			return w.Flush()
		},
	}
}

// newAuditClearCmd returns a command that removes all audit entries for a set.
func newAuditClearCmd() *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "clear <set>",
		Short: "Clear all audit log entries for an environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			if !confirm {
				return fmt.Errorf("pass --confirm to delete audit entries for %q", setName)
			}

			log, err := config.LoadAudit()
			if err != nil {
				return fmt.Errorf("loading audit log: %w", err)
			}

			// Remove entries for the named set.
			delete(log.Entries, setName)

			if err := config.SaveAudit(log); err != nil {
				return fmt.Errorf("saving audit log: %w", err)
			}

			fmt.Printf("Audit entries for set %q cleared.\n", setName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion of audit entries")
	return cmd
}
