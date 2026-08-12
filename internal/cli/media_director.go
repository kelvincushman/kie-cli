// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"kie-pp-cli/internal/config"
	"kie-pp-cli/internal/media"
)

func init() {
	registerNovelCommand(func(root *cobra.Command, flags *rootFlags) {
		addNovelCommandIfAbsent(root, newMediaCreateCmd(flags))
		addNovelCommandIfAbsent(root, newMediaCmd(flags))
	})
}

type mediaCreateOptions struct {
	workflow        string
	lesson          string
	briefID         string
	answer          string
	mediaType       string
	purpose         string
	platform        string
	aspectRatio     string
	durationSeconds int
	resolution      string
	audioMode       string
	videoMode       string
	outputFormat    string
	returnLastFrame bool
	webSearch       bool
	style           string
	references      []string
	referenceVideos []string
	referenceAudio  []string
	firstFrame      string
	lastFrame       string
	identityIDs     []string
	model           string
	previewModel    string
	productionMode  string
	planOnly        bool
	preview         bool
	approvePreview  bool
	rejectPreview   bool
	submit          bool
	wait            bool
	waitInterval    time.Duration
	waitTimeout     time.Duration
}

// pp:data-source auto
func newMediaCreateCmd(flags *rootFlags) *cobra.Command {
	var options mediaCreateOptions
	cmd := &cobra.Command{
		Use:   "create [what you want to make]",
		Short: "Create an image or video through a guided, resumable media brief",
		Long: "Qualify a media request one question at a time, preserve the brief locally, and submit it to Kie.ai only after review. Video requires a generated still preview and explicit approval before final submission. " +
			"Agents should use --agent and resume with --brief <id> --answer <value>.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := ""
			if len(args) == 1 {
				request = args[0]
			}
			return runMediaCreate(cmd, flags, options, request)
		},
	}
	cmd.Flags().StringVar(&options.briefID, "brief", "", "Resume an existing local media brief")
	cmd.Flags().StringVar(&options.workflow, "workflow", "", "Kie workflow: generate, academy, brandkit, marketplace-cards, product-photoshoot, identity, video-explainer, websites, or youtube-thumbnail")
	cmd.Flags().StringVar(&options.lesson, "lesson", "", "Academy lesson key from 'lesson recommend', such as course-slug/lesson-slug")
	cmd.Flags().StringVar(&options.answer, "answer", "", "Answer the brief's current question (use with --brief)")
	cmd.Flags().StringVar(&options.mediaType, "type", "", "Media type: image or video")
	cmd.Flags().StringVar(&options.purpose, "purpose", "", "Intended use, such as social post or website hero")
	cmd.Flags().StringVar(&options.platform, "platform", "", "Destination platform, such as website, Instagram, TikTok, or YouTube")
	cmd.Flags().StringVar(&options.aspectRatio, "aspect-ratio", "", "Output aspect ratio: 16:9, 9:16, 1:1, 4:3, or 3:4")
	cmd.Flags().IntVar(&options.durationSeconds, "duration", 0, "Video duration in seconds (4-30 single-shot, up to 600 storyboard, or -1 automatic)")
	cmd.Flags().StringVar(&options.resolution, "resolution", "", "SeedDance resolution: 480p or 720p")
	cmd.Flags().StringVar(&options.audioMode, "audio", "", "SeedDance generated audio: on or off")
	cmd.Flags().StringVar(&options.videoMode, "video-mode", "", "SeedDance mode: text, first-frame, first-last, or multimodal")
	cmd.Flags().StringVar(&options.outputFormat, "output-format", "", "SeedDance output: mp4 or mov")
	cmd.Flags().BoolVar(&options.returnLastFrame, "return-last-frame", false, "Ask SeedDance to return the generated last frame")
	cmd.Flags().BoolVar(&options.webSearch, "web-search", false, "Allow SeedDance web search for current context")
	cmd.Flags().StringVar(&options.style, "style", "", "Visual style, mood, camera, and lighting direction")
	cmd.Flags().StringSliceVar(&options.references, "reference", nil, "Reference image path, URL, or ref:<id>; repeat or comma-separate")
	cmd.Flags().StringSliceVar(&options.referenceVideos, "reference-video", nil, "SeedDance reference video path, URL, or ref:<id>; repeat or comma-separate")
	cmd.Flags().StringSliceVar(&options.referenceAudio, "reference-audio", nil, "SeedDance reference audio path, URL, or ref:<id>; repeat or comma-separate")
	cmd.Flags().StringVar(&options.firstFrame, "first-frame", "", "SeedDance first-frame image path, URL, or ref:<id>")
	cmd.Flags().StringVar(&options.lastFrame, "last-frame", "", "SeedDance last-frame image path, URL, or ref:<id>")
	cmd.Flags().StringSliceVar(&options.identityIDs, "identity", nil, "Local identity:<id> likeness bundle; repeat or comma-separate")
	cmd.Flags().StringVar(&options.model, "model", "", "Override the recommended Kie.ai model")
	cmd.Flags().StringVar(&options.previewModel, "preview-model", "", "Still gate model: GPT Image 2 or a supported Nano Banana route")
	cmd.Flags().StringVar(&options.productionMode, "production-mode", "", "Video production: single-shot or storyboard")
	cmd.Flags().BoolVar(&options.planOnly, "plan-only", false, "Build and save the brief without submitting a live generation")
	cmd.Flags().BoolVar(&options.preview, "preview", false, "Generate the review image required before a video can be submitted")
	cmd.Flags().BoolVar(&options.approvePreview, "approve-preview", false, "Explicitly approve the current video preview after viewing it")
	cmd.Flags().BoolVar(&options.rejectPreview, "reject-preview", false, "Reject the current video preview so it can be revised and regenerated")
	cmd.Flags().BoolVar(&options.submit, "submit", false, "Submit a ready brief to Kie.ai")
	cmd.Flags().BoolVar(&options.wait, "wait", false, "Wait for a submitted generation to finish")
	cmd.Flags().DurationVar(&options.waitInterval, "wait-interval", 3*time.Second, "Polling interval used with --wait")
	cmd.Flags().DurationVar(&options.waitTimeout, "wait-timeout", 10*time.Minute, "Maximum time to wait for a generation")
	return cmd
}

