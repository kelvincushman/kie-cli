// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestMediaModelsUsesCompleteCanonicalCatalog(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home, "--agent", "media", "models")
	var envelope struct {
		Results struct {
			Source            string `json:"source"`
			CatalogModelCount int    `json:"catalog_model_count"`
			ModelCount        int    `json:"model_count"`
			Models            []struct {
				Model       string   `json:"model"`
				InputFields []string `json:"input_fields"`
				Source      string   `json:"source"`
			} `json:"models"`
		} `json:"results"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("catalog output is not JSON: %v\n%s", err, output)
	}
	if envelope.Results.Source != "embedded_snapshot" || envelope.Results.CatalogModelCount != 129 || envelope.Results.ModelCount != 129 {
		t.Fatalf("unexpected catalog metadata: %s", output)
	}
	required := map[string]bool{
		"kling/v2-5-turbo-image-to-video-pro": false,
		"qwen3/pro-text-to-image":             false,
		"bytedance/seedance-2-5":              false,
		"seedream/5-pro-layer-decomposition":  false,
		"wan/2-7-text-to-video":               false,
	}
	seen := map[string]bool{}
	for _, model := range envelope.Results.Models {
		if seen[model.Model] {
			t.Fatalf("duplicate model ID %q in catalog", model.Model)
		}
		if model.Source == "" {
			t.Fatalf("model %q lacks an official source", model.Model)
		}
		seen[model.Model] = true
		if _, ok := required[model.Model]; ok {
			required[model.Model] = true
		}
	}
	for model, present := range required {
		if !present {
			t.Fatalf("required model %q missing from catalog", model)
		}
	}
	for _, model := range envelope.Results.Models {
		if model.Model == "wan/2-7-text-to-video" && len(model.InputFields) != 10 {
			t.Fatalf("Wan 2.7 input fields = %v; want all 10 documented settings", model.InputFields)
		}
	}
}

func TestMediaModelsFiltersWithoutLoadingASecondCatalog(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home, "--agent", "media", "models", "wan", "--family", "video")
	if !strings.Contains(string(output), `"model": "wan/2-7-text-to-video"`) || strings.Contains(string(output), `"model": "gpt-image`) {
		t.Fatalf("filtered catalog output = %s", output)
	}
}

func TestMediaVideoHelpLabelsAdvancedGateBypassAndFreshConfirmation(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := strings.ToLower(string(executeMediaCommand(t, home, "media", "video", "--help")))
	for _, phrase := range []string{"advanced direct generation", "bypasses", "still/proof approval", "scoped paid-confirmation", "fresh explicit user confirmation", "--dry-run"} {
		if !strings.Contains(output, phrase) {
			t.Fatalf("video help does not contain %q:\n%s", phrase, output)
		}
	}
}

