package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/config"
	"github.com/yourorg/envoy-cli/internal/env"
)

// newPromoteCmd returns the parent 'promote' command with subcommands.
func newPromoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote environment variables between sets",
		Long: `Promote copies variables from a source set into a destination set.

Use 'promote all' to copy every variable, or 'promote keys' to copy a
specific subset of keys.`,
	}

	cmd.AddCommand(newPromoteAllCmd())
	cmd.AddCommand(newPromoteKeysCmd())
	return cmd
}

// newPromoteAllCmd promotes all variables from src to dst.
func newPromoteAllCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "all <src> <dst>",
		Short: "Promote all variables from one set to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName, dstName := args[0], args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			src := cfg.GetSet(srcName)
			if src == nil {
				return fmt.Errorf("source set %q not found", srcName)
			}

			dst := cfg.GetSet(dstName)
			if dst == nil {
				return fmt.Errorf("destination set %q not found", dstName)
			}

			promoted, skipped, err := env.Promote(src, dst, overwrite)
			if err != nil {
				return fmt.Errorf("promoting variables: %w", err)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Printf("Promoted %d variable(s) from %q to %q", promoted, srcName, dstName)
			if skipped > 0 {
				fmt.Printf(" (%d skipped, already present)", skipped)
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"Overwrite existing keys in the destination set")
	return cmd
}

// newPromoteKeysCmd promotes a specific list of keys from src to dst.
func newPromoteKeysCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "keys <src> <dst> <key> [key...]",
		Short: "Promote specific keys from one set to another",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName, dstName := args[0], args[1]
			keys := args[2:]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			src := cfg.GetSet(srcName)
			if src == nil {
				return fmt.Errorf("source set %q not found", srcName)
			}

			dst := cfg.GetSet(dstName)
			if dst == nil {
				return fmt.Errorf("destination set %q not found", dstName)
			}

			promoted, skipped, err := env.PromoteKeys(src, dst, keys, overwrite)
			if err != nil {
				return fmt.Errorf("promoting keys: %w", err)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Printf("Promoted %d key(s) from %q to %q", promoted, srcName, dstName)
			if skipped > 0 {
				fmt.Printf(" (%d skipped, already present)", skipped)
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"Overwrite existing keys in the destination set")
	return cmd
}
