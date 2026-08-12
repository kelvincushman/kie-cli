// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package kiecatalog

import "testing"

func TestCatalogContainsCompleteCurrentContracts(t *testing.T) {
	registry, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if registry.ModelCount < 129 {
		t.Fatalf("model count = %d, want at least 129", registry.ModelCount)
	}
	model, err := Get("bytedance/seedance-2-5")
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"prompt", "reference_image_urls", "reference_video_urls", "reference_audio_urls", "resolution", "duration", "generate_audio", "output_format"} {
		if _, ok := model.InputSchema["properties"].(map[string]any)[field]; !ok {
			t.Errorf("Seedance 2.5 input schema missing %q", field)
		}
	}
}

func TestCatalogUsesCorrectedUpstreamModelIDs(t *testing.T) {
	for _, id := range []string{"qwen2/text-to-image", "kling/v2-5-turbo-image-to-video-pro"} {
		if _, err := Get(id); err != nil {
			t.Errorf("Get(%q): %v", id, err)
		}
	}
}

func TestValidateChecksRequiredEnumAndType(t *testing.T) {
	issues, err := Validate("qwen2/text-to-image", map[string]any{"image_size": "not-a-size", "seed": "wrong"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"input.prompt": false, "input.image_size": false, "input.seed": false}
	for _, item := range issues {
		if _, ok := want[item.Path]; ok {
			want[item.Path] = true
		}
	}
	for path, found := range want {
		if !found {
			t.Errorf("missing validation issue for %s: %#v", path, issues)
		}
	}
}

func TestExampleReturnsModelInputOnly(t *testing.T) {
	example, err := Example("bytedance/seedance-2-5")
	if err != nil {
		t.Fatal(err)
	}
	if _, hasModel := example["model"]; hasModel {
		t.Fatalf("example unexpectedly contains top-level model: %#v", example)
	}
	if len(example) == 0 {
		t.Fatal("example is empty")
	}
}