func TestMediaVideoBuildsExactWan27RequestInDryRun(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home,
		"--agent", "--dry-run", "media", "video", "A red kite crosses a dawn sky",
		"--duration", "7", "--ratio", "9:16", "--negative-prompt", "no text", "--seed", "42",
	)
	var envelope struct {
		Results struct {
			DryRun bool `json:"dry_run"`
			Body   struct {
				Model string `json:"model"`
				Input struct {
					Prompt       string `json:"prompt"`
					Resolution   string `json:"resolution"`
					Duration     int    `json:"duration"`
					Ratio        string `json:"ratio"`
					PromptExtend bool   `json:"prompt_extend"`
					Watermark    bool   `json:"watermark"`
					Seed         int    `json:"seed"`
				} `json:"input"`
			} `json:"body"`
		} `json:"results"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("video dry-run output is not JSON: %v\n%s", err, output)
	}
	if !envelope.Results.DryRun || envelope.Results.Body.Model != "wan/2-7-text-to-video" ||
		envelope.Results.Body.Input.Prompt != "A red kite crosses a dawn sky" ||
		envelope.Results.Body.Input.Resolution != "1080p" || envelope.Results.Body.Input.Duration != 7 ||
		envelope.Results.Body.Input.Ratio != "9:16" || !envelope.Results.Body.Input.PromptExtend ||
		envelope.Results.Body.Input.Watermark || envelope.Results.Body.Input.Seed != 42 {
		t.Fatalf("video dry-run request = %s", output)
	}
}

func TestMediaVideoValidatesSelectedModelInputInDryRun(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home,
		"--agent", "--dry-run", "media", "video", "--model", "wan/2-6-text-to-video",
		"--input", `{"prompt":"A red kite crosses a dawn sky","duration":"5"}`,
	)
	if !strings.Contains(string(output), `"model": "wan/2-6-text-to-video"`) || !strings.Contains(string(output), `"duration": "5"`) {
		t.Fatalf("selected-model dry run = %s", output)
	}
}

func TestMediaVideoRejectsInvalidInputsBeforeLiveCalls(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "duration", args: []string{"prompt", "--duration", "16"}, want: "--duration must be from 2 through 15"},
		{name: "audio URL", args: []string{"prompt", "--audio-url", "not a URI"}, want: "--audio-url must be an absolute URI"},
		{name: "raw Wan input", args: []string{"--input", `{"prompt":"prompt"}`}, want: "Wan 2.7 uses its validated flags"},
		{name: "custom model enum", args: []string{"--model", "wan/2-6-text-to-video", "--input", `{"prompt":"prompt","duration":"7"}`}, want: "documented wan/2-6-text-to-video contract"},
		{name: "unknown model", args: []string{"--model", "future/video-model", "--input", `{"prompt":"prompt"}`}, want: "is not in the embedded catalog"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var flags rootFlags
			cmd := newRootCmd(&flags)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(append([]string{"--home", filepath.Join(t.TempDir(), "kie-home"), "--agent", "--dry-run", "media", "video"}, test.args...))
			_, err := cmd.ExecuteC()
			if err == nil || ExitCode(err) != 2 || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v; want usage error containing %q", err, test.want)
			}
		})
	}
}

func TestMediaVideoHelpDocumentsWan27Settings(t *testing.T) {
	var flags rootFlags
	cmd := newRootCmd(&flags)
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"media", "video", "--help"})
	if _, err := cmd.ExecuteC(); err != nil {
		t.Fatalf("video help failed: %v", err)
	}
	for _, want := range []string{"--audio-url", "--duration", "--negative-prompt", "--prompt-extend", "--ratio", "--resolution", "--seed", "--watermark", "--nsfw-checker", "--wait"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("video help missing %s:\n%s", want, stdout.String())
		}
	}
}

type queuedMarketTaskClient struct {
	responses []json.RawMessage
	calls     int
	path      string
	params    map[string]string
}

func (c *queuedMarketTaskClient) Get(_ context.Context, path string, params map[string]string) (json.RawMessage, error) {
	c.path = path
	c.params = params
	response := c.responses[c.calls]
	c.calls++
	return response, nil
}

func TestWaitForMarketTaskPollsToTerminalResult(t *testing.T) {
	client := &queuedMarketTaskClient{responses: []json.RawMessage{
		json.RawMessage(`{"data":{"state":"processing"}}`),
		json.RawMessage(`{"data":{"state":"success","resultJson":"{\"resultUrls\":[\"https://cdn.example.test/video.mp4\"]}"}}`),
	}}
	data, err := waitForMarketTask(&cobra.Command{}, client, "task_wan_1", time.Nanosecond, time.Second)
	if err != nil {
		t.Fatalf("polling failed: %v", err)
	}
	if client.calls != 2 || client.path != "/api/v1/jobs/recordInfo" || client.params["taskId"] != "task_wan_1" {
		t.Fatalf("unexpected polls: calls=%d path=%q params=%#v", client.calls, client.path, client.params)
	}
	if state := marketTaskState(data); state != "success" {
		t.Fatalf("terminal state = %q", state)
	}
	urls := marketResultURLs(data)
	if len(urls) != 1 || urls[0] != "https://cdn.example.test/video.mp4" {
		t.Fatalf("terminal URLs = %#v", urls)
	}
}

func TestMarketResultURLsParsesNestedResultFields(t *testing.T) {
	data := json.RawMessage(`{"data":{"resultJson":"{\"resultUrls\":[\"https://cdn.example.test/video.mp4\"],\"firstFrameUrl\":\"https://cdn.example.test/first.png\"}","resultObject":{"lastFrameUrl":"https://cdn.example.test/last.png"}}}`)
	urls := marketResultURLs(data)
	seen := map[string]bool{}
	for _, resultURL := range urls {
		seen[resultURL] = true
	}
	for _, want := range []string{"https://cdn.example.test/video.mp4", "https://cdn.example.test/first.png", "https://cdn.example.test/last.png"} {
		if !seen[want] {
			t.Fatalf("result URLs = %#v; missing %s", urls, want)
		}
	}
}
