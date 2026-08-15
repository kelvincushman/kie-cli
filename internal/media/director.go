// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"kie-pp-cli/internal/academy"
	"kie-pp-cli/internal/kiecatalog"
	"kie-pp-cli/internal/leaderboard"
)

type BriefInput struct {
	Workflow           string
	Lesson             string
	Request            string
	MediaType          string
	Purpose            string
	Platform           string
	AspectRatio        string
	DurationSeconds    int
	Resolution         string
	AudioMode          string
	VideoMode          string
	OutputFormat       string
	ReturnLastFrame    bool
	WebSearch          bool
	Style              string
	References         []string
	ReferenceVideos    []string
	ReferenceAudio     []string
	FirstFrame         string
	LastFrame          string
	IdentityIDs        []string
	Model              string
	PreviewModel       string
	ProductionMode     string
	RightsAcknowledged bool
}

func NewBrief(input BriefInput) (*Brief, error) {
	now := time.Now().UTC()
	b := &Brief{
		ID:                 newID("brief"),
		Workflow:           normalizeWorkflow(input.Workflow),
		Lesson:             strings.Trim(strings.TrimSpace(input.Lesson), "/"),
		Request:            strings.TrimSpace(input.Request),
		MediaType:          normalizeMediaType(input.MediaType),
		Purpose:            strings.TrimSpace(input.Purpose),
		Platform:           normalizePlatform(input.Platform),
		AspectRatio:        strings.TrimSpace(input.AspectRatio),
		DurationSeconds:    input.DurationSeconds,
		Resolution:         normalizeResolution(input.Resolution),
		AudioMode:          normalizeAudioMode(input.AudioMode),
		VideoMode:          normalizeVideoMode(input.VideoMode),
		OutputFormat:       normalizeOutputFormat(input.OutputFormat),
		ReturnLastFrame:    input.ReturnLastFrame,
		WebSearch:          input.WebSearch,
		Style:              strings.TrimSpace(input.Style),
		References:         cleanStrings(input.References),
		ReferenceVideos:    cleanStrings(input.ReferenceVideos),
		ReferenceAudio:     cleanStrings(input.ReferenceAudio),
		FirstFrame:         normalizeDroppedReference(input.FirstFrame),
		LastFrame:          normalizeDroppedReference(input.LastFrame),
		IdentityIDs:        cleanStrings(input.IdentityIDs),
		IdentityComplete:   true,
		Model:              strings.TrimSpace(input.Model),
		PreviewModel:       strings.TrimSpace(input.PreviewModel),
		ProductionMode:     normalizeProductionMode(input.ProductionMode),
		RightsAcknowledged: input.RightsAcknowledged,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if b.Workflow == "academy" && b.Lesson == "" && b.Request != "" {
		if recommendations, err := academy.Recommend(b.Request, 1); err == nil && len(recommendations) > 0 {
			b.Lesson = recommendations[0].Lesson.Key
		}
	}
	if b.ProductionMode == "" {
		b.ProductionMode = ProductionModeSingleShot
	}
	if len(input.References)+len(input.ReferenceVideos)+len(input.ReferenceAudio) > 0 {
		b.ReferencesComplete = true
	}
	if b.MediaType == "video" && b.VideoMode == "" {
		switch {
		case b.FirstFrame != "" && b.LastFrame != "":
			b.VideoMode = "first-last"
		case b.FirstFrame != "":
			b.VideoMode = "first-frame"
		case len(b.References)+len(b.ReferenceVideos)+len(b.ReferenceAudio)+len(b.IdentityIDs) > 0:
			b.VideoMode = "multimodal"
		}
	}
	if err := ValidateBrief(b); err != nil {
		return nil, err
	}
	Refresh(b)
	return b, nil
}

func Refresh(b *Brief) {
	if b == nil {
		return
	}
	if b.Status == StatusSubmitted && strings.TrimSpace(b.GenerationID) != "" {
		return
	}
	if b.ProofBriefHash != "" && b.ProofBriefHash != proofBriefFingerprint(b) {
		invalidateProof(b)
	}
	if q := NextQuestion(b); q == nil {
		b.Status = StatusReady
		b.Plan = BuildPlan(b)
	} else if b.Status != StatusSubmitted {
		b.Status = StatusDraft
		b.Plan = nil
	}
	b.UpdatedAt = time.Now().UTC()
}

func TurnFor(b *Brief) Turn {
	if b != nil && b.Status == StatusSubmitted && strings.TrimSpace(b.GenerationID) != "" {
		return newTurn(b, nil, false, false, "check_generation_status")
	}
	q := NextQuestion(b)
	next := "answer_question"
	canSubmit := false
	if q == nil {
		switch {
		case b == nil:
			next = "answer_question"
		case b.MediaType != "video":
			next = "review_then_submit"
			canSubmit = true
		case b.ProductionMode == ProductionModeStoryboard:
			switch {
			case strings.TrimSpace(b.ScriptID) == "":
				next = "draft_script"
			case b.ScriptBriefHash != creativeBriefFingerprint(b):
				next = "revise_script"
			case b.ScriptStatus != StatusApproved:
				next = "review_script"
			case strings.TrimSpace(b.StoryboardID) == "":
				next = "draft_storyboard"
			case b.StoryboardBriefHash != creativeBriefFingerprint(b) || b.StoryboardScriptHash != b.ScriptHash:
				next = "revise_storyboard"
			case b.StoryboardStatus != StatusApproved:
				next = "review_storyboard"
			default:
				next = "generate_shot_previews"
			}
		case VideoPreviewApproved(b):
			switch {
			case b.ProofStatus == StatusSubmitted:
				next = "check_proof_status"
			case generationFailed(b.ProofStatus):
				next = "regenerate_proof"
			case strings.TrimSpace(b.ProofURL) != "" && mediaGenerationComplete(b.ProofStatus) && b.ProofApprovedAt == nil:
				next = "review_proof"
			case VideoProofApproved(b) || VideoProofSkipped(b):
				next = "review_then_submit"
				canSubmit = true
			default:
				option := ResolveProofOption(BuildPlan(b).Model)
				if option.Supported {
					next = "offer_proof"
				} else {
					next = "review_then_submit"
					canSubmit = true
				}
			}
		case !videoPreviewMatchesBrief(b) || strings.TrimSpace(b.PreviewGenerationID) == "":
			next = "generate_preview"
		case generationFailed(b.PreviewStatus):
			next = "regenerate_preview"
		case strings.TrimSpace(b.PreviewURL) != "" && mediaGenerationComplete(b.PreviewStatus):
			next = "review_preview"
		default:
			next = "check_preview_status"
		}
	}
	return newTurn(b, q, q == nil, canSubmit, next)
}

func newTurn(b *Brief, question *Question, ready, canSubmit bool, nextAction string) Turn {
	turn := Turn{Brief: b, NextQuestion: question, PaidAction: paidActionReview(b, nextAction), Ready: ready, CanSubmit: canSubmit, NextAction: nextAction, GateState: nextAction}
	if b != nil && b.ID != "" {
		turn.ResumeCommand = "kie-pp-cli create --brief " + b.ID
	}
	return turn
}

func NextQuestion(b *Brief) *Question {
	if b == nil {
		return nil
	}
	if strings.TrimSpace(b.Request) == "" {
		return &Question{Key: "request", Prompt: "What do you want to create? Describe the subject and desired outcome."}
	}
	if b.MediaType == "" {
		return &Question{Key: "media_type", Prompt: "Should this be an image or a video?", Recommendation: "image", RecommendationReason: "A still image is the lowest-cost starting point when motion is not essential.", Options: []string{"image", "video"}}
	}
	if strings.TrimSpace(b.Purpose) == "" {
		return &Question{Key: "purpose", Prompt: "What will this media be used for?", Recommendation: "social post", RecommendationReason: "This keeps the first draft broadly reusable."}
	}
	if strings.TrimSpace(b.Platform) == "" {
		return &Question{Key: "platform", Prompt: "Where will it be used?", Recommendation: "general", RecommendationReason: "General avoids platform-specific framing when no destination is known.", Options: []string{"general", "website", "instagram", "tiktok", "youtube", "linkedin"}}
	}
	if strings.TrimSpace(b.AspectRatio) == "" {
		recommended := recommendedAspect(b.Platform, b.MediaType)
		return &Question{Key: "aspect_ratio", Prompt: "Which aspect ratio should I use?", Recommendation: recommended, RecommendationReason: "This is the usual framing for the selected destination and media type.", Options: []string{"16:9", "9:16", "1:1", "4:3", "3:4"}}
	}
	if b.MediaType == "video" && b.DurationSeconds == 0 {
		return &Question{Key: "duration_seconds", Prompt: "How many seconds should the video run?", Recommendation: "5", RecommendationReason: "Five seconds is a compact, lower-cost first shot.", Options: []string{"5", "10", "15", "30", "-1 (automatic)"}}
	}
	if b.MediaType == "video" && b.AudioMode == "" {
		return &Question{Key: "audio_mode", Prompt: "Should SeedDance generate synchronized audio?", Recommendation: "off", RecommendationReason: "Audio can increase generation cost and is best enabled only when needed.", Options: []string{"on", "off"}}
	}
	if b.MediaType == "video" && b.VideoMode == "" {
		return &Question{Key: "video_mode", Prompt: "How should SeedDance guide the video?", Recommendation: "text", RecommendationReason: "Text mode is simplest when no reference assets have been supplied.", Options: []string{"text", "first-frame", "first-last", "multimodal"}}
	}
	if strings.TrimSpace(b.Style) == "" {
		return &Question{Key: "style", Prompt: "What visual style or mood should it have?", Recommendation: "cinematic, polished, natural lighting", RecommendationReason: "This is a neutral production-quality default that can be overridden."}
	}
	if requiresRightsAcknowledgement(b) && !b.RightsAcknowledged {
		return &Question{Key: "rights_acknowledged", Prompt: "Do you confirm you have the rights and consent needed for the likeness or voice used in this generation?", Options: []string{"yes", "no"}}
	}
	if b.MediaType == "video" {
		switch b.VideoMode {
		case "text":
			return nil
		case "first-frame":
			if b.FirstFrame == "" {
				return &Question{Key: "first_frame", Prompt: "Add the starting frame image path or URL."}
			}
			return nil
		case "first-last":
			if b.FirstFrame == "" {
				return &Question{Key: "first_frame", Prompt: "Add the starting frame image path or URL."}
			}
			if b.LastFrame == "" {
				return &Question{Key: "last_frame", Prompt: "Add the ending frame image path or URL."}
			}
			return nil
		}
	}
	if !b.ReferencesComplete {
		prompt := "Add a reference image path or URL, or type skip when you are done."
		if b.MediaType == "video" && b.VideoMode == "multimodal" {
			prompt = "Add a reference path or URL. Prefix video: or audio: for those types; unprefixed references are images. Type skip when done."
		}
		return &Question{Key: "reference", Prompt: prompt, Recommendation: "skip", AllowSkip: true}
	}
	return nil
}

// WrapUpBrief applies visible, conservative defaults when the user asks to
// use sensible defaults or stop the interview. It never submits a paid task.
func WrapUpBrief(b *Brief) error {
	if b == nil || strings.TrimSpace(b.Request) == "" {
		return fmt.Errorf("describe what you want to create before using sensible defaults")
	}
	inferBriefFromRequest(b)
	if b.MediaType == "" {
		b.MediaType = "image"
	}
	if b.Purpose == "" {
		b.Purpose = "general creative use"
	}
	if b.Platform == "" {
		b.Platform = "general"
	}
	if b.AspectRatio == "" {
		b.AspectRatio = recommendedAspect(b.Platform, b.MediaType)
	}
	if b.MediaType == "video" {
		if b.DurationSeconds == 0 {
			b.DurationSeconds = 5
		}
		if b.AudioMode == "" {
			b.AudioMode = "off"
		}
		if b.VideoMode == "" {
			b.VideoMode = "text"
		}
	}
	if b.Style == "" {
		b.Style = "cinematic, polished, natural lighting"
	}
	b.ReferencesComplete = true
	if err := ValidateBrief(b); err != nil {
		return err
	}
	Refresh(b)
	return nil
}

func ApplyNextAnswer(b *Brief, answer string) error {
	q := NextQuestion(b)
	if q == nil {
		return fmt.Errorf("brief %s is already ready", b.ID)
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = q.Recommendation
	}
	switch q.Key {
	case "request":
		b.Request = answer
	case "media_type":
		b.MediaType = normalizeMediaType(answer)
	case "purpose":
		b.Purpose = answer
	case "platform":
		b.Platform = normalizePlatform(answer)
	case "aspect_ratio":
		b.AspectRatio = answer
	case "duration_seconds":
		n, err := strconv.Atoi(answer)
		if err != nil {
			return fmt.Errorf("duration must be a whole number of seconds: %w", err)
		}
		b.DurationSeconds = n
	case "audio_mode":
		b.AudioMode = normalizeAudioMode(answer)
	case "video_mode":
		b.VideoMode = normalizeVideoMode(answer)
	case "style":
		b.Style = answer
	case "first_frame":
		b.FirstFrame = normalizeDroppedReference(answer)
	case "last_frame":
		b.LastFrame = normalizeDroppedReference(answer)
	case "reference":
		if isSkipAnswer(answer) {
			b.ReferencesComplete = true
		} else {
			mediaType, source := splitReferenceAnswer(answer)
			switch mediaType {
			case "video":
				b.ReferenceVideos = append(b.ReferenceVideos, source)
			case "audio":
				b.ReferenceAudio = append(b.ReferenceAudio, source)
			default:
				b.References = append(b.References, source)
			}
		}
	case "rights_acknowledged":
		if !strings.EqualFold(answer, "yes") && !strings.EqualFold(answer, "true") && !strings.EqualFold(answer, "confirm") {
			return fmt.Errorf("likeness and voice generation remains blocked until rights and consent are explicitly confirmed")
		}
		b.RightsAcknowledged = true
	default:
		return fmt.Errorf("unsupported question %q", q.Key)
	}
	if err := ValidateBrief(b); err != nil {
		return err
	}
	Refresh(b)
	return nil
}

// InferBrief extracts only facts explicitly present in the request. Unlike
// WrapUpBrief it does not invent missing defaults, so the next material
// question remains visible when the request is incomplete.
func InferBrief(b *Brief) error {
	if b == nil {
		return fmt.Errorf("media brief is required")
	}
	inferBriefFromRequest(b)
	if err := ValidateBrief(b); err != nil {
		return err
	}
	Refresh(b)
	return nil
}

// normalizeDroppedReference handles the forms macOS/Linux terminals commonly
// insert when a file is dragged into an interactive prompt. Flag values have
// already been decoded by the user's shell and pass through unchanged.
func normalizeDroppedReference(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			value = value[1 : len(value)-1]
		}
	}
	return strings.ReplaceAll(value, `\ `, " ")
}