func runMediaCreate(cmd *cobra.Command, flags *rootFlags, options mediaCreateOptions, request string) error {
	actionCount := 0
	for _, active := range []bool{options.preview, options.approvePreview, options.rejectPreview, options.submit} {
		if active {
			actionCount++
		}
	}
	if actionCount > 1 {
		return usageErr(fmt.Errorf("--preview, --approve-preview, --reject-preview, and --submit are separate actions"))
	}
	if options.wait && !options.submit && !options.preview {
		return usageErr(fmt.Errorf("--wait requires --preview or --submit; use 'media generation status <id> --wait' for an existing generation"))
	}
	if (options.preview || options.approvePreview || options.rejectPreview) && strings.TrimSpace(options.briefID) == "" {
		return usageErr(fmt.Errorf("video preview actions require --brief <id>"))
	}
	store, err := media.DefaultStore()
	if err != nil {
		return err
	}
	var brief *media.Brief
	if strings.TrimSpace(options.briefID) != "" {
		brief, err = store.GetBrief(options.briefID)
		if err != nil {
			return fmt.Errorf("loading media brief %q: %w", options.briefID, err)
		}
		if brief.Status == media.StatusSubmitted && (strings.TrimSpace(options.answer) != "" || mediaCreateHasOverrides(options, request)) {
			return fmt.Errorf("brief %s was already submitted as generation %s; create a new brief to change or generate it again", brief.ID, brief.GenerationID)
		}
		if strings.TrimSpace(options.answer) != "" {
			if err := media.ApplyNextAnswer(brief, options.answer); err != nil {
				return err
			}
		}
		applyMediaOverrides(brief, options, request)
		if err := media.ValidateBrief(brief); err != nil {
			return err
		}
		media.Refresh(brief)
	} else {
		if strings.TrimSpace(options.answer) != "" {
			return fmt.Errorf("--answer requires --brief <id>")
		}
		brief, err = media.NewBrief(media.BriefInput{
			Workflow: options.workflow, Lesson: options.lesson, Request: request, MediaType: options.mediaType, Purpose: options.purpose,
			Platform: options.platform, AspectRatio: options.aspectRatio,
			DurationSeconds: options.durationSeconds, Resolution: options.resolution,
			AudioMode: options.audioMode, VideoMode: options.videoMode, OutputFormat: options.outputFormat,
			ReturnLastFrame: options.returnLastFrame, WebSearch: options.webSearch, Style: options.style,
			References: options.references, ReferenceVideos: options.referenceVideos,
			ReferenceAudio: options.referenceAudio, FirstFrame: options.firstFrame,
			LastFrame: options.lastFrame, IdentityIDs: options.identityIDs, Model: options.model,
			PreviewModel: options.previewModel, ProductionMode: options.productionMode,
		})
		if err != nil {
			return err
		}
	}
	if err := store.VaultBriefReferences(brief); err != nil {
		return err
	}
	if options.rejectPreview {
		if err := media.RejectVideoPreview(brief); err != nil {
			return err
		}
	}
	if options.approvePreview {
		if err := media.ApproveVideoPreview(brief); err != nil {
			return err
		}
	}
	if err := store.SaveBrief(brief); err != nil {
		return err
	}

	interactive := !flags.noInput && !flags.agent && isTerminal(cmd.OutOrStdout())
	if interactive {
		if err := runMediaInterview(cmd, store, brief); err != nil {
			return err
		}
	}

	turn := media.TurnFor(brief)
	if !turn.Ready || options.planOnly {
		if options.submit && brief.Status == media.StatusSubmitted {
			return fmt.Errorf("brief %s was already submitted as generation %s; check its status instead", brief.ID, brief.GenerationID)
		}
		return printMediaValue(cmd, flags, turn)
	}
	if options.approvePreview || options.rejectPreview {
		return printMediaValue(cmd, flags, turn)
	}

	preview := options.preview
	if interactive && brief.MediaType == "video" && (turn.NextAction == "generate_preview" || turn.NextAction == "regenerate_preview") {
		if err := printMediaPreviewPlan(cmd.OutOrStdout(), brief); err != nil {
			return err
		}
		confirmed, err := confirmMediaPreview(cmd)
		if err != nil {
			return err
		}
		preview = confirmed
	}
	if preview {
		if flags.dryRun {
			return printMediaValue(cmd, flags, map[string]any{
				"brief": brief, "plan": media.BuildPreviewPlan(brief), "kind": media.GenerationKindPreview,
				"dry_run": true, "submitted": false,
			})
		}
		client, err := flags.newClient()
		if err != nil {
			return err
		}
		service := &media.Service{API: client, Store: store}
		generation, err := service.SubmitPreview(cmd.Context(), brief)
		if err != nil {
			return classifyAPIError(err, flags)
		}
		if options.wait || interactive {
			generation, err = waitForMediaGeneration(cmd, service, generation.ID, options.waitInterval, options.waitTimeout)
			if err != nil {
				return err
			}
		}
		if !interactive {
			return printMediaValue(cmd, flags, generation)
		}
		brief, err = store.GetBrief(brief.ID)
		if err != nil {
			return err
		}
		if strings.TrimSpace(brief.PreviewURL) == "" {
			return printMediaValue(cmd, flags, media.TurnFor(brief))
		}
		fmt.Fprintf(cmd.OutOrStdout(), "\nReview this generated video preview:\n%s\n", brief.PreviewURL)
		approved, err := confirmMediaPreviewApproval(cmd)
		if err != nil {
			return err
		}
		if !approved {
			return printMediaValue(cmd, flags, media.TurnFor(brief))
		}
		if err := media.ApproveVideoPreview(brief); err != nil {
			return err
		}
		if err := store.SaveBrief(brief); err != nil {
			return err
		}
		turn = media.TurnFor(brief)
	}
	if brief.MediaType == "video" && !turn.CanSubmit {
		if options.submit {
			return fmt.Errorf("video brief %s cannot be submitted yet; next action is %s", brief.ID, turn.NextAction)
		}
		return printMediaValue(cmd, flags, turn)
	}

	submit := options.submit
	if interactive && !submit {
		if err := printMediaPlan(cmd.OutOrStdout(), brief); err != nil {
			return err
		}
		confirmed, err := confirmMediaSubmit(cmd)
		if err != nil {
			return err
		}
		submit = confirmed
	}
	if !submit {
		return printMediaValue(cmd, flags, turn)
	}
	if flags.dryRun {
		return printMediaValue(cmd, flags, map[string]any{
			"brief": brief, "plan": brief.Plan, "dry_run": true, "submitted": false,
		})
	}
	client, err := flags.newClient()
	if err != nil {
		return err
	}
	service := &media.Service{API: client, Store: store}
	generation, err := service.Submit(cmd.Context(), brief)
	if err != nil {
		return classifyAPIError(err, flags)
	}
	if options.wait {
		generation, err = waitForMediaGeneration(cmd, service, generation.ID, options.waitInterval, options.waitTimeout)
		if err != nil {
			return err
		}
	}
	return printMediaValue(cmd, flags, generation)
}

