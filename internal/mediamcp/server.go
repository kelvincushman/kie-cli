// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

// Package mediamcp exposes the durable media-director workflow through the
// official Model Context Protocol Go SDK. Tool calls carry local brief,
// reference, and generation handles, so no conversational server state is
// required between calls.
package mediamcp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"kie-pp-cli/internal/academy"
	"kie-pp-cli/internal/cli"
	"kie-pp-cli/internal/client"
	"kie-pp-cli/internal/cliutil"
	"kie-pp-cli/internal/config"
	"kie-pp-cli/internal/kiecatalog"
	"kie-pp-cli/internal/leaderboard"
	"kie-pp-cli/internal/media"
)

const defaultRateLimit = 2

var ToolNames = []string{
	"media_setup_get",
	"media_grill_start",
	"media_grill_answer",
	"media_grill_wrap_up",
	"media_brief_start",
	"media_brief_answer",
	"media_brief_get",
	"media_workflow_list",
	"media_workflow_get",
	"media_lesson_list",
	"media_lesson_get",
	"media_lesson_recommend",
	"media_leaderboard_get",
	"media_model_list",
	"media_model_get",
	"media_model_example",
	"media_model_validate",
	"media_capability_list",
	"media_capability_get",
	"media_reference_add",
	"media_reference_list",
	"media_identity_create",
	"media_identity_get",
	"media_identity_list",
	"media_paid_confirm",
	"media_preview_generate",
	"media_preview_approve",
	"media_preview_reject",
	"media_proof_generate",
	"media_proof_approve",
	"media_proof_reject",
	"media_proof_skip",
	"media_script_set",
	"media_script_get",
	"media_script_decide",
	"media_storyboard_set",
	"media_storyboard_get",
	"media_storyboard_decide",
	"media_generate",
	"media_generation_status",
}

type Dependencies struct {
	Store       func() (*media.Store, error)
	LoadConfig  func() (*config.Config, error)
	LiveService func(context.Context, *media.Store) (*media.Service, func(), error)
}

type setupOutput struct {
	AuthConfigured      bool   `json:"auth_configured"`
	NextStep            string `json:"next_step"`
	GetAPIKey           string `json:"get_api_key,omitempty"`
	GetAPIKeyLinkType   string `json:"get_api_key_link_type,omitempty"`
	AffiliateDisclosure string `json:"affiliate_disclosure,omitempty"`
}

type briefStartInput struct {
	Workflow        string   `json:"workflow,omitempty" jsonschema:"optional Kie-native workflow name from media_workflow_list"`
	Lesson          string   `json:"lesson,omitempty" jsonschema:"optional Academy course-slug/lesson-slug from media_lesson_recommend"`
	Request         string   `json:"request" jsonschema:"what the user wants to create"`
	MediaType       string   `json:"media_type,omitempty" jsonschema:"optional image or video"`
	Purpose         string   `json:"purpose,omitempty" jsonschema:"optional intended use"`
	Platform        string   `json:"platform,omitempty" jsonschema:"optional destination such as website, Instagram, TikTok, or YouTube"`
	AspectRatio     string   `json:"aspect_ratio,omitempty" jsonschema:"optional 16:9, 9:16, 1:1, 4:3, or 3:4"`
	DurationSeconds int      `json:"duration_seconds,omitempty" jsonschema:"optional duration: 4 to 30 seconds single-shot, up to 600 storyboard, or -1 automatic"`
	Resolution      string   `json:"resolution,omitempty" jsonschema:"optional SeedDance resolution: 480p or 720p"`
	AudioMode       string   `json:"audio_mode,omitempty" jsonschema:"optional SeedDance generated audio: on or off"`
	VideoMode       string   `json:"video_mode,omitempty" jsonschema:"optional SeedDance mode: text, first-frame, first-last, or multimodal"`
	OutputFormat    string   `json:"output_format,omitempty" jsonschema:"optional SeedDance output: mp4 or mov"`
	ReturnLastFrame bool     `json:"return_last_frame,omitempty" jsonschema:"return the generated SeedDance last frame"`
	WebSearch       bool     `json:"web_search,omitempty" jsonschema:"allow SeedDance to use web search"`
	Style           string   `json:"style,omitempty" jsonschema:"optional visual style and mood"`
	References      []string `json:"references,omitempty" jsonschema:"optional local image paths, URLs, or ref:<id> handles"`
	ReferenceVideos []string `json:"reference_videos,omitempty" jsonschema:"optional SeedDance video paths, URLs, or ref:<id> handles"`
	ReferenceAudio  []string `json:"reference_audio,omitempty" jsonschema:"optional SeedDance audio paths, URLs, or ref:<id> handles"`
	FirstFrame      string   `json:"first_frame,omitempty" jsonschema:"optional SeedDance first frame path, URL, or ref:<id>"`
	LastFrame       string   `json:"last_frame,omitempty" jsonschema:"optional SeedDance last frame path, URL, or ref:<id>"`
	IdentityIDs     []string `json:"identity_ids,omitempty" jsonschema:"optional local identity:<id> likeness bundles"`
	Model           string   `json:"model,omitempty" jsonschema:"optional explicit Kie.ai model override"`
	PreviewModel    string   `json:"preview_model,omitempty" jsonschema:"optional still gate model: gpt-image-2-text-to-image, gpt-image-2-image-to-image, nano-banana-2, nano-banana-2-lite, or nano-banana-pro"`
	ProductionMode  string   `json:"production_mode,omitempty" jsonschema:"optional single-shot or storyboard video production"`
	RightsConfirmed bool     `json:"rights_confirmed,omitempty" jsonschema:"confirm rights and consent for likeness or voice generation"`
}

