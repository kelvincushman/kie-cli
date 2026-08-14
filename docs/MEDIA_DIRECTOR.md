# Kie Agent Media Factory: Director Contract

This document defines the local-first media creation workflow implemented by
`kie-pp-cli`. It is written for terminal users and local MCP/skill agents in
Codex, Claude Code, Cursor, and compatible hosts. The media director is the
control plane for an open-source agent media factory; it is not a browser-hosted
editor or a drop-in clone of Higgsfield's cloud product.

The concise product entry point is:

```bash
kie-pp-cli grill-me "what you want to make"
```

The director infers explicit prompt facts, qualifies the remainder one material
question at a time, shows its recommendation/rationale/overrides, saves durable
local handles, and stores reusable references in a private vault. Video has a
hard human-in-the-middle still gate and an optional complete-shot proof at the
model's lowest faithful resolution. Every live preview, proof, and final call
requires a separate fresh paid confirmation. Storyboard mode applies this
contract independently to every shot.

## Product Commands

Use the product commands for normal operation:

```bash
kie-pp-cli media setup --agent
kie-pp-cli media workflow list --agent
kie-pp-cli media workflow show <workflow> --agent
kie-pp-cli lesson recommend "<what you want to create>" --agent
kie-pp-cli media leaderboard <task> --agent
kie-pp-cli grill-me "a polished product image for a website hero" --agent
kie-pp-cli create "a polished product image for a website hero" --agent
kie-pp-cli create --workflow product-photoshoot "a polished product image" --agent
kie-pp-cli create --brief <brief_id> --answer <value> --agent
kie-pp-cli create --brief <brief_id> --wrap-up --agent
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
kie-pp-cli create --brief <brief_id> --approve-preview --agent
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
kie-pp-cli create --brief <brief_id> --approve-proof --agent
kie-pp-cli create --brief <brief_id> --reject-proof --agent
kie-pp-cli create --brief <brief_id> --skip-proof --agent
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media brief list --agent
kie-pp-cli media reference add <path-or-url> --name <name> --agent
kie-pp-cli media reference list --agent
kie-pp-cli media identity create <name> --reference <image> --consent --agent
kie-pp-cli media identity show <identity-id> --agent
kie-pp-cli media identity list --agent
kie-pp-cli media generation status <generation_id> --agent
kie-pp-cli media generation status <generation_id> --wait --agent
kie-pp-cli media capability list --capability kie-video --agent
kie-pp-cli media capability show <model-id> --agent
kie-pp-cli media models --family video --agent
```

Agents must still inspect current runtime truth before relying on copied command details:

```bash
kie-pp-cli doctor --json
kie-pp-cli agent-context --pretty
kie-pp-cli create --help
kie-pp-cli media --help
```

Raw generated endpoints such as `kie-pp-cli kie-ai-jobs market-create-task` and
`kie-pp-cli kie-ai-jobs market-query-task` are advanced/fallback surfaces. Do
not prefer them for the media-director workflow.

`kie-pp-cli media video` is a validated direct-generation shortcut for advanced
automation. Its built-in flags cover `wan/2-7-text-to-video`; another captured
model can be selected with `--model` and its exact input object supplied with
`--input`. Both paths validate the complete embedded per-model contract before
making a live call. This shortcut does not create a brief or enforce the still
preview/show/approve gate, so agents must not substitute it for the normal
director workflow when human confirmation is required.

## Design Goals

- Keep creative direction and reference selection local until the user approves live generation.
- Infer explicit facts, then ask exactly one material media question at a time,
  inspired by Matt Pocock's grilling structure but adapted to avoid an
  over-the-top creative intake. Include a recommended answer and reason.
- Use durable local IDs for briefs, references, identities, scripts,
  storyboards, and generations so agents pass handles instead of full context.
- Accept typed image, video, and audio paths, URLs, or `ref:<id>` handles from the local reference vault. Images support JPEG, PNG, GIF, WebP, BMP, and TIFF; video supports MP4, MOV, and MKV; audio supports MP3, WAV, AAC, M4A, and OGG. Local safety limits are 30 MiB for images, 200 MiB for video, and 15 MiB for audio.
- Keep reusable likenesses as consented local identity-reference bundles; do not claim biometric training or a cross-model Soul equivalent.
- Upload local reference files only during explicit live preview or final generation.
- Use existing Kie.ai authentication only: `kie-pp-cli auth setup` in an interactive terminal, or `KIE_BEARER_AUTH` through an environment secret store.
- When authentication is missing, CLI and MCP setup metadata may recommend the maintainer's Kie.ai referral link to support continued development. Agents must show its adjacent affiliate disclosure and must not describe it as a neutral link.
- Never estimate cost by default.
- For video briefs, treat `create --brief <id> --preview --confirm-paid --agent` as a separate live action that generates a still visual anchor and may consume credits. Still-image briefs use `--submit`, not `--preview`.
- For video, require `create --brief <id> --approve-preview --agent` after the preview has been shown. Silence or a ready brief is not approval.
- Offer a complete-shot proof at the selected model's lowest documented
  faithful tier after still approval. Proof generation is optional but the
  user must approve, reject, or skip it before the final gate.
