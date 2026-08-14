// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func paidConfirmationForTest(t *testing.T, store *Store, brief *Brief, plan *Plan, scope, kind string) string {
	t.Helper()
	confirmation, err := NewPaidConfirmation(brief, plan, PaidConfirmationRequest{
		Scope: scope, GenerationKind: kind, ConfirmedBy: "test-user", Acknowledged: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SavePaidConfirmation(confirmation); err != nil {
		t.Fatal(err)
	}
	return confirmation.ID
}

func TestPaidConfirmationIsExactExpiringAndSingleUse(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{Request: "Product hero", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "studio", References: []string{"https://example.test/product.png"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	confirmation, err := NewPaidConfirmation(brief, brief.Plan, PaidConfirmationRequest{
		Scope: PaidScopeFinal, GenerationKind: GenerationKindFinal, ConfirmedBy: "cli-flag", Acknowledged: true, Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SavePaidConfirmation(confirmation); err != nil {
		t.Fatal(err)
	}
	if _, err := store.UsePaidConfirmation(confirmation.ID, brief, brief.Plan, PaidScopePreview, GenerationKindPreview, now); err == nil {
		t.Fatal("wrong-scope confirmation was accepted")
	}
	if _, err := store.UsePaidConfirmation(confirmation.ID, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal, now.Add(11*time.Minute)); err == nil {
		t.Fatal("expired confirmation was accepted")
	}
	if _, err := store.UsePaidConfirmation(confirmation.ID, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal, now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.UsePaidConfirmation(confirmation.ID, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal, now); err == nil || !strings.Contains(err.Error(), "already used") {
		t.Fatalf("reused confirmation error = %v", err)
	}
}

func TestPaidConfirmationClaimIsAtomicAcrossStoreInstances(t *testing.T) {
	root := filepath.Join(t.TempDir(), "media")
	brief, err := NewBrief(BriefInput{Request: "Product hero", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "studio", References: []string{"https://example.test/product.png"}})
	if err != nil {
		t.Fatal(err)
	}
	store := NewStore(root)
	confirmationID := paidConfirmationForTest(t, store, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal)
	start := make(chan struct{})
	var successes atomic.Int32
	var reused atomic.Int32
	var wait sync.WaitGroup
	for range 8 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			_, useErr := NewStore(root).UsePaidConfirmation(confirmationID, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal, time.Now())
			switch {
			case useErr == nil:
				successes.Add(1)
			case strings.Contains(useErr.Error(), "already used"):
				reused.Add(1)
			default:
				t.Errorf("unexpected claim error: %v", useErr)
			}
		}()
	}
	close(start)
	wait.Wait()
	if successes.Load() != 1 || reused.Load() != 7 {
		t.Fatalf("atomic confirmation claims = %d success, %d reused; want 1 and 7", successes.Load(), reused.Load())
	}
}

func TestPaidConfirmationRejectsUnhashablePlan(t *testing.T) {
	brief, err := NewBrief(BriefInput{Request: "Product hero", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "studio", References: []string{"https://example.test/product.png"}})
	if err != nil {
		t.Fatal(err)
	}
	plan := &Plan{Model: "test/model", Input: map[string]any{"invalid": func() {}}}
	if _, err := NewPaidConfirmation(brief, plan, PaidConfirmationRequest{Scope: PaidScopeFinal, GenerationKind: GenerationKindFinal, ConfirmedBy: "test", Acknowledged: true}); err == nil {
		t.Fatal("unhashable paid plan was accepted")
	}
}

func TestTurnExposesExactPaidActionReviewForEveryDirectorSpend(t *testing.T) {
	image, err := NewBrief(BriefInput{Request: "Product hero", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "studio", References: []string{"https://example.test/product.png"}})
	if err != nil {
		t.Fatal(err)
	}
	assertPaidActionMatchesPlan(t, TurnFor(image).PaidAction, BuildPlan(image), PaidScopeFinal, GenerationKindFinal)

	video, err := NewBrief(BriefInput{Request: "Product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "studio"})
	if err != nil {
		t.Fatal(err)
	}
	assertPaidActionMatchesPlan(t, TurnFor(video).PaidAction, BuildPreviewPlan(video), PaidScopePreview, GenerationKindPreview)
	video.PreviewGenerationID = "gen_preview"
	video.PreviewStatus = "success"
	video.PreviewURL = "https://cdn.example.test/preview.png"
	video.PreviewBriefHash = previewBriefFingerprint(video)
	if err := ApproveVideoPreview(video); err != nil {
		t.Fatal(err)
	}
	proofPlan, _, err := BuildProofPlan(video)
	if err != nil {
		t.Fatal(err)
	}
	assertPaidActionMatchesPlan(t, TurnFor(video).PaidAction, proofPlan, PaidScopeProof, GenerationKindProof)
	if err := SkipVideoProof(video); err != nil {
		t.Fatal(err)
	}
	assertPaidActionMatchesPlan(t, TurnFor(video).PaidAction, BuildPlan(video), PaidScopeFinal, GenerationKindFinal)
}

func TestFinalSubmissionRequiresProofDecisionAtServiceBoundary(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "Product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "studio",
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
	if turn := TurnFor(brief); turn.NextAction != "offer_proof" || turn.PaidAction == nil || turn.PaidAction.Scope != PaidScopeProof {
		t.Fatalf("preview-approved turn = %#v", turn)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	confirmationID := paidConfirmationForTest(t, store, brief, BuildPlan(brief), PaidScopeFinal, GenerationKindFinal)
	api := &fakeAPI{}
	if _, err := (&Service{API: api, Store: store}).Submit(context.Background(), brief, confirmationID); err == nil || !strings.Contains(err.Error(), "proof decision") {
		t.Fatalf("final submission before proof decision err=%v", err)
	}
	if api.posts != 0 {
		t.Fatalf("final submission made %d live calls before proof decision", api.posts)
	}
}

func assertPaidActionMatchesPlan(t *testing.T, action *PaidActionReview, plan *Plan, scope, kind string) {
	t.Helper()
	wantHash, err := PlanFingerprint(plan)
	if err != nil {
		t.Fatal(err)
	}
	if action == nil || action.Scope != scope || action.GenerationKind != kind || action.Model != plan.Model || action.PlanHash != wantHash || action.BlockedReason != "" {
		t.Fatalf("paid action = %#v, want scope=%s kind=%s model=%s hash=%s", action, scope, kind, plan.Model, wantHash)
	}
}

func TestPaidConfirmationBecomesStaleWhenCreativeStateChanges(t *testing.T) {
	brief, err := NewBrief(BriefInput{Request: "A runner", MediaType: "video", Purpose: "ad", Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "cinematic"})
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPreviewPlan(brief)
	confirmation, err := NewPaidConfirmation(brief, plan, PaidConfirmationRequest{Scope: PaidScopePreview, GenerationKind: GenerationKindPreview, ConfirmedBy: "mcp-user", Acknowledged: true})
	if err != nil {
		t.Fatal(err)
	}
	brief.Style = "documentary"
	if err := ValidatePaidConfirmation(confirmation, brief, BuildPreviewPlan(brief), PaidScopePreview, GenerationKindPreview, time.Now()); err == nil {
		t.Fatal("stale creative confirmation was accepted")
	}
}

func TestProofOptionUsesModelLowestDocumentedTier(t *testing.T) {
	option := ResolveProofOption("bytedance/seedance-2-5")
	if !option.Supported || option.ResolutionField != "resolution" || option.ResolutionValue != "480p" || !option.SameModel {
		t.Fatalf("Seedance proof option = %#v", option)
	}
	if option.ExactCostKnown || !strings.Contains(strings.ToLower(option.Disclosure), "paid") {
		t.Fatalf("proof cost disclosure = %#v", option)
	}
	unsupported := ResolveProofOption("gpt-image-2-text-to-image")
	if unsupported.Supported || unsupported.UnsupportedReason == "" {
		t.Fatalf("image proof unexpectedly supported: %#v", unsupported)
	}
}

func TestOldBriefFixtureLoadsWithConservativeNewState(t *testing.T) {
	root := filepath.Join(t.TempDir(), "media")
	store := NewStore(root)
	created := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	legacy := map[string]any{
		"id": "brief_legacy", "request": "Legacy image", "media_type": "image", "purpose": "site",
		"platform": "website", "aspect_ratio": "16:9", "style": "clean", "references_complete": true,
		"identity_complete": true, "status": "ready", "created_at": created, "updated_at": created,
	}
	data, _ := json.Marshal(legacy)
	path := filepath.Join(root, "briefs", "brief_legacy.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	brief, err := store.GetBrief("brief_legacy")
	if err != nil {
		t.Fatal(err)
	}
	Refresh(brief)
	if brief.Status != StatusReady || brief.Plan == nil || brief.GenerationID != "" || brief.PreviewApprovedAt != nil {
		t.Fatalf("legacy brief migration = %#v", brief)
	}
}

func TestInferenceAndWrapUpStayConciseAndExposeRoute(t *testing.T) {
	brief, err := NewBrief(BriefInput{Request: "Create a 5 second vertical silent text-to-video TikTok video ad in a cinematic documentary style"})
	if err != nil {
		t.Fatal(err)
	}
	if err := InferBrief(brief); err != nil {
		t.Fatal(err)
	}
	if question := NextQuestion(brief); question != nil {
		t.Fatalf("complete request asked an unnecessary question: %#v", question)
	}
	if brief.Plan == nil || brief.Plan.ProductionSkill != "kie-video" || brief.Plan.CapabilitySkill != "kie-video" || brief.Plan.CostStatus == "" || len(brief.Plan.OverrideOptions) < 2 {
		t.Fatalf("inferred route = %#v", brief.Plan)
	}

	incomplete, err := NewBrief(BriefInput{Request: "A product launch visual"})
	if err != nil {
		t.Fatal(err)
	}
	if question := NextQuestion(incomplete); question == nil || question.Recommendation == "" || question.RecommendationReason == "" {
		t.Fatalf("material question lacks recommendation rationale: %#v", question)
	}
	if err := WrapUpBrief(incomplete); err != nil {
		t.Fatal(err)
	}
	if NextQuestion(incomplete) != nil || incomplete.Plan == nil {
		t.Fatalf("wrap up did not produce an inspectable plan: %#v", incomplete)
	}
}

func TestCompleteShotProofLifecycleAndSeparateFinalConfirmation(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A complete product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "studio",
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
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	plan, option, err := BuildProofPlan(brief)
	if err != nil || option.ResolutionValue != "480p" || plan.Input["resolution"] != "480p" {
		t.Fatalf("proof plan = %#v option=%#v err=%v", plan, option, err)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	confirmationID := paidConfirmationForTest(t, store, brief, plan, PaidScopeProof, GenerationKindProof)
	proof, err := service.SubmitProof(context.Background(), brief, confirmationID)
	if err != nil {
		t.Fatal(err)
	}
	if proof.Kind != GenerationKindProof || api.posts != 1 {
		t.Fatalf("proof=%#v posts=%d", proof, api.posts)
	}
	if _, err := service.RefreshGeneration(context.Background(), proof.ID); err != nil {
		t.Fatal(err)
	}
	brief, _ = store.GetBrief(brief.ID)
	if err := ApproveVideoProof(brief); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	if turn := TurnFor(brief); !turn.CanSubmit || turn.NextAction != "review_then_submit" {
		t.Fatalf("approved proof turn = %#v", turn)
	}
	if _, err := service.Submit(context.Background(), brief); err == nil || api.posts != 1 {
		t.Fatalf("final without separate confirmation err=%v posts=%d", err, api.posts)
	}
	brief.Style = "handheld"
	Refresh(brief)
	if VideoProofApproved(brief) || brief.ProofStatus != "" {
		t.Fatalf("creative change did not invalidate proof: %#v", brief)
	}
}

func TestAvatarRouteRequiresExplicitRightsAcknowledgement(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "A presenter speaks", MediaType: "video", Purpose: "explainer", Platform: "website",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "on", VideoMode: "multimodal", Style: "studio",
		References: []string{"https://example.test/presenter.png"}, Model: "kling/ai-avatar-pro",
	})
	if err != nil {
		t.Fatal(err)
	}
	if question := NextQuestion(brief); question == nil || question.Key != "rights_acknowledged" {
		t.Fatalf("rights question = %#v", question)
	}
	if err := ApplyNextAnswer(brief, "no"); err == nil {
		t.Fatal("rights denial was accepted")
	}
	if err := ApplyNextAnswer(brief, "yes"); err != nil {
		t.Fatal(err)
	}
}