type briefAnswerInput struct {
	BriefID string `json:"brief_id" jsonschema:"brief ID returned by media_brief_start"`
	Answer  string `json:"answer" jsonschema:"the user's answer; use skip when the optional reference step is complete"`
}

type briefGetInput struct {
	BriefID string `json:"brief_id" jsonschema:"durable media brief ID"`
}

type referenceAddInput struct {
	Source string `json:"source" jsonschema:"local supported image, video, or audio path, or an http(s) URL"`
	Name   string `json:"name,omitempty" jsonschema:"optional friendly reusable name"`
	Type   string `json:"type,omitempty" jsonschema:"optional image, video, or audio; local files are auto-detected"`
}

type referenceAddOutput struct {
	Reference media.PublicReference `json:"reference"`
}

type emptyInput struct{}

type workflowGetInput struct {
	Workflow string `json:"workflow" jsonschema:"workflow name returned by media_workflow_list"`
}

type workflowListOutput struct {
	Workflows []media.Workflow `json:"workflows"`
}

type workflowGetOutput struct {
	Workflow media.Workflow `json:"workflow"`
}

type lessonListInput struct {
	Query  string `json:"query,omitempty" jsonschema:"optional request or topic to match"`
	Course string `json:"course,omitempty" jsonschema:"optional exact Academy course slug"`
	Stage  string `json:"stage,omitempty" jsonschema:"optional production stage such as script, asset-lock, storyboard, keyframe, motion, post-production, or review"`
	Limit  int    `json:"limit,omitempty" jsonschema:"maximum results; defaults to 10"`
}

type lessonGetInput struct {
	Key string `json:"key" jsonschema:"exact course-slug/lesson-slug returned by media_lesson_list or media_lesson_recommend"`
}

type lessonListOutput struct {
	Recommendations []academy.Recommendation `json:"recommendations"`
}

type lessonOutput struct {
	Lesson academy.Lesson `json:"lesson"`
}

type leaderboardGetInput struct {
	Task string `json:"task,omitempty" jsonschema:"optional text-to-image, image-edit, character-consistency, or text-to-video"`
}

type leaderboardOutput struct {
	Ledger *leaderboard.Ledger `json:"ledger,omitempty"`
	Task   *leaderboard.Task   `json:"task,omitempty"`
}

type modelListInput struct {
	Search   string `json:"search,omitempty" jsonschema:"optional substring filter over model IDs, names, categories, and descriptions"`
	Category string `json:"category,omitempty" jsonschema:"optional category substring filter"`
}

type modelGetInput struct {
	Model string `json:"model" jsonschema:"exact Kie Market model ID returned by media_model_list"`
}

type modelListOutput struct {
	Models []kiecatalog.Summary `json:"models"`
}

type modelOutput struct {
	Model kiecatalog.Model `json:"model"`
}

type modelExampleOutput struct {
	Model string         `json:"model"`
	Input map[string]any `json:"input"`
}

type modelValidateInput struct {
	Model string         `json:"model" jsonschema:"exact Kie Market model ID"`
	Input map[string]any `json:"input" jsonschema:"model-specific generation input object"`
}

type modelValidationOutput struct {
	Model  string                       `json:"model"`
	Valid  bool                         `json:"valid"`
	Issues []kiecatalog.ValidationIssue `json:"issues"`
}

type capabilityListInput struct {
	Capability string `json:"capability,omitempty" jsonschema:"optional kie-image, kie-video, kie-audio, kie-avatar, or kie-identity filter"`
	Model      string `json:"model,omitempty" jsonschema:"optional exact model ID filter"`
}

type capabilityListOutput struct {
	Models []kiecatalog.ModelCapability `json:"models"`
}

type capabilityOutput struct {
	Model kiecatalog.ModelCapability `json:"model"`
}

