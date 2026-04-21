package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/env"
)

// newDiffCmd returns the parent 'diff' command with its subcommands.
func newDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare environment variable sets",
		Long:  "Show differences between two environment variable sets or display a formatted comparison.",
	}

	cmd.AddCommand(newDiffSetsCmd())
	cmd.AddCommand(newDiffFormatCmd())

	return cmd
}

// newDiffSetsCmd returns the 'diff sets' subcommand that prints a unified diff
// between two named sets.
func newDiffSetsCmd() *cobra.Command {
	var colorize bool

	cmd := &cobra.Command{
		Use:   "sets <src> <dst>",
		Short: "Show a unified diff between two sets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName := args[0]
			dstName := args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			srcVars, ok := cfg.GetSet(srcName)
			if !ok {
				return fmt.Errorf("set %q not found", srcName)
			}

			dstVars, ok := cfg.GetSet(dstName)
			if !ok {
				return fmt.Errorf("set %q not found", dstName)
			}

			src, err := env.NewSet(srcName)
			if err != nil {
				return err
			}
			for k, v := range srcVars {
				if err := src.Put(k, v); err != nil {
					return err
				}
			}

			dst, err := env.NewSet(dstName)
			if err != nil {
				return err
			}
			for k, v := range dstVars {
				if err := dst.Put(k, v); err != nil {
					return err
				}
			}

			changes, err := env.Diff(src, dst)
			if err != nil {
				return fmt.Errorf("computing diff: %w", err)
			}

			if len(changes) == 0 {
				fmt.Println("Sets are identical — no differences found.")
				return nil
			}

			out := env.FormatDiff(changes, colorize)
			fmt.Fprint(os.Stdout, out)
			return nil
		},
	}

	cmd.Flags().BoolVar(&colorize, "color", false, "Colorize diff output (added=green, removed=red, changed=yellow)")
	return cmd
}

// newDiffFormatCmd returns the 'diff format' subcommand that prints a structured
// comparison table between two sets using the compare formatter.
func newDiffFormatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format <src> <dst>",
		Short: "Show a structured comparison between two sets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcName := args[0]
			dstName := args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			srcVars, ok := cfg.GetSet(srcName)
			if !ok {
				return fmt.Errorf("set %q not found", srcName)
			}

			dstVars, ok := cfg.GetSet(dstName)
			if !ok {
				return fmt.Errorf("set %q not found", dstName)
			}

			src, err := env.NewSet(srcName)
			if err != nil {
				return err
			}
			for k, v := range srcVars {
				if err := src.Put(k, v); err != nil {
					return err
				}
			}

			dst, err := env.NewSet(dstName)
			if err != nil {
				return err
			}
			for k, v := range dstVars {
				if err := dst.Put(k, v); err != nil {
					return err
				}
			}

			result, err := env.Compare(src, dst)
			if err != nil {
				return fmt.Errorf("comparing sets: %w", err)
			}

			out := env.FormatCompare(result)
			if out == "" {
				fmt.Println("Sets are identical — no differences found.")
				return nil
			}

			fmt.Fprint(os.Stdout, out)
			return nil
		},
	}

	return cmd
}
