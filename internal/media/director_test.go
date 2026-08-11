// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var validPNG = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}

func TestDirectorAsksOneQuestionAndBuildsImagePlan(t *testing.T) {
	brief, err := NewBrief(BriefInput{Request: "A launch image for a coffee brand"})
	if err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{"media_type", "purpose", "platform", "aspect_ratio", "style", "reference"}
	answers := []string{"image", "Instagram launch", "instagram", "", "warm editorial photography", "skip"}
	for i, wantKey := range wantKeys {
		question := NextQuestion(brief)
		if question == nil || question.Key != wantKey {
			t.Fatalf("question %d = %#v, want key %q", i, question, wantKey)
		}
		if err := ApplyNextAnswer(brief, answers[i]); err != nil {
			t.Fatalf("answering %s: %v", wantKey, err)
		}
	}
	if question := NextQuestion(brief); question != nil {
		t.Fatalf("ready brief still has question: %#v", question)
	}
	if brief.Status != StatusReady || brief.Plan == nil {
		t.Fatalf("brief status/plan = %q/%#v", brief.Status, brief.Plan)
	}
	if brief.AspectRatio != "3:4" {
		t.Fatalf("recommended Instagram aspect = %q, want 3:4", brief.AspectRatio)
	}
	if brief.Plan.Model != "gpt-image-2-text-to-image" {
		t.Fatalf("model = %q", brief.Plan.Model)
	}
}

func TestDirectorCollectsMultipleReferencesUntilSkip(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "Animate this product", MediaType: "video", Purpose: "ad",
		Platform: "tiktok", AspectRatio: "9:16", DurationSeconds: 5,
		AudioMode: "off", VideoMode: "multimodal", Style: "handheld UGC",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, answer := range []string{"./front.png", "https://example.test/side.png"} {
		if question := NextQuestion(brief); question == nil || question.Key != "reference" {
			t.Fatalf("reference question missing before %q: %#v", answer, question)
		}
		if err := ApplyNextAnswer(brief, answer); err != nil {
			t.Fatal(err)
		}
	}
	if err := ApplyNextAnswer(brief, "done"); err != nil {
		t.Fatal(err)
	}
	if len(brief.References) != 2 || brief.Plan.Model != "bytedance/seedance-2-5" {
		t.Fatalf("references/model = %#v/%q", brief.References, brief.Plan.Model)
	}
}

func TestDirectorNormalizesTerminalDroppedReferencePath(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "Use this image", MediaType: "image", Purpose: "post",
		Platform: "general", AspectRatio: "1:1", Style: "natural",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := ApplyNextAnswer(brief, `'/tmp/My Reference.png'`); err != nil {
		t.Fatal(err)
	}
	if err := ApplyNextAnswer(brief, `C:\My\ Reference.png`); err != nil {
		t.Fatal(err)
	}
	if got := brief.References; len(got) != 2 || got[0] != "/tmp/My Reference.png" || got[1] != `C:\My Reference.png` {
		t.Fatalf("normalized references = %#v", got)
	}
}