type referenceListOutput struct {
	References []media.PublicReference `json:"references"`
}

type identityCreateInput struct {
	Name       string   `json:"name" jsonschema:"friendly local identity name"`
	References []string `json:"references" jsonschema:"1 to 20 image paths, URLs, or ref:<id> handles"`
	Consent    bool     `json:"consent" jsonschema:"confirm permission to use these likeness references"`
}

type identityGetInput struct {
	IdentityID string `json:"identity_id" jsonschema:"identity ID returned by media_identity_create"`
}

type identityOutput struct {
	Identity media.IdentityProfile `json:"identity"`
}

type identityListOutput struct {
	Identities []media.IdentityProfile `json:"identities"`
}

type generateInput struct {
	BriefID        string `json:"brief_id" jsonschema:"ready media brief ID approved by the user"`
	ConfirmationID string `json:"confirmation_id" jsonschema:"fresh scoped ID returned by media_paid_confirm"`
}

type liveGenerateInput struct {
	BriefID        string `json:"brief_id" jsonschema:"ready video brief ID"`
	ConfirmationID string `json:"confirmation_id" jsonschema:"fresh scoped ID returned by media_paid_confirm"`
}

type paidConfirmInput struct {
	BriefID                string `json:"brief_id" jsonschema:"durable ready brief ID"`
	Scope                  string `json:"scope" jsonschema:"preview, proof, or final"`
	GenerationKind         string `json:"generation_kind" jsonschema:"preview, proof, or final"`
	ExpectedModel          string `json:"expected_model" jsonschema:"model the user reviewed"`
	ExpectedPlanHash       string `json:"expected_plan_hash" jsonschema:"plan hash the user reviewed"`
	DisclosureAcknowledged bool   `json:"disclosure_acknowledged" jsonschema:"must be true after the user accepts that this live action may consume credits"`
}

type paidConfirmOutput struct {
	ConfirmationID string    `json:"confirmation_id"`
	ExpiresAt      time.Time `json:"expires_at"`
	NextAction     string    `json:"next_action"`
}

type generationStatusInput struct {
	GenerationID string `json:"generation_id" jsonschema:"generation ID returned by media_preview_generate or media_generate"`
}

type scriptSetInput struct {
	BriefID string `json:"brief_id" jsonschema:"master storyboard video brief ID"`
	Title   string `json:"title,omitempty" jsonschema:"optional script title"`
	Logline string `json:"logline,omitempty" jsonschema:"optional one-line story summary"`
	Content string `json:"content" jsonschema:"complete local script content"`
}

type productionDecisionInput struct {
	BriefID  string `json:"brief_id" jsonschema:"master storyboard video brief ID"`
	Decision string `json:"decision" jsonschema:"approve or reject after the user has reviewed the local artifact"`
}

type scriptOutput struct {
	Script     media.Script `json:"script"`
	NextAction string       `json:"next_action"`
}

type scriptStateOutput struct {
	ScriptID   string `json:"script_id"`
	BriefID    string `json:"brief_id"`
	Hash       string `json:"hash"`
	Status     string `json:"status"`
	NextAction string `json:"next_action"`
}

type storyboardSetInput struct {
	BriefID string                      `json:"brief_id" jsonschema:"master storyboard video brief ID with an approved script"`
	Title   string                      `json:"title,omitempty" jsonschema:"optional storyboard title"`
	Shots   []media.StoryboardShotInput `json:"shots" jsonschema:"ordered 4 to 30 second shots; durations must total the master brief duration"`
}

type storyboardOutput struct {
	View media.StoryboardView `json:"view"`
}

func NewServer(version string, dependencies *Dependencies) *mcp.Server {
	deps := dependenciesWithDefaults(dependencies)
	server := mcp.NewServer(
		&mcp.Implementation{Name: "kie-media-mcp", Version: version},
		&mcp.ServerOptions{
			Instructions: "Start with media_grill_start so the shared director infers explicit prompt facts and asks only one material question at a time. Show its recommendation reason, selected production/capability skills, model rationale, settings, cost status, and overrides. Use media_grill_wrap_up when the user asks for sensible defaults. If authentication is missing, call media_setup_get and show its get_api_key URL with its affiliate disclosure. For video, generate and display the mandatory still preview, then record the user's decision. Offer the optional paid complete-shot proof at the selected model's lowest documented tier; approval or skip never authorizes final generation. Every live director preview, proof, or final needs a fresh scoped media_paid_confirm result and confirmation_id. Never infer approval or paid confirmation from earlier text. All state is local and addressed by explicit IDs; the server is stateless MCP 2026-07-28.",
			Capabilities: &mcp.ServerCapabilities{},
		},
	)
	registerTools(server, deps)
	return server
}

