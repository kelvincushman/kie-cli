// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"kie-pp-cli/internal/kiecatalog"
)

type API interface {
	Get(context.Context, string, map[string]string) (json.RawMessage, error)
	PostWithParams(context.Context, string, map[string]string, any) (json.RawMessage, int, error)
	PostMultipartWithParams(context.Context, string, map[string]string, map[string]string, map[string]string) (json.RawMessage, int, error)
}

type Service struct {
	API   API
	Store *Store
}

func (s *Service) Submit(ctx context.Context, b *Brief, confirmationIDs ...string) (*Generation, error) {
	if s == nil || s.API == nil || s.Store == nil {
		return nil, fmt.Errorf("media service is not configured")
	}
	if b == nil {
		return nil, fmt.Errorf("media brief is required")
	}
	requestedBrief := b
	release, err := s.Store.acquireSubmission(b.ID)
	if err != nil {
		return nil, err
	}
	defer release()
	stored, err := s.Store.GetBrief(b.ID)
	if err != nil {
		return nil, fmt.Errorf("loading brief %s before submission: %w", b.ID, err)
	}
	if stored.Status == StatusSubmitted || strings.TrimSpace(stored.GenerationID) != "" {
		return nil, fmt.Errorf("brief %s was already submitted as generation %s; create a new brief to generate again", stored.ID, stored.GenerationID)
	}
	if previous, findErr := s.Store.findGenerationByBriefIDAndKind(stored.ID, GenerationKindFinal); findErr != nil {
		return nil, findErr
	} else if previous != nil {
		stored.Status = StatusSubmitted
		stored.GenerationID = previous.ID
		if err := s.Store.SaveBrief(stored); err != nil {
			return nil, fmt.Errorf("recording existing final generation %s on brief %s: %w", previous.ID, stored.ID, err)
		}
		return nil, fmt.Errorf("brief %s was already submitted as generation %s; check its status instead", stored.ID, previous.ID)
	}
	b = stored
	Refresh(b)
	if err := s.requireApprovedStoryboardForShot(b); err != nil {
		return nil, err
	}
	if b.MediaType == "video" && b.ProductionMode == ProductionModeStoryboard {
		return nil, fmt.Errorf("storyboard master brief %s is not submitted as one prompt-only video; generate and approve each shot brief, then assemble the returned clips", b.ID)
	}
	if b.Status != StatusReady || b.Plan == nil {
		return nil, fmt.Errorf("brief %s is not ready; answer the next question first", b.ID)
	}
	if b.MediaType == "video" && !VideoPreviewApproved(b) {
		return nil, fmt.Errorf("video brief %s requires an approved preview image before final submission; generate, show, and explicitly approve the preview first", b.ID)
	}
	plan := BuildPlan(b)
	if b.MediaType == "video" && ResolveProofOption(plan.Model).Supported && !VideoProofApproved(b) && !VideoProofSkipped(b) {
		return nil, fmt.Errorf("video brief %s requires an explicit proof decision before final submission; approve the current proof or explicitly skip the proof first", b.ID)
	}
	confirmationID, err := requiredConfirmationID(confirmationIDs)
	if err != nil {
		return nil, err
	}
	if _, err := s.Store.UsePaidConfirmation(confirmationID, b, plan, PaidScopeFinal, GenerationKindFinal, time.Now()); err != nil {
		return nil, err
	}
	imageSources, err := s.imageSourcesForBrief(b, false)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(plan.Model, "seedream/5-pro-") && len(imageSources) > 10 {
		return nil, fmt.Errorf("Seedream 5 Pro accepts at most 10 image references; got %d", len(imageSources))
	}
	if b.MediaType == "video" {
		if err := s.prepareVideoPlanReferences(ctx, b, plan, imageSources); err != nil {
			return nil, err
		}
	} else if len(imageSources) > 0 {
		key := imageInputKey(plan.Model)
		if err := s.resolvePlanReferences(ctx, b.ID, plan.Input, key, imageSources, "image"); err != nil {
			return nil, err
		}
	}
	if err := validatePaidPlan(plan); err != nil {
		return nil, err
	}
	body := map[string]any{"model": plan.Model, "input": plan.Input}
	data, status, err := s.API.PostWithParams(ctx, "/api/v1/jobs/createTask", map[string]string{}, body)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("create task returned HTTP %d", status)
	}
	taskID := firstStringForKeys(data, "taskId", "task_id")
	if taskID == "" {
		return nil, fmt.Errorf("create task response did not include a task id")
	}
	now := time.Now().UTC()
	g := &Generation{ID: newID("gen"), BriefID: b.ID, Kind: GenerationKindFinal, Fingerprint: previewBriefFingerprint(b), TaskID: taskID, Model: plan.Model, Status: "submitted", CreatedAt: now, UpdatedAt: now, Remote: data}
	if err := s.Store.SaveGeneration(g); err != nil {
		return nil, err
	}
	b.Status = StatusSubmitted
	b.Plan = plan
	b.GenerationID = g.ID
	b.UpdatedAt = now
	if err := s.Store.SaveBrief(b); err != nil {
		return nil, err
	}
	*requestedBrief = *b
	return g, nil
}

