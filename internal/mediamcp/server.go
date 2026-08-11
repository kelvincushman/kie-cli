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
	"kie-pp-cli/internal/cli"
	"kie-pp-cli/internal/client"
	"kie-pp-cli/internal/config"
	"kie-pp-cli/internal/kiecatalog"
	"kie-pp-cli/internal/media"
)

const defaultRateLimit = 2

var ToolNames = []string{
	"media_brief_start",
	"media_brief_answer",
	"media_brief_get",
	"media_workflow_list",
	"media_workflow_get",
	"media_model_list",
	"media_model_get",
	"media_model_example",
	"media_model_validate",
	"media_reference_add",
	"media_reference_list",
	"media_identity_create",
	"media_identity_get",
	"media_identity_list",
	"media_preview_generate",
	"media_preview_approve",
	"media_preview_reject",
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
	LiveService func(context.Context, *media.Store) (*media.Service, func(), error)
}

type briefStartInput struct {
	Workflow        string   `json:"workflow,omitempty" jsonschema:"optional Kie-native workflow name from media_workflow_list"`
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
	ProductionMode  string   `json:"production_mode,omitempty" jsonschema:"optional single-shot or storyboard video production"`
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
	BriefID string `json:"brief_id" jsonschema:"ready media brief ID approved by the user"`
}

type previewInput struct {
	BriefID string `json:"brief_id" jsonschema:"ready video brief ID"`
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
			Instructions: "Qualify image and video requests one question at a time. Use media_model_list and media_model_get to select a model and inspect every documented input setting, then media_model_validate before a paid generation. For storyboard video, save and explicitly approve a local script and storyboard, then use each returned shot_brief_id. For every video shot, call media_preview_generate, poll with media_generation_status, display the returned preview image URL, and call media_preview_approve only after explicit approval. media_generate rejects video briefs without that approval and rejects storyboard master briefs. Preview and final generation are separate live actions that may consume Kie.ai credits.",
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
		Name: "media_brief_start", Description: "Start a durable local image/video brief. Returns exactly one next question for the agent to ask the user.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefStartInput) (*mcp.CallToolResult, media.Turn, error) {
		brief, err := media.NewBrief(media.BriefInput{
			Workflow: input.Workflow, Request: input.Request, MediaType: input.MediaType, Purpose: input.Purpose,
			Platform: input.Platform, AspectRatio: input.AspectRatio, DurationSeconds: input.DurationSeconds,
			Resolution: input.Resolution, AudioMode: input.AudioMode, VideoMode: input.VideoMode,
			OutputFormat: input.OutputFormat, ReturnLastFrame: input.ReturnLastFrame, WebSearch: input.WebSearch,
			Style: input.Style, References: input.References, ReferenceVideos: input.ReferenceVideos,
			ReferenceAudio: input.ReferenceAudio, FirstFrame: input.FirstFrame, LastFrame: input.LastFrame,
			IdentityIDs: input.IdentityIDs, Model: input.Model,
			ProductionMode: input.ProductionMode,
		})
		if err != nil {
			return nil, media.Turn{}, err
		}
		store, err := deps.Store()
		if err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.VaultBriefReferences(brief); err != nil {
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
		Name: "media_brief_answer", Description: "Answer the current question for a durable media brief. Returns the next single question or a ready generation plan.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input briefAnswerInput) (*mcp.CallToolResult, media.Turn, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Turn{}, err
		}
		if err := media.ApplyNextAnswer(brief, input.Answer); err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.VaultBriefReferences(brief); err != nil {
			return nil, media.Turn{}, err
		}
		if err := store.SaveBrief(brief); err != nil {
			return nil, media.Turn{}, err
		}
		return nil, media.TurnFor(brief), nil
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
		Name: "media_preview_generate", Description: "Generate the required review image for a ready video brief. This is a separate live action that may consume credits. Poll the returned generation, then display its result URL to the user before requesting approval.", Annotations: liveMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input previewInput) (*mcp.CallToolResult, media.Generation, error) {
		store, brief, err := loadBrief(deps, input.BriefID)
		if err != nil {
			return nil, media.Generation{}, err
		}
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.SubmitPreview(ctx, brief)
		if err != nil {
			return nil, media.Generation{}, err
		}
		return nil, *generation, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "media_preview_approve", Description: "Record explicit user approval of the current preview image. Call only after the host has displayed preview_url and the user has affirmatively approved it.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input previewInput) (*mcp.CallToolResult, media.Turn, error) {
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
		Name: "media_preview_reject", Description: "Reject the current video preview and clear its approval state so the brief can be revised and a new preview generated.", Annotations: localMutation,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input previewInput) (*mcp.CallToolResult, media.Turn, error) {
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
		service, cleanup, err := deps.LiveService(ctx, store)
		if err != nil {
			return nil, media.Generation{}, err
		}
		defer cleanup()
		generation, err := service.Submit(ctx, brief)
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
