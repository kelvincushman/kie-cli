// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package kiecatalog

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

//go:embed capabilities.json
var capabilitiesJSON []byte

type CapabilityRegistry struct {
	SchemaVersion               int                   `json:"schema_version"`
	CatalogSourceSHA256         string                `json:"catalog_source_sha256"`
	CoverageSourceSHA256        string                `json:"coverage_source_sha256"`
	CatalogModelCount           int                   `json:"catalog_model_count"`
	DocumentedOperationVariants int                   `json:"documented_operation_variants"`
	Models                      []ModelCapability     `json:"models"`
	Operations                  []OperationCapability `json:"operations"`
}

type ModelCapability struct {
	ModelID               string          `json:"model_id"`
	PrimaryCapability     string          `json:"primary_capability"`
	SecondaryCapabilities []string        `json:"secondary_capabilities,omitempty"`
	ProductionFit         []string        `json:"production_fit,omitempty"`
	Creative              bool            `json:"creative"`
	Proof                 ProofCapability `json:"proof"`
	RoutingNotes          string          `json:"routing_notes,omitempty"`
}

type ProofCapability struct {
	ResolutionFields   []string `json:"resolution_fields,omitempty"`
	LowestFaithfulTier string   `json:"lowest_faithful_tier,omitempty"`
	AlternateAllowed   bool     `json:"alternate_allowed"`
}

type OperationCapability struct {
	OperationID       string `json:"operation_id"`
	Method            string `json:"method"`
	Path              string `json:"path"`
	VariantCount      int    `json:"variant_count"`
	PrimaryCapability string `json:"primary_capability"`
	Creative          bool   `json:"creative"`
	Plumbing          bool   `json:"plumbing"`
	Reason            string `json:"reason"`
}

var (
	capabilitiesOnce sync.Once
	capabilities     CapabilityRegistry
	capabilitiesErr  error
	capabilityByID   map[string]int
)

func LoadCapabilities() (*CapabilityRegistry, error) {
	capabilitiesOnce.Do(func() {
		if err := json.Unmarshal(capabilitiesJSON, &capabilities); err != nil {
			capabilitiesErr = fmt.Errorf("decode embedded Kie capability classifier: %w", err)
			return
		}
		if err := ValidateCapabilities(&capabilities, catalogJSON, nil); err != nil {
			capabilitiesErr = err
			return
		}
		capabilityByID = make(map[string]int, len(capabilities.Models))
		for i := range capabilities.Models {
			capabilityByID[capabilities.Models[i].ModelID] = i
		}
	})
	if capabilitiesErr != nil {
		return nil, capabilitiesErr
	}
	return &capabilities, nil
}

func GetCapability(modelID string) (*ModelCapability, error) {
	registry, err := LoadCapabilities()
	if err != nil {
		return nil, err
	}
	index, ok := capabilityByID[strings.TrimSpace(modelID)]
	if !ok {
		return nil, fmt.Errorf("Kie model %q is unclassified; refresh and review capabilities.json before routing it", modelID)
	}
	return &registry.Models[index], nil
}

// ValidateCapabilities fails closed when catalog models, source hashes, or
// documented operation coverage drift. coverageJSON may be nil at runtime;
// tests and refresh tooling pass it to verify the second checked-in source.
func ValidateCapabilities(registry *CapabilityRegistry, catalogSource, coverageSource []byte) error {
	if registry == nil || registry.SchemaVersion != 1 {
		return fmt.Errorf("unsupported Kie capability classifier schema")
	}
	if registry.CatalogSourceSHA256 != sumHex(catalogSource) {
		return fmt.Errorf("Kie capability classifier catalog hash is stale")
	}
	var catalog Registry
	if err := json.Unmarshal(catalogSource, &catalog); err != nil {
		return fmt.Errorf("decode catalog while validating capabilities: %w", err)
	}
	if registry.CatalogModelCount != catalog.ModelCount || len(registry.Models) != catalog.ModelCount {
		return fmt.Errorf("Kie capability classifier covers %d models, catalog contains %d", len(registry.Models), catalog.ModelCount)
	}
	wanted := make(map[string]struct{}, len(catalog.Models))
	for _, model := range catalog.Models {
		wanted[model.ID] = struct{}{}
	}
	seen := make(map[string]struct{}, len(registry.Models))
	allowed := map[string]bool{"kie-image": true, "kie-video": true, "kie-audio": true, "kie-avatar": true, "kie-identity": true}
	for _, item := range registry.Models {
		if _, ok := wanted[item.ModelID]; !ok {
			return fmt.Errorf("capability classifier contains unknown model %q", item.ModelID)
		}
		if _, duplicate := seen[item.ModelID]; duplicate {
			return fmt.Errorf("capability classifier contains duplicate model %q", item.ModelID)
		}
		if !allowed[item.PrimaryCapability] {
			return fmt.Errorf("model %q has invalid primary capability %q", item.ModelID, item.PrimaryCapability)
		}
		seen[item.ModelID] = struct{}{}
	}
	if len(coverageSource) > 0 {
		if registry.CoverageSourceSHA256 != sumHex(coverageSource) {
			return fmt.Errorf("Kie capability classifier coverage hash is stale")
		}
		var coverage struct {
			DocumentedOperationVariants int `json:"documented_operation_variants"`
			UniqueOperations            int `json:"unique_operations"`
		}
		if err := json.Unmarshal(coverageSource, &coverage); err != nil {
			return fmt.Errorf("decode API coverage while validating capabilities: %w", err)
		}
		variants := 0
		operationIDs := make(map[string]struct{}, len(registry.Operations))
		allowedOperations := map[string]bool{
			"plumbing": true, "creative-router": true,
			"kie-image": true, "kie-video": true, "kie-audio": true, "kie-avatar": true, "kie-identity": true,
		}
		for _, operation := range registry.Operations {
			if operation.OperationID == "" || operation.VariantCount < 1 || operation.Creative == operation.Plumbing {
				return fmt.Errorf("invalid operation capability row for %q", operation.OperationID)
			}
			if !allowedOperations[operation.PrimaryCapability] || (operation.Plumbing && operation.PrimaryCapability != "plumbing") || (operation.Creative && operation.PrimaryCapability == "plumbing") {
				return fmt.Errorf("operation %q has invalid primary capability %q", operation.OperationID, operation.PrimaryCapability)
			}
			if _, duplicate := operationIDs[operation.OperationID]; duplicate {
				return fmt.Errorf("duplicate operation capability %q", operation.OperationID)
			}
			operationIDs[operation.OperationID] = struct{}{}
			variants += operation.VariantCount
		}
		if len(operationIDs) != coverage.UniqueOperations || variants != coverage.DocumentedOperationVariants || variants != registry.DocumentedOperationVariants {
			return fmt.Errorf("operation classifier drift: %d operations/%d variants, want %d/%d", len(operationIDs), variants, coverage.UniqueOperations, coverage.DocumentedOperationVariants)
		}
	}
	return nil
}

func CapabilityNames() ([]string, error) {
	registry, err := LoadCapabilities()
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	for _, item := range registry.Models {
		seen[item.PrimaryCapability] = struct{}{}
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values, nil
}

func sumHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
