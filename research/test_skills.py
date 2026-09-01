import pathlib
import re
import unittest


ROOT = pathlib.Path(__file__).resolve().parents[1]
SKILLS = ROOT / "skills"
REQUIRED_NEW = {
    "kie-grill-me",
    "kie-grilling",
    "kie-image",
    "kie-video",
    "kie-audio",
    "kie-avatar",
    "kie-film",
}


def frontmatter(text: str) -> dict[str, str]:
    match = re.match(r"\A---\n(.*?)\n---\n", text, re.DOTALL)
    if not match:
        return {}
    values = {}
    for line in match.group(1).splitlines():
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        values[key.strip()] = value.strip().strip('"')
    return values


class SkillSuiteTests(unittest.TestCase):
    def test_every_skill_has_valid_metadata_and_agent_descriptor(self):
        skill_dirs = sorted(path for path in SKILLS.iterdir() if path.is_dir())
        self.assertTrue(REQUIRED_NEW.issubset({path.name for path in skill_dirs}))
        for directory in skill_dirs:
            with self.subTest(skill=directory.name):
                skill_file = directory / "SKILL.md"
                agent_file = directory / "agents" / "openai.yaml"
                self.assertTrue(skill_file.is_file())
                self.assertTrue(agent_file.is_file())
                text = skill_file.read_text()
                metadata = frontmatter(text)
                self.assertEqual(metadata.get("name"), directory.name)
                self.assertGreaterEqual(len(metadata.get("description", "")), 40)
                agent = agent_file.read_text()
                self.assertRegex(agent, r"(?m)^interface:\s*$")
                self.assertRegex(agent, r"(?m)^  display_name: \".+\"$")
                self.assertRegex(agent, r"(?m)^  short_description: \".+\"$")
                self.assertRegex(agent, r"(?m)^  default_prompt: \".+\"$")

    def test_entry_and_capability_skills_share_the_grilling_primitive(self):
        grill_me = (SKILLS / "kie-grill-me" / "SKILL.md").read_text()
        self.assertIn("$kie-grilling", grill_me)
        for name in ("kie-image", "kie-video", "kie-audio", "kie-avatar", "kie-film"):
            with self.subTest(skill=name):
                self.assertIn("$kie-grilling", (SKILLS / name / "SKILL.md").read_text())

    def test_video_skills_preserve_visual_and_paid_gates(self):
        for name in ("kie-grilling", "kie-video", "kie-film"):
            text = (SKILLS / name / "SKILL.md").read_text().lower()
            with self.subTest(skill=name):
                self.assertIn("still", text)
                self.assertIn("proof", text)
                self.assertIn("fresh", text)
                self.assertIn("paid", text)

    def test_identity_video_skill_preserves_the_optional_proof_branch(self):
        text = (SKILLS / "kie-identity" / "SKILL.md").read_text()
        for decision in ("--approve-proof", "--reject-proof", "--skip-proof"):
            self.assertIn(decision, text)

    def test_mcp_generation_skills_poll_before_display(self):
        image = (SKILLS / "kie-image" / "SKILL.md").read_text()
        image_generate = image.index("media_generate")
        self.assertGreater(image.index("media_generation_status", image_generate), image_generate)

        create = (SKILLS / "kie-create" / "SKILL.md").read_text()
        proof_generate = create.index("media_proof_generate")
        self.assertGreater(create.index("media_generation_status", proof_generate), proof_generate)
        final_generate = create.rindex("media_generate")
        self.assertGreater(create.index("media_generation_status", final_generate), final_generate)

    def test_capability_skills_use_runtime_schema_introspection(self):
        for name in ("kie-image", "kie-video", "kie-audio", "kie-avatar"):
            text = (SKILLS / name / "SKILL.md").read_text()
            with self.subTest(skill=name):
                self.assertIn("media capability", text)
                self.assertIn("models show", text)
        self.assertIn("$kie-video", (SKILLS / "kie-film" / "SKILL.md").read_text())

    def test_skills_do_not_copy_generated_model_schemas(self):
        for skill_file in sorted(SKILLS.glob("*/SKILL.md")):
            text = skill_file.read_text()
            with self.subTest(skill=skill_file.parent.name):
                self.assertNotIn("input_schema", text)
                self.assertNotRegex(text, r'(?m)^\s*"(?:properties|required)"\s*:')

    def test_kie_generate_is_only_the_advanced_manual_escape_hatch(self):
        text = (SKILLS / "kie-generate" / "SKILL.md").read_text().lower()
        self.assertIn("advanced/manual", text)
        for preferred in ("$kie-image", "$kie-video", "$kie-audio", "$kie-avatar", "$kie-film"):
            self.assertIn(preferred, text)


if __name__ == "__main__":
    unittest.main()