func mediaCreateHasOverrides(options mediaCreateOptions, request string) bool {
	return strings.TrimSpace(request) != "" || strings.TrimSpace(options.workflow) != "" || strings.TrimSpace(options.lesson) != "" || strings.TrimSpace(options.mediaType) != "" ||
		strings.TrimSpace(options.purpose) != "" || strings.TrimSpace(options.platform) != "" ||
		strings.TrimSpace(options.aspectRatio) != "" || options.durationSeconds != 0 ||
		strings.TrimSpace(options.resolution) != "" || strings.TrimSpace(options.audioMode) != "" ||
		strings.TrimSpace(options.videoMode) != "" || strings.TrimSpace(options.outputFormat) != "" ||
		options.returnLastFrame || options.webSearch || strings.TrimSpace(options.style) != "" ||
		len(options.references)+len(options.referenceVideos)+len(options.referenceAudio)+len(options.identityIDs) > 0 ||
		strings.TrimSpace(options.firstFrame) != "" || strings.TrimSpace(options.lastFrame) != "" || strings.TrimSpace(options.model) != "" ||
		strings.TrimSpace(options.previewModel) != "" ||
		strings.TrimSpace(options.productionMode) != ""
}

func applyMediaOverrides(brief *media.Brief, options mediaCreateOptions, request string) {
	if value := strings.TrimSpace(options.workflow); value != "" {
		if workflow, err := media.GetWorkflow(value); err == nil {
			brief.Workflow = workflow.Name
		} else {
			brief.Workflow = value
		}
	}
	if value := strings.Trim(strings.TrimSpace(options.lesson), "/"); value != "" {
		brief.Lesson = value
	}
	if value := strings.TrimSpace(request); value != "" {
		brief.Request = value
	}
	if value := strings.TrimSpace(options.mediaType); value != "" {
		brief.MediaType = strings.ToLower(value)
	}
	if value := strings.TrimSpace(options.purpose); value != "" {
		brief.Purpose = value
	}
	if value := strings.TrimSpace(options.platform); value != "" {
		brief.Platform = strings.ToLower(value)
	}
	if value := strings.TrimSpace(options.aspectRatio); value != "" {
		brief.AspectRatio = value
	}
	if options.durationSeconds != 0 {
		brief.DurationSeconds = options.durationSeconds
	}
	if value := strings.TrimSpace(options.resolution); value != "" {
		brief.Resolution = strings.ToLower(value)
	}
	if value := strings.TrimSpace(options.audioMode); value != "" {
		brief.AudioMode = strings.ToLower(value)
	}
	if value := strings.TrimSpace(options.videoMode); value != "" {
		brief.VideoMode = strings.ToLower(value)
	}
	if value := strings.TrimSpace(options.outputFormat); value != "" {
		brief.OutputFormat = strings.ToLower(value)
	}
	if options.returnLastFrame {
		brief.ReturnLastFrame = true
	}
	if options.webSearch {
		brief.WebSearch = true
	}
	if value := strings.TrimSpace(options.style); value != "" {
		brief.Style = value
	}
	if len(options.references) > 0 {
		brief.References = append([]string(nil), options.references...)
		brief.ReferencesComplete = true
	}
	if len(options.referenceVideos) > 0 {
		brief.ReferenceVideos = append([]string(nil), options.referenceVideos...)
		brief.ReferencesComplete = true
	}
	if len(options.referenceAudio) > 0 {
		brief.ReferenceAudio = append([]string(nil), options.referenceAudio...)
		brief.ReferencesComplete = true
	}
	if value := strings.TrimSpace(options.firstFrame); value != "" {
		brief.FirstFrame = value
	}
	if value := strings.TrimSpace(options.lastFrame); value != "" {
		brief.LastFrame = value
	}
	if len(options.identityIDs) > 0 {
		brief.IdentityIDs = append([]string(nil), options.identityIDs...)
		brief.IdentityComplete = true
	}
	if value := strings.TrimSpace(options.model); value != "" {
		brief.Model = value
	}
	if value := strings.TrimSpace(options.previewModel); value != "" {
		brief.PreviewModel = value
	}
	if value := strings.TrimSpace(options.productionMode); value != "" {
		brief.ProductionMode = strings.ToLower(value)
	}
}

