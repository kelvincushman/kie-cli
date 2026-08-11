---
name: kie-youtube-thumbnail
description: Concept, generate, edit, and validate truthful YouTube thumbnail variants with Kie.ai imagery, consented identity references, and deterministic local text/logo composition. Use for thumbnail strategy, 16:9 variants, faces, products, logos, short headlines, mobile readability, or A/B-ready delivery.
---

# Kie YouTube Thumbnail

Read [references/thumbnail-frameworks.md](references/thumbnail-frameworks.md) before choosing a concept.

Prefer `kie-pp-cli media workflow show youtube-thumbnail --agent` and one `kie-pp-cli create --workflow youtube-thumbnail ... --agent` per variant. MCP hosts use the same workflow name through `media_workflow_get` and `media_brief_start`.

## Qualify the promise

Ask one question at a time:

1. What is the video's truthful promise and target viewer?
2. What moment, object, transformation, or tension can represent it?
3. Are zero to three people needed, and are their likeness references authorized?
4. Is there a style donor, product, logo, or brand lock?
5. Should the thumbnail include a two-to-four-word headline?
6. How many distinct variants are required?

Do not invent outcomes, quotes, reactions, wealth, rankings, or before/after evidence absent from the video.

## Develop and generate

Brainstorm at least five meaningfully different concepts. Select the strongest concepts by clarity, curiosity, truthfulness, face/object recognition, and small-size readability.

Preserve reference order: identity/face references first, then product and logo. Use a consented local identity through `$kie-identity` when needed.

Generate each variant as a separate 16:9 task; do not ask one task to produce a contact sheet:

```bash
kie-pp-cli create --workflow youtube-thumbnail "<one clean thumbnail concept without final text>" --type image \
  --purpose "YouTube thumbnail" --platform youtube --aspect-ratio 16:9 \
  --reference <ordered-ref> --model gpt-image-2-image-to-image --agent
```

Use Seedream 5 Pro image-to-image for reference-heavy concepts. Use a fresh edit brief for narrow corrections rather than silently overwriting the source result.

Add exact headline and logo locally with a licensed font and deterministic canvas/image tool when requested. Keep a clean image-only version. Never rely on model-rendered small text as the final deliverable.

## Inspect at delivery size

Review each thumbnail at full size and around 120 pixels wide. Reject variants with weak silhouette, unreadable expression/object, false promise, clutter, accidental text, identity drift, distorted products, unsafe crop, or insufficient contrast.

Deliver passing variants with semantic labels, clean and text-baked files, aspect and dimensions, brief/generation/task IDs, identity/reference provenance, and a short truthfulness/readability note. Do not claim a click-through-rate prediction without actual experiment data.