func NewHTTPHandler(server *mcp.Server) http.Handler {
	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{
			Stateless:                    true,
			JSONResponse:                 true,
			PropagateRequestCancellation: true,
		},
	)
	return http.NewCrossOriginProtection().Handler(handler)
}

func dependenciesWithDefaults(input *Dependencies) Dependencies {
	deps := Dependencies{}
	if input != nil {
		deps = *input
	}
	if deps.Store == nil {
		deps.Store = media.DefaultStore
	}
	if deps.LoadConfig == nil {
		deps.LoadConfig = func() (*config.Config, error) { return config.Load("") }
	}
	if deps.LiveService == nil {
		deps.LiveService = newLiveService
	}
	return deps
}

func registerTools(server *mcp.Server, deps Dependencies) {
	localMutation := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPointer(false), OpenWorldHint: boolPointer(false)}
	readLocal := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, DestructiveHint: boolPointer(false), OpenWorldHint: boolPointer(false)}
	liveMutation := &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPointer(false), OpenWorldHint: boolPointer(true)}
	readRemote := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, DestructiveHint: boolPointer(false), OpenWorldHint: boolPointer(true)}

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_setup_get", Description: "Check local Kie authentication without exposing credentials. When setup is required, returns the maintainer referral URL and its required affiliate disclosure.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, setupOutput, error) {
		cfg, err := deps.LoadConfig()
		if err != nil {
			return nil, setupOutput{}, fmt.Errorf("loading Kie configuration: %w", err)
		}
		if cfg.CredentialConfigured() {
			return nil, setupOutput{
				AuthConfigured: true,
				NextStep:       "Start or resume a media brief",
			}, nil
		}
		return nil, setupOutput{
			AuthConfigured:      false,
			NextStep:            "Run kie-pp-cli auth setup in an interactive terminal",
			GetAPIKey:           cliutil.KieAPIKeyURL,
			GetAPIKeyLinkType:   "affiliate",
			AffiliateDisclosure: cliutil.KieAffiliateDisclosure,
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_brief_start", Description: "Start a durable local image/video brief. Returns exactly one next question for the agent to ask the user.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefStartInput) (*mcp.CallToolResult, media.Turn, error) {
		turn, err := startBrief(deps, input, false)
		return nil, turn, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_grill_start", Description: "Start the concise shared media interview. Infers only explicit request facts, returns one material question with rationale, or a ready overridable route.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefStartInput) (*mcp.CallToolResult, media.Turn, error) {
		turn, err := startBrief(deps, input, true)
		return nil, turn, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_grill_answer", Description: "Answer exactly one material grill question and return the next question or ready route.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefAnswerInput) (*mcp.CallToolResult, media.Turn, error) {
		turn, err := answerBrief(deps, input)
		return nil, turn, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_grill_wrap_up", Description: "Apply visible sensible defaults to remaining questions and return an inspectable plan without making a live call.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Turn{}, err
		}
		if err := media.WrapUpBrief(brief); err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.SaveBrief(brief); err != nil {
			return nil, media.Turn{}, err
		}
		return nil, media.TurnFor(brief), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_workflow_list", Description: "List compact Kie-native skill workflow metadata so agents can route media work without loading every skill document.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, workflowListOutput, error) {
		return nil, workflowListOutput{Workflows: media.ListWorkflows()}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_workflow_get", Description: "Get the stages, supported media, skill name, and honest capability gaps for one Kie-native workflow.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input workflowGetInput) (*mcp.CallToolResult, workflowGetOutput, error) {
		workflow, err := media.GetWorkflow(input.Workflow)
		if err != nil {
			return nil, workflowGetOutput{}, err
		}
		return nil, workflowGetOutput{Workflow: workflow}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_lesson_list", Description: "Search the checked-in map of public Academy lesson titles and original Kie-native production methods. No copied lesson scripts or prompts are stored.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input lessonListInput) (*mcp.CallToolResult, lessonListOutput, error) {
		results, err := academy.Search(input.Query, input.Course, input.Stage, input.Limit)
		if err != nil {
			return nil, lessonListOutput{}, err
		}
		return nil, lessonListOutput{Recommendations: results}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_lesson_get", Description: "Get one public lesson link plus its original Kie method, prompt focus, candidate models, and capability boundary.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input lessonGetInput) (*mcp.CallToolResult, lessonOutput, error) {
		lesson, err := academy.GetLesson(input.Key)
		if err != nil {
			return nil, lessonOutput{}, err
		}
		return nil, lessonOutput{Lesson: *lesson}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_lesson_recommend", Description: "Recommend source-linked production lessons for what the user wants to create. Use the selected key in media_brief_start.lesson with workflow academy.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input lessonListInput) (*mcp.CallToolResult, lessonListOutput, error) {
		results, err := academy.Recommend(input.Query, input.Limit)
		if err != nil {
			return nil, lessonListOutput{}, err
		}
		return nil, lessonListOutput{Recommendations: results}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_leaderboard_get", Description: "Inspect the dated, task-specific external model evidence ledger and Kie route availability. Scores from different sources are never combined.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input leaderboardGetInput) (*mcp.CallToolResult, leaderboardOutput, error) {
		if strings.TrimSpace(input.Task) == "" {
			ledger, err := leaderboard.Load()
			if err != nil {
				return nil, leaderboardOutput{}, err
			}
			return nil, leaderboardOutput{Ledger: ledger}, nil
		}
		task, err := leaderboard.GetTask(input.Task)
		if err != nil {
			return nil, leaderboardOutput{}, err
		}
		return nil, leaderboardOutput{Task: task}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_model_list", Description: "List compact summaries for every embedded Kie Market model, including documented input field names. This is local, token-efficient, and spends no credits.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input modelListInput) (*mcp.CallToolResult, modelListOutput, error) {
		models, err := kiecatalog.List(input.Search, input.Category)
		if err != nil {
			return nil, modelListOutput{}, err
		}
		return nil, modelListOutput{Models: models}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_model_get", Description: "Get one model's complete official request schema, input settings, required fields, enums, defaults, limits, examples, and source page.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input modelGetInput) (*mcp.CallToolResult, modelOutput, error) {
		model, err := kiecatalog.Get(input.Model)
		if err != nil {
			return nil, modelOutput{}, err
		}
		return nil, modelOutput{Model: *model}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_model_example", Description: "Return a documented starter input object for one Kie Market model without making a network request.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input modelGetInput) (*mcp.CallToolResult, modelExampleOutput, error) {
		example, err := kiecatalog.Example(input.Model)
		if err != nil {
			return nil, modelExampleOutput{}, err
		}
		return nil, modelExampleOutput{Model: input.Model, Input: example}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_model_validate", Description: "Validate a proposed generation input locally against the selected model's documented required fields, types, enums, and ranges before spending credits.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input modelValidateInput) (*mcp.CallToolResult, modelValidationOutput, error) {
		issues, err := kiecatalog.Validate(input.Model, input.Input)
		if err != nil {
			return nil, modelValidationOutput{}, err
		}
		return nil, modelValidationOutput{Model: input.Model, Valid: len(issues) == 0, Issues: issues}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_capability_list", Description: "List compact, locally classified model routes by capability without loading full model schemas or spending credits.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input capabilityListInput) (*mcp.CallToolResult, capabilityListOutput, error) {
		registry, err := kiecatalog.LoadCapabilities()
		if err != nil {
			return nil, capabilityListOutput{}, err
		}
		modelFilter := strings.TrimSpace(input.Model)
		capabilityFilter := strings.TrimSpace(input.Capability)
		models := make([]kiecatalog.ModelCapability, 0, len(registry.Models))
		for _, item := range registry.Models {
			if modelFilter != "" && item.ModelID != modelFilter {
				continue
			}
			if capabilityFilter != "" && item.PrimaryCapability != capabilityFilter {
				continue
			}
			models = append(models, item)
		}
		return nil, capabilityListOutput{Models: models}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_capability_get", Description: "Get one model's primary capability, production fit, proof settings, and routing note. Use media_model_get for its full settings schema.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input modelGetInput) (*mcp.CallToolResult, capabilityOutput, error) {
		item, err := kiecatalog.GetCapability(input.Model)
		if err != nil {
			return nil, capabilityOutput{}, err
		}
		return nil, capabilityOutput{Model: *item}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_brief_answer", Description: "Answer the current question for a durable media brief. Returns the next single question or a ready generation plan.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefAnswerInput) (*mcp.CallToolResult, media.Turn, error) {
		turn, err := answerBrief(deps, input)
		return nil, turn, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_brief_get", Description: "Load a durable local media brief and its current next action.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		_, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Turn{}, err
		}
		return nil, media.TurnFor(brief), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_script_set", Description: "Save or revise a storyboard video's script in local private state. This never calls Kie.ai; the returned next action requires separate user review and approval.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input scriptSetInput) (*mcp.CallToolResult, scriptStateOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, scriptStateOutput{}, err
		}
		script, err := store.SetScript(input.BriefID, media.ScriptInput{Title: input.Title, Logline: input.Logline, Content: input.Content})
		if err != nil {
			return nil, scriptStateOutput{}, err
		}
		return nil, scriptStateOutput{ScriptID: script.ID, BriefID: script.BriefID, Hash: script.Hash, Status: script.Status, NextAction: "review_script"}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_script_get", Description: "Load the current local script for a storyboard video brief.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, scriptOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, scriptOutput{}, err
		}
		script, err := store.GetScript(input.BriefID)
		if err != nil {
			return nil, scriptOutput{}, err
		}
		brief, err := store.GetBrief(input.BriefID)
		if err != nil {
			return nil, scriptOutput{}, err
		}
		return nil, scriptOutput{Script: *script, NextAction: media.TurnFor(brief).NextAction}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_script_decide", Description: "Approve or reject the current local script after the user reviews it. Approval is invalid if the creative brief changed.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input productionDecisionInput) (*mcp.CallToolResult, scriptStateOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, scriptStateOutput{}, err
		}
		script, err := store.DecideScript(input.BriefID, input.Decision)
		if err != nil {
			return nil, scriptStateOutput{}, err
		}
		brief, err := store.GetBrief(input.BriefID)
		if err != nil {
			return nil, scriptStateOutput{}, err
		}
		return nil, scriptStateOutput{ScriptID: script.ID, BriefID: script.BriefID, Hash: script.Hash, Status: script.Status, NextAction: media.TurnFor(brief).NextAction}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_storyboard_set", Description: "Save an ordered local storyboard after script approval and create one normal gated child brief per shot. This never calls Kie.ai.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input storyboardSetInput) (*mcp.CallToolResult, storyboardOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		if _, err := store.SetStoryboard(input.BriefID, media.StoryboardInput{Title: input.Title, Shots: input.Shots}); err != nil {
			return nil, storyboardOutput{}, err
		}
		view, err := store.StoryboardView(input.BriefID)
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		return nil, storyboardOutput{View: *view}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_storyboard_get", Description: "Load a storyboard plus every shot_brief_id, current preview/generation turn, and the next aggregate action.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, storyboardOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		view, err := store.StoryboardView(input.BriefID)
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		return nil, storyboardOutput{View: *view}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_storyboard_decide", Description: "Approve or reject the current local storyboard after user review. Approval unlocks per-shot preview work but makes no Kie.ai call.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input productionDecisionInput) (*mcp.CallToolResult, storyboardOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		if _, err := store.DecideStoryboard(input.BriefID, input.Decision); err != nil {
			return nil, storyboardOutput{}, err
		}
		view, err := store.StoryboardView(input.BriefID)
		if err != nil {
			return nil, storyboardOutput{}, err
		}
		return nil, storyboardOutput{View: *view}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_reference_add", Description: "Vault a reusable local reference image or remember a reference URL. Local files are copied into a private vault and are not uploaded until an explicit live preview or final generation call.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input referenceAddInput) (*mcp.CallToolResult, referenceAddOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, referenceAddOutput{}, err
		}
		ref, err := store.AddReferenceTyped(input.Source, input.Name, input.Type)
		if err != nil {
			return nil, referenceAddOutput{}, err
		}
		return nil, referenceAddOutput{Reference: ref.Public()}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_reference_list", Description: "List reusable local media reference handles without exposing private filesystem paths.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, referenceListOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, referenceListOutput{}, err
		}
		refs, err := store.ListReferences()
		if err != nil {
			return nil, referenceListOutput{}, err
		}
		return nil, referenceListOutput{References: media.PublicReferences(refs)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_identity_create", Description: "Save a consented local likeness-reference bundle. This does not train or upload a biometric model.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input identityCreateInput) (*mcp.CallToolResult, identityOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, identityOutput{}, err
		}
		identity, err := store.CreateIdentity(input.Name, input.References, input.Consent)
		if err != nil {
			return nil, identityOutput{}, err
		}
		return nil, identityOutput{Identity: *identity}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_identity_get", Description: "Load one local identity-reference bundle by durable handle.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input identityGetInput) (*mcp.CallToolResult, identityOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, identityOutput{}, err
		}
		identity, err := store.GetIdentity(input.IdentityID)
		if err != nil {
			return nil, identityOutput{}, err
		}
		return nil, identityOutput{Identity: *identity}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_identity_list", Description: "List local consented likeness-reference bundles.", Annotations: readLocal,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, identityListOutput, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, identityListOutput{}, err
		}
		identities, err := store.ListIdentities()
		if err != nil {
			return nil, identityListOutput{}, err
		}
		return nil, identityListOutput{Identities: identities}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_paid_confirm", Description: "Create a 10-minute, single-use confirmation for exactly one reviewed live preview, proof, or final action. This local tool does not call Kie.ai.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input paidConfirmInput) (*mcp.CallToolResult, paidConfirmOutput, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, paidConfirmOutput{}, err
		}
		turn := media.TurnFor(brief)
		review := turn.PaidAction
		if review == nil {
			return nil, paidConfirmOutput{}, fmt.Errorf("brief %s has no paid action available; complete the current %s step first", brief.ID, turn.NextAction)
		}
		if review.BlockedReason != "" {
			return nil, paidConfirmOutput{}, fmt.Errorf("brief %s paid action is blocked: %s", brief.ID, review.BlockedReason)
		}
		if input.Scope != review.Scope || input.GenerationKind != review.GenerationKind || input.ExpectedModel != review.Model || input.ExpectedPlanHash != review.PlanHash {
			return nil, paidConfirmOutput{}, fmt.Errorf("requested paid action is not the brief's current reviewed action; refresh the turn and confirm its exact paid_action")
		}
		plan, kind, err := media.PaidPlanForScope(brief, review.Scope)
		if err != nil {
			return nil, paidConfirmOutput{}, err
		}
		planHash, err := media.PlanFingerprint(plan)
		if err != nil {
			return nil, paidConfirmOutput{}, err
		}
		if kind != review.GenerationKind || plan.Model != review.Model || planHash != review.PlanHash {
			return nil, paidConfirmOutput{}, fmt.Errorf("reviewed paid action no longer matches the current kind, model, or plan hash")
		}
		confirmation, err := media.NewPaidConfirmation(brief, plan, media.PaidConfirmationRequest{
			Scope: review.Scope, GenerationKind: kind, ConfirmedBy: "mcp-user",
			Acknowledged: input.DisclosureAcknowledged,
		})
		if err != nil {
			return nil, paidConfirmOutput{}, err
		}
		if err := store.SavePaidConfirmation(confirmation); err != nil {
			return nil, paidConfirmOutput{}, err
		}
		return nil, paidConfirmOutput{ConfirmationID: confirmation.ID, ExpiresAt: confirmation.ExpiresAt, NextAction: paidNextAction(review.Scope)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_preview_generate", Description: "Generate the required review image for a ready video brief. This is a separate live action that may consume credits. Poll the returned generation, then display its result URL to the user before requesting approval.", Annotations: liveMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input liveGenerateInput) (*mcp.CallToolResult, media.Generation, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		confirmationID, err := requireMCPPaidConfirmation(input.ConfirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.SubmitPreview(ctx, brief, confirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		return nil, *generation, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_proof_generate", Description: "Generate the optional paid complete intended shot at the selected model's lowest documented faithful tier. Requires a fresh proof confirmation. Poll the returned generation, then display its result URL before requesting a decision.", Annotations: liveMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input liveGenerateInput) (*mcp.CallToolResult, media.Generation, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		confirmationID, err := requireMCPPaidConfirmation(input.ConfirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.SubmitProof(ctx, brief, confirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		return nil, *generation, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_preview_approve", Description: "Record explicit user approval of the current preview image. Call only after the host has displayed preview_url and the user has affirmatively approved it.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Turn{}, err
		}
		if err := media.ApproveVideoPreview(brief); err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.SaveBrief(brief); err != nil {
			return nil, media.Turn{}, err
		}
		return nil, media.TurnFor(brief), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_proof_approve", Description: "Record the user's approval after the host displays the completed proof URL. This local decision never authorizes final spend.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		return decideProof(deps, input.BriefID, "approve")
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_proof_reject", Description: "Reject the current proof and return the brief to proof revision without a live call.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		return decideProof(deps, input.BriefID, "reject")
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_proof_skip", Description: "Skip the optional proof. This local decision does not authorize or submit the final generation.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		return decideProof(deps, input.BriefID, "skip")
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_preview_reject", Description: "Reject the current video preview and clear its approval state so the brief can be revised and a new preview generated.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefGetInput) (*mcp.CallToolResult, media.Turn, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Turn{}, err
		}
		if err := media.RejectVideoPreview(brief); err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.SaveBrief(brief); err != nil {
			return nil, media.Turn{}, err
		}
		return nil, media.TurnFor(brief), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_generate", Description: "Submit a ready brief to Kie.ai. This explicit live action may consume credits. Video is rejected unless the current displayed preview has separate explicit user approval.", Annotations: liveMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input generateInput) (*mcp.CallToolResult, media.Generation, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		confirmationID, err := requireMCPPaidConfirmation(input.ConfirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.Submit(ctx, brief, confirmationID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		return nil, *generation, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_generation_status", Description: "Refresh one Kie.ai media generation and return its status plus any result URLs.", Annotations: readRemote,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input generationStatusInput) (*mcp.CallToolResult, media.Generation, error) {
		store, err := deps.Store()
		if err != nil {
			return nil, media.Generation{}, err
		}
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.RefreshGeneration(ctx, strings.TrimSpace(input.GenerationID))
		if err != nil {
			return nil, media.Generation{}, err
		}
		return nil, *generation, nil
	})
}

