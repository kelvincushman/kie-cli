// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func (s *Store) SetScript(briefID string, input ScriptInput) (*Script, error) {
	brief, err := s.GetBrief(briefID)
	if err != nil {
		return nil, err
	}
	if brief.MediaType != "video" {
		return nil, fmt.Errorf("scripts apply only to video briefs")
	}
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, fmt.Errorf("script content is required")
	}
	brief.ProductionMode = ProductionModeStoryboard
	now := time.Now().UTC()
	script := &Script{
		ID: newID("script"), BriefID: brief.ID, Title: strings.TrimSpace(input.Title),
		Logline: strings.TrimSpace(input.Logline), Content: content, Status: StatusDraft,
		CreatedAt: now, UpdatedAt: now,
	}
	if strings.TrimSpace(brief.ScriptID) != "" {
		if previous, loadErr := s.GetScript(brief.ID); loadErr == nil {
			script.ID = previous.ID
			script.CreatedAt = previous.CreatedAt
		}
	}
	script.Hash = scriptFingerprint(script)
	script.BriefHash = creativeBriefFingerprint(brief)
	if err := s.SaveScript(script); err != nil {
		return nil, err
	}
	brief.ScriptID = script.ID
	brief.ScriptHash = script.Hash
	brief.ScriptBriefHash = script.BriefHash
	brief.ScriptStatus = script.Status
	brief.StoryboardID = ""
	brief.StoryboardScriptHash = ""
	brief.StoryboardBriefHash = ""
	brief.StoryboardStatus = ""
	brief.StoryboardShotsTotal = 0
	Refresh(brief)
	if err := s.SaveBrief(brief); err != nil {
		return nil, err
	}
	return script, nil
}

func (s *Store) SaveScript(script *Script) error {
	if script == nil || strings.TrimSpace(script.ID) == "" || strings.TrimSpace(script.BriefID) == "" {
		return fmt.Errorf("script id and brief id are required")
	}
	return s.writeJSON(s.productionPath("scripts", script.BriefID), script)
}

func (s *Store) GetScript(briefID string) (*Script, error) {
	var script Script
	if err := s.readJSON(s.productionPath("scripts", briefID), &script); err != nil {
		return nil, err
	}
	return &script, nil
}

func (s *Store) DecideScript(briefID, decision string) (*Script, error) {
	brief, err := s.GetBrief(briefID)
	if err != nil {
		return nil, err
	}
	script, err := s.GetScript(briefID)
	if err != nil {
		return nil, err
	}
	decision = strings.ToLower(strings.TrimSpace(decision))
	now := time.Now().UTC()
	switch decision {
	case "approve", "approved":
		if script.Hash != scriptFingerprint(script) {
			return nil, fmt.Errorf("script content has changed; save it again before approval")
		}
		if script.BriefHash != creativeBriefFingerprint(brief) {
			return nil, fmt.Errorf("the creative brief changed after this script was written; revise the script first")
		}
		script.Status = StatusApproved
		script.ApprovedAt = &now
	case "reject", "rejected":
		script.Status = StatusRejected
		script.ApprovedAt = nil
	default:
		return nil, fmt.Errorf("script decision must be approve or reject")
	}
	script.UpdatedAt = now
	if err := s.SaveScript(script); err != nil {
		return nil, err
	}
	brief.ScriptStatus = script.Status
	brief.ScriptHash = script.Hash
	brief.ScriptBriefHash = script.BriefHash
	if script.Status != StatusApproved {
		brief.StoryboardStatus = ""
		brief.StoryboardID = ""
		brief.StoryboardShotsTotal = 0
	}
	Refresh(brief)
	if err := s.SaveBrief(brief); err != nil {
		return nil, err
	}
	return script, nil
}