func runMediaInterview(cmd *cobra.Command, store *media.Store, brief *media.Brief) error {
	reader := bufio.NewReader(cmd.InOrStdin())
	for question := media.NextQuestion(brief); question != nil; question = media.NextQuestion(brief) {
		fmt.Fprintln(cmd.OutOrStdout())
		fmt.Fprintln(cmd.OutOrStdout(), question.Prompt)
		if len(question.Options) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Options: %s\n", strings.Join(question.Options, ", "))
		}
		if question.Recommendation != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Recommendation: %s (press Enter to accept)\n", question.Recommendation)
		}
		fmt.Fprint(cmd.OutOrStdout(), "> ")
		answer, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if errors.Is(err, io.EOF) && strings.TrimSpace(answer) == "" {
			return fmt.Errorf("interactive input ended before brief %s was ready", brief.ID)
		}
		if err := media.ApplyNextAnswer(brief, answer); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Please try again: %v\n", err)
			continue
		}
		if err := store.VaultBriefReferences(brief); err != nil {
			return err
		}
		if err := store.SaveBrief(brief); err != nil {
			return err
		}
	}
	return nil
}

func confirmMediaSubmit(cmd *cobra.Command) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Submit this live Kie.ai generation now? [y/N] ")
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func confirmMediaPreview(cmd *cobra.Command) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Generate the paid review image before creating the video? [y/N] ")
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func confirmMediaPreviewApproval(cmd *cobra.Command) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Approve this image as the visual anchor for the final video? [y/N] ")
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

