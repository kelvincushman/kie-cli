---
name: kie-product-photoshoot
description: Direct consistent Kie.ai product photography and campaign imagery from local product, packaging, brand, environment, and identity references. Use for studio shots, lifestyle scenes, detail images, moodboards, hero banners, carousels, ad packs, virtual try-ons, conceptual product art, or restyling.
---

# Kie Product Photoshoot

Start with `$kie-grilling`, then route each still through `$kie-image` and any
motion shot through `$kie-video`. Reuse the approved product references and
continuity lock; do not repeat the questionnaire for every asset.

Read [references/shoot-modes.md](references/shoot-modes.md), then ask one question at a time until the product, mode, audience, placement, references, visual direction, and number of variants are clear.

Prefer `kie-pp-cli media workflow show product-photoshoot --agent` and `kie-pp-cli create --workflow product-photoshoot ... --agent` for token-efficient, durable runs. MCP hosts use `workflow: "product-photoshoot"`.

Vault authoritative references before generation:

```bash
kie-pp-cli media reference add <product-image> --type image --name "product front" --agent
kie-pp-cli media reference add <logo-image> --type image --name "brand logo" --agent
```

Use Seedream 5 Pro image-to-image for reference-heavy product scenes and up to its currently documented ten references; verify the live schema before relying on that limit. Use GPT Image 2 image-to-image when its edit behavior better fits the task. Use Recraft background removal only for deterministic post-processing.

Create one durable brief per shot or variant:

```bash
kie-pp-cli create --workflow product-photoshoot "<single shot brief>" --type image --purpose "product campaign" \
  --platform <destination> --aspect-ratio <ratio> --style "<locked shoot direction>" \
  --reference ref:<product-ref> --reference ref:<logo-ref> \
  --model seedream/5-pro-image-to-image --agent
```

Keep a shot ledger with shot ID, mode, prompt, references, model, task ID, output, approval state, and retry lineage. Do not batch conceptually different shots into one task.

## Protect product truth

- Preserve shape, proportions, color, finish, label hierarchy, closure, and included quantity.
- Treat logos and small packaging copy as reference geometry; repair exact text locally.
- Repair exact copy with an available deterministic compositor such as SVG/HTML canvas or ImageMagick, and retain the editable source.
- Do not invent performance, ingredients, certifications, accessories, or usage claims.
- Distinguish a concept render from an accurate catalog image.
- Obtain consent before using a saved human identity; use `$kie-identity` for reusable likeness bundles.
- For virtual try-on, state that fit, drape, scale, and material remain illustrative unless verified.

Review the plan and obtain a fresh paid confirmation before each submission,
regeneration, edit, or upscale. Inspect results at full size and target placement
size. Deliver only passing variants with semantic names, task/result provenance,
and any fidelity caveat.
