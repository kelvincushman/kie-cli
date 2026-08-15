// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"context"
	"encoding/json"
	"fmt"
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
	if _, err := store.UsePaidConfirmation(confirmation.ID, brief, brief.Plan, PaidScopeFinal, GenerationKindFinal, now.Add(paidConfirmationTTL+time.Minute)); err == nil {
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
	storedConfirmation, err := store.GetPaidConfirmation(confirmationID)
	if err != nil {
		t.Fatal(err)
	}
	if storedConfirmation.UsedAt != nil {
		t.Fatal("proof-decision rejection consumed the final paid confirmation")
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

func TestPaidConfirmationRejectsEveryCreativeAndTransactionMismatch(t *testing.T) {
	base, err := NewBrief(BriefInput{
		Request: "A runner crosses a neon finish line", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "cinematic",
	})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		mutate func(*Brief, *Plan, *PaidConfirmation, *string, *string, *time.Time)
	}{
		{name: "prompt", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) { b.Request += " at sunrise" }},
		{name: "model", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) {
			b.Model = "pixverse-v6/text-to-video"
		}},
		{name: "duration", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) { b.DurationSeconds = 10 }},
		{name: "resolution", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) { b.Resolution = "480p" }},
		{name: "references", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) {
			b.References = []string{"https://example.test/runner.png"}
		}},
		{name: "identity", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) {
			b.IdentityIDs = []string{"identity_1"}
		}},
		{name: "consent", mutate: func(b *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) { b.RightsAcknowledged = true }},
		{name: "plan hash", mutate: func(_ *Brief, p *Plan, _ *PaidConfirmation, _, _ *string, _ *time.Time) { p.Input["duration"] = 15 }},
		{name: "wrong model record", mutate: func(_ *Brief, _ *Plan, c *PaidConfirmation, _, _ *string, _ *time.Time) { c.Model = "other/model" }},
		{name: "wrong brief", mutate: func(_ *Brief, _ *Plan, c *PaidConfirmation, _, _ *string, _ *time.Time) { c.BriefID = "brief_other" }},
		{name: "wrong scope", mutate: func(_ *Brief, _ *Plan, _ *PaidConfirmation, scope, _ *string, _ *time.Time) { *scope = PaidScopeFinal }},
		{name: "wrong generation kind", mutate: func(_ *Brief, _ *Plan, _ *PaidConfirmation, _, kind *string, _ *time.Time) {
			*kind = GenerationKindFinal
		}},
		{name: "expired", mutate: func(_ *Brief, _ *Plan, _ *PaidConfirmation, _, _ *string, at *time.Time) {
			*at = now.Add(paidConfirmationTTL)
		}},
		{name: "reused", mutate: func(_ *Brief, _ *Plan, c *PaidConfirmation, _, _ *string, _ *time.Time) {
			used := now
			c.UsedAt = &used
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			brief := cloneBriefForContractTest(t, base)
			plan := BuildPreviewPlan(brief)
			confirmation, err := NewPaidConfirmation(brief, plan, PaidConfirmationRequest{
				Scope: PaidScopePreview, GenerationKind: GenerationKindPreview, ConfirmedBy: "test", Acknowledged: true, Now: now,
			})
			if err != nil {
				t.Fatal(err)
			}
			scope, kind, at := PaidScopePreview, GenerationKindPreview, now
			test.mutate(brief, plan, confirmation, &scope, &kind, &at)
			if err := ValidatePaidConfirmation(confirmation, brief, plan, scope, kind, at); err == nil {
				t.Fatal("mismatched paid confirmation was accepted")
			}
		})
	}
}

