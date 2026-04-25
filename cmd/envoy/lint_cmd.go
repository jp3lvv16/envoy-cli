package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/env"
)

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint environment variable sets for style and correctness issues",
	}
	cmd.AddCommand(newLintRunCmd())
	return cmd
}

func newLintRunCmd() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "run <set-name>",
		Short: "Run lint checks on a named environment set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			envSet, ok := cfg.GetSet(setName)
			if !ok {
				return fmt.Errorf("set %q not found", setName)
			}

			results, err := env.Lint(envSet)
			if err != nil {
				return fmt.Errorf("lint: %w", err)
			}

			fmt.Println(env.FormatLint(results))

			if strict && len(results) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Exit with non-zero status if any lint issues are found")
	return cmd
}
