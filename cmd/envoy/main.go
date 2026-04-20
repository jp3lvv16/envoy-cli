// Package main is the entry point for the envoy-cli tool.
// It wires together the internal packages and exposes a
// Cobra-based command hierarchy for managing environment
// variable sets across deployment targets.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/internal/config"
	"github.com/envoy-cli/internal/env"
)

// configPath is the default location of the envoy config file.
const configPath = ".envoy.json"

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd builds and returns the top-level cobra command.
func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envoy",
		Short: "Manage environment variable sets across deployment targets",
		Long: `envoy-cli lets you create, edit, export, and import named
environment variable sets so you can switch contexts quickly
when deploying to different targets (dev, staging, prod, …).`,
	}

	root.AddCommand(
		newListCmd(),
		newSetCmd(),
		newGetCmd(),
		newDeleteCmd(),
		newExportCmd(),
		newImportCmd(),
	)

	return root
}

// newListCmd returns the "list" sub-command that prints all known set names.
func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all environment variable sets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			sets := cfg.Sets()
			if len(sets) == 0 {
				fmt.Println("No sets defined.")
				return nil
			}
			for _, name := range sets {
				fmt.Println(name)
			}
			return nil
		},
	}
}

// newSetCmd returns the "set" sub-command that writes a key=value pair into a
// named environment variable set, creating the set if it does not exist.
func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <KEY=VALUE>",
		Short: "Set a variable in a named set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName, pair := args[0], args[1]
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			s := cfg.GetSet(setName)
			if s == nil {
				s, err = env.NewSet(setName)
				if err != nil {
					return err
				}
				cfg.AddOrUpdateSet(s)
			}
			key, value, ok := splitPair(pair)
			if !ok {
				return fmt.Errorf("argument must be in KEY=VALUE format, got %q", pair)
			}
			if err := s.Put(key, value); err != nil {
				return err
			}
			return config.Save(configPath, cfg)
		},
	}
}

// newGetCmd returns the "get" sub-command that prints the value of a key inside
// a named set.
func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <name> <KEY>",
		Short: "Get a variable from a named set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName, key := args[0], args[1]
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			s := cfg.GetSet(setName)
			if s == nil {
				return fmt.Errorf("set %q not found", setName)
			}
			val, err := s.Get(key)
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}
}

// newDeleteCmd returns the "delete" sub-command that removes an entire named
// set from the config.
func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named environment variable set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			cfg.DeleteSet(setName)
			return config.Save(configPath, cfg)
		},
	}
}

// newExportCmd returns the "export" sub-command that serialises a named set to
// stdout in the requested format (dotenv | shell | json).
func newExportCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export a named set to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			s := cfg.GetSet(setName)
			if s == nil {
				return fmt.Errorf("set %q not found", setName)
			}
			out, err := env.Export(s, format)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, shell, json")
	return cmd
}

// newImportCmd returns the "import" sub-command that reads a file and populates
// (or creates) a named set from it.
func newImportCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "import <name> <file>",
		Short: "Import variables from a file into a named set",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName, filePath := args[0], args[1]
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			s := cfg.GetSet(setName)
			if s == nil {
				s, err = env.NewSet(setName)
				if err != nil {
					return err
				}
				cfg.AddOrUpdateSet(s)
			}
			if err := env.Import(s, string(data), format); err != nil {
				return err
			}
			return config.Save(configPath, cfg)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Input format: dotenv, json")
	return cmd
}

// splitPair splits a "KEY=VALUE" string into its components.
// The value may itself contain '=' characters.
func splitPair(pair string) (key, value string, ok bool) {
	for i, ch := range pair {
		if ch == '=' {
			return pair[:i], pair[i+1:], true
		}
	}
	return "", "", false
}
