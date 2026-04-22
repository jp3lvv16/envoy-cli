package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/env"
)

func newScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Resolve environment variables across prioritised scopes",
	}
	cmd.AddCommand(newScopeResolveCmd())
	cmd.AddCommand(newScopeResolveAllCmd())
	return cmd
}

// newScopeResolveCmd resolves a single key across scopes supplied as
// name=priority pairs, e.g.: --scope dev=1 --scope prod=10
func newScopeResolveCmd() *cobra.Command {
	var scopeFlags []string
	cmd := &cobra.Command{
		Use:   "resolve KEY",
		Short: "Print the winning value for KEY across the given scopes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			scopes, err := buildScopes(cfg, scopeFlags)
			if err != nil {
				return err
			}
			r, err := env.NewScopeResolver(scopes)
			if err != nil {
				return err
			}
			val, winner, err := r.Resolve(key)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s  (from scope: %s)\n", val, winner)
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&scopeFlags, "scope", nil, "scope in name=priority format (repeatable)")
	_ = cmd.MarkFlagRequired("scope")
	return cmd
}

// newScopeResolveAllCmd merges all scopes into a single flat set and prints it.
func newScopeResolveAllCmd() *cobra.Command {
	var scopeFlags []string
	var outName string
	cmd := &cobra.Command{
		Use:   "resolve-all",
		Short: "Merge all scopes into a single resolved set and display it",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			scopes, err := buildScopes(cfg, scopeFlags)
			if err != nil {
				return err
			}
			r, err := env.NewScopeResolver(scopes)
			if err != nil {
				return err
			}
			out, err := r.ResolveAll(outName)
			if err != nil {
				return err
			}
			keys, err := env.SortedKeys(out, true)
			if err != nil {
				return err
			}
			for _, k := range keys {
				v, _ := out.Get(k)
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&scopeFlags, "scope", nil, "scope in name=priority format (repeatable)")
	cmd.Flags().StringVar(&outName, "name", "resolved", "name for the merged result set")
	_ = cmd.MarkFlagRequired("scope")
	return cmd
}

// buildScopes converts name=priority flag strings into []*env.Scope by loading
// each named set from the config store.
func buildScopes(cfg *config.Config, flags []string) ([]*env.Scope, error) {
	var scopes []*env.Scope
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("scope flag %q must be in name=priority format", f)
		}
		name := parts[0]
		prio, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("scope %q: invalid priority %q", name, parts[1])
		}
		s, err := cfg.GetSet(name)
		if err != nil {
			return nil, fmt.Errorf("scope %q: %w", name, err)
		}
		scopes = append(scopes, &env.Scope{Name: name, Priority: prio, Set: s})
	}
	return scopes, nil
}