func splitReferenceAnswer(value string) (string, string) {
	value = strings.TrimSpace(value)
	for _, mediaType := range []string{"image", "video", "audio"} {
		prefix := mediaType + ":"
		if strings.HasPrefix(strings.ToLower(value), prefix) && !strings.HasPrefix(strings.ToLower(value), "http") {
			return mediaType, normalizeDroppedReference(strings.TrimSpace(value[len(prefix):]))
		}
	}
	return "image", normalizeDroppedReference(value)
}

func ValidateBrief(b *Brief) error {
	if b.ProductionMode != "" && b.ProductionMode != ProductionModeSingleShot && b.ProductionMode != ProductionModeStoryboard {
		return fmt.Errorf("production mode must be single-shot or storyboard")
	}
	if b.ProductionMode == ProductionModeStoryboard && b.MediaType != "" && b.MediaType != "video" {
		return fmt.Errorf("storyboard production mode applies only to video briefs")
	}
	if b.Workflow != "" {
		workflow, err := GetWorkflow(b.Workflow)
		if err != nil {
			return err
		}
		if b.MediaType != "" && !containsString(workflow.MediaTypes, b.MediaType) {
			return fmt.Errorf("workflow %s does not support media type %s", b.Workflow, b.MediaType)
		}
	}
	if b.Lesson != "" {
		if _, err := academy.GetLesson(b.Lesson); err != nil {
			return err
		}
	}
	if b.MediaType != "" && b.MediaType != "image" && b.MediaType != "video" {
		return fmt.Errorf("media type must be image or video")
	}
	if b.AspectRatio != "" {
		allowed := map[string]bool{"16:9": true, "9:16": true, "1:1": true, "4:3": true, "3:4": true, "21:9": true, "adaptive": true}
		if !allowed[b.AspectRatio] {
			return fmt.Errorf("unsupported aspect ratio %q", b.AspectRatio)
		}
	}
	if b.MediaType == "video" && b.DurationSeconds != 0 && b.DurationSeconds != -1 {
		maxDuration := 30
		if b.ProductionMode == ProductionModeStoryboard {
			maxDuration = 600
		}
		if b.DurationSeconds < 4 || b.DurationSeconds > maxDuration {
			return fmt.Errorf("video duration must be -1 (automatic) or between 4 and %d seconds for %s production", maxDuration, defaultString(b.ProductionMode, ProductionModeSingleShot))
		}
	}
	if b.AudioMode != "" && b.AudioMode != "on" && b.AudioMode != "off" {
		return fmt.Errorf("audio mode must be on or off")
	}
	if b.VideoMode != "" && b.VideoMode != "text" && b.VideoMode != "first-frame" && b.VideoMode != "first-last" && b.VideoMode != "multimodal" {
		return fmt.Errorf("video mode must be text, first-frame, first-last, or multimodal")
	}
	if b.Resolution != "" && b.Resolution != "480p" && b.Resolution != "720p" {
		return fmt.Errorf("resolution must be 480p or 720p")
	}
	if b.OutputFormat != "" && b.OutputFormat != "mp4" && b.OutputFormat != "mov" {
		return fmt.Errorf("output format must be mp4 or mov")
	}
	if len(b.References) > 30 {
		return fmt.Errorf("SeedDance 2.5 accepts at most 30 reference images")
	}
	if b.MediaType == "video" && b.VideoMode == "multimodal" && len(b.References) > 29 {
		return fmt.Errorf("SeedDance 2.5 multimodal video accepts at most 29 user reference images because the approved preview occupies one of 30 image slots")
	}
	if len(b.ReferenceAudio) > 10 {
		return fmt.Errorf("SeedDance 2.5 accepts at most 10 reference audio files")
	}
	if len(b.ReferenceVideos) > 10 {
		return fmt.Errorf("SeedDance 2.5 accepts at most 10 reference video files")
	}
	if b.VideoMode == "text" && (b.FirstFrame != "" || b.LastFrame != "" || len(b.References)+len(b.ReferenceVideos)+len(b.ReferenceAudio)+len(b.IdentityIDs) > 0) {
		return fmt.Errorf("text video mode cannot include frame, identity, or multimodal references")
	}
	if b.VideoMode == "multimodal" && (b.FirstFrame != "" || b.LastFrame != "") {
		return fmt.Errorf("SeedDance multimodal mode cannot include first or last frames; use first-frame or first-last mode")
	}
	if (b.VideoMode == "first-frame" || b.VideoMode == "first-last") && (len(b.References)+len(b.ReferenceVideos)+len(b.ReferenceAudio)+len(b.IdentityIDs) > 0) {
		return fmt.Errorf("SeedDance first-frame modes cannot be combined with multimodal or identity references")
	}
	if b.LastFrame != "" && b.FirstFrame == "" {
		return fmt.Errorf("last frame requires a first frame")
	}
	if b.PreviewModel != "" && !supportedPreviewModel(b.PreviewModel) {
		return fmt.Errorf("preview model %q is not supported by the guided still contract", b.PreviewModel)
	}
	if len(previewImageSources(b)) > 0 && b.PreviewModel == "gpt-image-2-text-to-image" {
		return fmt.Errorf("gpt-image-2-text-to-image cannot receive preview references; use gpt-image-2-image-to-image or a Nano Banana reference model")
	}
	return nil
}