func (s *Service) SubmitPreview(ctx context.Context, b *Brief, confirmationIDs ...string) (*Generation, error) {
	if s == nil || s.API == nil || s.Store == nil {
		return nil, fmt.Errorf("media service is not configured")
	}
	if b == nil {
		return nil, fmt.Errorf("media brief is required")
	}
	requestedBrief := b
	release, err := s.Store.acquirePreviewSubmission(b.ID)
	if err != nil {
		return nil, err
	}
	defer release()
	stored, err := s.Store.GetBrief(b.ID)
	if err != nil {
		return nil, fmt.Errorf("loading brief %s before preview: %w", b.ID, err)
	}
	Refresh(stored)
	if err := s.requireApprovedStoryboardForShot(stored); err != nil {
		return nil, err
	}
	if stored.MediaType != "video" {
		return nil, fmt.Errorf("preview generation applies only to video briefs")
	}
	if stored.ProductionMode == ProductionModeStoryboard {
		return nil, fmt.Errorf("storyboard master brief %s does not have one preview; generate previews for its shot brief ids", stored.ID)
	}
	if stored.Status != StatusReady || stored.Plan == nil {
		return nil, fmt.Errorf("brief %s is not ready; answer the next question first", stored.ID)
	}
	fingerprint := previewBriefFingerprint(stored)
	if previous, findErr := s.Store.findActiveGenerationByFingerprint(stored.ID, GenerationKindPreview, fingerprint); findErr != nil {
		return nil, findErr
	} else if previous != nil {
		stored.PreviewGenerationID = previous.ID
		stored.PreviewStatus = previous.Status
		stored.PreviewBriefHash = fingerprint
		if len(previous.ResultURLs) > 0 {
			stored.PreviewURL = previous.ResultURLs[0]
		}
		if err := s.Store.SaveBrief(stored); err != nil {
			return nil, fmt.Errorf("recording existing preview generation %s on brief %s: %w", previous.ID, stored.ID, err)
		}
		return nil, fmt.Errorf("brief %s already has preview generation %s; refresh its status, then approve or reject it", stored.ID, previous.ID)
	}
	plan := BuildPreviewPlan(stored)
	confirmationID, err := requiredConfirmationID(confirmationIDs)
	if err != nil {
		return nil, err
	}
	if _, err := s.Store.UsePaidConfirmation(confirmationID, stored, plan, PaidScopePreview, GenerationKindPreview, time.Now()); err != nil {
		return nil, err
	}
	imageSources, err := s.imageSourcesForBrief(stored, true)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(plan.Model, "seedream/5-pro-") && len(imageSources) > 10 {
		return nil, fmt.Errorf("Seedream 5 Pro accepts at most 10 preview image references; got %d", len(imageSources))
	}
	if len(imageSources) > 0 {
		if err := s.resolvePlanReferences(ctx, stored.ID, plan.Input, imageInputKey(plan.Model), imageSources, "image"); err != nil {
			return nil, err
		}
	}
	if err := validatePaidPlan(plan); err != nil {
		return nil, err
	}
	body := map[string]any{"model": plan.Model, "input": plan.Input}
	data, status, err := s.API.PostWithParams(ctx, "/api/v1/jobs/createTask", map[string]string{}, body)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("create preview task returned HTTP %d", status)
	}
	taskID := firstStringForKeys(data, "taskId", "task_id")
	if taskID == "" {
		return nil, fmt.Errorf("create preview task response did not include a task id")
	}
	now := time.Now().UTC()
	g := &Generation{
		ID: newID("gen"), BriefID: stored.ID, Kind: GenerationKindPreview, Fingerprint: fingerprint,
		TaskID: taskID, Model: plan.Model, Status: "submitted", CreatedAt: now, UpdatedAt: now, Remote: data,
	}
	if err := s.Store.SaveGeneration(g); err != nil {
		return nil, err
	}
	stored.PreviewGenerationID = g.ID
	stored.PreviewStatus = g.Status
	stored.PreviewURL = ""
	stored.PreviewBriefHash = fingerprint
	stored.PreviewApprovedAt = nil
	stored.UpdatedAt = now
	stored.Plan = BuildPlan(stored)
	if err := s.Store.SaveBrief(stored); err != nil {
		return nil, err
	}
	*requestedBrief = *stored
	return g, nil
}

