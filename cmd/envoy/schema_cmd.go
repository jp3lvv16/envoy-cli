package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/env"
)

func newSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate an env set against a JSON schema file",
	}
	cmd.AddCommand(newSchemaValidateCmd())
	return cmd
}

func newSchemaValidateCmd() *cobra.Command {
	var schemaFile string

	cmd := &cobra.Command{
		Use:   "validate <set-name>",
		Short: "Validate a named env set against a schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			vars, ok := cfg.GetSet(setName)
			if !ok {
				return fmt.Errorf("set %q not found", setName)
			}

			s, err := env.NewSet(setName)
			if err != nil {
				return err
			}
			for k, v := range vars {
				if putErr := s.Put(k, v); putErr != nil {
					return putErr
				}
			}

			schema, err := loadSchemaFile(schemaFile)
			if err != nil {
				return fmt.Errorf("load schema: %w", err)
			}

			if valErr := env.ValidateSchema(s, schema); valErr != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "validation failed: %v\n", valErr)
				os.Exit(1)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "set %q is valid against schema %q\n", setName, schema.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&schemaFile, "file", "f", "", "path to JSON schema file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

// loadSchemaFile reads a JSON file with the structure:
//
//	{"name": "...", "fields": [{"key": "...", "required": true, "pattern": "..."}]}
func loadSchemaFile(path string) (*env.Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schema env.Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parse schema JSON: %w", err)
	}
	if schema.Name == "" {
		return nil, fmt.Errorf("schema must have a non-empty name")
	}
	return &schema, nil
}