func (s *Store) SetStoryboard(briefID string, input StoryboardInput) (*Storyboard, error) {
	brief, err := s.GetBrief(briefID)
	if err != nil {
		return nil, err
	}
	script, err := s.GetScript(briefID)
	if err != nil {
		return nil, fmt.Errorf("load the script before storyboarding: %w", err)
	}
	if script.Status != StatusApproved || script.Hash != brief.ScriptHash || script.BriefHash != creativeBriefFingerprint(brief) {
		return nil, fmt.Errorf("the current script must be approved and match the latest creative brief before storyboarding")
	}
	if script.Hash != scriptFingerprint(script) {
		return nil, fmt.Errorf("the approved script content changed on disk; save and approve it again before storyboarding")
	}
	if err := validateStoryboardInput(brief, input); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	storyboard := &Storyboard{
		ID: newID("storyboard"), BriefID: brief.ID, ScriptID: script.ID, ScriptHash: script.Hash,
		BriefHash: creativeBriefFingerprint(brief), Title: strings.TrimSpace(input.Title),
		Status: StatusDraft, CreatedAt: now, UpdatedAt: now,
	}
	for i, source := range input.Shots {
		shotID := strings.TrimSpace(source.ID)
		if shotID == "" {
			shotID = fmt.Sprintf("shot-%02d", i+1)
		}
		child, err := newShotBrief(brief, storyboard.ID, shotID, i+1, source)
		if err != nil {
			return nil, fmt.Errorf("shot %d: %w", i+1, err)
		}
		if err := s.VaultBriefReferences(child); err != nil {
			return nil, fmt.Errorf("shot %d references: %w", i+1, err)
		}
		if err := s.SaveBrief(child); err != nil {
			return nil, err
		}
		shotReferences := []string(nil)
		if len(child.References) > len(brief.References) {
			shotReferences = append(shotReferences, child.References[len(brief.References):]...)
		}
		storyboard.Shots = append(storyboard.Shots, StoryboardShot{
			ID: shotID, Number: i + 1, Title: strings.TrimSpace(source.Title),
			DurationSeconds: source.DurationSeconds, Visual: strings.TrimSpace(source.Visual),
			Camera: strings.TrimSpace(source.Camera), Narration: strings.TrimSpace(source.Narration),
			Dialogue: strings.TrimSpace(source.Dialogue), Transition: strings.TrimSpace(source.Transition),
			References: shotReferences, BriefID: child.ID,
		})
	}
	storyboard.Hash = storyboardFingerprint(storyboard)
	if err := s.SaveStoryboard(storyboard); err != nil {
		return nil, err
	}
	brief.StoryboardID = storyboard.ID
	brief.StoryboardScriptHash = storyboard.ScriptHash
	brief.StoryboardBriefHash = storyboard.BriefHash
	brief.StoryboardStatus = storyboard.Status
	brief.StoryboardShotsTotal = len(storyboard.Shots)
	Refresh(brief)
	if err := s.SaveBrief(brief); err != nil {
		return nil, err
	}
	return storyboard, nil
}

func (s *Store) SaveStoryboard(storyboard *Storyboard) error {
	if storyboard == nil || strings.TrimSpace(storyboard.ID) == "" || strings.TrimSpace(storyboard.BriefID) == "" {
		return fmt.Errorf("storyboard id and brief id are required")
	}
	return s.writeJSON(s.productionPath("storyboards", storyboard.BriefID), storyboard)
}

func (s *Store) GetStoryboard(briefID string) (*Storyboard, error) {
	var storyboard Storyboard
	if err := s.readJSON(s.productionPath("storyboards", briefID), &storyboard); err != nil {
		return nil, err
	}
	return &storyboard, nil
}

func (s *Store) DecideStoryboard(briefID, decision string) (*Storyboard, error) {
	brief, err := s.GetBrief(briefID)
	if err != nil {
		return nil, err
	}
	storyboard, err := s.GetStoryboard(briefID)
	if err != nil {
		return nil, err
	}
	script, err := s.GetScript(briefID)
	if err != nil {
		return nil, err
	}
	decision = strings.ToLower(strings.TrimSpace(decision))
	now := time.Now().UTC()
	switch decision {
	case "approve", "approved":
		if script.Hash != scriptFingerprint(script) || script.Hash != brief.ScriptHash || script.Status != StatusApproved {
			return nil, fmt.Errorf("the approved script content changed; save and approve it again before approving the storyboard")
		}
		if storyboard.Hash != storyboardFingerprint(storyboard) {
			return nil, fmt.Errorf("the storyboard content changed on disk; save it again before approval")
		}
		if storyboard.ScriptHash != brief.ScriptHash || storyboard.BriefHash != creativeBriefFingerprint(brief) {
			return nil, fmt.Errorf("the script or creative brief changed after this storyboard was written; revise the storyboard first")
		}
		if brief.ScriptStatus != StatusApproved {
			return nil, fmt.Errorf("the current script must be approved before the storyboard")
		}
		storyboard.Status = StatusApproved
		storyboard.ApprovedAt = &now
	case "reject", "rejected":
		storyboard.Status = StatusRejected
		storyboard.ApprovedAt = nil
	default:
		return nil, fmt.Errorf("storyboard decision must be approve or reject")
	}
	storyboard.UpdatedAt = now
	if err := s.SaveStoryboard(storyboard); err != nil {
		return nil, err
	}
	brief.StoryboardStatus = storyboard.Status
	Refresh(brief)
	if err := s.SaveBrief(brief); err != nil {
		return nil, err
	}
	return storyboard, nil
}

