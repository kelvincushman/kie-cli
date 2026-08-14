import json
import pathlib
import re
import unittest


ROOT = pathlib.Path(__file__).resolve().parents[1]


class DocumentationTests(unittest.TestCase):
    def test_local_markdown_links_resolve(self):
        markdown_files = [ROOT / "README.md", *sorted((ROOT / "docs").glob("*.md"))]
        link_pattern = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
        for source in markdown_files:
            for target in link_pattern.findall(source.read_text()):
                if target.startswith(("http://", "https://", "#", "mailto:")):
                    continue
                clean = target.split("#", 1)[0]
                if not clean:
                    continue
                with self.subTest(source=source.name, target=target):
                    self.assertTrue((source.parent / clean).resolve().exists())

    def test_readme_names_current_public_contract(self):
        readme = (ROOT / "README.md").read_text()
        for expected in (
            "17 agent skills",
            "40 media tools",
            "129 typed models",
            "70 API operations",
            "--proof --confirm-paid",
            "docs/SKILLS.md",
            "docs/PROOF_AND_PAID_CONFIRMATION.md",
        ):
            with self.subTest(expected=expected):
                self.assertIn(expected, readme)

    def test_printing_press_patch_preserves_new_surfaces(self):
        patch_path = ROOT / ".printing-press-patches" / "media-director-local-extension.json"
        patch = json.loads(patch_path.read_text())
        files = set(patch["files"])
        required = {
            "docs/SKILLS.md",
            "docs/PROOF_AND_PAID_CONFIRMATION.md",
            "internal/kiecatalog/capabilities.json",
            "internal/media/paid.go",
            "internal/media/proof.go",
            "research/generate_capabilities.py",
            "research/test_skills.py",
            "skills/kie-grill-me/SKILL.md",
            "skills/kie-grilling/SKILL.md",
            "skills/kie-image/SKILL.md",
            "skills/kie-video/SKILL.md",
            "skills/kie-audio/SKILL.md",
            "skills/kie-avatar/SKILL.md",
            "skills/kie-film/SKILL.md",
        }
        self.assertTrue(required.issubset(files))
        for relative in files:
            with self.subTest(relative=relative):
                self.assertTrue((ROOT / relative).exists())

    def test_guides_keep_still_and_optional_video_proof_paths_distinct(self):
        root_skill = (ROOT / "SKILL.md").read_text()
        director = (ROOT / "docs" / "MEDIA_DIRECTOR.md").read_text()
        storyboard = (ROOT / "docs" / "SCRIPT_AND_STORYBOARD.md").read_text()
        create_skill = (ROOT / "skills" / "kie-create" / "SKILL.md").read_text()
        self.assertIn("Still-image submission", root_skill)
        self.assertIn("Video-only preview, proof, and final gates", root_skill)
        self.assertIn("Decline the optional proof without making a live proof call", director)
        self.assertIn("Accept the optional proof only after a fresh paid confirmation", director)
        self.assertIn("--skip-proof", storyboard)
        self.assertIn("--reject-proof", storyboard)
        self.assertIn("new proof-scoped `media_paid_confirm`", create_skill)


if __name__ == "__main__":
    unittest.main()