func (s *Service) SubmitProof(ctx context.Context, b *Brief, confirmationIDs ...string) (*Generation, error) {
	if s == nil || s.API == nil || s.Store == nil {
		return nil, fmt.Errorf("media service is not configured")
	}
	if b == nil {
		return nil, fmt.Errorf("media brief is required")
	}
	requestedBrief := b
	release, err := s.Store.acquireProofSubmission(b.ID)
	if err != nil {
		return nil, err
	}
	defer release()
	stored, err := s.Store.GetBrief(b.ID)
	if err != nil {
		return nil, fmt.Errorf("loading brief %s before proof: %w", b.ID, err)
	}
	Refresh(stored)
	if err := s.requireApprovedStoryboardForShot(stored); err != nil {
		return nil, err
	}
	if stored.MediaType != "video" || !VideoPreviewApproved(stored) {
		return nil, fmt.Errorf("complete-shot proof requires a current approved video preview")
	}
	if stored.ProductionMode == ProductionModeStoryboard {
		return nil, fmt.Errorf("generate proofs for storyboard shot brief ids, not the master brief")
	}
	fingerprint := proofBriefFingerprint(stored)
	if previous, findErr := s.Store.findActiveGenerationByFingerprint(stored.ID, GenerationKindProof, fingerprint); findErr != nil {
		return nil, findErr
	} else if previous != nil {
		return nil, fmt.Errorf("brief %s already has proof generation %s; refresh and review it", stored.ID, previous.ID)
	}
	plan, option, err := BuildProofPlan(stored)
	if err != nil {
		return nil, err
	}
	confirmationID, err := requiredConfirmationID(confirmationIDs)
	if err != nil {
		return nil, err
	}
	if _, err := s.Store.UsePaidConfirmation(confirmationID, stored, plan, PaidScopeProof, GenerationKindProof, time.Now()); err != nil {
		return nil, err
	}
	imageSources, err := s.imageSourcesForBrief(stored, false)
	if err != nil {
		return nil, err
	}
	if err := s.prepareVideoPlanReferences(ctx, stored, plan, imageSources); err != nil {
		return nil, err
	}
	if err := validatePaidPlan(plan); err != nil {
		return nil, err
	}
	data, status, err := s.API.PostWithParams(ctx, "/api/v1/jobs/createTask", map[string]string{}, map[string]any{"model": plan.Model, "input": plan.Input})
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("create proof task returned HTTP %d", status)
	}
	taskID := firstStringForKeys(data, "taskId", "task_id")
	if taskID == "" {
		return nil, fmt.Errorf("create proof task response did not include a task id")
	}
	now := time.Now().UTC()
	g := &Generation{ID: newID("gen"), BriefID: stored.ID, Kind: GenerationKindProof, Fingerprint: fingerprint, TaskID: taskID, Model: plan.Model, Status: "submitted", CreatedAt: now, UpdatedAt: now, Remote: data}
	if err := s.Store.SaveGeneration(g); err != nil {
		return nil, err
	}
	stored.ProofOption = &option
	stored.ProofGenerationID = g.ID
	stored.ProofStatus = g.Status
	stored.ProofURL = ""
	stored.ProofBriefHash = fingerprint
	stored.ProofApprovedAt = nil
	stored.ProofSkippedAt = nil
	stored.UpdatedAt = now
	if err := s.Store.SaveBrief(stored); err != nil {
		return nil, err
	}
	*requestedBrief = *stored
	return g, nil
}

