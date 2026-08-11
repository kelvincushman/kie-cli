// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package mediamcp

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"kie-pp-cli/internal/media"
)

func TestHTTPNegotiatesLatestStatelessProtocolAndListsMediaTools(t *testing.T) {
	store := media.NewStore(filepath.Join(t.TempDir(), "media"))
	server := NewServer("test", &Dependencies{Store: func() (*media.Store, error) { return store, nil }})
	httpServer := httptest.NewServer(NewHTTPHandler(server))
	defer httpServer.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	if got := session.InitializeResult().ProtocolVersion; got != "2026-07-28" {
		t.Fatalf("protocol version = %q, want 2026-07-28", got)
	}

	listed, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	gotNames := make([]string, 0, len(listed.Tools))
	for _, tool := range listed.Tools {
		gotNames = append(gotNames, tool.Name)
		if tool.Name == "media_generate" || tool.Name == "media_preview_generate" {
			if tool.Annotations == nil || tool.Annotations.ReadOnlyHint || tool.Annotations.OpenWorldHint == nil || !*tool.Annotations.OpenWorldHint {
				t.Fatalf("media_generate annotations = %#v", tool.Annotations)
			}
		}
	}
	sort.Strings(gotNames)
	wantNames := append([]string(nil), ToolNames...)
	sort.Strings(wantNames)
	if stringSliceJSON(gotNames) != stringSliceJSON(wantNames) {
		t.Fatalf("tools = %v, want %v", gotNames, wantNames)
	}

	workflowResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "media_workflow_get", Arguments: map[string]any{"workflow": "product-photoshoot"},
	})
	if err != nil || workflowResult.IsError {
		t.Fatalf("media_workflow_get error=%v result=%#v", err, workflowResult)
	}
	var workflowOutput workflowGetOutput
	decodeStructured(t, workflowResult.StructuredContent, &workflowOutput)
	if workflowOutput.Workflow.Skill != "kie-product-photoshoot" {
		t.Fatalf("workflow output = %#v", workflowOutput)
	}

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "media_brief_start",
		Arguments: map[string]any{
			"request": "A polished coffee launch image", "media_type": "image",
			"purpose": "social launch", "platform": "instagram", "aspect_ratio": "9:16",
			"style": "warm editorial photography", "references": []string{"https://example.test/coffee.png"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("media_brief_start returned an error: %#v", result.Content)
	}
	var turn media.Turn
	decodeStructured(t, result.StructuredContent, &turn)
	if !turn.Ready || turn.Brief == nil || turn.Brief.ID == "" || turn.Brief.Plan == nil {
		t.Fatalf("ready turn = %#v", turn)
	}
}

