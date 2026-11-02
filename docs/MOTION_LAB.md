# Motion Lab: deterministic motion + Kie.ai media

This is an **optional, review-led extension** to the existing local-first Kie media director, not a replacement for its paid-generation and creative-approval gates.

The working idea is taken from production practice: a compelling result usually requires an explicit shot brief, measurable timing, repeatable renders and several small review rounds. The output is not "a one-prompt miracle". It is a versioned creative asset with evidence.

## Two production lanes

| Decision | Code-motion lane | Kie.ai lane |
| --- | --- | --- |
| Best for | Exact branding, typography, 2D/3D camera moves, interface animations, precise beats and loops | Cinematic shots, detailed surfaces, real-world scenes, visual ideas and image-to-video |
| Source | HTML/CSS/Canvas/Three.js/GSAP/Remotion or a pure-time renderer | The current Kie model registry and approved prompts |
| Edit | Change a value and re-render any frame | Regenerate/extend a shot, subject to provider behaviour |
| QA | Deterministic frame capture, transition stills, contact sheet, endpoint seam test | Human preview approval, full-shot proof if accepted, shot consistency inspection |
| Constraints | Local dependencies, render capacity | Credit spend, model availability, provider retention |

A **hybrid** asset starts with the code-motion lane, asks Kie for only those shots or visual layers that truly benefit from generation, then composites the approved outputs through a deterministic timeline. This prevents generative drift from affecting exact logos, captions, diagrams, UI, subtitles and required product details.

## Contract: one frame from one time

Each animation must expose `window.renderAt(t)` or an equivalent pure function where `t` is seconds. Returning to the same `t` must reproduce the same pixels. No wall-clock time, timers, unseeded random data, network-dependent visuals or autonomous requestAnimationFrame loops during frame export.

- Fix canvas size, duration, frame rate, language, typefaces and brand-safe areas up front.
- Use `N = round(duration * fps)`; render only frames `0..N-1` to the encoded sequence.
- For a seamless loop, explicitly render a **separate** diagnostic at `t=duration` and compare it with the frame at `t=0`. Never encode that diagnostic as frame N.
- Periodic functions, closed-path camera movement and deterministic noise must return to their starting state at the seam.
- Font and external asset loading must complete before frame 0; fail the render on missing dependencies.
- All timing-critical state, including particle positions, comes from `t`, not previous frame state.
- Frame checks include boundary frames before and after each major key, the title hold, and the loop seam.
- On audio, use absolute sample positions, not accumulation of rounded beat lengths; inspect decoded media, not only the source WAV.

## Required creative brief

```yaml
goal: "What must a viewer understand or feel?"
format: "16:9 / 1:1 / 4:5 / 9:16"
canvas: "1280x720"
duration_seconds: 6
fps: 24
loop: true
brand:
  palette: ["#070b12", "#f5b301", "#f4f5f7"]
  typography: "Inter + Geist Mono or explicit licensed fallbacks"
  logo_lock: "use provided exact logo file; never ask a model to redraw it"
look: "physical materials, lighting, shot framing, texture constraints"
shots:
  - range: "0-36"
    hero: "subject / action / precise outcome"
    cue: "main hit on an exact frame"
  - range: "36-96"
    hero: "second action"
  - range: "96-144"
    hero: "return to the starting composition"
motion:
  easing: "named cubic-bezier values"
  spring_damping: 0.75
  stagger_frames: 3
audio: "optional, absolute frame events + loudness/headroom checks"
banned: ["generic neon glow", "random particle decorations", "arbitrary brand text"]
qa:
  contact_sheet: true
  inspect_boundary_frames: true
  exact_endpoint_seam: true
  sign_off: "human approval before publication"
```

Use numbers as initial *creative choices*, not universal rules. The brief should name excluded visual clichés and what must not change.

## Render → inspect → critique → revise

