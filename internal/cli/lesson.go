// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"kie-pp-cli/internal/academy"
	"kie-pp-cli/internal/leaderboard"
)

func init() {
	registerNovelCommand(func(root *cobra.Command, flags *rootFlags) {
		addNovelCommandIfAbsent(root, newLessonCmd(flags))
	})
}

// pp:data-source local
func newLessonCmd(flags *rootFlags) *cobra.Command {
	parent := &cobra.Command{
		Use:   "lesson",
		Short: "Choose a source-linked production lesson and start a guided Kie workflow",
		Long: "Browse the checked-in public Academy map, receive an original Kie-native method recommendation, and start a durable brief. " +
			"The catalog stores public titles and links, not copied course scripts or prompts. Agent hosts may expose the matching kie-lesson skill as /kie-lesson.",
	}

	var course, stage string
	var limit int
	list := &cobra.Command{
		Use:   "list [query]",
		Short: "List courses or search the source-linked lesson methods",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) == 1 {
				query = args[0]
			}
			if query == "" && course == "" && stage == "" {
				courses, err := academy.ListCourses()
				if err != nil {
					return err
				}
				return printMediaValue(cmd, flags, courses)
			}
			results, err := academy.Search(query, course, stage, limit)
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, results)
		},
	}
	list.Flags().StringVar(&course, "course", "", "Limit results to one course slug")
	list.Flags().StringVar(&stage, "stage", "", "Limit results to a production stage")
	list.Flags().IntVar(&limit, "limit", 20, "Maximum lesson results")
	parent.AddCommand(list)

	parent.AddCommand(&cobra.Command{
		Use:   "show <course-or-lesson-key>",
		Short: "Show a course or one original source-linked lesson method",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.Trim(args[0], "/")
			if strings.Contains(key, "/") {
				lesson, err := academy.GetLesson(key)
				if err != nil {
					return err
				}
				return printMediaValue(cmd, flags, lesson)
			}
			course, err := academy.GetCourse(key)
			if err != nil {
				if lesson, lessonErr := academy.GetLesson(key); lessonErr == nil {
					return printMediaValue(cmd, flags, lesson)
				}
				return err
			}
			return printMediaValue(cmd, flags, course)
		},
	})

	var recommendLimit int
	recommend := &cobra.Command{
		Use:   "recommend <what-you-want-to-create>",
		Short: "Rank original lesson-method adaptations for a media request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := academy.Recommend(args[0], recommendLimit)
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, map[string]any{
				"query": args[0], "recommendations": results,
				"next_step": "Run kie-pp-cli lesson start <request> --lesson <course/lesson> or invoke $kie-lesson.",
			})
		},
	}
	recommend.Flags().IntVar(&recommendLimit, "limit", 5, "Maximum recommendations")
	parent.AddCommand(recommend)

	var selectedLesson string
	start := &cobra.Command{
		Use:   "start <what-you-want-to-create>",
		Short: "Start a storyboard-first Academy workflow with an explicit or recommended lesson",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lessonKey := strings.Trim(selectedLesson, "/")
			if lessonKey == "" {
				recommendations, err := academy.Recommend(args[0], 1)
				if err != nil {
					return err
				}
				if len(recommendations) == 0 {
					return fmt.Errorf("no Academy lesson matched %q", args[0])
				}
				lessonKey = recommendations[0].Lesson.Key
			}
			if _, err := academy.GetLesson(lessonKey); err != nil {
				return err
			}
			return runMediaCreate(cmd, flags, mediaCreateOptions{
				workflow: "academy", lesson: lessonKey, mediaType: "video",
				productionMode: "storyboard",
			}, args[0])
		},
	}
	start.Flags().StringVar(&selectedLesson, "lesson", "", "Use an exact course-slug/lesson-slug instead of automatic recommendation")
	parent.AddCommand(start)

	return parent
}

// pp:data-source local
func newMediaLeaderboardCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "leaderboard [task]",
		Short: "Inspect current task-specific model evidence and Kie route availability",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				ledger, err := leaderboard.Load()
				if err != nil {
					return err
				}
				return printMediaValue(cmd, flags, ledger)
			}
			task, err := leaderboard.GetTask(args[0])
			if err != nil {
				return err
			}
			return printMediaValue(cmd, flags, task)
		},
	}
}
