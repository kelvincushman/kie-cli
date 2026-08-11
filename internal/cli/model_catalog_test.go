// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestModelsShowReturnsSeedanceSettings(t *testing.T) {
	cmd := RootCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "show", "bytedance/seedance-2-5", "--json", "--no-learn"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"reference_image_urls", "reference_video_urls", "generate_audio", "output_format"} {
		if !strings.Contains(out.String(), field) {
			t.Errorf("models show output missing %q", field)
		}
	}
}

func TestModelsValidateRejectsInvalidInput(t *testing.T) {
	cmd := RootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "validate", "qwen2/text-to-image", "--input", `{}`, "--json", "--no-learn"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("error = %v", err)
	}
}

func TestMarketCreateValidatesKnownModelBeforeRequest(t *testing.T) {
	cmd := RootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"kie-ai-jobs", "market-create-task", "--model", "qwen2/text-to-image", "--input", `{}`, "--dry-run", "--json", "--no-learn"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "documented qwen2/text-to-image contract") {
		t.Fatalf("error = %v", err)
	}
}

func TestMarketCreateValidatesKnownModelFromStdinBeforeRequest(t *testing.T) {
	cmd := RootCmd()
	cmd.SetIn(strings.NewReader(`{"model":"qwen2/text-to-image","input":{}}`))
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"kie-ai-jobs", "market-create-task", "--stdin", "--dry-run", "--json", "--no-learn"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "documented qwen2/text-to-image contract") {
		t.Fatalf("error = %v", err)
	}
}
