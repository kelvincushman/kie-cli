from __future__ import annotations

import json
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "research"))

import build_academy_map
import refresh_model_leaderboard


class AcademyMapTests(unittest.TestCase):
    def test_checked_in_catalog_is_source_linked_and_copy_safe(self) -> None:
        catalog = json.loads((ROOT / "internal/academy/catalog.json").read_text())
        kie_models = {
            model["id"]
            for model in json.loads((ROOT / "internal/kiecatalog/catalog.json").read_text())["models"]
        }
        self.assertEqual(16, catalog["course_count"])
        self.assertEqual(171, catalog["lesson_count"])
        self.assertNotIn("compiledMdx", json.dumps(catalog))
        for course in catalog["courses"]:
            self.assertEqual(course["lesson_count"], len(course["lessons"]))
            for lesson in course["lessons"]:
                self.assertTrue(lesson["source_url"].startswith(build_academy_map.COURSES_URL + "/"))
                self.assertTrue(lesson["kie_method"])
                self.assertTrue(lesson["prompt_focus"])
                self.assertTrue(set(lesson["recommended_kie_models"]).issubset(kie_models))

    def test_classifier_uses_title_before_course_description(self) -> None:
        self.assertEqual(
            "reference-review",
            build_academy_map.classify("Watch the film", "build reusable character assets"),
        )
        self.assertEqual(
            "asset-lock",
            build_academy_map.classify("Build reusable character assets", "watch an example"),
        )


class LeaderboardTests(unittest.TestCase):
    def test_checked_in_ledger_does_not_transfer_seedance_proxy_score(self) -> None:
        ledger = json.loads((ROOT / "internal/leaderboard/leaderboard.json").read_text())
        kie_models = {
            model["id"]
            for model in json.loads((ROOT / "internal/kiecatalog/catalog.json").read_text())["models"]
        }
        for task in ledger["tasks"]:
            for entry in task["entries"]:
                self.assertEqual(entry["kie_model"] in kie_models, entry["kie_available"])
        video = next(task for task in ledger["tasks"] if task["id"] == "text-to-video")
        seedance = video["entries"][0]
        self.assertEqual("bytedance/seedance-2-5", seedance["kie_model"])
        self.assertIsNone(seedance["source_rank"])
        self.assertIsNone(seedance["source_score"])
        self.assertFalse(seedance["direct_match"])
        self.assertIn("do not transfer", seedance["proxy"]["warning"].lower())

    def test_table_parser_extracts_public_rankings(self) -> None:
        rows = refresh_model_leaderboard.table_rows(
            "<table><tr><td>1</td><td>1</td><td>Lab</td><td>Model</td><td>1,234</td></tr></table>"
        )
        self.assertEqual(["1", "1", "Lab", "Model", "1,234"], rows[0])
        self.assertEqual(1234, refresh_model_leaderboard.parse_number(rows[0][4]))


if __name__ == "__main__":
    unittest.main()
