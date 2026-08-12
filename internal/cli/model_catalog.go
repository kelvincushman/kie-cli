// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"kie-pp-cli/internal/kiecatalog"
)

func init() {
	registerNovelCommand(func(root *cobra.Command, flags *rootFlags) {
		addNovelCommandIfAbsent(root, newModelsCmd(flags))
		decorateMarketCreateWithModelValidation(root)
	})
}

// pp:data-source computed
func newModelsCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{
		Use:   "models",
		Short: "Inspect every embedded Kie Market model input and setting",
		Long: "Browse the complete model registry captured from Kie.ai's official OpenAPI pages. " +
			"These commands are local, token-efficient, and do not require an API key or spend credits.",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE:        parentNoSubcommandRunE(flags),
	}
	parent.AddCommand(newModelsListCmd(flags))
	parent.AddCommand(newModelsShowCmd(flags))
	parent.AddCommand(newModelsExampleCmd(flags))
	parent.AddCommand(newModelsValidateCmd(flags))
	return parent
}

func newModelsListCmd(flags *rootFlags) *cobra.Command {
	var query string
	var category string
	cmd := &cobra.Command{
		Use:         "list",
		Short:       "List compact model summaries and documented input field names",
		Args:        cobra.NoArgs,
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			models, err := kiecatalog.List(query, category)
			if err != nil {
				return err
			}
			if flags.asJSON || flags.agent || !isTerminal(cmd.OutOrStdout()) {
				return printJSONFiltered(cmd.OutOrStdout(), models, flags)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Kie Market models (%d)\n\n", len(models))
			for _, model := range models {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-48s %-24s %s\n", model.ID, model.Category, strings.Join(model.InputFields, ", "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&query, "search", "", "Filter IDs, names, categories, and descriptions")
	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	return cmd
}

func newModelsShowCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:         "show <model-id>",
		Aliases:     []string{"settings", "schema"},
		Short:       "Show one model's full request schema, input settings, constraints, and examples",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := kiecatalog.Get(args[0])
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, model)
		},
	}
}

func newModelsExampleCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:         "example <model-id>",
		Short:       "Print a documented starter input payload for one model",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			example, err := kiecatalog.Example(args[0])
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, example)
		},
	}
}

func newModelsValidateCmd(flags *rootFlags) *cobra.Command {
	var inputJSON string
	var stdin bool
	cmd := &cobra.Command{
		Use:         "validate <model-id>",
		Short:       "Validate documented types, required fields, enums, and ranges before generation",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := readModelInput(cmd, inputJSON, stdin)
			if err != nil {
				return usageErr(err)
			}
			issues, err := kiecatalog.Validate(args[0], input)
			if err != nil {
				return usageErr(err)
			}
			result := map[string]any{"model": args[0], "valid": len(issues) == 0, "issues": issues}
			if err := printMediaValue(cmd, flags, result); err != nil {
				return err
			}
			if len(issues) > 0 {
				return usageErr(fmt.Errorf("input does not match the documented %s contract (%d issue(s))", args[0], len(issues)))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&inputJSON, "input", "", "Model input JSON object")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read the model input JSON object from stdin")
	return cmd
}

func readModelInput(cmd *cobra.Command, value string, stdin bool) (map[string]any, error) {
	if stdin && strings.TrimSpace(value) != "" {
		return nil, fmt.Errorf("--input and --stdin are mutually exclusive")
	}
	var data []byte
	var err error
	if stdin {
		data, err = io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
	} else {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("--input <json> or --stdin is required")
		}
		data = []byte(value)
	}
	var input map[string]any
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("parsing model input JSON: %w", err)
	}
	return input, nil
}

func decorateMarketCreateWithModelValidation(root *cobra.Command) {
	jobs := findDirectChild(root, "kie-ai-jobs")
	if jobs == nil {
		return
	}
	create := findDirectChild(jobs, "market-create-task")
	if create == nil || create.Flags().Lookup("model") == nil || create.Flags().Lookup("input") == nil {
		return
	}
	var allowUnknown bool
	create.Flags().BoolVar(&allowUnknown, "allow-unknown-model", false, "Allow a newly released model absent from this embedded catalog")
	originalPreRunE := create.PreRunE
	create.PreRunE = func(cmd *cobra.Command, args []string) error {
		stdin, _ := cmd.Flags().GetBool("stdin")
		modelID, _ := cmd.Flags().GetString("model")
		inputJSON, _ := cmd.Flags().GetString("input")
		if stdin {
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return usageErr(fmt.Errorf("reading stdin: %w", err))
			}
			cmd.SetIn(bytes.NewReader(data))
			var request struct {
				Model string         `json:"model"`
				Input map[string]any `json:"input"`
			}
			if err := json.Unmarshal(data, &request); err != nil {
				return usageErr(fmt.Errorf("parsing stdin JSON: %w", err))
			}
			modelID = request.Model
			inputData, _ := json.Marshal(request.Input)
			inputJSON = string(inputData)
		}
		if strings.TrimSpace(modelID) != "" {
			if _, err := kiecatalog.Get(modelID); err != nil {
				if !allowUnknown {
					return usageErr(err)
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v; sending because --allow-unknown-model was set\n", err)
			} else {
				input := map[string]any{}
				if strings.TrimSpace(inputJSON) != "" {
					if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
						return usageErr(fmt.Errorf("parsing model input JSON: %w", err))
					}
				}
				issues, err := kiecatalog.Validate(modelID, input)
				if err != nil {
					return usageErr(err)
				}
				if len(issues) > 0 {
					data, _ := json.Marshal(issues)
					return usageErr(fmt.Errorf("model input does not match the documented %s contract: %s", modelID, data))
				}
			}
		}
		if originalPreRunE != nil {
			return originalPreRunE(cmd, args)
		}
		return nil
	}
}

func findDirectChild(parent *cobra.Command, name string) *cobra.Command {
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}
