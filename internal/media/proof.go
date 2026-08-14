// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"kie-pp-cli/internal/kiecatalog"
)

type ProofOption struct {
	Supported          bool     `json:"supported"`
	ProofModel         string   `json:"proof_model"`
	FinalModel         string   `json:"final_model"`
	ResolutionField    string   `json:"resolution_field,omitempty"`
	ResolutionValue    string   `json:"resolution_value,omitempty"`
	LowestTier         string   `json:"lowest_tier,omitempty"`
	SameModel          bool     `json:"same_model"`
	AlternateLabel     string   `json:"alternate_label,omitempty"`
	PriceEstimate      string   `json:"price_estimate,omitempty"`
	ExactCostKnown     bool     `json:"exact_cost_known"`
	Disclosure         string   `json:"disclosure"`
	UnsupportedReason  string   `json:"unsupported_reason,omitempty"`
	SourceSchemaFields []string `json:"source_schema_fields,omitempty"`
}

func ResolveProofOption(modelID string) ProofOption {
	option := ProofOption{
		ProofModel: modelID, FinalModel: modelID, SameModel: true,
		Disclosure: "A complete-shot proof is a live paid generation. Exact cost is not available in the local Kie catalog and must be confirmed before submission.",
	}
	capability, err := kiecatalog.GetCapability(modelID)
	if err != nil {
		option.UnsupportedReason = err.Error()
		return option
	}
	if capability.PrimaryCapability != "kie-video" && capability.PrimaryCapability != "kie-avatar" {
		option.UnsupportedReason = "complete-shot proofs apply only to video or avatar generation models"
		return option
	}
	model, err := kiecatalog.Get(modelID)
	if err != nil {
		option.UnsupportedReason = err.Error()
		return option
	}
	properties, _ := model.InputSchema["properties"].(map[string]any)
	fields := append([]string(nil), capability.Proof.ResolutionFields...)
	if len(fields) == 0 {
		option.UnsupportedReason = "the selected model has no documented lower-fidelity setting for a faithful proof"
		return option
	}
	for _, field := range fields {
		schema, _ := properties[field].(map[string]any)
		values := stringValues(schema["enum"])
		if len(values) == 0 {
			continue
		}
		option.ResolutionField = field
		option.ResolutionValue = lowestProofValue(values)
		option.LowestTier = option.ResolutionValue
		option.SourceSchemaFields = fields
		option.Supported = option.ResolutionValue != ""
		if option.Supported && proofPixels(option.ResolutionValue) >= 720 {
			option.Disclosure += " No cheaper faithful same-model tier is documented; this proof uses the model's lowest supported tier."
		}
		return option
	}
	option.UnsupportedReason = "resolution-like fields exist but do not publish selectable tiers"
	return option
}

func stringValues(value any) []string {
	items, _ := value.([]any)
	values := make([]string, 0, len(items))
	for _, item := range items {
		switch typed := item.(type) {
		case string:
			values = append(values, typed)
		case float64:
			values = append(values, strconv.FormatFloat(typed, 'f', -1, 64))
		}
	}
	return values
}

func lowestProofValue(values []string) string {
	values = append([]string(nil), values...)
	rank := func(value string) int {
		if pixels := proofPixels(value); pixels > 0 {
			return pixels
		}
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "low", "draft":
			return 1
		case "std", "standard":
			return 2
		case "medium":
			return 3
		case "high":
			return 4
		case "pro":
			return 5
		default:
			return 100000
		}
	}
	sort.SliceStable(values, func(i, j int) bool { return rank(values[i]) < rank(values[j]) })
	if len(values) == 0 || rank(values[0]) == 100000 {
		return ""
	}
	return values[0]
}

func proofPixels(value string) int {
	value = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(value)), "p")
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func BuildProofPlan(brief *Brief) (*Plan, ProofOption, error) {
	if brief == nil || brief.MediaType != "video" {
		return nil, ProofOption{}, fmt.Errorf("complete-shot proof applies only to video briefs")
	}
	final := BuildPlan(brief)
	option := ResolveProofOption(final.Model)
	if !option.Supported {
		return nil, option, fmt.Errorf("proof is unavailable: %s", option.UnsupportedReason)
	}
	data := cloneMap(final.Input)
	data[option.ResolutionField] = option.ResolutionValue
	return &Plan{
		Model: option.ProofModel, Input: data, Rationale: option.Disclosure,
		ProductionSkill: final.ProductionSkill, CapabilitySkill: final.CapabilitySkill,
		CostStatus: "unknown_until_provider_quote", OverrideOptions: append([]string(nil), final.OverrideOptions...),
		ProductionStage: "proof",
	}, option, nil
}

func cloneMap(input map[string]any) map[string]any {
	data := make(map[string]any, len(input))
	for key, value := range input {
		data[key] = value
	}
	return data
}