func TestStoreVaultsPrivateReferenceAndPersistsBrief(t *testing.T) {
	root := t.TempDir()
	store := NewStore(filepath.Join(root, "media"))
	source := filepath.Join(root, "logo.png")
	if err := os.WriteFile(source, validPNG, 0o600); err != nil {
		t.Fatal(err)
	}
	ref, err := store.AddReference(source, "logo")
	if err != nil {
		t.Fatal(err)
	}
	if ref.Kind != "file" || ref.StoredPath == source {
		t.Fatalf("reference was not vaulted: %#v", ref)
	}
	info, err := os.Stat(ref.StoredPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("vaulted reference mode = %o, want 600", info.Mode().Perm())
	}
	brief, err := NewBrief(BriefInput{Request: "Use my logo", MediaType: "image", Purpose: "hero", Platform: "website", AspectRatio: "16:9", Style: "minimal", References: []string{"ref:" + ref.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.GetBrief(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != brief.ID || loaded.References[0] != "ref:"+ref.ID {
		t.Fatalf("loaded brief = %#v", loaded)
	}
}

type fakeAPI struct {
	uploads int
	posts   int
	body    map[string]any
	bodies  []map[string]any
}

func (f *fakeAPI) Get(context.Context, string, map[string]string) (json.RawMessage, error) {
	return json.RawMessage(`{"data":{"status":"success","resultUrls":["https://cdn.example.test/final.png"]}}`), nil
}

func (f *fakeAPI) PostWithParams(_ context.Context, _ string, _ map[string]string, body any) (json.RawMessage, int, error) {
	f.posts++
	f.body, _ = body.(map[string]any)
	f.bodies = append(f.bodies, f.body)
	return json.RawMessage(`{"data":{"taskId":"task_123"}}`), 200, nil
}

func (f *fakeAPI) PostMultipartWithParams(context.Context, string, map[string]string, map[string]string, map[string]string) (json.RawMessage, int, error) {
	f.uploads++
	return json.RawMessage(`{"data":{"downloadUrl":"https://uploads.example.test/reference.png"}}`), 200, nil
}

func TestServiceUploadsLocalReferenceSubmitsAndRefreshes(t *testing.T) {
	root := t.TempDir()
	store := NewStore(filepath.Join(root, "media"))
	path := filepath.Join(root, "reference.png")
	if err := os.WriteFile(path, validPNG, 0o600); err != nil {
		t.Fatal(err)
	}
	brief, err := NewBrief(BriefInput{
		Request: "A precise product image", MediaType: "image", Purpose: "shop",
		Platform: "website", AspectRatio: "1:1", Style: "studio light", References: []string{path},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	generation, err := service.Submit(context.Background(), brief)
	if err != nil {
		t.Fatal(err)
	}
	if api.uploads != 1 || generation.TaskID != "task_123" {
		t.Fatalf("uploads/task = %d/%q", api.uploads, generation.TaskID)
	}
	input := api.body["input"].(map[string]any)
	urls := input["input_urls"].([]string)
	if len(urls) != 1 || urls[0] != "https://uploads.example.test/reference.png" {
		t.Fatalf("submitted input_urls = %#v", urls)
	}
	refreshed, err := service.RefreshGeneration(context.Background(), generation.ID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Status != "success" || len(refreshed.ResultURLs) != 1 {
		t.Fatalf("refreshed generation = %#v", refreshed)
	}
	if _, err := service.Submit(context.Background(), brief); err == nil {
		t.Fatal("submitted brief was accepted for a second paid task")
	}
	turn := TurnFor(brief)
	if turn.Ready || turn.NextAction != "check_generation_status" {
		t.Fatalf("submitted turn = %#v", turn)
	}
}

func TestVideoRequiresGeneratedPreviewAndExplicitApproval(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "A runner launches a new shoe", MediaType: "video", Purpose: "campaign",
		Platform: "instagram", AspectRatio: "9:16", DurationSeconds: 5,
		AudioMode: "off", VideoMode: "text", Style: "cinematic sunrise",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	turn := TurnFor(brief)
	if !turn.Ready || turn.CanSubmit || turn.NextAction != "generate_preview" {
		t.Fatalf("initial video turn = %#v", turn)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	if _, err := service.Submit(context.Background(), brief); err == nil || !strings.Contains(err.Error(), "approved preview") {
		t.Fatalf("final video without preview error = %v", err)
	}
	if api.posts != 0 {
		t.Fatalf("blocked final video submitted %d paid tasks", api.posts)
	}

	preview, err := service.SubmitPreview(context.Background(), brief)
	if err != nil {
		t.Fatal(err)
	}
	if preview.Kind != GenerationKindPreview || preview.Model != "gpt-image-2-text-to-image" || api.posts != 1 {
		t.Fatalf("preview = %#v, posts=%d", preview, api.posts)
	}
	if got := TurnFor(brief).NextAction; got != "check_preview_status" {
		t.Fatalf("next action after preview submit = %q", got)
	}
	if _, err := service.RefreshGeneration(context.Background(), preview.ID); err != nil {
		t.Fatal(err)
	}
	brief, err = store.GetBrief(brief.ID)
	if err != nil {
		t.Fatal(err)
	}
	turn = TurnFor(brief)
	if turn.CanSubmit || turn.NextAction != "review_preview" || brief.PreviewURL != "https://cdn.example.test/final.png" {
		t.Fatalf("review turn = %#v", turn)
	}
	if err := ApproveVideoPreview(brief); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	turn = TurnFor(brief)
	if !turn.CanSubmit || turn.NextAction != "review_then_submit" {
		t.Fatalf("approved turn = %#v", turn)
	}
	final, err := service.Submit(context.Background(), brief)
	if err != nil {
		t.Fatal(err)
	}
	if final.Kind != GenerationKindFinal || api.posts != 2 {
		t.Fatalf("final generation = %#v, posts=%d", final, api.posts)
	}
	finalInput := api.bodies[1]["input"].(map[string]any)
	if got := finalInput["first_frame_url"]; got != brief.PreviewURL {
		t.Fatalf("final visual anchor = %#v, want %q", got, brief.PreviewURL)
	}
}

func TestVideoPreviewApprovalBecomesStaleAndRejectedPreviewCanRegenerate(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	brief, err := NewBrief(BriefInput{
		Request: "Product reveal", MediaType: "video", Purpose: "launch",
		Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 5,
		AudioMode: "off", VideoMode: "text", Style: "clean studio",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{}
	service := &Service{API: api, Store: store}
	preview, err := service.SubmitPreview(context.Background(), brief)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.RefreshGeneration(context.Background(), preview.ID); err != nil {
		t.Fatal(err)
	}
	brief, _ = store.GetBrief(brief.ID)
	if err := ApproveVideoPreview(brief); err != nil {
		t.Fatal(err)
	}
	brief.Style = "handheld documentary"
	Refresh(brief)
	if turn := TurnFor(brief); turn.CanSubmit || turn.NextAction != "generate_preview" {
		t.Fatalf("stale approval turn = %#v", turn)
	}
	brief.Style = "clean studio"
	if err := RejectVideoPreview(brief); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SubmitPreview(context.Background(), brief); err != nil {
		t.Fatalf("regenerating rejected preview: %v", err)
	}
	if api.posts != 2 {
		t.Fatalf("preview regenerations = %d, want 2", api.posts)
	}
}

func TestStoreRejectsNonImageLocalReference(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".env")
	if err := os.WriteFile(path, []byte("KIE_API_KEY=secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewStore(filepath.Join(root, "media")).AddReference(path, ""); err == nil {
		t.Fatal("non-image local reference was accepted")
	}
}

func TestVaultedBriefAndPublicReferenceRedactLocalPaths(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "private-logo.png")
	if err := os.WriteFile(path, validPNG, 0o600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(filepath.Join(root, "media"))
	brief, err := NewBrief(BriefInput{Request: "Use my logo", MediaType: "image", Purpose: "site", Platform: "website", AspectRatio: "16:9", Style: "clean", References: []string{path}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.VaultBriefReferences(brief); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(brief)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), root) || !strings.Contains(string(data), "ref:") {
		t.Fatalf("vaulted brief leaked path or omitted handle: %s", data)
	}
	refs, err := store.ListReferences()
	if err != nil || len(refs) != 1 {
		t.Fatalf("references = %#v, err=%v", refs, err)
	}
	publicData, _ := json.Marshal(refs[0].Public())
	if strings.Contains(string(publicData), root) {
		t.Fatalf("public reference leaked path: %s", publicData)
	}
}

func TestSubmissionLockRejectsConcurrentSubmit(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	release, err := store.acquireSubmission("brief_test")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.acquireSubmission("brief_test"); err == nil {
		t.Fatal("concurrent submission lock was accepted")
	}
	release()
	releaseAgain, err := store.acquireSubmission("brief_test")
	if err != nil {
		t.Fatalf("released submission lock remained blocked: %v", err)
	}
	releaseAgain()
}

func TestSeedDancePlanUsesExclusiveMultimodalInputs(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Request: "Creator presents a new camera", MediaType: "video", Purpose: "launch",
		Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 10,
		AudioMode: "on", VideoMode: "multimodal", Style: "cinematic studio",
		References:      []string{"https://example.test/person.jpg"},
		ReferenceVideos: []string{"https://example.test/motion.mp4"},
		ReferenceAudio:  []string{"https://example.test/voice.mp3"},
		Resolution:      "720p", OutputFormat: "mp4", ReturnLastFrame: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if brief.Status != StatusReady || brief.Plan == nil {
		t.Fatalf("brief not ready: %#v", brief)
	}
	if brief.Plan.Model != "bytedance/seedance-2-5" {
		t.Fatalf("model = %q", brief.Plan.Model)
	}
	if _, ok := brief.Plan.Input["first_frame_url"]; ok {
		t.Fatalf("multimodal plan included first_frame_url: %#v", brief.Plan.Input)
	}
	if got := brief.Plan.Input["generate_audio"]; got != true {
		t.Fatalf("generate_audio = %#v", got)
	}
	if got := brief.Plan.Input["reference_video_urls"].([]string); len(got) != 1 {
		t.Fatalf("reference videos = %#v", got)
	}
}

func TestSeedDanceRejectsMixedInputModesAndExcessVideoReferences(t *testing.T) {
	base := BriefInput{
		Request: "Animate the campaign", MediaType: "video", Purpose: "launch",
		Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 10,
		AudioMode: "off", VideoMode: "multimodal", Style: "cinematic",
	}
	withFrame := base
	withFrame.FirstFrame = "https://example.test/start.png"
	withFrame.References = []string{"https://example.test/reference.png"}
	if _, err := NewBrief(withFrame); err == nil || !strings.Contains(err.Error(), "cannot include first or last frames") {
		t.Fatalf("mixed SeedDance modes error = %v", err)
	}
	tooManyVideos := base
	for i := 0; i < 11; i++ {
		tooManyVideos.ReferenceVideos = append(tooManyVideos.ReferenceVideos, "https://example.test/reference.mp4")
	}
	if _, err := NewBrief(tooManyVideos); err == nil || !strings.Contains(err.Error(), "at most 10 reference video") {
		t.Fatalf("video reference limit error = %v", err)
	}
}

func TestStoreCreatesConsentedIdentityAndExpandsItOnSubmit(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "media"))
	if _, err := store.CreateIdentity("Creator", []string{"https://example.test/front.jpg"}, false); err == nil {
		t.Fatal("identity profile was accepted without consent")
	}
	identity, err := store.CreateIdentity("Creator", []string{
		"https://example.test/front.jpg", "https://example.test/profile.jpg",
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(identity.ImageReferences) != 2 || !strings.HasPrefix(identity.ImageReferences[0], "ref:") {
		t.Fatalf("identity references = %#v", identity.ImageReferences)
	}
	brief, err := NewBrief(BriefInput{
		Request: "Creator portrait", MediaType: "image", Purpose: "blog",
		Platform: "website", AspectRatio: "1:1", Style: "editorial",
		IdentityIDs: []string{identity.ID}, References: []string{"https://example.test/location.jpg"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{}
	if _, err := (&Service{API: api, Store: store}).Submit(context.Background(), brief); err != nil {
		t.Fatal(err)
	}
	urls := api.body["input"].(map[string]any)["input_urls"].([]string)
	if len(urls) != 3 {
		t.Fatalf("expanded identity urls = %#v", urls)
	}
}

func TestStoreRejectsReferenceTypeMismatch(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "image.png")
	if err := os.WriteFile(path, validPNG, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewStore(filepath.Join(root, "media")).AddReferenceTyped(path, "", "audio"); err == nil {
		t.Fatal("image was accepted as an audio reference")
	}
}

func TestStoreAcceptsCurrentSeedDanceVideoAndAudioFormats(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name      string
		extension string
		contents  []byte
		mediaType string
	}{
		{name: "Matroska video", extension: ".mkv", contents: []byte{0x1a, 0x45, 0xdf, 0xa3}, mediaType: "video"},
		{name: "AAC audio", extension: ".aac", contents: []byte{0xff, 0xf1, 0x50, 0x80}, mediaType: "audio"},
		{name: "M4A audio", extension: ".m4a", contents: []byte{0, 0, 0, 24, 'f', 't', 'y', 'p', 'M', '4', 'A', ' '}, mediaType: "audio"},
		{name: "Ogg audio", extension: ".ogg", contents: []byte("OggS\x00\x02"), mediaType: "audio"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(root, "reference"+test.extension)
			if err := os.WriteFile(path, test.contents, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := NewStore(filepath.Join(root, "media-"+strings.TrimPrefix(test.extension, "."))).AddReferenceTyped(path, "", test.mediaType); err != nil {
				t.Fatalf("documented %s reference rejected: %v", test.mediaType, err)
			}
		})
	}
}

func TestWorkflowCatalogAndBriefRouting(t *testing.T) {
	workflows := ListWorkflows()
	if len(workflows) != 8 {
		t.Fatalf("workflow count = %d, want 8", len(workflows))
	}
	for _, workflow := range workflows {
		for _, mediaType := range workflow.MediaTypes {
			if mediaType != "image" && mediaType != "video" {
				t.Fatalf("workflow %s advertises unsupported director media type %q", workflow.Name, mediaType)
			}
		}
	}
	workflow, err := GetWorkflow("kie-product-photoshoot")
	if err != nil || workflow.Name != "product-photoshoot" {
		t.Fatalf("workflow = %#v, err=%v", workflow, err)
	}
	brief, err := NewBrief(BriefInput{
		Workflow: "youtube-thumbnail", Request: "A truthful camera review thumbnail",
		MediaType: "image", Purpose: "YouTube", Platform: "youtube",
		AspectRatio: "16:9", Style: "clean, high contrast", References: []string{"https://example.test/camera.jpg"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if brief.Workflow != "youtube-thumbnail" || !strings.Contains(brief.Plan.Input["prompt"].(string), "Workflow: youtube-thumbnail") {
		t.Fatalf("workflow brief = %#v", brief)
	}
	if _, err := NewBrief(BriefInput{Workflow: "unknown", Request: "x"}); err == nil {
		t.Fatal("unknown workflow was accepted")
	}
	if _, err := NewBrief(BriefInput{
		Workflow: "video-explainer", Request: "Reusable style key", MediaType: "image",
		Purpose: "explainer continuity", Platform: "youtube", AspectRatio: "16:9", Style: "paper cutout",
		References: []string{"https://example.test/style.png"},
	}); err != nil {
		t.Fatalf("video-explainer style-key image was rejected: %v", err)
	}
}

func TestProductWorkflowSelectsSeedreamReferenceSchema(t *testing.T) {
	brief, err := NewBrief(BriefInput{
		Workflow: "product-photoshoot", Request: "Bottle hero", MediaType: "image",
		Purpose: "website", Platform: "website", AspectRatio: "16:9", Style: "studio",
		References: []string{"https://example.test/bottle.jpg"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if brief.Plan.Model != "seedream/5-pro-image-to-image" {
		t.Fatalf("model = %q", brief.Plan.Model)
	}
	if _, ok := brief.Plan.Input["image_urls"]; !ok {
		t.Fatalf("Seedream plan input = %#v", brief.Plan.Input)
	}
}
