// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestApprovalFingerprintsSurviveJSONRoundTripWithEmptyReferences(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "A single shot", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text", Style: "clean",
	})
	if err != nil {
		t.Fatal(err)
	}
	brief.PreviewGenerationID = "gen_preview"
	brief.PreviewStatus = "success"
	brief.PreviewURL = "https://cdn.example.test/preview.png"
	brief.PreviewBriefHash = previewBriefFingerprint(brief)
	if err := ApproveVideoPreview(brief); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(brief)
	if err != nil {
		t.Fatal(err)
	}
	var loaded Brief
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if !VideoPreviewApproved(&loaded) {
		t.Fatal("preview approval became stale after empty slices were omitted during JSON round trip")
	}
}

func TestStoryboardProductionCreatesGatedShotBriefs(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A founder introduces a new camera", MediaType: "video", Purpose: "launch film",
		Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 10, AudioMode: "on",
		VideoMode: "text", Style: "cinematic documentary", ProductionMode: ProductionModeStoryboard,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	if got := TurnFor(brief).NextAction; got != "draft_script" {
		t.Fatalf("initial next action = %q", got)
	}
	script, err := store.SetScript(brief.ID, ScriptInput{Title: "Launch", Content: "The founder opens the case, then reveals the camera."})
	if err != nil {
		t.Fatal(err)
	}
	if script.Status != StatusDraft {
		t.Fatalf("script status = %q", script.Status)
	}
	if _, err := store.DecideScript(brief.ID, "approve"); err != nil {
		t.Fatal(err)
	}
	brief, _ = store.GetBrief(brief.ID)
	if got := TurnFor(brief).NextAction; got != "draft_storyboard" {
		t.Fatalf("post-script next action = %q", got)
	}

	storyboard, err := store.SetStoryboard(brief.ID, StoryboardInput{Title: "Two shots", Shots: []StoryboardShotInput{
		{Title: "Open", DurationSeconds: 5, Visual: "Founder opens a hard case on a studio table", Camera: "slow dolly in"},
		{Title: "Reveal", DurationSeconds: 5, Visual: "Camera rises into a clean hero composition", Camera: "controlled orbit"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(storyboard.Shots) != 2 || storyboard.Shots[0].BriefID == "" || storyboard.Shots[0].BriefID == brief.ID {
		t.Fatalf("storyboard shots = %#v", storyboard.Shots)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	unapprovedShot, err := store.GetBrief(storyboard.Shots[0].BriefID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SubmitPreview(context.Background(), unapprovedShot); err == nil {
		t.Fatal("shot preview was generated before storyboard approval")
	}
	if api.posts != 0 {
		t.Fatalf("unapproved storyboard made %d live calls", api.posts)
	}
	if _, err := store.DecideStoryboard(brief.ID, "approve"); err != nil {
		t.Fatal(err)
	}
	view, err := store.StoryboardView(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.NextAction != "generate_shot_previews" || len(view.Shots) != 2 {
		t.Fatalf("approved storyboard view = %#v", view)
	}
	for _, shot := range view.Shots {
		if shot.Turn.NextAction != "generate_preview" || shot.Turn.CanSubmit {
			t.Fatalf("shot did not inherit preview gate: %#v", shot)
		}
		if shot.Turn.Brief.MasterBriefID != brief.ID || shot.Turn.Brief.ProductionMode != ProductionModeSingleShot {
			t.Fatalf("shot lineage = %#v", shot.Turn.Brief)
		}
	}

	if _, err := service.Submit(context.Background(), brief); err == nil {
		t.Fatal("storyboard master was submitted as one prompt-only video")
	}
	if api.posts != 0 {
		t.Fatalf("blocked storyboard master made %d live calls", api.posts)
	}
}

func TestScriptAndCreativeEditsInvalidateDownstreamApproval(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A product story", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text",
		Style: "studio", ProductionMode: ProductionModeStoryboard,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetScript(brief.ID, ScriptInput{Content: "First version"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DecideScript(brief.ID, "approve"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetStoryboard(brief.ID, StoryboardInput{Shots: []StoryboardShotInput{{DurationSeconds: 5, Visual: "Product on a plinth"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DecideStoryboard(brief.ID, "approve"); err != nil {
		t.Fatal(err)
	}

	brief, _ = store.GetBrief(brief.ID)
	brief.Style = "handheld"
	Refresh(brief)
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	if got := TurnFor(brief).NextAction; got != "revise_script" {
		t.Fatalf("creative edit next action = %q", got)
	}
	if _, err := store.DecideStoryboard(brief.ID, "approve"); err == nil {
		t.Fatal("stale storyboard approval was accepted")
	}

	if _, err := store.SetScript(brief.ID, ScriptInput{Content: "Second version"}); err != nil {
		t.Fatal(err)
	}
	brief, _ = store.GetBrief(brief.ID)
	if brief.StoryboardID != "" || brief.StoryboardStatus != "" || TurnFor(brief).NextAction != "review_script" {
		t.Fatalf("script edit did not invalidate storyboard: %#v", brief)
	}
}

func TestStoryboardShotDurationsMustMatchMaster(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A short story", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 10, AudioMode: "off", VideoMode: "text",
		Style: "cinematic", ProductionMode: ProductionModeStoryboard,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = store.SaveBrief(brief)
	_, _ = store.SetScript(brief.ID, ScriptInput{Content: "A ten second story"})
	_, _ = store.DecideScript(brief.ID, "approve")
	_, err = store.SetStoryboard(brief.ID, StoryboardInput{Shots: []StoryboardShotInput{{DurationSeconds: 5, Visual: "Only half"}}})
	if err == nil {
		t.Fatal("duration mismatch was accepted")
	}
}

func TestModifiedApprovedArtifactsLockStoryboardShots(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A product story", MediaType: "video", Purpose: "launch", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text",
		Style: "studio", ProductionMode: ProductionModeStoryboard,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	_, _ = store.SetScript(brief.ID, ScriptInput{Content: "Approved script"})
	_, _ = store.DecideScript(brief.ID, "approve")
	storyboard, err := store.SetStoryboard(brief.ID, StoryboardInput{Shots: []StoryboardShotInput{{DurationSeconds: 5, Visual: "Approved visual"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.DecideStoryboard(brief.ID, "approve"); err != nil {
		t.Fatal(err)
	}

	script, _ := store.GetScript(brief.ID)
	script.Content = "Changed behind the approval state"
	if err := store.SaveScript(script); err != nil {
		t.Fatal(err)
	}
	view, err := store.StoryboardView(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.NextAction != "revise_storyboard" {
		t.Fatalf("modified script next action = %q", view.NextAction)
	}
	child, _ := store.GetBrief(storyboard.Shots[0].BriefID)
	api := &fakeAPI{}
	if _, err := (&Service{API: api, Store: store}).SubmitPreview(context.Background(), child); err == nil {
		t.Fatal("shot preview accepted modified approved script")
	}
	if api.posts != 0 {
		t.Fatalf("modified script caused %d live calls", api.posts)
	}

	// Restore the script through the public workflow, then mutate the approved
	// storyboard content without updating its approval hash.
	_, _ = store.SetScript(brief.ID, ScriptInput{Content: "Approved script again"})
	_, _ = store.DecideScript(brief.ID, "approve")
	storyboard, _ = store.SetStoryboard(brief.ID, StoryboardInput{Shots: []StoryboardShotInput{{DurationSeconds: 5, Visual: "Approved visual again"}}})
	_, _ = store.DecideStoryboard(brief.ID, "approve")
	storyboard.Shots[0].Visual = "Changed behind the approval state"
	if err := store.SaveStoryboard(storyboard); err != nil {
		t.Fatal(err)
	}
	view, err = store.StoryboardView(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.NextAction != "revise_storyboard" {
		t.Fatalf("modified storyboard next action = %q", view.NextAction)
	}
	if _, err := store.DecideStoryboard(brief.ID, "approve"); err == nil {
		t.Fatal("modified storyboard was re-approved without being saved through SetStoryboard")
	}
}

func TestFailedPreviewWithResultURLCannotBeApproved(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "A failed preview", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text", Style: "clean",
	})
	if err != nil {
		t.Fatal(err)
	}
	brief.PreviewGenerationID = "gen_failed"
	brief.PreviewStatus = "failed"
	brief.PreviewURL = "https://cdn.example.test/error-output.png"
	brief.PreviewBriefHash = previewBriefFingerprint(brief)
	if err := ApproveVideoPreview(brief); err == nil {
		t.Fatal("failed preview with a result URL was approved")
	}
	if VideoPreviewApproved(brief) || TurnFor(brief).NextAction != "regenerate_preview" {
		t.Fatalf("failed preview turn = %#v", TurnFor(brief))
	}
}

func TestPreviewDedupeScansAllMatchingFingerprintGenerations(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A preview", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text", Style: "clean",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	matching := &Generation{
		ID: "gen_matching", BriefID: brief.ID, Kind: GenerationKindPreview,
		Fingerprint: previewBriefFingerprint(brief), TaskID: "task_matching", Model: "gpt-image-2-text-to-image",
		Status: "submitted", CreatedAt: now.Add(-time.Minute), UpdatedAt: now.Add(-time.Minute),
	}
	newerDifferent := &Generation{
		ID: "gen_newer", BriefID: brief.ID, Kind: GenerationKindPreview,
		Fingerprint: "different", TaskID: "task_newer", Model: "gpt-image-2-text-to-image",
		Status: "submitted", CreatedAt: now, UpdatedAt: now,
	}
	if err := store.SaveGeneration(matching); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveGeneration(newerDifferent); err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{}
	if _, err := (&Service{API: api, Store: store}).SubmitPreview(context.Background(), brief); err == nil {
		t.Fatal("duplicate preview was submitted while an older matching generation was active")
	}
	if api.posts != 0 {
		t.Fatalf("dedupe failure made %d live calls", api.posts)
	}
}