- Require a fresh explicit user confirmation immediately before every live
  generation. The director additionally enforces a scoped, expiring, single-use
  record for preview, proof, and final calls. Broad generated endpoint tools are
  advanced manual surfaces and must not be described as protected by that
  director record.
- Apply the same separation to MCP: `media_paid_confirm`,
  `media_preview_generate`, display/approve, optional `media_proof_generate`,
  display/approve or skip, then a new confirmation and `media_generate`.
- Return machine-readable brief, reference, generation, task state, and result URL data.
- Keep `media models` backed by the canonical complete model registry; never
  maintain a second summary catalog that can drift from captured inputs and
  settings.
- Keep provider gaps explicit. Do not describe local identity bundles as trained
  Soul models or qualitative review as proprietary virality prediction.

## IDs and Local State

The implementation stores media state under the CLI's resolved data directory. Use `kie-pp-cli media setup --agent` or `kie-pp-cli agent-context --pretty` to discover the resolved store path instead of hardcoding it.

State is organized around opaque handles:

- `brief_...`: durable media brief ID returned by `kie-pp-cli create`.
- `ref_...`: reusable local reference ID returned by `kie-pp-cli media reference add`.
- `ref:<id>`: reference handle accepted by `kie-pp-cli create --reference`.
- `identity_...`: consented local likeness-reference bundle.
- `script_...`: versioned master script for a storyboard production.
- `storyboard_...`: ordered multi-shot production plan.
- `gen_...`: local generation ID returned after submission.

The local store contains brief answers, selected model, the generated plan, reference handles, provider task IDs, generation status, and result URLs. It must not contain credentials, bearer tokens, cookies, or auth headers. Public CLI/MCP reference responses omit local source and vault paths.

The workflow catalog keeps repeated routing metadata in Go instead of agent prompts. `media workflow list/show` returns the nine Kie-native skill ports, their compact stages, supported media, and unsupported provider gaps. Pass the selected name through `create --workflow <name>` or MCP `media_brief_start.workflow`.

## Workflow Catalog

| Workflow | Media | Purpose |
| --- | --- | --- |
| `generate` | Image, video | General model routing; audio and music use dedicated CLI commands |
| `academy` | Image, video | Source-linked lesson selection, original Kie method, storyboard, and per-shot approvals |
| `brandkit` | Image | Approval-led brand concepts and local brand-system handoff |
| `marketplace-cards` | Image | Truthful listing visuals and local exact-copy composition |
| `product-photoshoot` | Image | Reference-led product campaign imagery |
| `identity` | Image, video | Consented likeness references without biometric training |
| `video-explainer` | Image, video | Scripts, visual blocks, narration handoff, and local assembly |
| `websites` | Image, video | Site/app media assets with a separate local build/deploy step |
| `youtube-thumbnail` | Image | 16:9 concepts with local exact-text composition |

The 17 installable skills add a thin `kie-grill-me` entry, shared
`kie-grilling` primitive, capability skills for image/video/audio/avatar,
`kie-film`, the core `kie-create` director, identity and lesson skills, and the
outcome workflows above. Workflow metadata is intentionally compact for token
savings; model schemas stay in runtime introspection. See [SKILLS.md](SKILLS.md).

## Lesson Director

In hosts that expose installed skills as slash commands, `/kie-lesson` invokes
the `kie-lesson` skill. The stable cross-host terminal contract is:

```bash
kie-pp-cli lesson list --agent
kie-pp-cli lesson recommend "<what you want to create>" --agent
kie-pp-cli lesson show <course-slug/lesson-slug> --agent
kie-pp-cli lesson start "<what you want to create>" \
  --lesson <course-slug/lesson-slug> --agent
```

`lesson start` creates an `academy` workflow video brief in storyboard mode and
returns exactly one next question. An agent may instead pass the selected
`lesson` with `--workflow academy --type image` for an image-only job.

The checked-in map contains public lesson titles, ordering, duration, difficulty,
and source URLs plus original Kie-native method abstractions. It does not copy
course scripts, prompts, prose, videos, or downloads. See
[ACADEMY_METHODS.md](ACADEMY_METHODS.md) for all 171 mapped lessons.

