---
name: kie-motion-lab
description: Build inspectable code-generated motion graphics, combine them with separately approved Kie.ai image/video inserts, and ship reproducible animations with contact sheets, frame-level revision notes, export QA and provenance. Use for branded launch loops, code motion, creative media pipelines and AI-assisted studio-style animation.
---

# Kie Motion Lab — code and generated media

Read `docs/MOTION_LAB.md` first. This skill adds a **code-motion lane** alongside the existing `kie-create`, `kie-image`, `kie-video`, `kie-film` and `kie-audio` skills. It is not a replacement for them and must never bypass paid-action approvals.

## Route

1. Identify the deliverable: exact 2D/3D typographic motion, cinematic generated footage, or a hybrid. Prefer deterministic code for logos, titles, precise UI, diagrams and beat-synchronised loops; use Kie for genuinely generative visuals.
2. Inspect the existing repository/framework, brand tokens, approved reference media and required output sizes. Keep the code motion build local.
3. Write a short, numerical production brief: `GOAL`, `OUTPUT`, `TIME`, `LOOK`, `SHOTS`, `MOTION`, `BANNED`, `QA`, `AUDIO`.
4. Require `renderAt(t)` (or equivalent). Render frame `i` at `t=i/fps` for `i=0..N-1`; do not render duplicate frame N into a loop. Separately render `t=duration` for an exact seam comparison.
5. Review a 12-still contact sheet and key transitions, writing each correction as a frame-numbered, testable change. Prefer two to four **small** correction rounds over a total re-prompt.
6. If Kie is useful, run `kie-pp-cli media capability list`, `models show` and the local director's brief route. Never guess current schemas or costs. Do not expose secrets in HTML, the browser, media prompt or commit.
7. For **every** paid image generation, ask for a fresh explicit yes after showing model/settings/known cost and reference uploads. For **each** video shot, obtain fresh confirmation for the mandatory preview still; show the actual output and wait for explicit still approval, then offer an optional separately confirmed complete-shot proof, then get another fresh yes for final rendering. Use `--confirm-paid` only after permission. Preserve the director's brief IDs and approval state.
8. Copy provider results before expiry, record provider task IDs and media provenance, composite only approved inputs into a code-owned timeline, export and review the finished encoded MP4/GIF/still.
9. Produce a manifest with FPS, duration, frame count, deterministic seam check, contact sheet path, versions, model/task IDs where used, audio sync and review status. Mark unavailable evidence `not_tested`; never invent successful Kie generations.

## Local QA

The optional `examples/motion-lab/qa.py` uses Pillow to inspect existing PNG frames. Install `Pillow` locally before running. It writes a contact sheet and structured QA report; it does **not** render with Kie or make provider calls.

```bash
python examples/motion-lab/qa.py --frames ./out/frames --fps 24 --duration 6 \
  --endpoint ./out/seam-endpoint.png --out ./out/qa
```

When the environment cannot run Playwright, headless Chromium or a GPU renderer, report the limitation. A code preview is not a verified Kie render and neither is publication approval.

## Example first-party brief

> Six seconds, 24fps, 1280x720, seamless KelvinLee.io black/gold motion identity. Code controls all text, exact geometry, timeline and loop. One central gold sculptural element, four production-stage cards (BRIEF, CODE, KIE.AI, QA + EXPORT), no purple gradients, generic neon, misdrawn logo or invented claims. `renderAt(t)` is periodic in six seconds. Render keyframes, contact sheet, audio plan, then one visual revision at a time. Kie may supply an approved visual plate only after the paid-gate workflow; never use generated text as the final company identity.

When the user asks for a blog thumbnail alongside the video, create an exact title and metadata in the code layer; reserve Kie for the unlettered illustration plate, and preserve enough negative space for different responsive crops.
