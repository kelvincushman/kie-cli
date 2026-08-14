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


if __name__ == "__main__":
    unittest.main()