func BuildPlan(b *Brief) *Plan {
	model := strings.TrimSpace(b.Model)
	if model == "" {
		switch {
		case b.MediaType == "video":
			model = "bytedance/seedance-2-5"
		case len(b.References)+len(b.IdentityIDs) > 0 && (b.Workflow == "product-photoshoot" || b.Workflow == "marketplace-cards"):
			model = "seedream/5-pro-image-to-image"
		case len(b.References)+len(b.IdentityIDs) > 0:
			model = rankedModel("character-consistency", "gpt-image-2-image-to-image")
		default:
			model = rankedModel("text-to-image", "gpt-image-2-text-to-image")
		}
	}
	promptParts := []string{b.Request}
	if b.Workflow != "" {
		promptParts = append(promptParts, "Workflow: "+b.Workflow)
	}
	if b.Purpose != "" {
		promptParts = append(promptParts, "Created for "+b.Purpose)
	}
	if b.Platform != "" && b.Platform != "general" {
		promptParts = append(promptParts, "Optimized for "+b.Platform)
	}
	if b.Style != "" {
		promptParts = append(promptParts, "Visual direction: "+b.Style)
	}
	lesson := selectedLesson(b)
	if lesson != nil {
		promptParts = append(promptParts, "Production priorities: "+lesson.PromptFocus)
	}
	input := map[string]any{"prompt": strings.Join(promptParts, ". ") + "."}
	if b.MediaType == "video" && model == "bytedance/seedance-2-5" {
		input["aspect_ratio"] = b.AspectRatio
		input["duration"] = b.DurationSeconds
		input["resolution"] = defaultString(b.Resolution, "720p")
		input["generate_audio"] = b.AudioMode == "on"
		input["return_last_frame"] = b.ReturnLastFrame
		input["output_format"] = defaultString(b.OutputFormat, "mp4")
		input["web_search"] = b.WebSearch
		switch b.VideoMode {
		case "multimodal":
			imageRefs := append([]string(nil), b.References...)
			for _, id := range b.IdentityIDs {
				imageRefs = append(imageRefs, "identity:"+strings.TrimPrefix(id, "identity:"))
			}
			if VideoPreviewApproved(b) {
				imageRefs = append([]string{b.PreviewURL}, imageRefs...)
			}
			if len(imageRefs) > 0 {
				input["reference_image_urls"] = imageRefs
			}
			if len(b.ReferenceVideos) > 0 {
				input["reference_video_urls"] = append([]string(nil), b.ReferenceVideos...)
			}
			if len(b.ReferenceAudio) > 0 {
				input["reference_audio_urls"] = append([]string(nil), b.ReferenceAudio...)
			}
		default:
			firstFrame := b.FirstFrame
			if VideoPreviewApproved(b) {
				firstFrame = b.PreviewURL
			}
			if firstFrame != "" {
				input["first_frame_url"] = firstFrame
			}
			if b.LastFrame != "" {
				input["last_frame_url"] = b.LastFrame
			}
		}
	} else if b.MediaType == "video" {
		applyDocumentedVideoDefaults(input, b, model)
		imageRefs := append([]string(nil), b.References...)
		for _, id := range b.IdentityIDs {
			imageRefs = append(imageRefs, "identity:"+strings.TrimPrefix(id, "identity:"))
		}
		if VideoPreviewApproved(b) {
			imageRefs = append([]string{b.PreviewURL}, imageRefs...)
		}
		setDocumentedMediaSources(input, model, "image", imageRefs)
		setDocumentedMediaSources(input, model, "video", b.ReferenceVideos)
		setDocumentedMediaSources(input, model, "audio", b.ReferenceAudio)
	} else {
		input["aspect_ratio"] = b.AspectRatio
	}
	if b.MediaType != "video" && len(b.References)+len(b.IdentityIDs) > 0 {
		key := imageInputKey(model)
		values := append([]string(nil), b.References...)
		for _, id := range b.IdentityIDs {
			values = append(values, "identity:"+strings.TrimPrefix(id, "identity:"))
		}
		input[key] = values
	}
	plan := &Plan{Model: model, Input: input, Rationale: planRationale(b, model)}
	plan.ProductionSkill, plan.CapabilitySkill = routeSkills(b, model)
	plan.CostStatus = "unknown_until_provider_quote"
	plan.OverrideOptions = routeOverrides(b, model)
	applyLessonToPlan(plan, lesson)
	return plan
}

