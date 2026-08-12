---
name: kie-brandkit
description: Build or extend a reusable brand identity with staged approvals, local brand-lock state, deterministic logo and typography artifacts, and Kie.ai concept imagery or mockups. Use for palettes, logos, type systems, brand books, social templates, packaging, signage, merchandise, or campaign asset systems.
---

# Kie Brand Kit

Keep the brand source of truth local. Use Kie for concepts, imagery, mockups, removal, and upscale; use deterministic local tooling for exact SVG logos, typography, copy, measurements, PDF, and presentation output.

Prefer `kie-pp-cli media workflow show brandkit --agent` and `kie-pp-cli create --workflow brandkit ... --agent` for compact, resumable execution. MCP hosts call `media_workflow_get` and pass `workflow: "brandkit"` to `media_brief_start`.

Read [references/brand-lock.md](references/brand-lock.md) before creating or modifying a brand system.

## Establish the project

Ask one question at a time:

1. What does the brand do and for whom?
2. What should it feel like, and what must it never feel like?
3. Which existing assets are authoritative?
4. Which deliverables and channels are required?
5. Which accessibility, cultural, legal, or trademark constraints apply?

Choose a local output folder and create a human-readable `brand-lock.json` or `brand-lock.yaml`. Record provenance for every supplied asset. Never overwrite approved assets silently.

## Approve in stages

Use this dependency order:

1. strategy and attributes;
2. palette with contrast checks;
3. three distinct logo directions;
4. selected logo geometry and local vector reconstruction;
5. typography with licensed/local font choices;
6. imagery and motion direction;
7. applications and brand book.

Pause after each stage for approval. If an upstream choice changes, mark dependent artifacts stale and regenerate them.

Use Kie to create concept boards and mockups through the durable director:

```bash
kie-pp-cli create --workflow brandkit "<brand concept or application>" --type image --purpose "brand system" \
  --style "<approved visual direction>" --reference <approved-asset> --agent
```

Prefer GPT Image 2 for broad concepts, Seedream 5 Pro image-to-image for reference-heavy mockups, and Recraft background removal/upscale only as post-processing. Kie does not currently expose Recraft vector generation; reconstruct final logo geometry locally as SVG and inspect it at small and large sizes.

## Preserve exactness

- Render exact names, taglines, legal copy, dimensions, and templates locally.
- Do not trust raster model text as the final logo or layout.
- Check WCAG contrast for digital colors.
- Keep font names, licenses, weights, fallbacks, and download/source records.
- Export monochrome, reversed, small-size, and clear-space logo variants.
- Distinguish generated mockups from production-ready artwork.

Review each Kie plan before `--submit`. Deliver the locked source file, asset manifest, editable vector/logo files, color/type tokens, application templates, and brand book. State which items are approved, generated, deterministic, stale, or awaiting review.