func requiredConfirmationID(values []string) (string, error) {
	if len(values) != 1 || strings.TrimSpace(values[0]) == "" {
		return "", fmt.Errorf("a fresh scoped paid confirmation is required for this live action")
	}
	return strings.TrimSpace(values[0]), nil
}

func validatePaidPlan(plan *Plan) error {
	if plan == nil {
		return fmt.Errorf("generation plan is required")
	}
	// The catalog validator receives JSON-shaped values from CLI/MCP callers.
	// Director plans use native []string slices, so normalize through JSON before
	// applying the exact same contract.
	encoded, err := json.Marshal(plan.Input)
	if err != nil {
		return fmt.Errorf("encoding %s input for local validation: %w", plan.Model, err)
	}
	normalized := map[string]any{}
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		return fmt.Errorf("normalizing %s input for local validation: %w", plan.Model, err)
	}
	issues, err := kiecatalog.Validate(plan.Model, normalized)
	if err != nil {
		return fmt.Errorf("validating %s input before paid generation: %w", plan.Model, err)
	}
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, strings.TrimSpace(issue.Path+" "+issue.Message))
	}
	return fmt.Errorf("%s input failed local contract validation: %s", plan.Model, strings.Join(parts, "; "))
}

func (s *Service) requireApprovedStoryboardForShot(brief *Brief) error {
	if brief == nil || strings.TrimSpace(brief.MasterBriefID) == "" {
		return nil
	}
	master, err := s.Store.GetBrief(brief.MasterBriefID)
	if err != nil {
		return fmt.Errorf("loading storyboard master %s: %w", brief.MasterBriefID, err)
	}
	if master.StoryboardStatus != StatusApproved || master.StoryboardID != brief.StoryboardID ||
		master.StoryboardBriefHash != creativeBriefFingerprint(master) || master.StoryboardScriptHash != master.ScriptHash ||
		master.ScriptStatus != StatusApproved || master.ScriptBriefHash != creativeBriefFingerprint(master) {
		return fmt.Errorf("shot brief %s is locked until its current master script and storyboard are explicitly approved", brief.ID)
	}
	storyboard, err := s.Store.GetStoryboard(master.ID)
	if err != nil {
		return fmt.Errorf("loading storyboard %s: %w", master.StoryboardID, err)
	}
	script, err := s.Store.GetScript(master.ID)
	if err != nil {
		return fmt.Errorf("loading script %s: %w", master.ScriptID, err)
	}
	if script.Hash != scriptFingerprint(script) || script.Hash != master.ScriptHash ||
		script.Status != StatusApproved || script.BriefHash != creativeBriefFingerprint(master) {
		return fmt.Errorf("shot brief %s belongs to a stale or modified script", brief.ID)
	}
	if storyboard.ID != master.StoryboardID || storyboard.Status != StatusApproved ||
		storyboard.Hash != storyboardFingerprint(storyboard) || storyboard.ScriptHash != master.ScriptHash || storyboard.BriefHash != creativeBriefFingerprint(master) {
		return fmt.Errorf("shot brief %s belongs to a stale or unapproved storyboard", brief.ID)
	}
	for _, shot := range storyboard.Shots {
		if shot.ID == brief.ShotID && shot.BriefID == brief.ID {
			return nil
		}
	}
	return fmt.Errorf("shot brief %s is not part of the current approved storyboard", brief.ID)
}

