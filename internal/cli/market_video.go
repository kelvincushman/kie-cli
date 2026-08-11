// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"kie-pp-cli/internal/kiecatalog"
)

type mediaModelSummary struct {
	Model               string   `json:"model"`
	Name                string   `json:"name"`
	Category            string   `json:"category"`
	InputFields         []string `json:"input_fields"`
	RequiredInputFields []string `json:"required_input_fields"`
	Source              string   `json:"source"`
}

// pp:data-source local
func newMediaModelsCmd(flags *rootFlags) *cobra.Command {
	var family string
	cmd := &cobra.Command{
		Use:   "models [query]",
		Short: "List compact summaries from the complete embedded Kie model registry",
		Long: "List a token-efficient view of the same complete Kie Market model registry exposed by the top-level models command. " +
			"Every full input schema, setting, constraint, example, and official source remains available through 'kie-pp-cli models show <model-id>'.",
		Args:        cobra.MaximumNArgs(1),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDataSourceStrategy(flags, "local"); err != nil {
				return usageErr(err)
			}
			registry, err := kiecatalog.Load()
			if err != nil {
				return err
			}
			models, err := kiecatalog.List("", "")
			if err != nil {
				return err
			}
			query := strings.TrimSpace(family)
			if len(args) == 1 {
				query = strings.TrimSpace(query + " " + args[0])
			}
			words := strings.Fields(strings.ToLower(query))
			results := make([]mediaModelSummary, 0, len(models))
			for _, model := range models {
				haystack := strings.ToLower(strings.Join([]string{
					model.ID, model.Name, model.Category, strings.Join(model.InputFields, " "),
				}, " "))
				matches := true
				for _, word := range words {
					if !strings.Contains(haystack, word) {
						matches = false
						break
					}
				}
				if matches {
					results = append(results, mediaModelSummary{
						Model: model.ID, Name: model.Name, Category: model.Category,
						InputFields: model.InputFields, RequiredInputFields: model.RequiredInputFields,
						Source: model.Source,
					})
				}
			}
			return printMediaValue(cmd, flags, map[string]any{
				"source":              "embedded_snapshot",
				"catalog_source":      registry.SourceIndex,
				"catalog_model_count": registry.ModelCount,
				"model_count":         len(results),
				"refresh":             "scripts/weekly-refresh.sh, then rebuild the CLI",
				"models":              results,
			})
		},
	}
	cmd.Flags().StringVar(&family, "family", "", "Filter model ID, family, category, name, or input field (for example: video or wan)")
	return cmd
}

type marketVideoOptions struct {
	model        string
	input        string
	callbackURL  string
	prompt       string
	negative     string
	audioURL     string
	resolution   string
	ratio        string
	duration     int
	promptExtend bool
	watermark    bool
	seed         int64
	nsfwChecker  bool
	wait         bool
	waitInterval time.Duration
	waitTimeout  time.Duration
}

