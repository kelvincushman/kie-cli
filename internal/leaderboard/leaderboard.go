// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

// Package leaderboard exposes the checked-in task-specific media evidence
// ledger used by the director. It does not normalize unrelated benchmarks into
// a fictional universal score.
package leaderboard

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed leaderboard.json
var leaderboardData []byte

type Ledger struct {
	SchemaVersion int      `json:"schema_version"`
	GeneratedAt   string   `json:"generated_at"`
	Methodology   string   `json:"methodology"`
	Sources       []Source `json:"sources"`
	Tasks         []Task   `json:"tasks"`
}

type Source struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	EvidenceType string `json:"evidence_type"`
	FetchedAt    string `json:"fetched_at"`
}

type Task struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	SourceIDs   []string `json:"source_ids"`
	RankingKind string   `json:"ranking_kind"`
	Warning     string   `json:"warning,omitempty"`
	Entries     []Entry  `json:"entries"`
}

type Entry struct {
	SourceRank   *int     `json:"source_rank"`
	SourceScore  *int     `json:"source_score"`
	ScoreType    string   `json:"score_type"`
	SourceModel  string   `json:"source_model"`
	KieModel     string   `json:"kie_model"`
	KieAvailable bool     `json:"kie_available"`
	DirectMatch  bool     `json:"direct_match"`
	Proxy        *Proxy   `json:"proxy,omitempty"`
	Strengths    []string `json:"strengths,omitempty"`
}

type Proxy struct {
	SourceModel string `json:"source_model"`
	SourceRank  int    `json:"source_rank"`
	SourceScore int    `json:"source_score"`
	Warning     string `json:"warning"`
}

var (
	loadOnce sync.Once
	loaded   Ledger
	loadErr  error
)

func Load() (*Ledger, error) {
	loadOnce.Do(func() {
		loadErr = json.Unmarshal(leaderboardData, &loaded)
		if loadErr != nil {
			loadErr = fmt.Errorf("decoding embedded model leaderboard: %w", loadErr)
		}
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return &loaded, nil
}

func ListTasks() ([]Task, error) {
	ledger, err := Load()
	if err != nil {
		return nil, err
	}
	result := make([]Task, len(ledger.Tasks))
	for i, task := range ledger.Tasks {
		task.Entries = nil
		result[i] = task
	}
	return result, nil
}

func GetTask(id string) (*Task, error) {
	ledger, err := Load()
	if err != nil {
		return nil, err
	}
	id = normalizeTask(id)
	for _, task := range ledger.Tasks {
		if task.ID == id {
			copy := task
			copy.Entries = append([]Entry(nil), task.Entries...)
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("unknown leaderboard task %q", id)
}

func BestAvailable(id string) (*Entry, error) {
	task, err := GetTask(id)
	if err != nil {
		return nil, err
	}
	for _, entry := range task.Entries {
		if entry.KieAvailable && entry.DirectMatch {
			copy := entry
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("leaderboard task %q has no directly matched Kie route", task.ID)
}

func normalizeTask(id string) string {
	id = strings.ToLower(strings.TrimSpace(id))
	switch id {
	case "image", "text-image", "t2i":
		return "text-to-image"
	case "edit", "image-editing", "i2i":
		return "image-edit"
	case "character", "identity", "consistency":
		return "character-consistency"
	case "video", "text-video", "t2v":
		return "text-to-video"
	default:
		return id
	}
}