func TestInvalidPaidConfirmationsMakeZeroAPICalls(t *testing.T) {
	mutations := map[string]func(*Brief, *PaidConfirmation){
		"stale brief": func(brief *Brief, _ *PaidConfirmation) { brief.Style = "documentary" },
		"wrong model": func(_ *Brief, confirmation *PaidConfirmation) { confirmation.Model = "other/model" },
		"wrong plan":  func(_ *Brief, confirmation *PaidConfirmation) { confirmation.PlanHash = "wrong" },
		"wrong scope": func(_ *Brief, confirmation *PaidConfirmation) { confirmation.Scope = PaidScopePreview },
		"wrong brief": func(_ *Brief, confirmation *PaidConfirmation) { confirmation.BriefID = "brief_other" },
		"expired":     func(_ *Brief, confirmation *PaidConfirmation) { confirmation.ExpiresAt = time.Now().Add(-time.Minute) },
		"reused":      func(_ *Brief, confirmation *PaidConfirmation) { used := time.Now().UTC(); confirmation.UsedAt = &used },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			store := NewStore(filepath.Join(t.TempDir(), "media"))
			brief, err := NewBrief(BriefInput{Request: "Product hero", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "studio"})
			if err != nil {
				t.Fatal(err)
			}
			plan := BuildPlan(brief)
			confirmation, err := NewPaidConfirmation(brief, plan, PaidConfirmationRequest{
				Scope: PaidScopeFinal, GenerationKind: GenerationKindFinal, ConfirmedBy: "test", Acknowledged: true,
			})
			if err != nil {
				t.Fatal(err)
			}
			mutate(brief, confirmation)
			Refresh(brief)
			if err := store.SaveBrief(brief); err != nil {
				t.Fatal(err)
			}
			if err := store.SavePaidConfirmation(confirmation); err != nil {
				t.Fatal(err)
			}
			api := &fakeAPI{}
			if _, err := (&Service{API: api, Store: store}).Submit(context.Background(), brief, confirmation.ID); err == nil {
				t.Fatal("invalid confirmation was accepted")
			}
			if api.posts != 0 {
				t.Fatalf("invalid confirmation made %d API calls", api.posts)
			}
		})
	}
}

func cloneBriefForContractTest(t *testing.T, source *Brief) *Brief {
	t.Helper()
	data, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}
	var cloned Brief
	if err := json.Unmarshal(data, &cloned); err != nil {
		t.Fatal(err)
	}
	return &cloned
}

func TestProofOptionUsesEveryDocumentedTierShape(t *testing.T) {
	tests := []struct {
		model, field, value, tier string
		noCheaper                 bool
	}{
		{model: "pixverse-v6/text-to-video", field: "quality", value: "360p", tier: "360p"},
		{model: "bytedance/seedance-2-5", field: "resolution", value: "480p", tier: "480p"},
		{model: "bytedance/v1-pro-fast-image-to-video", field: "resolution", value: "720p", tier: "720p", noCheaper: true},
		{model: "kling-3.0/video", field: "mode", value: "std", tier: "std"},
		{model: "omnihuman-1-5", field: "output_resolution", value: "720", tier: "720", noCheaper: true},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			option := ResolveProofOption(test.model)
			if !option.Supported || option.ResolutionField != test.field || option.ResolutionValue != test.value || option.LowestTier != test.tier || !option.SameModel {
				t.Fatalf("proof option = %#v", option)
			}
			if option.ExactCostKnown || !strings.Contains(strings.ToLower(option.Disclosure), "paid") || !strings.Contains(strings.ToLower(option.Disclosure), "exact cost") {
				t.Fatalf("proof cost disclosure = %#v", option)
			}
			if test.noCheaper && !strings.Contains(strings.ToLower(option.Disclosure), "no cheaper faithful same-model tier") {
				t.Fatalf("minimum-tier disclosure = %#v", option)
			}
		})
	}
	option := ResolveProofOption("bytedance/v1-pro-fast-image-to-video")
	if !strings.Contains(strings.ToLower(option.AlternateLabel), "guidance only") || !strings.Contains(strings.ToLower(option.AlternateLabel), "differ") {
		t.Fatalf("alternate proof label = %#v", option)
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
	if brief.ProofOption != nil || brief.ProofStatus != "" || brief.ProofApprovedAt != nil || brief.ProofSkippedAt != nil {
		t.Fatalf("legacy image acquired non-conservative proof state: %#v", brief)
	}
}