func TestGenerateIsExplicitAndCannotResubmitBrief(t *testing.T) {
	store := media.NewStore(filepath.Join(t.TempDir(), "media"))
	api := &fakeAPI{}
	server := NewServer("test", &Dependencies{
		Store: func() (*media.Store, error) { return store, nil },
		LiveService: func(context.Context, *media.Store) (*media.Service, func(), error) {
			return &media.Service{API: api, Store: store}, func() {}, nil
		},
	})
	httpServer := httptest.NewServer(NewHTTPHandler(server))
	defer httpServer.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	brief, err := media.NewBrief(media.BriefInput{
		Request: "Product launch", MediaType: "image", Purpose: "website",
		Platform: "website", AspectRatio: "16:9", Style: "minimal", References: []string{"https://example.test/ref.png"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	call := func() *mcp.CallToolResult {
		result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "media_generate", Arguments: map[string]any{"brief_id": brief.ID}})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	first := call()
	if first.IsError || api.submissions != 1 {
		t.Fatalf("first generation result=%#v submissions=%d", first.Content, api.submissions)
	}
	second := call()
	if !second.IsError || api.submissions != 1 {
		t.Fatalf("resubmit result=%#v submissions=%d", second.Content, api.submissions)
	}
}

func TestMCPVideoGenerationRequiresDisplayedPreviewApproval(t *testing.T) {
	store := media.NewStore(filepath.Join(t.TempDir(), "media"))
	api := &fakeAPI{}
	server := NewServer("test", &Dependencies{
		Store: func() (*media.Store, error) { return store, nil },
		LiveService: func(context.Context, *media.Store) (*media.Service, func(), error) {
			return &media.Service{API: api, Store: store}, func() {}, nil
		},
	})
	httpServer := httptest.NewServer(NewHTTPHandler(server))
	defer httpServer.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	brief, err := media.NewBrief(media.BriefInput{
		Request: "Creator introduces a camera", MediaType: "video", Purpose: "launch",
		Platform: "youtube", AspectRatio: "16:9", DurationSeconds: 5,
		AudioMode: "off", VideoMode: "text", Style: "cinematic studio",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveBrief(brief); err != nil {
		t.Fatal(err)
	}
	call := func(name string, arguments map[string]any) *mcp.CallToolResult {
		result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	if result := call("media_generate", map[string]any{"brief_id": brief.ID}); !result.IsError || api.submissions != 0 {
		t.Fatalf("video generation bypassed preview: result=%#v submissions=%d", result.Content, api.submissions)
	}
	previewResult := call("media_preview_generate", map[string]any{"brief_id": brief.ID})
	if previewResult.IsError || api.submissions != 1 {
		t.Fatalf("preview result=%#v submissions=%d", previewResult.Content, api.submissions)
	}
	var preview media.Generation
	decodeStructured(t, previewResult.StructuredContent, &preview)
	if preview.Kind != media.GenerationKindPreview || preview.ID == "" {
		t.Fatalf("preview generation = %#v", preview)
	}
	statusResult := call("media_generation_status", map[string]any{"generation_id": preview.ID})
	if statusResult.IsError {
		t.Fatalf("preview status = %#v", statusResult.Content)
	}
	brief, err = store.GetBrief(brief.ID)
	if err != nil || brief.PreviewURL == "" {
		t.Fatalf("preview URL was not persisted: brief=%#v err=%v", brief, err)
	}
	approveResult := call("media_preview_approve", map[string]any{"brief_id": brief.ID})
	if approveResult.IsError {
		t.Fatalf("preview approval = %#v", approveResult.Content)
	}
	var approved media.Turn
	decodeStructured(t, approveResult.StructuredContent, &approved)
	if !approved.CanSubmit || approved.NextAction != "review_then_submit" {
		t.Fatalf("approved turn = %#v", approved)
	}
	finalResult := call("media_generate", map[string]any{"brief_id": brief.ID})
	if finalResult.IsError || api.submissions != 2 {
		t.Fatalf("final result=%#v submissions=%d", finalResult.Content, api.submissions)
	}
}

func TestMCPStoryboardCreatesLocalGatedShotBriefs(t *testing.T) {
	store := media.NewStore(filepath.Join(t.TempDir(), "media"))
	api := &fakeAPI{}
	server := NewServer("test", &Dependencies{
		Store: func() (*media.Store, error) { return store, nil },
		LiveService: func(context.Context, *media.Store) (*media.Service, func(), error) {
			return &media.Service{API: api, Store: store}, func() {}, nil
		},
	})
	httpServer := httptest.NewServer(NewHTTPHandler(server))
	defer httpServer.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	call := func(name string, arguments map[string]any) *mcp.CallToolResult {
		result, callErr := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: arguments})
		if callErr != nil {
			t.Fatal(callErr)
		}
		return result
	}
	start := call("media_brief_start", map[string]any{
		"request": "A two-shot launch film", "media_type": "video", "purpose": "launch",
		"platform": "youtube", "aspect_ratio": "16:9", "duration_seconds": 10,
		"audio_mode": "off", "video_mode": "text", "style": "cinematic",
		"production_mode": "storyboard",
	})
	if start.IsError {
		t.Fatalf("start = %#v", start.Content)
	}
	var turn media.Turn
	decodeStructured(t, start.StructuredContent, &turn)
	if turn.NextAction != "draft_script" {
		t.Fatalf("start turn = %#v", turn)
	}
	briefID := turn.Brief.ID
	if result := call("media_script_set", map[string]any{"brief_id": briefID, "content": "Show the product, then reveal the feature."}); result.IsError {
		t.Fatalf("script set = %#v", result.Content)
	}
	if result := call("media_script_decide", map[string]any{"brief_id": briefID, "decision": "approve"}); result.IsError {
		t.Fatalf("script approve = %#v", result.Content)
	}
	set := call("media_storyboard_set", map[string]any{
		"brief_id": briefID,
		"shots": []map[string]any{
			{"duration_seconds": 5, "visual": "Product on a dark plinth", "camera": "slow dolly"},
			{"duration_seconds": 5, "visual": "Feature in a bright hero frame", "camera": "controlled orbit"},
		},
	})
	if set.IsError {
		t.Fatalf("storyboard set = %#v", set.Content)
	}
	approved := call("media_storyboard_decide", map[string]any{"brief_id": briefID, "decision": "approve"})
	if approved.IsError {
		t.Fatalf("storyboard approve = %#v", approved.Content)
	}
	var output storyboardOutput
	decodeStructured(t, approved.StructuredContent, &output)
	if output.View.NextAction != "generate_shot_previews" || len(output.View.Shots) != 2 || output.View.Shots[0].Shot.BriefID == "" {
		t.Fatalf("storyboard output = %#v", output)
	}
	if api.submissions != 0 {
		t.Fatalf("local script/storyboard tools made %d live calls", api.submissions)
	}
	if result := call("media_generate", map[string]any{"brief_id": briefID}); !result.IsError || api.submissions != 0 {
		t.Fatalf("master generation bypassed shot workflow: result=%#v submissions=%d", result.Content, api.submissions)
	}
}

type fakeAPI struct{ submissions int }

func (*fakeAPI) Get(context.Context, string, map[string]string) (json.RawMessage, error) {
	return json.RawMessage(`{"data":{"status":"success","resultUrls":["https://cdn.example.test/preview.png"]}}`), nil
}

func (f *fakeAPI) PostWithParams(context.Context, string, map[string]string, any) (json.RawMessage, int, error) {
	f.submissions++
	return json.RawMessage(`{"data":{"taskId":"task_123"}}`), 200, nil
}

func (*fakeAPI) PostMultipartWithParams(context.Context, string, map[string]string, map[string]string, map[string]string) (json.RawMessage, int, error) {
	return json.RawMessage(`{"data":{"downloadUrl":"https://uploads.example.test/ref.png"}}`), 200, nil
}

func decodeStructured(t *testing.T, value any, target any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("decoding structured content: %v (%s)", err, data)
	}
}

func stringSliceJSON(value []string) string {
	data, _ := json.Marshal(value)
	return string(data)
}