// BuildPreviewPlan creates a still-image preflight for a ready video brief.
// The final video service refuses submission until the resulting image has
// been returned and explicitly approved by the user.
func BuildPreviewPlan(b *Brief) *Plan {
	imageSources := previewImageSources(b)
	model := rankedModel("text-to-image", "gpt-image-2-text-to-image")
	if len(imageSources) > 0 {
		if b.Workflow == "product-photoshoot" || b.Workflow == "marketplace-cards" {
			model = "seedream/5-pro-image-to-image"
		} else {
			model = rankedModel("character-consistency", "gpt-image-2-image-to-image")
		}
	}
	if strings.TrimSpace(b.PreviewModel) != "" {
		model = strings.TrimSpace(b.PreviewModel)
	}
	promptParts := []string{
		"Create one production-ready still image as the visual anchor and proposed first frame for the planned video",
		b.Request,
	}
	if b.Workflow != "" {
		promptParts = append(promptParts, "Workflow: "+b.Workflow)
	}
	if b.Purpose != "" {
		promptParts = append(promptParts, "Created for "+b.Purpose)
	}
	if b.Platform != "" && b.Platform != "general" {
		promptParts = append(promptParts, "Optimized for "+b.Platform)
	}
	if b.Style != "" {
		promptParts = append(promptParts, "Visual direction: "+b.Style)
	}
	lesson := selectedLesson(b)
	if lesson != nil {
		promptParts = append(promptParts, "Production priorities: "+lesson.PromptFocus)
	}
	promptParts = append(promptParts, "Use a clear composition, stable subject identity, and production-ready lighting; avoid motion blur and do not add unrequested text")
	input := map[string]any{
		"prompt":       strings.Join(promptParts, ". ") + ".",
		"aspect_ratio": b.AspectRatio,
	}
	if len(imageSources) > 0 {
		input[imageInputKey(model)] = append([]string(nil), imageSources...)
	}
	plan := &Plan{
		Model:     model,
		Input:     input,
		Rationale: "A leaderboard-informed still-image preflight lets the user inspect composition, subject, identity, product fidelity, and style before spending credits on the final video.",
	}
	plan.ProductionSkill, plan.CapabilitySkill = routeSkills(b, model)
	plan.CostStatus = "unknown_until_provider_quote"
	plan.OverrideOptions = routeOverrides(&Brief{MediaType: "image"}, model)
	applyLessonToPlan(plan, lesson)
	return plan
}

