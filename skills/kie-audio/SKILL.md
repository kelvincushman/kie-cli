---
name: kie-audio
description: Plan and create Kie.ai narration, dialogue, music, or sound effects with local capability discovery, exact model input validation, voice/likeness rights checks, and fresh per-call paid confirmation. Use for audio-only assets or film audio production.
---

# Kie Audio

Start with `$kie-grilling`. Establish the audience, duration, language, delivery, voice qualities, script/lyrics, music or sound direction, output format, and synchronization target. Do not imitate a real person's voice without explicit authority; do not claim a synthetic voice is authentic.

Use compact local discovery before loading one full schema:

```bash
kie-pp-cli media capability list --capability kie-audio --agent
kie-pp-cli media capability show <model-id> --agent
kie-pp-cli models show <model-id> --agent
kie-pp-cli models example <model-id> --agent
kie-pp-cli models validate <model-id> --input '<json>' --agent
```

Use the generated CLI/MCP endpoint named by the chosen model's current documentation. Do not guess payload fields. Treat the direct model command as an advanced paid surface: immediately before each narration, dialogue, music, sound, extension, or regeneration, show the exact model/settings and cost uncertainty, then ask for a fresh explicit yes. Never reuse an earlier consent or paid confirmation.

For a film, attach approved audio to `$kie-film` shot/scene records and preserve timing, pronunciation notes, stems, licenses, and provenance. Return the provider task ID, status, URLs, exact settings, and resume/poll command.