func (s *Store) StoryboardView(briefID string) (*StoryboardView, error) {
	brief, err := s.GetBrief(briefID)
	if err != nil {
		return nil, err
	}
	storyboard, err := s.GetStoryboard(briefID)
	if err != nil {
		return nil, err
	}
	script, err := s.GetScript(briefID)
	if err != nil {
		return nil, err
	}
	view := &StoryboardView{Storyboard: storyboard}
	if script.Hash != scriptFingerprint(script) || script.Hash != brief.ScriptHash || script.Status != StatusApproved || script.BriefHash != creativeBriefFingerprint(brief) ||
		storyboard.Hash != storyboardFingerprint(storyboard) || storyboard.ScriptHash != brief.ScriptHash || storyboard.BriefHash != creativeBriefFingerprint(brief) {
		view.NextAction = "revise_storyboard"
	} else if storyboard.Status != StatusApproved {
		view.NextAction = "review_storyboard"
	}
	priorities := map[string]int{
		"complete_shot_briefs": 1, "generate_shot_previews": 2, "check_shot_previews": 3,
		"review_shot_previews": 4, "generate_shot_videos": 5, "check_shot_videos": 6,
		"revise_failed_shots": 7, "assemble_locally": 8,
	}
	derivedAction := "assemble_locally"
	for _, shot := range storyboard.Shots {
		child, loadErr := s.GetBrief(shot.BriefID)
		if loadErr != nil {
			return nil, fmt.Errorf("loading shot %s brief: %w", shot.ID, loadErr)
		}
		shotView := StoryboardShotView{Shot: shot, Turn: TurnFor(child)}
		action := shotAction(shotView.Turn)
		if child.GenerationID != "" {
			if generation, generationErr := s.GetGeneration(child.GenerationID); generationErr == nil {
				shotView.Generation = generation
				if generationFailed(generation.Status) {
					action = "revise_failed_shots"
				} else if mediaGenerationComplete(generation.Status) {
					action = "assemble_locally"
				}
			}
		}
		view.Shots = append(view.Shots, shotView)
		if priorities[action] < priorities[derivedAction] {
			derivedAction = action
		}
	}
	if view.NextAction == "" {
		view.NextAction = derivedAction
	}
	return view, nil
}

func newShotBrief(master *Brief, storyboardID, shotID string, number int, shot StoryboardShotInput) (*Brief, error) {
	prompt := []string{master.Request, fmt.Sprintf("Storyboard shot %d: %s", number, strings.TrimSpace(shot.Visual))}
	if shot.Camera != "" {
		prompt = append(prompt, "Camera: "+strings.TrimSpace(shot.Camera))
	}
	if shot.Narration != "" {
		prompt = append(prompt, "Narration: "+strings.TrimSpace(shot.Narration))
	}
	if shot.Dialogue != "" {
		prompt = append(prompt, "Dialogue: "+strings.TrimSpace(shot.Dialogue))
	}
	if shot.Transition != "" {
		prompt = append(prompt, "Transition: "+strings.TrimSpace(shot.Transition))
	}
	videoMode := master.VideoMode
	firstFrame, lastFrame := master.FirstFrame, master.LastFrame
	if len(shot.References) > 0 && (videoMode == "text" || videoMode == "first-frame" || videoMode == "first-last") {
		videoMode, firstFrame, lastFrame = "multimodal", "", ""
	}
	child, err := NewBrief(BriefInput{
		Workflow: master.Workflow, Request: strings.Join(prompt, ". "), MediaType: "video",
		Purpose: master.Purpose, Platform: master.Platform, AspectRatio: master.AspectRatio,
		DurationSeconds: shot.DurationSeconds, Resolution: master.Resolution, AudioMode: master.AudioMode,
		VideoMode: videoMode, OutputFormat: master.OutputFormat, ReturnLastFrame: master.ReturnLastFrame,
		WebSearch: master.WebSearch, Style: master.Style,
		References:      append(append([]string(nil), master.References...), shot.References...),
		ReferenceVideos: master.ReferenceVideos, ReferenceAudio: master.ReferenceAudio,
		FirstFrame: firstFrame, LastFrame: lastFrame, IdentityIDs: master.IdentityIDs,
		Model: master.Model, ProductionMode: ProductionModeSingleShot,
	})
	if err != nil {
		return nil, err
	}
	child.MasterBriefID = master.ID
	child.StoryboardID = storyboardID
	child.ShotID = shotID
	child.ReferencesComplete = true
	Refresh(child)
	return child, nil
}

