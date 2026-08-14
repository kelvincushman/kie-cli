---
name: kie-websites
description: Build a local website, generative media app, or browser game while using Kie.ai for approved visual, video, audio, and marketing assets. Use for site concepts, landing pages, product sites, media-generation apps, game art, covers, animation references, textures, or asset pipelines where code and deployment remain local or use the user's chosen platform.
---

# Kie Websites

Start with `$kie-grilling`, then route still assets through `$kie-image` and
motion assets through `$kie-video`. This skill owns local code and asset
integration; Kie does not build or deploy the website.

Kie is the media provider, not a hosted full-stack website platform. Do not claim Higgsfield subdomains, repositories, Quanta/fnf SDKs, Cloudflare resources, multiplayer rooms, community publishing, contests, or deployment are available through Kie.

Prefer `kie-pp-cli media workflow show websites --agent` and `kie-pp-cli create --workflow websites ... --agent` for compact asset jobs. MCP hosts use `workflow: "websites"`; keep code/build/deploy operations in the local coding environment.

Read [references/project-routing.md](references/project-routing.md) before choosing the build lane.

## Define the product

Ask one question at a time:

1. Is this a website, a media-generation app, or a browser game?
2. What is the primary user outcome and conversion/action?
3. Which repository, framework, runtime, and deployment target are already in scope?
4. Which generated assets are required?
5. Which references, brand locks, identity permissions, and accessibility constraints apply?

Inspect the existing repository and its instructions before editing. Preserve its framework and design system unless the user requests a change.

## Separate code from media

Build code, auth, storage, API proxying, deployment, and security with the local coding environment and the user's authorized infrastructure. Never expose a Kie API key to browser JavaScript; call Kie from a trusted server route or local agent process.

Use the director for hero images, illustrations, product shots, backgrounds, covers, and short videos:

```bash
kie-pp-cli create --workflow websites "<site-specific asset brief>" --type image --purpose "website asset" \
  --platform website --aspect-ratio 16:9 --reference <brand-ref> --agent
```

Use `$kie-product-photoshoot`, `$kie-brandkit`, `$kie-identity`, or `$kie-video-explainer` when those specialized constraints apply.

For every short video asset, obtain a fresh paid confirmation, generate the
director preview, render the returned still in the review surface, and obtain
explicit approval with `--approve-preview`. Offer the low-resolution
complete-shot proof before the final. Use `--reject-preview` and revise when the
hero composition, crop, product fidelity, identity, or safe text area is wrong.
Every preview, proof, and final is a separate paid action.

For games, use Kie for concept art, spritesheet references, textures, UI art, backgrounds, sound, and music. Kie does not currently document true 3D mesh/rig/animation generation; use user-supplied models or an explicitly approved external 3D tool. Treat generated textures and animation references as source material that still needs local engine validation.

## Validate the result

- Run repository lint, typecheck, tests, and production build.
- Inspect responsive breakpoints, keyboard access, contrast, reduced motion, loading states, and asset fallbacks.
- Optimize generated media and retain licenses/provenance.
- Test server-side secret isolation and authorization boundaries.
- Verify deployment only on the platform the user authorized.
- Do not publish to a community feed or contest unless that separate external action is explicitly requested.

Deliver the repository changes, local preview/build evidence, generated asset manifest with Kie task IDs, deployment status if authorized, and explicit unsupported or unverified platform features.