func ApproveVideoPreview(b *Brief) error {
	if b == nil || b.MediaType != "video" {
		return fmt.Errorf("preview approval applies only to video briefs")
	}
	if !videoPreviewMatchesBrief(b) {
		return fmt.Errorf("the current preview does not match the latest brief; generate a new preview")
	}
	if strings.TrimSpace(b.PreviewURL) == "" {
		return fmt.Errorf("preview image is not ready; refresh generation %s first", b.PreviewGenerationID)
	}
	if !mediaGenerationComplete(b.PreviewStatus) {
		return fmt.Errorf("preview generation %s is not successfully complete (status %q)", b.PreviewGenerationID, b.PreviewStatus)
	}
	now := time.Now().UTC()
	b.PreviewApprovedAt = &now
	invalidateProof(b)
	b.UpdatedAt = now
	b.Plan = BuildPlan(b)
	return nil
}

func RejectVideoPreview(b *Brief) error {
	if b == nil || b.MediaType != "video" {
		return fmt.Errorf("preview rejection applies only to video briefs")
	}
	b.PreviewGenerationID = ""
	b.PreviewStatus = ""
	b.PreviewURL = ""
	b.PreviewBriefHash = ""
	b.PreviewApprovedAt = nil
	b.PreviewRevision++
	invalidateProof(b)
	b.UpdatedAt = time.Now().UTC()
	b.Plan = BuildPlan(b)
	return nil
}

func VideoPreviewApproved(b *Brief) bool {
	return b != nil && b.MediaType == "video" && b.PreviewApprovedAt != nil && strings.TrimSpace(b.PreviewURL) != "" && mediaGenerationComplete(b.PreviewStatus) && videoPreviewMatchesBrief(b)
}

func videoPreviewMatchesBrief(b *Brief) bool {
	return b != nil && strings.TrimSpace(b.PreviewBriefHash) != "" && b.PreviewBriefHash == previewBriefFingerprint(b)
}

