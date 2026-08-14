// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package kiecatalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCapabilityClassifierCoversCatalogAndOperations(t *testing.T) {
	coverage, err := os.ReadFile(filepath.Join("..", "..", "research", "kie-api-coverage.json"))
	if err != nil {
		t.Fatal(err)
	}
	var registry CapabilityRegistry
	if err := json.Unmarshal(capabilitiesJSON, &registry); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCapabilities(&registry, catalogJSON, coverage); err != nil {
		t.Fatal(err)
	}
	if len(registry.Models) != 129 {
		t.Fatalf("classified models = %d, want 129", len(registry.Models))
	}
	if got, err := GetCapability("bytedance/seedance-2-5"); err != nil || got.PrimaryCapability != "kie-video" {
		t.Fatalf("Seedance capability = %#v, err=%v", got, err)
	}
}

func TestCapabilityClassifierFailsClosedOnSourceDrift(t *testing.T) {
	var registry CapabilityRegistry
	if err := json.Unmarshal(capabilitiesJSON, &registry); err != nil {
		t.Fatal(err)
	}
	registry.CatalogSourceSHA256 = "stale"
	if err := ValidateCapabilities(&registry, catalogJSON, nil); err == nil {
		t.Fatal("stale classifier hash was accepted")
	}
	if _, err := GetCapability("future/unclassified-model"); err == nil {
		t.Fatal("unclassified model was accepted")
	}
}

func TestOperationClassifierSeparatesCreativeMediaFromLanguagePlumbing(t *testing.T) {
	registry, err := LoadCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"create-music-video":       "kie-video",
		"suno-voice-generate":      "kie-audio",
		"market-create-task":       "creative-router",
		"gpt-codex-responses":      "plumbing",
		"gpt-5-2-chat-completions": "plumbing",
	}
	for _, operation := range registry.Operations {
		if expected, ok := want[operation.OperationID]; ok {
			if operation.PrimaryCapability != expected {
				t.Errorf("%s capability = %q, want %q", operation.OperationID, operation.PrimaryCapability, expected)
			}
			delete(want, operation.OperationID)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing classified operations: %#v", want)
	}
}