## Qualification Protocol

Start with `grill-me` or MCP `media_grill_start`. The runtime fills facts stated
explicitly in the initial request. Ask the returned `next_question.prompt`, show
its recommendation and reason briefly, wait for the answer, update the brief,
then ask the next question. Do not batch an intake form or hard-code a question
order: conditional route, rights, identity, frame, audio, and reference
questions can appear. Use `--wrap-up`/`media_grill_wrap_up` only when the user
asks the agent to choose sensible remaining defaults.

Agent resume flow:

```bash
kie-pp-cli create --brief <brief_id> --answer <value> --agent
```

The command returns either the next single `next_question` or a qualified plan.
Plans include production/capability skills, route rationale, cost status, and
override options. Images next require plan review and paid confirmation. Videos
first return `generate_preview`; later state offers `generate_proof` or records
proof skip before allowing final submission.

## References

References can be supported local image/video/audio paths, HTTP(S) URLs, or `ref:<id>` handles. Local files use type-aware limits and symbolic links are rejected. A local path supplied to a brief is immediately replaced with a private `ref:<id>` handle before the brief is saved or returned.

Add reusable references to the local vault:

```bash
kie-pp-cli media reference add ./reference.png --name "hero product" --agent
kie-pp-cli media reference list --agent
```

Use vaulted references in briefs with the `ref:<id>` form:

```bash
kie-pp-cli create --brief <brief_id> --answer ref:<ref_id> --agent
```

or as an override:

```bash
kie-pp-cli create --brief <brief_id> --reference ref:<ref_id> --agent
```

Local files in the vault remain local until explicit live generation. On submission, the service resolves references and uploads local files to Kie as needed. URLs are remembered and passed through when accepted by the selected model path.

## Create, Preview, Approve, Submit, and Poll

Use `--agent` for every agent-run command. It produces parseable JSON and non-interactive defaults.

Start or resume a brief:

```bash
kie-pp-cli create "a cinematic product image for a landing page" --agent
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media brief list --agent
```

When the brief is ready, review the returned plan. The plan includes the selected `model`, provider `input`, and `rationale`.

For a video, generate and wait for the separate review image:

```bash
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
```

The completed preview generation returns `kind: "preview"` and a `result_urls` image. The terminal prints the URL; an agent or chat host with image rendering must display the image itself to the user. Do not approve from prompt text, metadata, silence, or an agent's private judgment.

If the image is wrong, reject it, revise the brief, and regenerate:

```bash
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --style "revised direction" --agent
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
```

If the user explicitly approves the displayed image:

```bash
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

The approved still becomes the final video's visual anchor. SeedDance
text/first-frame routes use it as `first_frame_url`; multimodal routes include
it in `reference_image_urls`. Editing a creative brief field invalidates the
approval and requires a new preview.

Offer the optional complete-shot proof at `proof_option.lowest_faithful_tier`.
For Seedance 2.5 this resolves to 480p. Decline the optional proof without
making a live proof call:

```bash
kie-pp-cli create --brief <brief_id> --skip-proof --agent
```

Accept the optional proof only after a fresh paid confirmation. Generate it,
show the whole clip, and then record approval or rejection:

```bash
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
kie-pp-cli create --brief <brief_id> --approve-proof --agent
# or --reject-proof
```

Then ask for a new confirmation for the exact final render:

```bash
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

Poll an existing generation:

```bash
kie-pp-cli media generation status <generation_id> --agent
kie-pp-cli media generation status <generation_id> --wait --agent
```

Report:

- `brief_id`
- `generation_id`
- provider `task_id`
- selected model
- current status
- result URLs when available
- any warning or failure reason

## Script and Storyboard Production

Use `--production-mode storyboard` for multi-shot video. The master brief holds
the production intent but is never submitted as one prompt-only video.