func previewBriefFingerprint(b *Brief) string {
	if b == nil {
		return ""
	}
	data := creativeBriefFingerprint(b) + ":" + strconv.Itoa(b.PreviewRevision)
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func proofBriefFingerprint(b *Brief) string {
	if b == nil {
		return ""
	}
	data := creativeBriefFingerprint(b) + ":" + strconv.Itoa(b.PreviewRevision) + ":" + strconv.Itoa(b.ProofRevision)
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func VideoProofApproved(b *Brief) bool {
	return b != nil && b.ProofApprovedAt != nil && strings.TrimSpace(b.ProofURL) != "" && mediaGenerationComplete(b.ProofStatus) && b.ProofBriefHash == proofBriefFingerprint(b)
}

func VideoProofSkipped(b *Brief) bool {
	return b != nil && b.ProofStatus == "skipped" && b.ProofSkippedAt != nil && b.ProofBriefHash == proofBriefFingerprint(b)
}

func ApproveVideoProof(b *Brief) error {
	if b == nil || b.MediaType != "video" || !VideoPreviewApproved(b) {
		return fmt.Errorf("proof approval requires a current approved video preview")
	}
	if b.ProofBriefHash != proofBriefFingerprint(b) || strings.TrimSpace(b.ProofURL) == "" || !mediaGenerationComplete(b.ProofStatus) {
		return fmt.Errorf("the current proof is not ready or no longer matches the brief")
	}
	now := time.Now().UTC()
	b.ProofApprovedAt = &now
	b.ProofSkippedAt = nil
	b.UpdatedAt = now
	return nil
}

func RejectVideoProof(b *Brief) error {
	if b == nil || b.MediaType != "video" {
		return fmt.Errorf("proof rejection applies only to video briefs")
	}
	b.ProofGenerationID = ""
	b.ProofStatus = StatusRejected
	b.ProofURL = ""
	b.ProofBriefHash = ""
	b.ProofApprovedAt = nil
	b.ProofSkippedAt = nil
	b.ProofRevision++
	b.UpdatedAt = time.Now().UTC()
	return nil
}

func SkipVideoProof(b *Brief) error {
	if b == nil || b.MediaType != "video" || !VideoPreviewApproved(b) {
		return fmt.Errorf("proof skip requires a current approved video preview")
	}
	now := time.Now().UTC()
	b.ProofStatus = "skipped"
	b.ProofApprovedAt = nil
	b.ProofSkippedAt = &now
	b.ProofBriefHash = proofBriefFingerprint(b)
	b.UpdatedAt = now
	return nil
}

func invalidateProof(b *Brief) {
	if b == nil {
		return
	}
	b.ProofOption = nil
	b.ProofGenerationID = ""
	b.ProofStatus = ""
	b.ProofURL = ""
	b.ProofBriefHash = ""
	b.ProofApprovedAt = nil
	b.ProofSkippedAt = nil
}

func creativeBriefFingerprint(b *Brief) string {
	if b == nil {
		return ""
	}
	creativeState := struct {
		Workflow           string   `json:"workflow"`
		Lesson             string   `json:"lesson"`
		Request            string   `json:"request"`
		MediaType          string   `json:"media_type"`
		Purpose            string   `json:"purpose"`
		Platform           string   `json:"platform"`
		AspectRatio        string   `json:"aspect_ratio"`
		DurationSeconds    int      `json:"duration_seconds"`
		Resolution         string   `json:"resolution"`
		AudioMode          string   `json:"audio_mode"`
		VideoMode          string   `json:"video_mode"`
		OutputFormat       string   `json:"output_format"`
		ReturnLastFrame    bool     `json:"return_last_frame"`
		WebSearch          bool     `json:"web_search"`
		Style              string   `json:"style"`
		References         []string `json:"references"`
		ReferenceVideos    []string `json:"reference_videos"`
		ReferenceAudio     []string `json:"reference_audio"`
		FirstFrame         string   `json:"first_frame"`
		LastFrame          string   `json:"last_frame"`
		IdentityIDs        []string `json:"identity_ids"`
		RightsAcknowledged bool     `json:"rights_acknowledged,omitempty"`
		Model              string   `json:"model"`
		PreviewModel       string   `json:"preview_model"`
		ProductionMode     string   `json:"production_mode"`
	}{
		b.Workflow, b.Lesson, b.Request, b.MediaType, b.Purpose, b.Platform, b.AspectRatio,
		b.DurationSeconds, b.Resolution, b.AudioMode, b.VideoMode, b.OutputFormat,
		b.ReturnLastFrame, b.WebSearch, b.Style, cleanStrings(b.References), cleanStrings(b.ReferenceVideos),
		cleanStrings(b.ReferenceAudio), b.FirstFrame, b.LastFrame, cleanStrings(b.IdentityIDs), b.RightsAcknowledged, b.Model, b.PreviewModel, b.ProductionMode,
	}
	data, err := json.Marshal(creativeState)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func normalizeProductionMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "":
		return ""
	case "single", "single-shot", "shot":
		return ProductionModeSingleShot
	case "storyboard", "multi-shot", "multishot", "narrative":
		return ProductionModeStoryboard
	default:
		return value
	}
}

func previewImageSources(b *Brief) []string {
	if b == nil {
		return nil
	}
	sources := append([]string(nil), b.References...)
	if strings.TrimSpace(b.FirstFrame) != "" {
		sources = append(sources, b.FirstFrame)
	}
	for _, id := range b.IdentityIDs {
		sources = append(sources, "identity:"+strings.TrimPrefix(id, "identity:"))
	}
	return sources
}

func generationFailed(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "fail", "failed", "error", "cancelled", "canceled":
		return true
	default:
		return false
	}
}