// pp:data-source live
func newMediaVideoCmd(flags *rootFlags) *cobra.Command {
	var options marketVideoOptions
	cmd := &cobra.Command{
		Use:   "video [prompt]",
		Short: "Create a direct Kie Market video task; defaults to Wan 2.7 text-to-video",
		Long: "Create one direct Market video task without a local brief or the director's preview gate. The built-in flags implement only Wan 2.7 text-to-video. " +
			"Select another captured Kie Market video model with --model and pass its documented input object through --input; the complete embedded contract is validated before any live call.",
		Example: "  kie-pp-cli media video 'A red kite crosses a dawn sky' --duration 5 --ratio 16:9\n" +
			"  kie-pp-cli media video --model wan/2-6-text-to-video --input '{\"prompt\":\"A red kite crosses a dawn sky\",\"duration\":\"5\"}'",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDataSourceStrategy(flags, "live"); err != nil {
				return usageErr(err)
			}
			if len(args) == 1 {
				options.prompt = args[0]
			}
			body, err := marketVideoRequest(cmd, options)
			if err != nil {
				return usageErr(err)
			}
			if flags.dryRun {
				return printMediaValue(cmd, flags, map[string]any{
					"dry_run": true, "path": "/api/v1/jobs/createTask", "body": body,
				})
			}
			client, err := flags.newClient()
			if err != nil {
				return err
			}
			data, status, err := client.PostWithParams(cmd.Context(), "/api/v1/jobs/createTask", map[string]string{}, body)
			if err != nil {
				return classifyAPIError(err, flags)
			}
			if status < 200 || status >= 300 {
				return fmt.Errorf("create Market video task returned HTTP %d", status)
			}
			taskID := marketTaskID(data)
			if taskID == "" {
				return fmt.Errorf("create Market video task response did not include a task id")
			}
			result := map[string]any{
				"task_id": taskID, "model": strings.TrimSpace(options.model),
				"state": "submitted", "response": json.RawMessage(data),
			}
			if options.wait {
				data, err = waitForMarketTask(cmd, client, taskID, options.waitInterval, options.waitTimeout)
				if err != nil {
					return classifyAPIError(err, flags)
				}
				result["task"] = json.RawMessage(data)
				result["state"] = marketTaskState(data)
				result["result_urls"] = marketResultURLs(data)
			} else {
				result["next_action"] = "kie-pp-cli kie-ai-jobs market-query-task --task-id " + taskID
			}
			return printMediaValue(cmd, flags, result)
		},
	}
	cmd.Flags().StringVar(&options.model, "model", "wan/2-7-text-to-video", "Kie Market model ID; built-in flags support Wan 2.7 only")
	cmd.Flags().StringVar(&options.input, "input", "", "Raw model-specific input JSON object; required for a model other than Wan 2.7")
	cmd.Flags().StringVar(&options.callbackURL, "call-back-url", "", "Optional Kie callback URL")
	cmd.Flags().StringVar(&options.negative, "negative-prompt", "", "Wan 2.7 negative prompt (up to 500 characters)")
	cmd.Flags().StringVar(&options.audioURL, "audio-url", "", "Wan 2.7 audio URL")
	cmd.Flags().StringVar(&options.resolution, "resolution", "1080p", "Wan 2.7 resolution: 720p or 1080p")
	cmd.Flags().StringVar(&options.ratio, "ratio", "16:9", "Wan 2.7 aspect ratio: 16:9, 9:16, 1:1, 4:3, or 3:4")
	cmd.Flags().IntVar(&options.duration, "duration", 5, "Wan 2.7 duration in seconds (2-15)")
	cmd.Flags().BoolVar(&options.promptExtend, "prompt-extend", true, "Ask Wan 2.7 to extend the prompt")
	cmd.Flags().BoolVar(&options.watermark, "watermark", false, "Add a Wan 2.7 watermark")
	cmd.Flags().Int64Var(&options.seed, "seed", 0, "Wan 2.7 seed (0-2147483647)")
	cmd.Flags().BoolVar(&options.nsfwChecker, "nsfw-checker", false, "Enable Wan 2.7 NSFW checking")
	cmd.Flags().BoolVar(&options.wait, "wait", false, "Poll until Kie reports a terminal state")
	cmd.Flags().DurationVar(&options.waitInterval, "wait-interval", 3*time.Second, "Polling interval used with --wait")
	cmd.Flags().DurationVar(&options.waitTimeout, "wait-timeout", 10*time.Minute, "Maximum wait time")
	return cmd
}

func marketVideoRequest(cmd *cobra.Command, options marketVideoOptions) (map[string]any, error) {
	model := strings.TrimSpace(options.model)
	if model == "" {
		return nil, fmt.Errorf("--model is required")
	}
	if options.callbackURL != "" {
		if err := validateMarketURI("--call-back-url", options.callbackURL); err != nil {
			return nil, err
		}
	}
	body := map[string]any{"model": model}
	if options.callbackURL != "" {
		body["callBackUrl"] = options.callbackURL
	}
	if strings.TrimSpace(options.input) != "" {
		if model == "wan/2-7-text-to-video" {
			return nil, fmt.Errorf("Wan 2.7 uses its validated flags; --input is for another model")
		}
		if strings.TrimSpace(options.prompt) != "" {
			return nil, fmt.Errorf("use either [prompt] with Wan 2.7 flags or --input, not both")
		}
		for _, name := range []string{"negative-prompt", "audio-url", "resolution", "ratio", "duration", "prompt-extend", "watermark", "seed", "nsfw-checker"} {
			if cmd.Flags().Changed(name) {
				return nil, fmt.Errorf("--input cannot be combined with --%s", name)
			}
		}
		var input map[string]any
		if err := json.Unmarshal([]byte(options.input), &input); err != nil {
			return nil, fmt.Errorf("parsing --input JSON: %w", err)
		}
		if input == nil {
			return nil, fmt.Errorf("--input must be a JSON object")
		}
		if err := validateMarketModelInput(model, input); err != nil {
			return nil, err
		}
		body["input"] = input
		return body, nil
	}
	if model != "wan/2-7-text-to-video" {
		return nil, fmt.Errorf("--model %q requires --input; built-in flags only cover wan/2-7-text-to-video", model)
	}
	prompt := strings.TrimSpace(options.prompt)
	if utf8.RuneCountInString(prompt) < 1 || utf8.RuneCountInString(prompt) > 5000 {
		return nil, fmt.Errorf("Wan 2.7 prompt must contain 1-5000 characters")
	}
	if utf8.RuneCountInString(options.negative) > 500 {
		return nil, fmt.Errorf("Wan 2.7 --negative-prompt must contain at most 500 characters")
	}
	if options.duration < 2 || options.duration > 15 {
		return nil, fmt.Errorf("Wan 2.7 --duration must be from 2 through 15")
	}
	if options.seed < 0 || options.seed > math.MaxInt32 {
		return nil, fmt.Errorf("Wan 2.7 --seed must be from 0 through 2147483647")
	}
	if options.resolution != "720p" && options.resolution != "1080p" {
		return nil, fmt.Errorf("Wan 2.7 --resolution must be 720p or 1080p")
	}
	ratioOK := map[string]bool{"16:9": true, "9:16": true, "1:1": true, "4:3": true, "3:4": true}
	if !ratioOK[options.ratio] {
		return nil, fmt.Errorf("Wan 2.7 --ratio must be 16:9, 9:16, 1:1, 4:3, or 3:4")
	}
	if options.audioURL != "" {
		if err := validateMarketURI("--audio-url", options.audioURL); err != nil {
			return nil, err
		}
	}
	input := map[string]any{
		"prompt": prompt, "resolution": options.resolution, "ratio": options.ratio,
		"duration": options.duration, "prompt_extend": options.promptExtend, "watermark": options.watermark,
	}
	if options.negative != "" {
		input["negative_prompt"] = options.negative
	}
	if options.audioURL != "" {
		input["audio_url"] = options.audioURL
	}
	if cmd.Flags().Changed("seed") {
		input["seed"] = options.seed
	}
	if cmd.Flags().Changed("nsfw-checker") {
		input["nsfw_checker"] = options.nsfwChecker
	}
	if err := validateMarketModelInput(model, input); err != nil {
		return nil, err
	}
	body["input"] = input
	return body, nil
}