func TestLegacyApprovedVideoPreviewResumesWithoutProofOrPaidAuthorization(t *testing.T) {
	root := filepath.Join(t.TempDir(), "media")
	store := NewStore(root)
	created := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	legacyBrief := &Brief{
		ID: "brief_legacy_video", Request: "Legacy product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "off", VideoMode: "text", Style: "studio",
		ReferencesComplete: true, IdentityComplete: true, Status: StatusReady, CreatedAt: created, UpdatedAt: created,
		PreviewGenerationID: "gen_legacy_preview", PreviewStatus: "success", PreviewURL: "https://cdn.example.test/legacy-preview.png",
	}
	approved := created.Add(time.Minute)
	legacyBrief.PreviewApprovedAt = &approved
	legacyBrief.PreviewBriefHash = previewBriefFingerprint(legacyBrief)
	data, err := json.Marshal(legacyBrief)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "briefs", legacyBrief.ID+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	brief, err := store.GetBrief(legacyBrief.ID)
	if err != nil {
		t.Fatal(err)
	}
	Refresh(brief)
	turn := TurnFor(brief)
	if !VideoPreviewApproved(brief) || turn.NextAction != "offer_proof" || turn.PaidAction == nil || turn.PaidAction.Scope != PaidScopeProof {
		t.Fatalf("legacy video resume turn = %#v", turn)
	}
	if brief.ProofStatus != "" || brief.ProofApprovedAt != nil || brief.ProofSkippedAt != nil {
		t.Fatalf("legacy video acquired proof approval: %#v", brief)
	}
	if _, err := store.GetPaidConfirmation("confirm_missing"); err == nil {
		t.Fatal("legacy video unexpectedly acquired a paid confirmation")
	}
	brief.Style = "documentary"
	Refresh(brief)
	if VideoPreviewApproved(brief) || TurnFor(brief).NextAction != "generate_preview" {
		t.Fatalf("changed legacy preview remained approved: %#v", TurnFor(brief))
	}
}

func TestLegacyIdentityDefaultsToBlockedRightsScope(t *testing.T) {
	root := filepath.Join(t.TempDir(), "media")
	identity := map[string]any{
		"id": "identity_legacy", "name": "Legacy presenter", "image_references": []string{"ref:legacy"},
		"created_at": time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC), "updated_at": time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC),
	}
	data, err := json.Marshal(identity)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "identities", "identity_legacy.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(root)
	loaded, err := store.GetIdentity("identity_legacy")
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.ConsentConfirmedAt.IsZero() {
		t.Fatalf("legacy identity acquired consent timestamp: %#v", loaded)
	}
	brief, err := NewBrief(BriefInput{
		Request: "Presenter explains the product", MediaType: "video", Purpose: "explainer", Platform: "website",
		AspectRatio: "16:9", DurationSeconds: 5, Resolution: "720p", AudioMode: "on", VideoMode: "multimodal", Style: "studio",
		IdentityIDs: []string{loaded.ID}, Model: "omnihuman-1-5",
	})
	if err != nil {
		t.Fatal(err)
	}
	turn := TurnFor(brief)
	if turn.NextQuestion == nil || turn.NextQuestion.Key != "rights_acknowledged" || turn.PaidAction != nil || turn.Ready {
		t.Fatalf("legacy identity rights gate = %#v", turn)
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
	brief, err = store.GetBrief(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
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
		References: []string{"https://example.test/presenter.png"}, ReferenceAudio: []string{"https://example.test/voice.mp3"}, Model: "kling/ai-avatar-pro",
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
	if brief.Plan == nil || brief.Plan.Input["image_url"] != "https://example.test/presenter.png" || brief.Plan.Input["audio_url"] != "https://example.test/voice.mp3" {
		t.Fatalf("avatar plan inputs = %#v", brief.Plan)
	}
}

func TestCapabilityRoutedVoiceModelRequiresRightsAcknowledgement(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "A narrated product walkthrough", MediaType: "video", Purpose: "explainer", Platform: "website",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "on", VideoMode: "text", Style: "studio",
		Model: "elevenlabs/text-to-speech-turbo-2-5",
	})
	if err != nil {
		t.Fatal(err)
	}
	if question := NextQuestion(brief); question == nil || question.Key != "rights_acknowledged" {
		t.Fatalf("voice rights question = %#v", question)
	}
}

