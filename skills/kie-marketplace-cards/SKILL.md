---
name: kie-marketplace-cards
description: Plan and produce truthful e-commerce marketplace image sets with Kie.ai product imagery and deterministic local layouts. Use for Amazon-style main images, galleries, infographics, A+ modules, lifestyle scenes, feature cards, comparison graphics, or listing refreshes that require evidence, exact copy, and platform-aware QA.
---

# Kie Marketplace Cards

Start with `$kie-grilling`, then route each approved card brief through
`$kie-image`. Keep product claims, copy, measurements, badges, and final layout
deterministic and truthful.

Use Kie to generate or edit visuals. Compose exact copy, diagrams, badges, dimensions, and final marketplace canvases locally. Never claim that Kie or this skill certifies marketplace compliance.

Prefer `kie-pp-cli media workflow show marketplace-cards --agent` and `kie-pp-cli create --workflow marketplace-cards ... --agent` to keep repeated orchestration in the CLI. MCP hosts use the same name through `media_workflow_get` and `media_brief_start`.

Read [references/asset-matrix.md](references/asset-matrix.md) before planning a set.

## Qualify the listing

Ask one question at a time. Capture:

1. marketplace, country, category, and requested scope;
2. product truth: materials, dimensions, quantity, ingredients, included items, and verified claims;
3. authoritative product, packaging, logo, and brand references;
4. target buyer, objections, differentiators, and tone;
5. prohibited claims, required disclaimers, and current platform rules;
6. output dimensions, language, and variant count.

Reject or soften claims that lack supplied evidence. Confirm current marketplace rules from official sources when compliance matters.

## Produce the set

Create a shot plan before generating. Use Seedream 5 Pro image-to-image for product/reference fidelity or GPT Image 2 image-to-image as an alternative. Generate each concept as its own durable brief so failed variants can be resumed independently:

```bash
kie-pp-cli create --workflow marketplace-cards "<asset-specific visual brief>" --type image --purpose "marketplace listing" \
  --aspect-ratio 1:1 --reference <product-ref> --reference <brand-ref> \
  --model seedream/5-pro-image-to-image --agent
```

For the main image, preserve the actual product and packaging. Do not invent accessories or alter quantity. For secondary cards, separate model-generated imagery from exact local copy and diagrams.

## Validate

Obtain a fresh paid confirmation before every variant, edit, background removal,
or upscale. An approval of the set direction is not payment authorization for
the next asset.

Inspect every final flattened asset for:

- product identity, color, count, and packaging fidelity;
- legibility at mobile thumbnail size;
- truthful, evidence-backed copy;
- required dimensions, safe areas, and file format;
- no unsupported badges, rankings, testimonials, medical claims, or guarantees;
- consistent palette, type, spacing, and reference provenance.

Submit live Kie tasks only after the user approves each ready plan. Deliver semantic filenames, the editable local layout source, generated source URLs/job IDs, a claims-evidence ledger, and a compliance caveat tied to the marketplace and review date.