```bash
kie-pp-cli create "a 30-second product story" --type video --duration 30 \
  --production-mode storyboard --agent
kie-pp-cli media script set <brief_id> --file script.md --agent
kie-pp-cli media script show <brief_id> --agent
kie-pp-cli media script approve <brief_id> --agent
kie-pp-cli media storyboard set <brief_id> --file storyboard.json --agent
kie-pp-cli media storyboard show <brief_id> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

Script and storyboard writes and decisions are local-only. A storyboard must
contain 1-60 ordered shots, each 4-30 seconds, and their durations must total
the master duration. Saving it creates one standard `shot_brief_id` per shot.
Every shot uses the normal preview/show/approve/final gate. Script edits clear
the current storyboard; creative master-brief edits make existing approvals
stale. Final clip ordering, captions, audio mix, and export remain a transparent
local ffmpeg, Remotion, or editor step.

## Defaults

The implemented defaults are:

- Image without references: `gpt-image-2-text-to-image`
- Image with references: `gpt-image-2-image-to-image`
- Video: `bytedance/seedance-2-5`
- Video duration recommendation: `5` seconds
- Video defaults to 720p, generated audio off, and MP4 output.
- SeedDance modes are text, first-frame, first+last-frame, or multimodal image/video/audio references.
- Platform aspect-ratio recommendations: Instagram image `3:4`, Instagram/TikTok video `9:16`, YouTube/website/LinkedIn `16:9`, general video `16:9`, general image `1:1`

The preview-still model can be selected independently with `--preview-model`.
Supported routes are GPT Image 2 text/image modes and Nano Banana 2, 2 Lite,
or Pro. Use `media leaderboard character-consistency --agent` to inspect dated
task evidence, then validate the actual plan input against the captured Kie
schema. Model override is available with `--model`, but agents should prefer
the implemented defaults unless the user requests a specific model or current
evidence and input requirements support a better route.

## Approval and MCP

Image generation needs plan approval plus a fresh paid confirmation. Video
requires separate still, optional proof, and final actions. First create and
approve the preview:

```bash
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
# Show result_urls[0] and obtain explicit user approval.
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

Decline the optional proof without making a live proof call:

```bash
kie-pp-cli create --brief <brief_id> --skip-proof --agent
```

Or accept the optional proof only after a fresh paid confirmation, then show
the whole result and record approval or rejection:

```bash
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
# Show the whole proof and approve or reject it.
kie-pp-cli create --brief <brief_id> --approve-proof --agent
```

Finally, obtain another fresh confirmation for the exact final render:

```bash
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --agent
```

Do not submit merely because a brief is qualified. The final service rejects every video brief whose current preview has not been returned and explicitly approved. Preview approval is fingerprinted to the current brief, so later creative changes make it stale.

A brief can be submitted only once. After submission it returns `next_action: "check_generation_status"`; create a new brief for another paid generation.

For MCP, brief, reference, identity, and approval/skip tools are local state
actions. `media_preview_generate`, `media_proof_generate`, and `media_generate`
are distinct live actions and may each consume credits. Precede every one with
a fresh `media_paid_confirm`, poll it with `media_generation_status`, render the
artifact in the host, then record the user's decision. Confirmation validation
occurs before a live API request.

The focused MCP binary is `kie-media-mcp`. It uses the official MCP Go SDK and
supports the finalized `2026-07-28` stateless protocol over loopback Streamable
HTTP:

```bash
# Local child-process transport (recommended)
claude mcp add kie-media -- kie-media-mcp

# Or local-only stateless HTTP
kie-media-mcp --transport http --addr 127.0.0.1:7780
```

The HTTP listener intentionally rejects wildcard and non-loopback addresses.
The 40 stable focused tools are:

- `media_setup_get`
- `media_grill_start`
- `media_grill_answer`
- `media_grill_wrap_up`
- `media_brief_start`
- `media_brief_answer`
- `media_brief_get`
- `media_workflow_list`
- `media_workflow_get`
- `media_lesson_list`
- `media_lesson_get`
- `media_lesson_recommend`
- `media_leaderboard_get`
- `media_model_list`
- `media_model_get`
- `media_model_example`
- `media_model_validate`
- `media_capability_list`
- `media_capability_get`
- `media_reference_add`
- `media_reference_list`
- `media_identity_create`
- `media_identity_get`
- `media_identity_list`
- `media_paid_confirm`
- `media_preview_generate`
- `media_preview_approve`
- `media_preview_reject`
- `media_proof_generate`
- `media_proof_approve`
- `media_proof_reject`
- `media_proof_skip`
- `media_script_set`
- `media_script_get`
- `media_script_decide`
- `media_storyboard_set`
- `media_storyboard_get`
- `media_storyboard_decide`
- `media_generate`
- `media_generation_status`

See [Proof, Approval, and Paid Confirmation](PROOF_AND_PAID_CONFIRMATION.md)
for transaction and invalidation details.

## Sources

- Higgsfield CLI: https://higgsfield.ai/cli
- Higgsfield skills repository: https://github.com/higgsfield-ai/skills
- Matt Pocock Grill With Docs skill: https://github.com/mattpocock/skills/tree/main/skills/engineering/grill-with-docs
- Kie API docs: https://docs.kie.ai/
- Kie File Upload API quickstart: https://docs.kie.ai/file-upload-api/quickstart
- Official MCP Go SDK v1.7.0: https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.7.0
- MCP 2026-07-28 release: https://blog.modelcontextprotocol.io/posts/2026-07-28/
