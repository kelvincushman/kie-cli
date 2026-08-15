---
name: kie-image
description: Create or edit Kie.ai still images through the local director with compact model discovery, references, consent, route rationale, model-specific validation, and fresh paid confirmation. Use for text-to-image, image-to-image, product, portrait, logo concept, background, poster, thumbnail, edit, remove, or upscale requests.
---

# Kie Image

Start with `$kie-grilling`; do not ask a second parallel questionnaire. Use its durable brief and selected route.

Inspect compact candidates first, then only the chosen schema:

```bash
kie-pp-cli media capability list --capability kie-image --agent
kie-pp-cli media capability show <model-id> --agent
kie-pp-cli models show <model-id> --agent
```

Prefer `kie-pp-cli create` for text-to-image and reference-driven image generation. Add local/URL references through `media reference add`, and use `$kie-identity` for recurring people. Keep exact logo geometry, long copy, dimensions, and typography deterministic outside a raster generator.

Show the user the route rationale, model settings, references, cost status, and overrides. Validate any manual input with `kie-pp-cli models validate <model-id> --input '<json>' --agent`.

Immediately before a live generation, ask for confirmation for that exact brief/model/plan. Only then:

```bash
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

For an MCP host, use `media_paid_confirm` followed by `media_generate`. Do not reuse a confirmation. Poll `media_generation_status` until it completes or fails; handle failure before asking for another paid action, and render the image only after completion. Collect requested revisions and create a fresh confirmation for every regeneration, edit, removal, or upscale. Report brief/generation/task IDs, model, settings, provenance, result URLs, and the exact resume command.