func printMediaPlan(w io.Writer, brief *media.Brief) error {
	data, err := json.MarshalIndent(brief.Plan, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\nReady brief %s\n%s\n", brief.ID, data)
	return err
}

func printMediaPreviewPlan(w io.Writer, brief *media.Brief) error {
	data, err := json.MarshalIndent(media.BuildPreviewPlan(brief), "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\nVideo preview plan for brief %s\n%s\n", brief.ID, data)
	return err
}

func printMediaValue(cmd *cobra.Command, flags *rootFlags, value any) error {
	if flags.asJSON || flags.agent || !isTerminal(cmd.OutOrStdout()) {
		return printJSONFiltered(cmd.OutOrStdout(), value, flags)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return err
}

func waitForMediaGeneration(cmd *cobra.Command, service *media.Service, id string, interval, timeout time.Duration) (*media.Generation, error) {
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
	for {
		generation, err := service.RefreshGeneration(cmd.Context(), id)
		if err != nil {
			return nil, err
		}
		if mediaGenerationTerminal(generation.Status) {
			return generation, nil
		}
		select {
		case <-cmd.Context().Done():
			return nil, cmd.Context().Err()
		case <-deadline.C:
			return generation, fmt.Errorf("timed out waiting for generation %s; resume with 'kie-pp-cli media generation status %s --wait'", id, id)
		case <-ticker.C:
		}
	}
}

func mediaGenerationTerminal(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "success", "succeeded", "completed", "complete", "fail", "failed", "error", "cancelled", "canceled":
		return true
	default:
		return false
	}
}

// pp:data-source computed
func newMediaCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{Use: "media", Short: "Inspect and resume local media briefs, references, and generations"}
	cmd.AddCommand(newMediaSetupCmd(flags))
	cmd.AddCommand(newMediaWorkflowCmd(flags))
	cmd.AddCommand(newMediaLeaderboardCmd(flags))
	cmd.AddCommand(newMediaBriefCmd(flags))
	cmd.AddCommand(newMediaReferenceCmd(flags))
	cmd.AddCommand(newMediaIdentityCmd(flags))
	cmd.AddCommand(newMediaScriptCmd(flags))
	cmd.AddCommand(newMediaStoryboardCmd(flags))
	cmd.AddCommand(newMediaGenerationCmd(flags))
	cmd.AddCommand(newMediaModelsCmd(flags))
	cmd.AddCommand(newMediaVideoCmd(flags))
	return cmd
}

