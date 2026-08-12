// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

// Package academy exposes the compact, source-linked production-method map
// embedded in the CLI. The catalog contains public titles and original
// Kie-native adaptations, never copied lesson scripts or prompts.
package academy

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

//go:embed catalog.json
var catalogData []byte

type Catalog struct {
	SchemaVersion   int      `json:"schema_version"`
	GeneratedAt     string   `json:"generated_at"`
	Source          string   `json:"source"`
	CopyrightPolicy string   `json:"copyright_policy"`
	CourseCount     int      `json:"course_count"`
	LessonCount     int      `json:"lesson_count"`
	Courses         []Course `json:"courses"`
}

type Course struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	SourceURL       string   `json:"source_url"`
	LessonCount     int      `json:"lesson_count"`
	DurationMinutes int      `json:"duration_minutes"`
	Difficulty      string   `json:"difficulty"`
	Lessons         []Lesson `json:"lessons,omitempty"`
}

type Lesson struct {
	Slug                 string   `json:"slug"`
	Title                string   `json:"title"`
	Position             int      `json:"position"`
	SourceURL            string   `json:"source_url"`
	Key                  string   `json:"key"`
	ProductionStage      string   `json:"production_stage"`
	KieMethod            string   `json:"kie_method"`
	PromptFocus          string   `json:"prompt_focus"`
	RecommendedKieModels []string `json:"recommended_kie_models,omitempty"`
	CapabilityBoundary   string   `json:"capability_boundary,omitempty"`
	CourseSlug           string   `json:"course_slug,omitempty"`
	CourseTitle          string   `json:"course_title,omitempty"`
}

type Recommendation struct {
	Lesson Lesson `json:"lesson"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

var (
	loadOnce sync.Once
	loaded   Catalog
	loadErr  error
)

func Load() (*Catalog, error) {
	loadOnce.Do(func() {
		loadErr = json.Unmarshal(catalogData, &loaded)
		if loadErr != nil {
			loadErr = fmt.Errorf("decoding embedded academy catalog: %w", loadErr)
			return
		}
		if loaded.CourseCount != len(loaded.Courses) {
			loadErr = fmt.Errorf("academy catalog course count is %d, found %d", loaded.CourseCount, len(loaded.Courses))
			return
		}
		lessons := 0
		for i := range loaded.Courses {
			course := &loaded.Courses[i]
			lessons += len(course.Lessons)
			for j := range course.Lessons {
				course.Lessons[j].CourseSlug = course.Slug
				course.Lessons[j].CourseTitle = course.Title
			}
		}
		if loaded.LessonCount != lessons {
			loadErr = fmt.Errorf("academy catalog lesson count is %d, found %d", loaded.LessonCount, lessons)
		}
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return &loaded, nil
}

func ListCourses() ([]Course, error) {
	catalog, err := Load()
	if err != nil {
		return nil, err
	}
	result := make([]Course, len(catalog.Courses))
	for i, course := range catalog.Courses {
		course.Lessons = nil
		result[i] = course
	}
	return result, nil
}

func GetCourse(slug string) (*Course, error) {
	catalog, err := Load()
	if err != nil {
		return nil, err
	}
	slug = strings.TrimSpace(strings.ToLower(slug))
	for _, course := range catalog.Courses {
		if course.Slug == slug {
			copy := course
			copy.Lessons = append([]Lesson(nil), course.Lessons...)
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("unknown academy course %q", slug)
}

func GetLesson(key string) (*Lesson, error) {
	catalog, err := Load()
	if err != nil {
		return nil, err
	}
	key = strings.Trim(strings.TrimSpace(strings.ToLower(key)), "/")
	for _, course := range catalog.Courses {
		for _, lesson := range course.Lessons {
			if lesson.Key == key || lesson.Slug == key {
				copy := lesson
				return &copy, nil
			}
		}
	}
	return nil, fmt.Errorf("unknown academy lesson %q; use course-slug/lesson-slug", key)
}

func Search(query, courseSlug, stage string, limit int) ([]Recommendation, error) {
	catalog, err := Load()
	if err != nil {
		return nil, err
	}
	query = strings.TrimSpace(strings.ToLower(query))
	courseSlug = strings.TrimSpace(strings.ToLower(courseSlug))
	stage = strings.TrimSpace(strings.ToLower(stage))
	if limit <= 0 {
		limit = 10
	}
	tokens := searchTokens(query)
	results := make([]Recommendation, 0)
	for _, course := range catalog.Courses {
		if courseSlug != "" && course.Slug != courseSlug {
			continue
		}
		for _, lesson := range course.Lessons {
			if stage != "" && lesson.ProductionStage != stage {
				continue
			}
			score := lessonScore(course, lesson, query, tokens)
			if query != "" && score == 0 {
				continue
			}
			reason := "matches the requested production stage"
			if query != "" {
				reason = fmt.Sprintf("matched %q across the public lesson title and original Kie method", query)
			}
			results = append(results, Recommendation{Lesson: lesson, Score: score, Reason: reason})
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Lesson.Key < results[j].Lesson.Key
		}
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func Recommend(query string, limit int) ([]Recommendation, error) {
	return Search(query, "", "", limit)
}

func lessonScore(course Course, lesson Lesson, query string, tokens []string) int {
	if query == "" {
		return 1
	}
	title := strings.ToLower(lesson.Title)
	courseTitle := strings.ToLower(course.Title)
	method := strings.ToLower(lesson.KieMethod + " " + lesson.PromptFocus + " " + lesson.ProductionStage)
	score := 0
	if strings.Contains(title, query) {
		score += 20
	}
	if strings.Contains(courseTitle, query) {
		score += 12
	}
	for _, token := range tokens {
		if strings.Contains(title, token) {
			score += 6
		}
		if strings.Contains(courseTitle, token) {
			score += 3
		}
		if strings.Contains(method, token) {
			score += 2
		}
	}
	return score
}

func searchTokens(query string) []string {
	seen := map[string]bool{}
	result := make([]string, 0)
	for _, token := range strings.FieldsFunc(query, func(r rune) bool {
		return r < '0' || (r > '9' && r < 'a') || r > 'z'
	}) {
		if len(token) < 3 || seen[token] {
			continue
		}
		seen[token] = true
		result = append(result, token)
	}
	return result
}