func TestInferenceLeavesUnsupportedDurationUnsetAndExtractsStyleTokens(t *testing.T) {
	single, err := NewBrief(BriefInput{
		Request: "Create a 60 second cinematic documentary product video", MediaType: "video", Purpose: "ad",
		Platform: "youtube", AspectRatio: "16:9", AudioMode: "off", VideoMode: "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := InferBrief(single); err != nil {
		t.Fatal(err)
	}
	if single.DurationSeconds != 0 {
		t.Fatalf("single-shot inferred duration = %d, want unset", single.DurationSeconds)
	}
	if question := NextQuestion(single); question == nil || question.Key != "duration_seconds" {
		t.Fatalf("next question = %#v, want duration", question)
	}
	if single.Style != "cinematic, documentary" {
		t.Fatalf("inferred style = %q", single.Style)
	}

	storyboard, err := NewBrief(BriefInput{
		Request: "Create a 60 second cinematic documentary product video", MediaType: "video", Purpose: "ad",
		Platform: "youtube", AspectRatio: "16:9", AudioMode: "off", VideoMode: "text", ProductionMode: ProductionModeStoryboard,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := InferBrief(storyboard); err != nil {
		t.Fatal(err)
	}
	if storyboard.DurationSeconds != 60 {
		t.Fatalf("storyboard inferred duration = %d, want 60", storyboard.DurationSeconds)
	}
	requestCount := strings.Count(storyboard.Plan.Input["prompt"].(string), storyboard.Request)
	if requestCount != 1 {
		t.Fatalf("request appears %d times in prompt %q", requestCount, storyboard.Plan.Input["prompt"])
	}
}

func TestProofTierSelectionSkipsUnrankableFields(t *testing.T) {
	properties := map[string]any{
		"quality":    map[string]any{"enum": []any{"balanced", "creative"}},
		"resolution": map[string]any{"enum": []any{"1080p", "720p"}},
	}
	field, value := selectProofTier([]string{"quality", "resolution"}, properties)
	if field != "resolution" || value != "720p" {
		t.Fatalf("selected proof tier = %s/%s", field, value)
	}
	field, value = selectProofTier([]string{"quality"}, properties)
	if field != "" || value != "" {
		t.Fatalf("unrankable proof tier = %s/%s", field, value)
	}
}

func TestProofUsesDocumentedModelReferenceInputs(t *testing.T) {
	t.Run("Kling expands identity images", func(t *testing.T) {
		store := NewStore(filepath.Join(t.TempDir(), "media"))
		identity, err := store.CreateIdentity("Presenter", []string{
			"https://example.test/front.png", "https://example.test/profile.png",
		}, true)
		if err != nil {
			t.Fatal(err)
		}
		brief, err := NewBrief(BriefInput{
			Request: "A consistent presenter reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
			AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "multimodal", Style: "studio",
			References: []string{"https://example.test/prop.png"}, IdentityIDs: []string{identity.ID}, Model: "kling-3.0/video",
		})
		if err != nil {
			t.Fatal(err)
		}
		approvePreviewForProofTest(t, brief)
		if err := store.SaveBrief(brief); err != nil {
			t.Fatal(err)
		}
		plan, _, err := BuildProofPlan(brief)
		if err != nil {
			t.Fatal(err)
		}
		api := &fakeAPI{}
		service := &Service{API: api, Store: store}
		confirmationID := paidConfirmationForTest(t, store, brief, plan, PaidScopeProof, GenerationKindProof)
		if _, err := service.SubmitProof(context.Background(), brief, confirmationID); err != nil {
			t.Fatal(err)
		}
		input := api.body["input"].(map[string]any)
		images, ok := input["image_urls"].([]string)
		wantImages := "https://cdn.example.test/preview.png,https://example.test/prop.png,https://example.test/front.png,https://example.test/profile.png"
		if !ok || strings.Join(images, ",") != wantImages {
			t.Fatalf("Kling proof images = %#v", input["image_urls"])
		}
	})

	t.Run("Wan remains text to video", func(t *testing.T) {
		store := NewStore(filepath.Join(t.TempDir(), "media"))
		brief, err := NewBrief(BriefInput{
			Request: "A landscape product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
			AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text", Style: "cinematic",
			Model: "wan/2-7-text-to-video",
		})
		if err != nil {
			t.Fatal(err)
		}
		approvePreviewForProofTest(t, brief)
		if err := store.SaveBrief(brief); err != nil {
			t.Fatal(err)
		}
		plan, _, err := BuildProofPlan(brief)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := plan.Input["image_urls"]; ok {
			t.Fatalf("Wan proof plan invented image_urls: %#v", plan.Input)
		}
		api := &fakeAPI{}
		service := &Service{API: api, Store: store}
		confirmationID := paidConfirmationForTest(t, store, brief, plan, PaidScopeProof, GenerationKindProof)
		if _, err := service.SubmitProof(context.Background(), brief, confirmationID); err != nil {
			t.Fatal(err)
		}
		input := api.body["input"].(map[string]any)
		if _, ok := input["image_urls"]; ok {
			t.Fatalf("Wan proof submitted image_urls: %#v", input)
		}
	})
}

func TestSeedanceProofRejectsMoreThanThirtyExpandedImages(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	identityImages := make([]string, 20)
	for i := range identityImages {
		identityImages[i] = fmt.Sprintf("https://example.test/identity-%02d.png", i)
	}
	identity, err := store.CreateIdentity("Presenter", identityImages, true)
	if err != nil {
		t.Fatal(err)
	}
	references := make([]string, 10)
	for i := range references {
		references[i] = fmt.Sprintf("https://example.test/reference-%02d.png", i)
	}
	brief, err := NewBrief(BriefInput{
		Request: "A reference-heavy film shot", MediaType: "video", Purpose: "film", Platform: "youtube",
		AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "multimodal", Style: "cinematic",
		References: references, IdentityIDs: []string{identity.ID},
	})
	if err != nil {
		t.Fatal(err)
	}
	approvePreviewForProofTest(t, brief)
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	plan, _, err := BuildProofPlan(brief)
	if err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	confirmationID := paidConfirmationForTest(t, store, brief, plan, PaidScopeProof, GenerationKindProof)
	if _, err := service.SubmitProof(context.Background(), brief, confirmationID); err == nil || !strings.Contains(err.Error(), "at most 30") {
		t.Fatalf("SeedDance expanded reference limit error = %v", err)
	}
	if api.posts != 0 {
		t.Fatalf("SeedDance over-limit proof made %d live calls", api.posts)
	}
}

func TestRefreshGenerationPreservesApprovalForUnchangedURL(t *testing.T) {
	for _, kind := range []string{GenerationKindPreview, GenerationKindProof} {
		t.Run(kind, func(t *testing.T) {
			store := NewStore(filepath.Join(t.TempDir(), "media"))
			brief, err := NewBrief(BriefInput{
				Request: "A product reveal", MediaType: "video", Purpose: "ad", Platform: "youtube",
				AspectRatio: "16:9", DurationSeconds: 5, AudioMode: "off", VideoMode: "text", Style: "studio",
			})
			if err != nil {
				t.Fatal(err)
			}
			now := time.Now().UTC()
			fingerprint := previewBriefFingerprint(brief)
			generation := &Generation{ID: newID("gen"), BriefID: brief.ID, Kind: kind, Fingerprint: fingerprint, TaskID: "task_123", Model: "bytedance/seedance-2-5", Status: "success", CreatedAt: now, UpdatedAt: now}
			if kind == GenerationKindPreview {
				brief.PreviewGenerationID = generation.ID
				brief.PreviewBriefHash = fingerprint
				brief.PreviewStatus = "success"
				brief.PreviewURL = "https://cdn.example.test/final.png"
				brief.PreviewApprovedAt = &now
			} else {
				approvePreviewForProofTest(t, brief)
				fingerprint = proofBriefFingerprint(brief)
				generation.Fingerprint = fingerprint
				brief.ProofGenerationID = generation.ID
				brief.ProofBriefHash = fingerprint
				brief.ProofStatus = "success"
				brief.ProofURL = "https://cdn.example.test/final.png"
				brief.ProofApprovedAt = &now
			}
			if err := store.SaveBrief(brief); err != nil {
				t.Fatal(err)
			}
			if err := store.SaveGeneration(generation); err != nil {
				t.Fatal(err)
			}
			service := &Service{API: &fakeAPI{}, Store: store}
			if _, err := service.RefreshGeneration(context.Background(), generation.ID); err != nil {
				t.Fatal(err)
			}
			refreshed, err := store.GetBrief(brief.ID)
			if err != nil {
				t.Fatal(err)
			}
			if kind == GenerationKindPreview && refreshed.PreviewApprovedAt == nil {
				t.Fatal("unchanged preview refresh cleared approval")
			}
			if kind == GenerationKindProof && refreshed.ProofApprovedAt == nil {
				t.Fatal("unchanged proof refresh cleared approval")
			}
		})
	}
}

func approvePreviewForProofTest(t *testing.T, brief *Brief) {
	t.Helper()
	brief.PreviewGenerationID = "gen_preview"
	brief.PreviewStatus = "success"
	brief.PreviewURL = "https://cdn.example.test/preview.png"
	brief.PreviewBriefHash = previewBriefFingerprint(brief)
	if err := ApproveVideoPreview(brief); err != nil {
		t.Fatal(err)
	}
}