func (s *Service) RefreshGeneration(ctx context.Context, id string) (*Generation, error) {
	g, err := s.Store.GetGeneration(id)
	if err != nil {
		return nil, err
	}
	data, err := s.API.Get(ctx, "/api/v1/jobs/recordInfo", map[string]string{"taskId": g.TaskID})
	if err != nil {
		return nil, err
	}
	if status := firstStringForKeys(data, "state", "status", "taskStatus", "task_status"); status != "" {
		g.Status = strings.ToLower(status)
	}
	g.ResultURLs = collectResultURLs(data)
	g.Remote = data
	g.UpdatedAt = time.Now().UTC()
	if err := s.Store.SaveGeneration(g); err != nil {
		return nil, err
	}
	if g.Kind == GenerationKindPreview || g.Kind == GenerationKindProof {
		brief, briefErr := s.Store.GetBrief(g.BriefID)
		if briefErr != nil {
			return nil, fmt.Errorf("loading brief %s for generated status: %w", g.BriefID, briefErr)
		}
		if g.Kind == GenerationKindPreview && brief.PreviewGenerationID == g.ID && brief.PreviewBriefHash == g.Fingerprint {
			brief.PreviewStatus = g.Status
			if len(g.ResultURLs) > 0 {
				if brief.PreviewURL != g.ResultURLs[0] {
					brief.PreviewURL = g.ResultURLs[0]
					brief.PreviewApprovedAt = nil
				}
			}
			brief.UpdatedAt = g.UpdatedAt
			brief.Plan = BuildPlan(brief)
		} else if g.Kind == GenerationKindProof && brief.ProofGenerationID == g.ID && brief.ProofBriefHash == g.Fingerprint {
			brief.ProofStatus = g.Status
			if len(g.ResultURLs) > 0 {
				if brief.ProofURL != g.ResultURLs[0] {
					brief.ProofURL = g.ResultURLs[0]
					brief.ProofApprovedAt = nil
				}
			}
			brief.UpdatedAt = g.UpdatedAt
		}
		if err := s.Store.SaveBrief(brief); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (s *Service) imageSourcesForBrief(b *Brief, includeFirstFrame bool) ([]string, error) {
	imageSources := append([]string(nil), b.References...)
	if includeFirstFrame && strings.TrimSpace(b.FirstFrame) != "" {
		imageSources = append(imageSources, b.FirstFrame)
	}
	for _, identityID := range b.IdentityIDs {
		identity, err := s.Store.GetIdentity(identityID)
		if err != nil {
			return nil, fmt.Errorf("loading identity %q: %w", identityID, err)
		}
		imageSources = append(imageSources, identity.ImageReferences...)
	}
	return imageSources, nil
}

func (s *Service) prepareVideoPlanReferences(ctx context.Context, b *Brief, plan *Plan, explicitImages []string) error {
	if plan.Model == "bytedance/seedance-2-5" {
		switch b.VideoMode {
		case "multimodal":
			images := append([]string{b.PreviewURL}, explicitImages...)
			if len(images) > 30 {
				return fmt.Errorf("SeedDance 2.5 accepts at most 30 total preview, image, and identity references; got %d", len(images))
			}
			if err := s.resolvePlanReferences(ctx, b.ID, plan.Input, "reference_image_urls", images, "image"); err != nil {
				return err
			}
			if err := s.resolvePlanReferences(ctx, b.ID, plan.Input, "reference_video_urls", b.ReferenceVideos, "video"); err != nil {
				return err
			}
			return s.resolvePlanReferences(ctx, b.ID, plan.Input, "reference_audio_urls", b.ReferenceAudio, "audio")
		default:
			firstFrame, err := s.resolveReference(ctx, b.ID, b.PreviewURL, "image")
			if err != nil {
				return err
			}
			plan.Input["first_frame_url"] = firstFrame
			if b.LastFrame != "" {
				lastFrame, err := s.resolveReference(ctx, b.ID, b.LastFrame, "image")
				if err != nil {
					return err
				}
				plan.Input["last_frame_url"] = lastFrame
			}
			return nil
		}
	}

	images := append([]string{b.PreviewURL}, explicitImages...)
	if _, _, ok := documentedMediaInput(plan.Model, "image"); ok {
		if err := s.resolveDocumentedMediaReferences(ctx, b.ID, plan, "image", images); err != nil {
			return err
		}
	} else if len(explicitImages) > 0 {
		return fmt.Errorf("model %s does not document image references; remove the image/identity references or choose a compatible video model", plan.Model)
	}
	if err := s.resolveDocumentedMediaReferences(ctx, b.ID, plan, "video", b.ReferenceVideos); err != nil {
		return err
	}
	return s.resolveDocumentedMediaReferences(ctx, b.ID, plan, "audio", b.ReferenceAudio)
}

func (s *Service) resolveDocumentedMediaReferences(ctx context.Context, briefID string, plan *Plan, mediaType string, sources []string) error {
	if len(sources) == 0 {
		return nil
	}
	field, array, ok := documentedMediaInput(plan.Model, mediaType)
	if !ok {
		return fmt.Errorf("model %s does not document %s references", plan.Model, mediaType)
	}
	if array {
		return s.resolvePlanReferences(ctx, briefID, plan.Input, field, sources, mediaType)
	}
	if len(sources) != 1 {
		return fmt.Errorf("model %s accepts one %s reference in %s; got %d", plan.Model, mediaType, field, len(sources))
	}
	resolved, err := s.resolveReference(ctx, briefID, sources[0], mediaType)
	if err != nil {
		return err
	}
	plan.Input[field] = resolved
	return nil
}

func (s *Service) resolvePlanReferences(ctx context.Context, briefID string, input map[string]any, key string, sources []string, mediaType string) error {
	if len(sources) == 0 {
		delete(input, key)
		return nil
	}
	urls := make([]string, 0, len(sources))
	for _, source := range sources {
		resolved, err := s.resolveReference(ctx, briefID, source, mediaType)
		if err != nil {
			return err
		}
		urls = append(urls, resolved)
	}
	input[key] = urls
	return nil
}

func (s *Service) resolveReference(ctx context.Context, briefID, source, mediaType string) (string, error) {
	source = strings.TrimSpace(source)
	if strings.HasPrefix(source, "ref:") {
		ref, err := s.Store.GetReference(source)
		if err != nil {
			return "", err
		}
		if ref.MediaType != "" && ref.MediaType != mediaType {
			return "", fmt.Errorf("reference %s is %s, not %s", source, ref.MediaType, mediaType)
		}
		if ref.URL != "" {
			return ref.URL, nil
		}
		source = ref.StoredPath
	}
	if parsed, err := url.Parse(source); err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") {
		return source, nil
	}
	if _, err := inspectLocalReference(source, mediaType); err != nil {
		return "", err
	}
	fields := map[string]string{"fileName": filepath.Base(source), "uploadPath": "media-director/" + briefID}
	data, status, err := s.API.PostMultipartWithParams(ctx, "/api/file-stream-upload", map[string]string{}, fields, map[string]string{"file": source})
	if err != nil {
		return "", fmt.Errorf("uploading reference %q: %w", source, err)
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("uploading reference %q returned HTTP %d", source, status)
	}
	downloadURL := firstStringForKeys(data, "downloadUrl", "download_url", "url")
	if downloadURL == "" {
		return "", fmt.Errorf("upload response for %q did not include a download URL", source)
	}
	return downloadURL, nil
}

func firstStringForKeys(data json.RawMessage, keys ...string) string {
	var value any
	if json.Unmarshal(data, &value) != nil {
		return ""
	}
	keySet := make(map[string]bool, len(keys))
	for _, key := range keys {
		keySet[key] = true
	}
	var walk func(any) string
	walk = func(v any) string {
		switch typed := v.(type) {
		case map[string]any:
			for key, child := range typed {
				if keySet[key] {
					if result, ok := child.(string); ok && strings.TrimSpace(result) != "" {
						return result
					}
				}
			}
			for _, child := range typed {
				if result := walk(child); result != "" {
					return result
				}
			}
		case []any:
			for _, child := range typed {
				if result := walk(child); result != "" {
					return result
				}
			}
		}
		return ""
	}
	return walk(value)
}

func collectResultURLs(data json.RawMessage) []string {
	var value any
	if json.Unmarshal(data, &value) != nil {
		return nil
	}
	seen := map[string]bool{}
	var results []string
	var walk func(any, string)
	walk = func(v any, parent string) {
		switch typed := v.(type) {
		case map[string]any:
			for key, child := range typed {
				walk(child, strings.ToLower(key))
			}
		case []any:
			for _, child := range typed {
				walk(child, parent)
			}
		case string:
			if !strings.HasPrefix(typed, "https://") && !strings.HasPrefix(typed, "http://") {
				return
			}
			if !strings.Contains(parent, "url") && !strings.Contains(parent, "result") && !strings.Contains(parent, "output") {
				return
			}
			if !seen[typed] {
				seen[typed] = true
				results = append(results, typed)
			}
		}
	}
	walk(value, "")
	return results
}
