// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestMediaCreateAgentReadyPlanAndDryRun(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	first := executeMediaCommand(t, home,
		"--agent", "create", "A polished coffee launch image",
		"--type", "image", "--purpose", "social launch", "--platform", "instagram",
		"--aspect-ratio", "9:16", "--style", "warm editorial photography",
		"--reference", "https://example.test/coffee.png",
	)
	var envelope struct {
		Results struct {
			Ready bool `json:"ready"`
			Brief struct {
				ID   string `json:"id"`
				Plan struct {
					Model string `json:"model"`
				} `json:"plan"`
			} `json:"brief"`
		} `json:"results"`
	}
	if err := json.Unmarshal(first, &envelope); err != nil {
		t.Fatalf("agent output is not JSON: %v\n%s", err, first)
	}
	if !envelope.Results.Ready || envelope.Results.Brief.ID == "" {
		t.Fatalf("ready brief missing from output: %s", first)
	}
	if got := envelope.Results.Brief.Plan.Model; got != "gpt-image-2-image-to-image" {
		t.Fatalf("ready model = %q", got)
	}

	dryRun := executeMediaCommand(t, home,
		"--agent", "--dry-run", "create", "--brief", envelope.Results.Brief.ID, "--submit",
	)
	var dryEnvelope struct {
		Results struct {
			DryRun    bool `json:"dry_run"`
			Submitted bool `json:"submitted"`
		} `json:"results"`
	}
	if err := json.Unmarshal(dryRun, &dryEnvelope); err != nil {
		t.Fatalf("dry-run output is not JSON: %v\n%s", err, dryRun)
	}
	if !dryEnvelope.Results.DryRun || dryEnvelope.Results.Submitted {
		t.Fatalf("dry run unexpectedly submitted: %s", dryRun)
	}
}