func imageInputKey(model string) string {
	switch {
	case strings.HasPrefix(model, "seedream/5-pro-"):
		return "image_urls"
	case model == "nano-banana-2" || model == "nano-banana-pro":
		return "image_input"
	case model == "nano-banana-2-lite":
		return "image_urls"
	case strings.HasPrefix(model, "ideogram/character"):
		return "reference_image_urls"
	default:
		return "input_urls"
	}
}

// documentedMediaInput returns the first model-schema field that can carry the
// requested reference type. The boolean records whether the field is an array;
// callers must not invent generic keys for models that do not document one.
func documentedMediaInput(model, mediaType string) (string, bool, bool) {
	item, err := kiecatalog.Get(strings.TrimSpace(model))
	if err != nil {
		return "", false, false
	}
	properties, _ := item.InputSchema["properties"].(map[string]any)
	candidates := map[string][]string{
		"image": {"reference_image_urls", "image_urls", "input_urls", "image_input", "image_url", "first_frame_url"},
		"video": {"reference_video_urls", "video_urls", "video_url"},
		"audio": {"reference_audio_urls", "audio_urls", "audio_url"},
	}[mediaType]
	for _, field := range candidates {
		raw, ok := properties[field]
		if !ok {
			continue
		}
		schema, _ := raw.(map[string]any)
		return field, schema["type"] == "array", true
	}
	return "", false, false
}

func setDocumentedMediaSources(input map[string]any, model, mediaType string, sources []string) {
	field, array, ok := documentedMediaInput(model, mediaType)
	if !ok || len(sources) == 0 {
		return
	}
	if array {
		input[field] = append([]string(nil), sources...)
	} else if len(sources) == 1 {
		input[field] = sources[0]
	}
}

func applyDocumentedVideoDefaults(input map[string]any, b *Brief, model string) {
	item, err := kiecatalog.Get(model)
	if err != nil {
		return
	}
	properties, _ := item.InputSchema["properties"].(map[string]any)
	if _, ok := properties["aspect_ratio"]; ok {
		input["aspect_ratio"] = b.AspectRatio
	} else if _, ok := properties["ratio"]; ok {
		input["ratio"] = b.AspectRatio
	}
	if raw, ok := properties["duration"].(map[string]any); ok {
		if raw["type"] == "string" {
			input["duration"] = strconv.Itoa(b.DurationSeconds)
		} else {
			input["duration"] = b.DurationSeconds
		}
	}
	if _, ok := properties["sound"]; ok {
		input["sound"] = b.AudioMode == "on"
	}
	if _, ok := properties["multi_shots"]; ok {
		input["multi_shots"] = false
	}
	if _, ok := properties["multi_prompt"]; ok {
		input["multi_prompt"] = []any{}
	}
	for _, field := range []string{"mode", "resolution"} {
		raw, ok := properties[field].(map[string]any)
		if !ok {
			continue
		}
		if field == "resolution" && enumContains(raw["enum"], b.Resolution) {
			input[field] = b.Resolution
			continue
		}
		if value, ok := raw["default"]; ok {
			input[field] = value
		}
	}
}

func enumContains(value any, candidate string) bool {
	if strings.TrimSpace(candidate) == "" {
		return false
	}
	values, _ := value.([]any)
	for _, raw := range values {
		if text, ok := raw.(string); ok && text == candidate {
			return true
		}
	}
	return false
}

func supportedPreviewModel(model string) bool {
	switch strings.TrimSpace(model) {
	case "gpt-image-2-text-to-image", "gpt-image-2-image-to-image", "nano-banana-2", "nano-banana-2-lite", "nano-banana-pro", "seedream/5-pro-image-to-image":
		return true
	default:
		return false
	}
}

func planRationale(b *Brief, model string) string {
	if strings.TrimSpace(b.Model) != "" {
		return "Uses the model explicitly selected in the brief."
	}
	if b.MediaType == "video" {
		return "SeedDance 2.5 is the flagship video route, with 4-30 second or automatic duration, optional audio, first/last-frame control, and multimodal image/video/audio references."
	}
	if model == "seedream/5-pro-image-to-image" {
		return "Seedream 5 Pro image-to-image is the reference-fidelity route for product and marketplace workflows."
	}
	if len(b.References)+len(b.IdentityIDs) > 0 {
		return "GPT Image 2 image-to-image preserves the supplied visual references."
	}
	return "GPT Image 2 text-to-image is the default high-fidelity image route."
}

func routeSkills(b *Brief, model string) (string, string) {
	production := "kie-create"
	if b != nil {
		switch {
		case b.ProductionMode == ProductionModeStoryboard:
			production = "kie-film"
		case b.Workflow != "" && b.Workflow != "generate":
			production = "kie-" + b.Workflow
		case b.MediaType == "video":
			production = "kie-video"
		case b.MediaType == "image":
			production = "kie-image"
		}
	}
	capability := "kie-generate"
	if classified, err := kiecatalog.GetCapability(model); err == nil {
		capability = classified.PrimaryCapability
	} else if b != nil && b.MediaType != "" {
		capability = "kie-" + b.MediaType
	}
	return production, capability
}