func validateMarketModelInput(model string, input map[string]any) error {
	issues, err := kiecatalog.Validate(model, input)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		return nil
	}
	data, _ := json.Marshal(issues)
	return fmt.Errorf("model input does not match the documented %s contract: %s", model, data)
}

func validateMarketURI(flag, value string) error {
	parsed, err := url.ParseRequestURI(value)
	if err != nil || !parsed.IsAbs() {
		return fmt.Errorf("%s must be an absolute URI", flag)
	}
	return nil
}

func marketTaskID(data json.RawMessage) string {
	var value any
	if json.Unmarshal(data, &value) != nil {
		return ""
	}
	return findMarketString(value, "taskId", "task_id")
}

func marketTaskState(data json.RawMessage) string {
	var value any
	if json.Unmarshal(data, &value) != nil {
		return ""
	}
	return strings.ToLower(findMarketString(value, "state", "status"))
}

func findMarketString(value any, keys ...string) string {
	wanted := make(map[string]bool, len(keys))
	for _, key := range keys {
		wanted[key] = true
	}
	var find func(any) string
	find = func(value any) string {
		switch value := value.(type) {
		case map[string]any:
			for key, child := range value {
				if wanted[key] {
					if value, ok := child.(string); ok && strings.TrimSpace(value) != "" {
						return value
					}
				}
			}
			for _, child := range value {
				if found := find(child); found != "" {
					return found
				}
			}
		case []any:
			for _, child := range value {
				if found := find(child); found != "" {
					return found
				}
			}
		}
		return ""
	}
	return find(value)
}

func marketResultURLs(data json.RawMessage) []string {
	var value any
	if json.Unmarshal(data, &value) != nil {
		return nil
	}
	seen := map[string]bool{}
	var urls []string
	var walk func(any, string)
	walk = func(value any, parent string) {
		switch value := value.(type) {
		case map[string]any:
			for key, child := range value {
				walk(child, strings.ToLower(key))
			}
		case []any:
			for _, child := range value {
				walk(child, parent)
			}
		case string:
			if parent == "resultjson" {
				var parsed any
				if json.Unmarshal([]byte(value), &parsed) == nil {
					walk(parsed, "result")
				}
				return
			}
			if (strings.Contains(parent, "url") || strings.Contains(parent, "result") || strings.Contains(parent, "output")) &&
				(strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://")) && !seen[value] {
				seen[value] = true
				urls = append(urls, value)
			}
		}
	}
	walk(value, "")
	return urls
}

func waitForMarketTask(cmd *cobra.Command, client interface {
	Get(ctx context.Context, path string, params map[string]string) (json.RawMessage, error)
}, taskID string, interval, timeout time.Duration) (json.RawMessage, error) {
	if interval <= 0 {
		interval = 3 * time.Second
	}
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		data, err := client.Get(ctx, "/api/v1/jobs/recordInfo", map[string]string{"taskId": taskID})
		if err != nil {
			return nil, err
		}
		switch marketTaskState(data) {
		case "success", "fail", "failed":
			return data, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline.C:
			return nil, fmt.Errorf("timed out waiting for task %s; resume with 'kie-pp-cli kie-ai-jobs market-query-task --task-id %s'", taskID, taskID)
		case <-ticker.C:
		}
	}
}