func startBrief(deps Dependencies, input briefStartInput, infer bool) (media.Turn, error) {
	brief, err := media.NewBrief(media.BriefInput{
		Workflow: input.Workflow, Lesson: input.Lesson, Request: input.Request, MediaType: input.MediaType, Purpose: input.Purpose,
		Platform: input.Platform, AspectRatio: input.AspectRatio, DurationSeconds: input.DurationSeconds,
		Resolution: input.Resolution, AudioMode: input.AudioMode, VideoMode: input.VideoMode,
		OutputFormat: input.OutputFormat, ReturnLastFrame: input.ReturnLastFrame, WebSearch: input.WebSearch,
		Style: input.Style, References: input.References, ReferenceVideos: input.ReferenceVideos,
		ReferenceAudio: input.ReferenceAudio, FirstFrame: input.FirstFrame, LastFrame: input.LastFrame,
		IdentityIDs: input.IdentityIDs, Model: input.Model, PreviewModel: input.PreviewModel,
		ProductionMode: input.ProductionMode, RightsAcknowledged: input.RightsConfirmed,
	})
	if err != nil {
		return media.Turn{}, err
	}
	if infer {
		if err := media.InferBrief(brief); err != nil {
			return media.Turn{}, err
		}
	}
	store, err := deps.Store()
	if err != nil {
		return media.Turn{}, err
	}
	if err := store.VaultBriefReferences(brief); err != nil {
		return media.Turn{}, err
	}
	if err := store.SaveBrief(brief); err != nil {
		return media.Turn{}, err
	}
	return media.TurnFor(brief), nil
}