func validateStoryboardInput(brief *Brief, input StoryboardInput) error {
	if len(input.Shots) == 0 || len(input.Shots) > 60 {
		return fmt.Errorf("storyboard requires 1 to 60 shots")
	}
	total := 0
	seen := map[string]bool{}
	for i, shot := range input.Shots {
		if strings.TrimSpace(shot.Visual) == "" {
			return fmt.Errorf("shot %d visual is required", i+1)
		}
		if shot.DurationSeconds < 4 || shot.DurationSeconds > 30 {
			return fmt.Errorf("shot %d duration must be between 4 and 30 seconds", i+1)
		}
		if len(shot.References) > 29 {
			return fmt.Errorf("shot %d accepts at most 29 explicit image references because its approved preview occupies one SeedDance image slot", i+1)
		}
		id := strings.TrimSpace(shot.ID)
		if id != "" {
			if seen[id] {
				return fmt.Errorf("storyboard shot id %q is duplicated", id)
			}
			seen[id] = true
		}
		total += shot.DurationSeconds
	}
	if total > 600 {
		return fmt.Errorf("storyboard duration must not exceed 600 seconds")
	}
	if brief.DurationSeconds > 0 && total != brief.DurationSeconds {
		return fmt.Errorf("storyboard shots total %d seconds; the master brief requires %d", total, brief.DurationSeconds)
	}
	return nil
}

func scriptFingerprint(script *Script) string {
	state := struct {
		Title   string `json:"title"`
		Logline string `json:"logline"`
		Content string `json:"content"`
	}{script.Title, script.Logline, script.Content}
	data, _ := json.Marshal(state)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func storyboardFingerprint(storyboard *Storyboard) string {
	if storyboard == nil {
		return ""
	}
	shots := append([]StoryboardShot(nil), storyboard.Shots...)
	for i := range shots {
		shots[i].References = cleanStrings(shots[i].References)
	}
	state := struct {
		ID         string           `json:"id"`
		BriefID    string           `json:"brief_id"`
		ScriptID   string           `json:"script_id"`
		ScriptHash string           `json:"script_hash"`
		BriefHash  string           `json:"brief_hash"`
		Title      string           `json:"title"`
		Shots      []StoryboardShot `json:"shots"`
	}{storyboard.ID, storyboard.BriefID, storyboard.ScriptID, storyboard.ScriptHash, storyboard.BriefHash, storyboard.Title, shots}
	data, _ := json.Marshal(state)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (s *Store) productionPath(kind, briefID string) string {
	return filepath.Join(s.root, kind, safeID(briefID)+".json")
}

func shotAction(turn Turn) string {
	switch turn.NextAction {
	case "answer_question":
		return "complete_shot_briefs"
	case "generate_preview", "regenerate_preview":
		return "generate_shot_previews"
	case "check_preview_status":
		return "check_shot_previews"
	case "review_preview":
		return "review_shot_previews"
	case "review_then_submit":
		return "generate_shot_videos"
	case "check_generation_status":
		return "check_shot_videos"
	default:
		return "complete_shot_briefs"
	}
}

func mediaGenerationComplete(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "success", "succeeded", "completed", "complete":
		return true
	default:
		return false
	}
}
