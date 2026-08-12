// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package academy

import "testing"

func TestCatalogCountsAndKnownLesson(t *testing.T) {
	catalog, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if catalog.CourseCount != 16 || catalog.LessonCount != 171 {
		t.Fatalf("unexpected Academy snapshot: %d courses, %d lessons", catalog.CourseCount, catalog.LessonCount)
	}
	lesson, err := GetLesson("blockbuster-4k/one-character-many-worlds-watch-the-film")
	if err != nil {
		t.Fatal(err)
	}
	if lesson.SourceURL == "" || lesson.KieMethod == "" || lesson.ProductionStage == "" {
		t.Fatalf("lesson is missing source-linked guidance: %#v", lesson)
	}
}

func TestRecommendFindsCharacterConsistency(t *testing.T) {
	results, err := Recommend("consistent character across many worlds", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one lesson recommendation")
	}
	found := false
	for _, result := range results {
		if result.Lesson.Key == "blockbuster-4k/one-character-many-worlds-watch-the-film" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected Blockbuster character lesson in top recommendations: %#v", results)
	}
}