func routeOverrides(b *Brief, selected string) []string {
	values := []string{selected}
	if b != nil && b.MediaType == "video" {
		values = append(values, "bytedance/seedance-2-5", "kling-3.0/video", "wan/2-7-text-to-video")
	} else {
		values = append(values, "gpt-image-2-text-to-image", "nano-banana-2", "seedream/5-pro-text-to-image")
	}
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func requiresRightsAcknowledgement(b *Brief) bool {
	if b == nil {
		return false
	}
	model := strings.TrimSpace(BuildPlan(b).Model)
	if classified, err := kiecatalog.GetCapability(model); err == nil {
		if classified.PrimaryCapability == "kie-avatar" || classified.PrimaryCapability == "kie-audio" {
			return true
		}
		for _, secondary := range classified.SecondaryCapabilities {
			if secondary == "kie-avatar" {
				return true
			}
		}
	}
	return containsAny(strings.ToLower(model), "avatar", "voice", "text-to-speech", "text-to-dialogue", "infinitalk", "omnihuman")
}

var durationPattern = regexp.MustCompile(`(?i)\b(\d{1,3})\s*(?:seconds?|secs?|s)\b`)

func inferBriefFromRequest(b *Brief) {
	if b == nil {
		return
	}
	request := strings.ToLower(b.Request)
	if b.MediaType == "" {
		switch {
		case containsAny(request, " video", "video ", "clip", "animate", "film", "trailer"):
			b.MediaType = "video"
		case containsAny(request, "image", "photo", "picture", "thumbnail", "poster", "logo", "illustration"):
			b.MediaType = "image"
		}
	}
	if b.Platform == "" {
		for _, platform := range []string{"tiktok", "instagram", "youtube", "linkedin", "website"} {
			if strings.Contains(request, platform) {
				b.Platform = platform
				break
			}
		}
	}
	if b.AspectRatio == "" {
		for _, ratio := range []string{"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"} {
			if strings.Contains(request, ratio) {
				b.AspectRatio = ratio
				break
			}
		}
		if b.AspectRatio == "" && strings.Contains(request, "vertical") {
			b.AspectRatio = "9:16"
		}
	}
	if b.DurationSeconds == 0 && b.MediaType == "video" {
		if match := durationPattern.FindStringSubmatch(request); len(match) == 2 {
			if seconds, err := strconv.Atoi(match[1]); err == nil && seconds >= 4 {
				maxDuration := 30
				if b.ProductionMode == ProductionModeStoryboard {
					maxDuration = 600
				}
				if seconds <= maxDuration {
					b.DurationSeconds = seconds
				}
			}
		}
	}
	if b.Purpose == "" {
		switch {
		case containsAny(request, "advert", " ad ", "campaign"):
			b.Purpose = "advertising campaign"
		case strings.Contains(request, "thumbnail"):
			b.Purpose = "thumbnail"
		case containsAny(request, "hero section", "website hero"):
			b.Purpose = "website hero"
		case containsAny(request, "social post", "social media"):
			b.Purpose = "social post"
		}
	}
	if b.Style == "" {
		var matched []string
		for _, token := range []string{"cinematic", "photorealistic", "documentary", "handheld", "editorial", "minimal", "studio", "anime", "watercolor"} {
			if strings.Contains(request, token) {
				matched = append(matched, token)
			}
		}
		if len(matched) > 0 {
			b.Style = strings.Join(matched, ", ")
		}
	}
	if b.MediaType == "video" {
		if b.AudioMode == "" {
			switch {
			case containsAny(request, "with audio", "synchronised audio", "synchronized audio", "dialogue", "speaking"):
				b.AudioMode = "on"
			case containsAny(request, "silent", "without audio", "no audio"):
				b.AudioMode = "off"
			}
		}
		if b.VideoMode == "" {
			switch {
			case len(b.ReferenceVideos)+len(b.ReferenceAudio)+len(b.IdentityIDs)+len(b.References) > 0:
				b.VideoMode = "multimodal"
			case b.FirstFrame != "" && b.LastFrame != "":
				b.VideoMode = "first-last"
			case b.FirstFrame != "":
				b.VideoMode = "first-frame"
			case containsAny(request, "text to video", "text-to-video"):
				b.VideoMode = "text"
			}
		}
	}
}

func containsAny(value string, candidates ...string) bool {
	for _, candidate := range candidates {
		if strings.Contains(value, candidate) {
			return true
		}
	}
	return false
}

func rankedModel(task, fallback string) string {
	entry, err := leaderboard.BestAvailable(task)
	if err != nil || strings.TrimSpace(entry.KieModel) == "" {
		return fallback
	}
	return entry.KieModel
}

func selectedLesson(b *Brief) *academy.Lesson {
	if b == nil || strings.TrimSpace(b.Lesson) == "" {
		return nil
	}
	lesson, err := academy.GetLesson(b.Lesson)
	if err != nil {
		return nil
	}
	return lesson
}

func applyLessonToPlan(plan *Plan, lesson *academy.Lesson) {
	if plan == nil || lesson == nil {
		return
	}
	plan.Lesson = lesson.Key
	plan.ProductionStage = lesson.ProductionStage
	plan.Method = lesson.KieMethod
	plan.PromptFocus = lesson.PromptFocus
}

func isSkipAnswer(answer string) bool {
	return strings.EqualFold(answer, "skip") || strings.EqualFold(answer, "done") || strings.EqualFold(answer, "none")
}

func normalizeAudioMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "yes", "true", "audio", "with audio", "on":
		return "on"
	case "no", "false", "silent", "without audio", "off":
		return "off"
	default:
		return value
	}
}

func normalizeVideoMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "first", "first frame", "image-to-video":
		return "first-frame"
	case "first last", "first-and-last", "first and last":
		return "first-last"
	case "reference", "references", "reference-to-video":
		return "multimodal"
	default:
		return value
	}
}

func normalizeResolution(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func normalizeOutputFormat(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func defaultString(value, fallback string) string {
	if value = strings.TrimSpace(value); value != "" {
		return value
	}
	return fallback
}

func recommendedAspect(platform, mediaType string) string {
	switch normalizePlatform(platform) {
	case "instagram":
		if mediaType == "image" {
			return "3:4"
		}
		return "9:16"
	case "tiktok":
		return "9:16"
	case "youtube", "website", "linkedin":
		return "16:9"
	default:
		if mediaType == "video" {
			return "16:9"
		}
		return "1:1"
	}
}

func normalizeMediaType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "photo", "picture", "graphic", "still":
		return "image"
	case "clip", "movie", "animation":
		return "video"
	default:
		return value
	}
}

func normalizePlatform(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "web" || value == "site" {
		return "website"
	}
	return value
}

func cleanStrings(values []string) []string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
