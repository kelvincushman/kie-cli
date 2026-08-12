// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package leaderboard

import "testing"

func TestBestAvailableImageRoutes(t *testing.T) {
	text, err := BestAvailable("image")
	if err != nil {
		t.Fatal(err)
	}
	if text.KieModel != "gpt-image-2-text-to-image" {
		t.Fatalf("unexpected text-to-image leader: %s", text.KieModel)
	}
	edit, err := BestAvailable("character")
	if err != nil {
		t.Fatal(err)
	}
	if edit.KieModel != "gpt-image-2-image-to-image" {
		t.Fatalf("unexpected character-reference leader: %s", edit.KieModel)
	}
}

func TestSeedance25DoesNotInheritProxyScore(t *testing.T) {
	task, err := GetTask("video")
	if err != nil {
		t.Fatal(err)
	}
	if len(task.Entries) == 0 || task.Entries[0].KieModel != "bytedance/seedance-2-5" {
		t.Fatalf("unexpected video routing entries: %#v", task.Entries)
	}
	entry := task.Entries[0]
	if entry.DirectMatch || entry.SourceScore != nil || entry.Proxy == nil {
		t.Fatalf("Seedance 2.5 must remain an unscored family proxy: %#v", entry)
	}
}