func TestMediaCreateWaitRequiresSubmit(t *testing.T) {
	var flags rootFlags
	cmd := newRootCmd(&flags)
	cmd.SetArgs([]string{"--home", filepath.Join(t.TempDir(), "kie-home"), "--agent", "create", "request", "--wait"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_, err := cmd.ExecuteC()
	if err == nil {
		t.Fatal("--wait without --submit unexpectedly succeeded")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("exit code = %d, want 2; err=%v", got, err)
	}
}

func TestMediaCreateVideoReturnsPreviewGateAndDryRunPlan(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home,
		"--agent", "create", "A cinematic trainer launch",
		"--type", "video", "--purpose", "campaign", "--platform", "instagram",
		"--aspect-ratio", "9:16", "--duration", "5", "--audio", "off",
		"--video-mode", "text", "--style", "sunrise tracking shot",
	)
	var envelope struct {
		Results struct {
			Ready      bool   `json:"ready"`
			CanSubmit  bool   `json:"can_submit"`
			NextAction string `json:"next_action"`
			Brief      struct {
				ID string `json:"id"`
			} `json:"brief"`
		} `json:"results"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("video output is not JSON: %v\n%s", err, output)
	}
	if !envelope.Results.Ready || envelope.Results.CanSubmit || envelope.Results.NextAction != "generate_preview" {
		t.Fatalf("video gate output = %s", output)
	}

	dryRun := executeMediaCommand(t, home,
		"--agent", "--dry-run", "create", "--brief", envelope.Results.Brief.ID, "--preview",
	)
	var dryEnvelope struct {
		Results struct {
			Kind string `json:"kind"`
			Plan struct {
				Model string `json:"model"`
			} `json:"plan"`
			DryRun bool `json:"dry_run"`
		} `json:"results"`
	}
	if err := json.Unmarshal(dryRun, &dryEnvelope); err != nil {
		t.Fatalf("preview dry-run output is not JSON: %v\n%s", err, dryRun)
	}
	if !dryEnvelope.Results.DryRun || dryEnvelope.Results.Kind != "preview" || dryEnvelope.Results.Plan.Model != "gpt-image-2-text-to-image" {
		t.Fatalf("preview dry run = %s", dryRun)
	}

	var flags rootFlags
	cmd := newRootCmd(&flags)
	cmd.SetArgs([]string{"--home", home, "--agent", "create", "--brief", envelope.Results.Brief.ID, "--submit"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if _, err := cmd.ExecuteC(); err == nil {
		t.Fatal("video submit bypassed the preview gate")
	}
}

func TestMediaGenerationTerminalIncludesKieFail(t *testing.T) {
	if !mediaGenerationTerminal("fail") {
		t.Fatal("Kie fail status must stop polling")
	}
}

func TestMediaWorkflowDiscoveryIsCompactAndLocal(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home, "--agent", "media", "workflow", "show", "youtube-thumbnail")
	var envelope struct {
		Results struct {
			Name   string   `json:"name"`
			Skill  string   `json:"skill"`
			Stages []string `json:"stages"`
		} `json:"results"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("workflow output is not JSON: %v\n%s", err, output)
	}
	if envelope.Results.Name != "youtube-thumbnail" || envelope.Results.Skill != "kie-youtube-thumbnail" || len(envelope.Results.Stages) == 0 {
		t.Fatalf("workflow output = %s", output)
	}
}

func TestMediaStoryboardCLIProducesGatedShotBriefs(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "kie-home")
	created := executeMediaCommand(t, home,
		"--agent", "create", "A two-shot product reveal",
		"--type", "video", "--purpose", "launch", "--platform", "youtube",
		"--aspect-ratio", "16:9", "--duration", "10", "--audio", "off",
		"--video-mode", "text", "--style", "cinematic studio", "--production-mode", "storyboard",
	)
	var createEnvelope struct {
		Results struct {
			NextAction string `json:"next_action"`
			Brief      struct {
				ID string `json:"id"`
			} `json:"brief"`
		} `json:"results"`
	}
	if err := json.Unmarshal(created, &createEnvelope); err != nil {
		t.Fatal(err)
	}
	briefID := createEnvelope.Results.Brief.ID
	if briefID == "" || createEnvelope.Results.NextAction != "draft_script" {
		t.Fatalf("storyboard master = %s", created)
	}

	scriptPath := filepath.Join(root, "script.md")
	if err := os.WriteFile(scriptPath, []byte("Open on the product. Reveal its defining feature."), 0o600); err != nil {
		t.Fatal(err)
	}
	executeMediaCommand(t, home, "--agent", "media", "script", "set", briefID, "--file", scriptPath, "--title", "Reveal")
	executeMediaCommand(t, home, "--agent", "media", "script", "approve", briefID)

	storyboardPath := filepath.Join(root, "storyboard.json")
	storyboardJSON := `{"title":"Reveal","shots":[{"duration_seconds":5,"visual":"Product on a dark plinth","camera":"slow dolly"},{"duration_seconds":5,"visual":"Feature illuminated in a hero frame","camera":"controlled orbit"}]}`
	if err := os.WriteFile(storyboardPath, []byte(storyboardJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	setOutput := executeMediaCommand(t, home, "--agent", "media", "storyboard", "set", briefID, "--file", storyboardPath)
	var setEnvelope struct {
		Results struct {
			Shots []struct {
				BriefID string `json:"shot_brief_id"`
			} `json:"shots"`
		} `json:"results"`
	}
	if err := json.Unmarshal(setOutput, &setEnvelope); err != nil {
		t.Fatal(err)
	}
	if len(setEnvelope.Results.Shots) != 2 || setEnvelope.Results.Shots[0].BriefID == "" {
		t.Fatalf("storyboard set output = %s", setOutput)
	}

	approved := executeMediaCommand(t, home, "--agent", "media", "storyboard", "approve", briefID)
	var approvedEnvelope struct {
		Results struct {
			NextAction string `json:"next_action"`
			Shots      []struct {
				Turn struct {
					NextAction string `json:"next_action"`
					CanSubmit  bool   `json:"can_submit"`
				} `json:"turn"`
			} `json:"shots"`
		} `json:"results"`
	}
	if err := json.Unmarshal(approved, &approvedEnvelope); err != nil {
		t.Fatal(err)
	}
	if approvedEnvelope.Results.NextAction != "generate_shot_previews" || len(approvedEnvelope.Results.Shots) != 2 ||
		approvedEnvelope.Results.Shots[0].Turn.NextAction != "generate_preview" || approvedEnvelope.Results.Shots[0].Turn.CanSubmit {
		t.Fatalf("approved storyboard = %s", approved)
	}
	var flags rootFlags
	cmd := newRootCmd(&flags)
	cmd.SetArgs([]string{"--home", home, "--agent", "--dry-run", "create", "--brief", briefID, "--submit"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_, err := cmd.ExecuteC()
	if err == nil || !strings.Contains(err.Error(), "next action is generate_shot_previews") {
		t.Fatalf("storyboard master submit error = %v", err)
	}
}

func TestMediaScriptReadErrorRedactsParentPath(t *testing.T) {
	var flags rootFlags
	cmd := newRootCmd(&flags)
	privatePath := filepath.Join(t.TempDir(), "private-project", "missing-script.md")
	cmd.SetArgs([]string{"--home", filepath.Join(t.TempDir(), "home"), "--agent", "media", "script", "set", "brief_missing", "--file", privatePath})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_, err := cmd.ExecuteC()
	if err == nil {
		t.Fatal("missing script unexpectedly succeeded")
	}
	if strings.Contains(err.Error(), filepath.Dir(privatePath)) || !strings.Contains(err.Error(), filepath.Base(privatePath)) {
		t.Fatalf("script read error leaked parent path: %v", err)
	}
}

func executeMediaCommand(t *testing.T, home string, args ...string) []byte {
	t.Helper()
	var flags rootFlags
	cmd := newRootCmd(&flags)
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(append([]string{"--home", home}, args...))
	if _, err := cmd.ExecuteC(); err != nil {
		t.Fatalf("media command failed: %v\nstderr: %s\nstdout: %s", err, stderr.String(), stdout.String())
	}
	return stdout.Bytes()
}

func TestMediaModelsExposesCurrentStandardTaskModelCatalog(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home, "--agent", "media", "models")
	var envelope struct {
		Results struct {
			Source string `json:"source"`
			Models []struct {
				Model  string `json:"model"`
				Source string `json:"source"`
			} `json:"models"`
		} `json:"results"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("catalog output is not JSON: %v\n%s", err, output)
	}
	if envelope.Results.Source != "embedded_snapshot" {
		t.Fatalf("catalog source = %q", envelope.Results.Source)
	}
	required := map[string]bool{
		"kling/v2-5-turbo-image-to-video-pro": false,
		"qwen3/image-to-image":                false,
		"qwen3/pro-image-to-image":            false,
		"qwen3/pro-text-to-image":             false,
		"qwen3/text-to-image":                 false,
		"seedream/5-pro-layer-decomposition":  false,
		"wan/2-7-text-to-video":               false,
	}
	seen := map[string]bool{}
	for _, model := range envelope.Results.Models {
		if seen[model.Model] {
			t.Fatalf("duplicate model ID %q in catalog", model.Model)
		}
		if model.Source == "" {
			t.Fatalf("model %q has no official source", model.Model)
		}
		seen[model.Model] = true
		if _, ok := required[model.Model]; ok {
			required[model.Model] = true
		}
	}
	for model, present := range required {
		if !present {
			t.Fatalf("required model %q missing from catalog: %s", model, output)
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

func TestMediaVideoPassesSelectedModelInputInDryRun(t *testing.T) {
	home := filepath.Join(t.TempDir(), "kie-home")
	output := executeMediaCommand(t, home,
		"--agent", "--dry-run", "media", "video", "--model", "wan/2-6-text-to-video",
		"--input", `{"prompt":"A red kite crosses a dawn sky","duration":"5"}`,
	)
	if !strings.Contains(string(output), `"model": "wan/2-6-text-to-video"`) || !strings.Contains(string(output), `"duration": "5"`) {
		t.Fatalf("selected-model dry run = %s", output)
	}
}

func TestMediaVideoRejectsInvalidWan27Input(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "duration", args: []string{"prompt", "--duration", "16"}, want: "--duration must be from 2 through 15"},
		{name: "audio URL", args: []string{"prompt", "--audio-url", "not a URI"}, want: "--audio-url must be an absolute URI"},
		{name: "raw Wan input", args: []string{"--input", `{"prompt":"prompt"}`}, want: "Wan 2.7 uses its validated flags"},
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

func TestMediaVideoHelpDocumentsWan27Schema(t *testing.T) {
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

func TestMarketResultURLsParsesResultJSONAndFutureURLFields(t *testing.T) {
	data := json.RawMessage(`{"data":{"resultJson":"{\"resultUrls\":[\"https://cdn.example.test/video.mp4\"],\"firstFrameUrl\":\"https://cdn.example.test/first.png\"}","resultObject":{"lastFrameUrl":"https://cdn.example.test/last.png"}}}`)
	urls := marketResultURLs(data)
	seen := map[string]bool{}
	for _, url := range urls {
		seen[url] = true
	}
	for _, want := range []string{"https://cdn.example.test/video.mp4", "https://cdn.example.test/first.png", "https://cdn.example.test/last.png"} {
		if !seen[want] {
			t.Fatalf("result URLs = %#v; missing %s", urls, want)
		}
	}
}