1. **Plan** one hero movement per shot with named keyframes, 10–30% anticipation or overshoot only where it serves the animation, and explicit overlap for supporting elements.
2. **Prototype with code** using low-resolution frames. Generate a 12-frame contact sheet, a seekable preview and a diagnostic frame for the loop endpoint.
3. **Inspect** keyframes and adjacent frames, including text readability, moving edges, empty margins and brand fidelity. Record the observed issue with a frame number.
4. **Revise minimally**: change the cause, not the whole brief. Preserve all other validated shot constraints. Keep a prompt/fix ledger.
5. **Choose the Kie inserts**: still background, texture, cinematic opening/closing shot or other assets that cannot be efficiently drawn in code. Use `kie-pp-cli media capability list` and `models show`; never assume a specific model, parameter or price.
6. **Obtain approval** before **every** Kie paid call. For video: paid still preview → display and explicit still approval → optional complete-shot proof and decision → separately confirmed final. The code-only preview has no authority to bypass the provider gates.
7. **Store immediately**: Kie.ai retains generated media temporarily. Copy approved result to the tenant/local asset store and retain provider task ID, model, settings, source reference, hash, rights notes and generation timestamps.
8. **Composite/export** in code: render the exact timeline, then add sound, captions, title and logo. Produce H.264 MP4, a brand-safe still and a shareable short loop.
9. **Audit and approve**: record frame count, endpoint seam result, contact sheet, audio sync notes, aspect ratio, encoded media probe, approval and publication status.

The scripts in `examples/motion-lab/` perform local frame QA. They **do not call Kie** and **do not approve media**.

## Production manifest

Use one manifest per creative version, and write a new version after every change:

```json
{
  "project_id": "kelvin-motion-lab",
  "version": 1,
  "stage": "review",
  "source": {"engine": "code", "commit": "replace-with-real-commit"},
  "timeline": {"width": 1280, "height": 720, "fps": 24, "duration_seconds": 6},
  "kie": [{"asset_id": null, "brief_id": null, "generation_id": null, "task_id": null, "status": "not_requested"}],
  "qa": {"frame_count_ok": null, "endpoint_seam_ok": null, "keyframe_review": "pending", "audio_review": "not_applicable"},
  "review": {"creative_approval": false, "paid_confirmation": false, "publish_approval": false},
  "deliverables": {"video": null, "thumbnail": null, "contact_sheet": null},
  "provenance": "Do not claim Kie generated code-rendered assets"
}
```

Never store credentials, personal likenesses, raw consent evidence or private reference files in a public manifest. Preserve a truthful origin field when code, user footage, a provider or AI images are mixed.

## Worked KelvinLee.io example

A six-second, 24 fps, 1280 × 720 black-and-gold identity animation:

- 0–1.5 s: "AI MOTION LAB" appears against a measured orbit, with no invented logo.
- 1.5–3 s: BRIEF → CODE is emphasised; a smooth camera-like glide connects the two stages.
- 3–4.5 s: a **previously approved** Kie visual may occupy an image window while the motion engine owns the overlay.
- 4.5–6 s: QA + EXPORT closes the cycle, preserving a mathematically periodic seam.
- Gold `#f5b301`, off-white type, near-black background; restrained grid; no competing gimmicks.
- Key title and company URL are composed in code; the model is not asked to typeset them.
- Export 144 frames `0..143`; separately check `renderAt(0)` against `renderAt(6)`; verify every key shot via contact sheet.

Example Kie **still** prompt (planning only; no paid submission):

> Editorial 16:9 abstract cinematic macro image illustrating deterministic animation meeting generated media: sculptural brushed gold orbit around a matte-black geometric core, restrained studio lighting, ink-dark background, carefully controlled reflections, sophisticated British engineering/editorial aesthetic. No words, no logos, no UI, no watermark, clear negative space for a code-composited title. Honour the reference palette #070b12 and #f5b301. No purple/blue neon, random particles or distracting glow.

Example Kie **motion insert** prompt (apply only to an approved still):

> Six-second gentle cinematic orbit of the supplied gold geometric sculpture; maintain subject proportions, core geometry, material finish and lighting direction. Restrained camera movement, no cuts, no lettering, no hallucinated branding. Provide start/end compositions suitable for a code-controlled crossfade; any final loop seam must be validated after export rather than assumed from provider output.

## Credit and publication safety

Follow `docs/PROOF_AND_PAID_CONFIRMATION.md`. A user approving the brief does **not** authorise a provider call. Preview, proof, regeneration and final each need a fresh consent transaction. A local code render has no credit cost but must still be marked as unapproved for publication until reviewed. Do not automatically post a blog or social clip.

## Further reading

- Florence, *Opus 5.5 Motion Lab: Make Claude Animate Like a Studio* (30 September 2026) — experiment and reporting approach; do not copy article or assets.
- Anthropic, Opus 5.5 model documentation (released 22 September 2026).
- Kie.ai API documentation — task creation, polling, expiry, and account limits.
- Three.js, Remotion, GSAP or HyperFrames documentation for frame-addressable code motion.

