// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"encoding/json"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusReady     = "ready"
	StatusSubmitted = "submitted"
	StatusApproved  = "approved"
	StatusRejected  = "rejected"

	ProductionModeSingleShot = "single-shot"
	ProductionModeStoryboard = "storyboard"

	GenerationKindFinal   = "final"
	GenerationKindPreview = "preview"
)

type Brief struct {
	ID                   string     `json:"id"`
	Workflow             string     `json:"workflow,omitempty"`
	Lesson               string     `json:"lesson,omitempty"`
	Request              string     `json:"request,omitempty"`
	MediaType            string     `json:"media_type,omitempty"`
	Purpose              string     `json:"purpose,omitempty"`
	Platform             string     `json:"platform,omitempty"`
	AspectRatio          string     `json:"aspect_ratio,omitempty"`
	DurationSeconds      int        `json:"duration_seconds,omitempty"`
	Resolution           string     `json:"resolution,omitempty"`
	AudioMode            string     `json:"audio_mode,omitempty"`
	VideoMode            string     `json:"video_mode,omitempty"`
	OutputFormat         string     `json:"output_format,omitempty"`
	ReturnLastFrame      bool       `json:"return_last_frame,omitempty"`
	WebSearch            bool       `json:"web_search,omitempty"`
	Style                string     `json:"style,omitempty"`
	References           []string   `json:"references,omitempty"`
	ReferenceVideos      []string   `json:"reference_videos,omitempty"`
	ReferenceAudio       []string   `json:"reference_audio,omitempty"`
	FirstFrame           string     `json:"first_frame,omitempty"`
	LastFrame            string     `json:"last_frame,omitempty"`
	ReferencesComplete   bool       `json:"references_complete"`
	IdentityIDs          []string   `json:"identity_ids,omitempty"`
	IdentityComplete     bool       `json:"identity_complete"`
	Model                string     `json:"model,omitempty"`
	PreviewModel         string     `json:"preview_model,omitempty"`
	Status               string     `json:"status"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Plan                 *Plan      `json:"plan,omitempty"`
	GenerationID         string     `json:"generation_id,omitempty"`
	PreviewGenerationID  string     `json:"preview_generation_id,omitempty"`
	PreviewStatus        string     `json:"preview_status,omitempty"`
	PreviewURL           string     `json:"preview_url,omitempty"`
	PreviewBriefHash     string     `json:"preview_brief_hash,omitempty"`
	PreviewApprovedAt    *time.Time `json:"preview_approved_at,omitempty"`
	PreviewRevision      int        `json:"preview_revision,omitempty"`
	ProductionMode       string     `json:"production_mode,omitempty"`
	MasterBriefID        string     `json:"master_brief_id,omitempty"`
	StoryboardID         string     `json:"storyboard_id,omitempty"`
	ShotID               string     `json:"shot_id,omitempty"`
	ScriptID             string     `json:"script_id,omitempty"`
	ScriptHash           string     `json:"script_hash,omitempty"`
	ScriptBriefHash      string     `json:"script_brief_hash,omitempty"`
	ScriptStatus         string     `json:"script_status,omitempty"`
	StoryboardScriptHash string     `json:"storyboard_script_hash,omitempty"`
	StoryboardBriefHash  string     `json:"storyboard_brief_hash,omitempty"`
	StoryboardStatus     string     `json:"storyboard_status,omitempty"`
	StoryboardShotsTotal int        `json:"storyboard_shots_total,omitempty"`
}

type Question struct {
	Key            string   `json:"key"`
	Prompt         string   `json:"prompt"`
	Recommendation string   `json:"recommendation,omitempty"`
	Options        []string `json:"options,omitempty"`
	AllowSkip      bool     `json:"allow_skip,omitempty"`
}

type Turn struct {
	Brief        *Brief    `json:"brief"`
	NextQuestion *Question `json:"next_question,omitempty"`
	Ready        bool      `json:"ready"`
	CanSubmit    bool      `json:"can_submit"`
	NextAction   string    `json:"next_action"`
}

type Plan struct {
	Model           string         `json:"model"`
	Input           map[string]any `json:"input"`
	Rationale       string         `json:"rationale"`
	Lesson          string         `json:"lesson,omitempty"`
	ProductionStage string         `json:"production_stage,omitempty"`
	Method          string         `json:"method,omitempty"`
	PromptFocus     string         `json:"prompt_focus,omitempty"`
}

type Reference struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Kind       string    `json:"kind"`
	MediaType  string    `json:"media_type"`
	Source     string    `json:"source"`
	StoredPath string    `json:"stored_path,omitempty"`
	URL        string    `json:"url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Generation struct {
	ID          string          `json:"id"`
	BriefID     string          `json:"brief_id"`
	Kind        string          `json:"kind"`
	Fingerprint string          `json:"fingerprint,omitempty"`
	TaskID      string          `json:"task_id"`
	Model       string          `json:"model"`
	Status      string          `json:"status"`
	ResultURLs  []string        `json:"result_urls,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Remote      json.RawMessage `json:"-"`
}

// Script is a durable, local production artifact. Its approval is bound to
// both the content hash and the creative fields of the master brief.
type Script struct {
	ID         string     `json:"id"`
	BriefID    string     `json:"brief_id"`
	Title      string     `json:"title,omitempty"`
	Logline    string     `json:"logline,omitempty"`
	Content    string     `json:"content"`
	Hash       string     `json:"hash"`
	BriefHash  string     `json:"brief_hash"`
	Status     string     `json:"status"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ScriptInput struct {
	Title   string `json:"title,omitempty"`
	Logline string `json:"logline,omitempty"`
	Content string `json:"content"`
}

// StoryboardShot describes one independently gated SeedDance clip. BriefID
// points to a normal single-shot brief, so every shot must pass the existing
// generated-still review and approval gate before video generation.
type StoryboardShot struct {
	ID              string   `json:"id"`
	Number          int      `json:"number"`
	Title           string   `json:"title,omitempty"`
	DurationSeconds int      `json:"duration_seconds"`
	Visual          string   `json:"visual"`
	Camera          string   `json:"camera,omitempty"`
	Narration       string   `json:"narration,omitempty"`
	Dialogue        string   `json:"dialogue,omitempty"`
	Transition      string   `json:"transition,omitempty"`
	References      []string `json:"references,omitempty"`
	BriefID         string   `json:"shot_brief_id"`
}

type StoryboardShotInput struct {
	ID              string   `json:"id,omitempty"`
	Title           string   `json:"title,omitempty"`
	DurationSeconds int      `json:"duration_seconds"`
	Visual          string   `json:"visual"`
	Camera          string   `json:"camera,omitempty"`
	Narration       string   `json:"narration,omitempty"`
	Dialogue        string   `json:"dialogue,omitempty"`
	Transition      string   `json:"transition,omitempty"`
	References      []string `json:"references,omitempty"`
}

type StoryboardInput struct {
	Title string                `json:"title,omitempty"`
	Shots []StoryboardShotInput `json:"shots"`
}

type Storyboard struct {
	ID         string           `json:"id"`
	BriefID    string           `json:"brief_id"`
	ScriptID   string           `json:"script_id"`
	ScriptHash string           `json:"script_hash"`
	BriefHash  string           `json:"brief_hash"`
	Hash       string           `json:"hash"`
	Title      string           `json:"title,omitempty"`
	Shots      []StoryboardShot `json:"shots"`
	Status     string           `json:"status"`
	ApprovedAt *time.Time       `json:"approved_at,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type StoryboardShotView struct {
	Shot       StoryboardShot `json:"shot"`
	Turn       Turn           `json:"turn"`
	Generation *Generation    `json:"generation,omitempty"`
}

type StoryboardView struct {
	Storyboard *Storyboard          `json:"storyboard"`
	Shots      []StoryboardShotView `json:"shots"`
	NextAction string               `json:"next_action"`
}

// PublicReference is the safe response form of a reusable reference. Local
// source and vault paths are deliberately excluded from CLI and MCP output.
type PublicReference struct {
	ID        string    `json:"id"`
	Handle    string    `json:"handle"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	MediaType string    `json:"media_type"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Reference) Public() PublicReference {
	return PublicReference{
		ID: r.ID, Handle: "ref:" + r.ID, Name: r.Name, Kind: r.Kind, MediaType: r.MediaType,
		URL: r.URL, CreatedAt: r.CreatedAt,
	}
}

// IdentityProfile is a reusable, local-only likeness reference bundle. It is
// not a trained biometric model: the stored handles are expanded into image
// references only when the user approves a live generation.
type IdentityProfile struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	ImageReferences    []string  `json:"image_references"`
	ConsentConfirmedAt time.Time `json:"consent_confirmed_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