func answerBrief(deps Dependencies, input briefAnswerInput) (media.Turn, error) {
	store, brief, err := loadBrief(deps, input.BriefID)
	if err != nil {
		return media.Turn{}, err
	}
	answer := strings.TrimSpace(input.Answer)
	if strings.EqualFold(answer, "wrap up") || strings.EqualFold(answer, "use sensible defaults") {
		err = media.WrapUpBrief(brief)
	} else {
		err = media.ApplyNextAnswer(brief, answer)
	}
	if err != nil {
		return media.Turn{}, err
	}
	if err := store.VaultBriefReferences(brief); err != nil {
		return media.Turn{}, err
	}
	if err := store.SaveBrief(brief); err != nil {
		return media.Turn{}, err
	}
	return media.TurnFor(brief), nil
}

func requireMCPPaidConfirmation(confirmationID string) (string, error) {
	confirmationID = strings.TrimSpace(confirmationID)
	if confirmationID == "" {
		return "", fmt.Errorf("this live action requires a fresh confirmation_id from media_paid_confirm")
	}
	return confirmationID, nil
}

func paidNextAction(scope string) string {
	switch scope {
	case media.PaidScopePreview:
		return "call media_preview_generate with confirmation_id"
	case media.PaidScopeProof:
		return "call media_proof_generate with confirmation_id"
	case media.PaidScopeFinal:
		return "call media_generate with confirmation_id"
	default:
		return "use the confirmation_id for the reviewed paid action"
	}
}