// pp:data-source local
func newMediaScriptCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "script", Short: "Write and approve a durable local video script"}
	var file, title, logline string
	set := &cobra.Command{
		Use:   "set <brief-id>",
		Short: "Save script text locally from a file or stdin; no Kie.ai call is made",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readMediaInput(cmd, file)
			if err != nil {
				return err
			}
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			script, err := store.SetScript(args[0], media.ScriptInput{Title: title, Logline: logline, Content: string(content)})
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, map[string]any{
				"script_id": script.ID, "brief_id": script.BriefID, "hash": script.Hash,
				"status": script.Status, "next_action": "review_script",
			})
		},
	}
	set.Flags().StringVar(&file, "file", "", "Script file path, or - to read from stdin")
	set.Flags().StringVar(&title, "title", "", "Script title")
	set.Flags().StringVar(&logline, "logline", "", "One-line story summary")
	parent.AddCommand(set)
	parent.AddCommand(&cobra.Command{Use: "show <brief-id>", Short: "Show the complete local script for review", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		script, err := store.GetScript(args[0])
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, script)
	}})
	for _, decision := range []string{"approve", "reject"} {
		decision := decision
		parent.AddCommand(&cobra.Command{Use: decision + " <brief-id>", Short: strings.ToUpper(decision[:1]) + decision[1:] + " the current reviewed script locally", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			script, err := store.DecideScript(args[0], decision)
			if err != nil {
				return err
			}
			brief, err := store.GetBrief(args[0])
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, map[string]any{
				"script_id": script.ID, "status": script.Status, "next_action": media.TurnFor(brief).NextAction,
			})
		}})
	}
	return parent
}

// pp:data-source local
func newMediaStoryboardCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "storyboard", Short: "Plan and approve local multi-shot video production"}
	var file string
	set := &cobra.Command{
		Use:   "set <brief-id>",
		Short: "Save storyboard JSON and create one gated child brief per shot; no Kie.ai call is made",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := readMediaInput(cmd, file)
			if err != nil {
				return err
			}
			var input media.StoryboardInput
			if err := json.Unmarshal(data, &input); err != nil {
				return fmt.Errorf("decoding storyboard JSON: %w", err)
			}
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			storyboard, err := store.SetStoryboard(args[0], input)
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, storyboard)
		},
	}
	set.Flags().StringVar(&file, "file", "", "Storyboard JSON file path, or - to read from stdin")
	parent.AddCommand(set)
	parent.AddCommand(&cobra.Command{Use: "show <brief-id>", Short: "Show the storyboard, shot brief IDs, and aggregate next action", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		view, err := store.StoryboardView(args[0])
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, view)
	}})
	for _, decision := range []string{"approve", "reject"} {
		decision := decision
		parent.AddCommand(&cobra.Command{Use: decision + " <brief-id>", Short: strings.ToUpper(decision[:1]) + decision[1:] + " the current reviewed storyboard locally", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			if _, err := store.DecideStoryboard(args[0], decision); err != nil {
				return err
			}
			view, err := store.StoryboardView(args[0])
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, view)
		}})
	}
	return parent
}

func readMediaInput(cmd *cobra.Command, path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, usageErr(fmt.Errorf("--file <path> is required; use --file - for stdin"))
	}
	if path == "-" {
		return io.ReadAll(cmd.InOrStdin())
	}
	data, err := os.ReadFile(path) // #nosec G304 -- the user explicitly selected the local script/storyboard file.
	if err != nil {
		reason := "could not be read"
		switch {
		case errors.Is(err, os.ErrNotExist):
			reason = "was not found"
		case errors.Is(err, os.ErrPermission):
			reason = "is not readable"
		}
		return nil, fmt.Errorf("selected file %q %s", filepath.Base(path), reason)
	}
	return data, nil
}

