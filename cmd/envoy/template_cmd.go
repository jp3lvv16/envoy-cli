package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/config"
	"github.com/yourorg/envoy-cli/internal/env"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Render templates using environment variable sets",
	}
	cmd.AddCommand(newTemplateRenderCmd())
	return cmd
}

func newTemplateRenderCmd() *cobra.Command {
	var setName string
	var strict bool
	var templateStr string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a template string using variables from a set",
		Example: `  envoy template render --set production --template "http://{{HOST}}:{{PORT}}"
  envoy template render --set staging --template "{{DB_URL}}" --strict`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if setName == "" {
				return fmt.Errorf("--set is required")
			}
			if templateStr == "" {
				return fmt.Errorf("--template is required")
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			s := cfg.GetSet(setName)
			if s == nil {
				return fmt.Errorf("set %q not found", setName)
			}

			if strict {
				out, err := env.RenderStrict(s, templateStr)
				if err != nil {
					return err
				}
				fmt.Fprintln(os.Stdout, out)
				return nil
			}

			res, err := env.Render(s, templateStr)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, res.Output)
			if len(res.Missing) > 0 {
				fmt.Fprintf(os.Stderr, "warning: unresolved placeholders: %v\n", res.Missing)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&setName, "set", "s", "", "Name of the environment set to use")
	cmd.Flags().StringVarP(&templateStr, "template", "t", "", "Template string with {{KEY}} placeholders")
	cmd.Flags().BoolVar(&strict, "strict", false, "Fail if any placeholder cannot be resolved")
	return cmd
}