func decideProof(deps Dependencies, briefID, decision string) (*mcp.CallToolResult, media.Turn, error) {
	store, brief, err := loadBrief(deps, briefID)
	if err != nil {
		return nil, media.Turn{}, err
	}
	switch decision {
	case "approve":
		err = media.ApproveVideoProof(brief)
	case "reject":
		err = media.RejectVideoProof(brief)
	case "skip":
		err = media.SkipVideoProof(brief)
	default:
		err = fmt.Errorf("unsupported proof decision %q", decision)
	}
	if err != nil {
		return nil, media.Turn{}, err
	}
	if err := store.SaveBrief(brief); err != nil {
		return nil, media.Turn{}, err
	}
	return nil, media.TurnFor(brief), nil
}

func loadBrief(deps Dependencies, id string) (*media.Store, *media.Brief, error) {
	store, err := deps.Store()
	if err != nil {
		return nil, nil, err
	}
	brief, err := store.GetBrief(strings.TrimSpace(id))
	if err != nil {
		return nil, nil, err
	}
	return store, brief, nil
}

func newLiveService(ctx context.Context, store *media.Store) (*media.Service, func(), error) {
	cfg, err := config.Load("")
	if err != nil {
		return nil, nil, fmt.Errorf("loading Kie configuration: %w", err)
	}
	if !cfg.CredentialConfigured() {
		return nil, nil, fmt.Errorf("Kie.ai authentication is not configured; run kie-pp-cli auth setup")
	}
	c := client.New(cfg, 60*time.Second, defaultRateLimit)
	c.NoCache = true
	session, err := cli.BindMCPClient(ctx, c)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		if session != nil {
			session.ZeroCredentials()
		}
	}
	return &media.Service{API: c, Store: store}, cleanup, nil
}

func boolPointer(value bool) *bool { return &value }