// pp:data-source computed
func newMediaWorkflowCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "workflow", Short: "Discover compact Kie-native skill workflow metadata"}
	parent.AddCommand(&cobra.Command{Use: "list", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		return printMediaValue(cmd, flags, media.ListWorkflows())
	}})
	parent.AddCommand(&cobra.Command{Use: "show <workflow>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		workflow, err := media.GetWorkflow(args[0])
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, workflow)
	}})
	return parent
}

// pp:data-source local
func newMediaSetupCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Check local media director setup without exposing credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(flags.configPath)
			if err != nil {
				return err
			}
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			configured := cfg.CredentialConfigured()
			result := map[string]any{
				"auth_configured": configured,
				"media_store":     store.Root(),
				"next_step":       "Run kie-pp-cli create \"your media request\"",
			}
			if !configured {
				result["next_step"] = "Run kie-pp-cli auth setup in an interactive terminal"
			}
			return printMediaValue(cmd, flags, result)
		},
	}
}

// pp:data-source local
func newMediaBriefCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "brief", Short: "Inspect durable media briefs"}
	parent.AddCommand(&cobra.Command{Use: "show <brief-id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		brief, err := store.GetBrief(args[0])
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, media.TurnFor(brief))
	}})
	parent.AddCommand(&cobra.Command{Use: "list", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		briefs, err := store.ListBriefs()
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, briefs)
	}})
	return parent
}

// pp:data-source local
func newMediaReferenceCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "reference", Short: "Manage reusable local image, video, and audio references"}
	var name, mediaType string
	add := &cobra.Command{Use: "add <path-or-url>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		ref, err := store.AddReferenceTyped(args[0], name, mediaType)
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, ref.Public())
	}}
	add.Flags().StringVar(&name, "name", "", "Friendly reusable reference name")
	add.Flags().StringVar(&mediaType, "type", "", "Reference media type: image, video, or audio (auto-detected for local files)")
	parent.AddCommand(add)
	parent.AddCommand(&cobra.Command{Use: "list", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		refs, err := store.ListReferences()
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, media.PublicReferences(refs))
	}})
	return parent
}

// pp:data-source local
func newMediaIdentityCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "identity", Short: "Manage consented local likeness reference bundles"}
	var references []string
	var consent bool
	create := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a local identity bundle without training or uploading a biometric model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := media.DefaultStore()
			if err != nil {
				return err
			}
			identity, err := store.CreateIdentity(args[0], references, consent)
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, identity)
		},
	}
	create.Flags().StringSliceVar(&references, "reference", nil, "Image path, URL, or ref:<id>; repeat or comma-separate (1-20)")
	create.Flags().BoolVar(&consent, "consent", false, "Confirm permission to use these likeness references")
	parent.AddCommand(create)
	parent.AddCommand(&cobra.Command{Use: "show <identity-id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		identity, err := store.GetIdentity(args[0])
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, identity)
	}})
	parent.AddCommand(&cobra.Command{Use: "list", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		identities, err := store.ListIdentities()
		if err != nil {
			return err
		}
		return printMediaValue(cmd, flags, identities)
	}})
	return parent
}

// pp:data-source live
func newMediaGenerationCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{Use: "generation", Short: "Inspect Kie.ai media generation jobs"}
	var wait bool
	var interval, timeout time.Duration
	status := &cobra.Command{Use: "status <generation-id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		store, err := media.DefaultStore()
		if err != nil {
			return err
		}
		client, err := flags.newClient()
		if err != nil {
			return err
		}
		service := &media.Service{API: client, Store: store}
		var generation *media.Generation
		if wait {
			generation, err = waitForMediaGeneration(cmd, service, args[0], interval, timeout)
		} else {
			generation, err = service.RefreshGeneration(cmd.Context(), args[0])
		}
		if err != nil {
			return classifyAPIError(err, flags)
		}
		return printMediaValue(cmd, flags, generation)
	}}
	status.Flags().BoolVar(&wait, "wait", false, "Wait until the generation reaches a terminal state")
	status.Flags().DurationVar(&interval, "wait-interval", 3*time.Second, "Polling interval")
	status.Flags().DurationVar(&timeout, "wait-timeout", 10*time.Minute, "Maximum wait time")
	parent.AddCommand(status)
	return parent
}
