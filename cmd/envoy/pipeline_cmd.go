package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/env"
)

func newPipelineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Run a transformation pipeline on an environment set",
	}
	cmd.AddCommand(newPipelineRunCmd())
	return cmd
}

func newPipelineRunCmd() *cobra.Command {
	var steps []string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "run <set-name>",
		Short: "Apply an ordered list of transforms to a set and print the result",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			setName := args[0]

			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			src, err := cfg.GetSet(setName)
			if err != nil {
				return fmt.Errorf("set %q not found: %w", setName, err)
			}

			p, err := env.NewPipeline(setName + "_pipeline")
			if err != nil {
				return err
			}

			for _, step := range steps {
				step := step // capture
				switch strings.ToLower(step) {
				case "uppercase":
					_ = p.AddStep(func(s *env.Set) (*env.Set, error) {
						return env.UppercaseValues(s)
					})
				case "prefix_app":
					_ = p.AddStep(func(s *env.Set) (*env.Set, error) {
						return env.PrefixValues(s, "APP_")
					})
				default:
					return fmt.Errorf("unknown pipeline step: %q", step)
				}
			}

			out, err := p.Run(src)
			if err != nil {
				return fmt.Errorf("pipeline run: %w", err)
			}

			data, err := env.Export(out, outputFormat)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&steps, "step", "s", nil,
		"Ordered transform steps (uppercase, prefix_app)")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "dotenv",
		"Output format: dotenv, shell, json")
	return cmd
}
