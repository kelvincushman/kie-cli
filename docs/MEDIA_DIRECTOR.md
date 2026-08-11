# Kie Agent Media Factory: Director Contract

This document defines the local-first media creation workflow implemented by
`kie-pp-cli`. It is written for terminal users and local MCP/skill agents in
Codex, Claude Code, Cursor, and compatible hosts. The media director is the
control plane for an open-source agent media factory; it is not a browser-hosted
editor or a drop-in clone of Higgsfield's cloud product.

The product entry point is:

```bash
kie-pp-cli create "what you want to make"
```

The director qualifies a brief one question at a time, saves durable local
handles, stores reusable references in a private vault, and returns generation
state plus output URLs. Video has a hard human-in-the-middle gate: generate a
still preview, show it to the user, record explicit approval, and only then
submit the final video. Storyboard mode applies this contract independently to
every shot.

## Product Commands

Use the product commands for normal operation:

```bash
kie-pp-cli media setup --agent
kie-pp-cli media workflow list --agent
kie-pp-cli media workflow show <workflow> --agent
kie-pp-cli create "a polished product image for a website hero" --agent
kie-pp-cli create --workflow product-photoshoot "a polished product image" --agent
kie-pp-cli create --brief <brief_id> --answer <value> --agent
kie-pp-cli create --brief <brief_id> --preview --wait --agent
kie-pp-cli create --brief <brief_id> --approve-preview --agent
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --submit --agent
kie-pp-cli create --brief <brief_id> --submit --wait --agent
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media brief list --agent
kie-pp-cli media reference add <path-or-url> --name <name> --agent
kie-pp-cli media reference list --agent
kie-pp-cli media identity create <name> --reference <image> --consent --agent
kie-pp-cli media identity show <identity-id> --agent
kie-pp-cli media identity list --agent
kie-pp-cli media generation status <generation_id> --agent
kie-pp-cli media generation status <generation_id> --wait --agent
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

## Design Goals

- Keep creative direction and reference selection local until the user approves live generation.
- Ask exactly one media-qualification question at a time, inspired by Matt Pocock's Grill With Docs workflow but adapted from software-domain alignment to media-brief alignment.
- Use durable local IDs for briefs, references, identities, scripts,
  storyboards, and generations so agents pass handles instead of full context.
- Accept typed image, video, and audio paths, URLs, or `ref:<id>` handles from the local reference vault. Images support JPEG, PNG, GIF, WebP, BMP, and TIFF; video supports MP4, MOV, and MKV; audio supports MP3, WAV, AAC, M4A, and OGG. Local safety limits are 30 MiB for images, 200 MiB for video, and 15 MiB for audio.
- Keep reusable likenesses as consented local identity-reference bundles; do not claim biometric training or a cross-model Soul equivalent.
- Upload local reference files only during explicit live preview or final generation.
- Use existing Kie.ai authentication only: `kie-pp-cli auth setup` in an interactive terminal, or `KIE_BEARER_AUTH` through an environment secret store.
- Never estimate cost by default.
- Treat `create --brief <id> --preview --agent` as a separate live image-generation action that may consume credits.
- For video, require `create --brief <id> --approve-preview --agent` after the preview has been shown. Silence or a ready brief is not approval.
- Treat `create --brief <id> --submit --agent` as the explicit final-generation action. `--wait` only changes whether the command polls after submitting.
- Apply the same separation to MCP: `media_preview_generate`, display/poll, `media_preview_approve`, then `media_generate`.
- Return machine-readable brief, reference, generation, task state, and result URL data.
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

The workflow catalog keeps repeated routing metadata in Go instead of agent prompts. `media workflow list/show` returns the eight Kie-native skill ports, their compact stages, supported media, and unsupported provider gaps. Pass the selected name through `create --workflow <name>` or MCP `media_brief_start.workflow`.

## Workflow Catalog

| Workflow | Media | Purpose |
| --- | --- | --- |
| `generate` | Image, video | General model routing; audio and music use dedicated CLI commands |
| `brandkit` | Image | Approval-led brand concepts and local brand-system handoff |
| `marketplace-cards` | Image | Truthful listing visuals and local exact-copy composition |
| `product-photoshoot` | Image | Reference-led product campaign imagery |
| `identity` | Image, video | Consented likeness references without biometric training |
| `video-explainer` | Image, video | Scripts, visual blocks, narration handoff, and local assembly |
| `websites` | Image, video | Site/app media assets with a separate local build/deploy step |
| `youtube-thumbnail` | Image | 16:9 concepts with local exact-text composition |

The nine installable skills are the core `kie-create` director plus one skill
for each workflow above. Workflow metadata is intentionally compact for token
savings; load a specialized skill only when its domain guidance is needed.

## Qualification Protocol

Ask one question, wait for the answer, update the brief, then ask the next question. Do not batch an intake form.

The implemented v1 question order is exactly:

1. `request`: "What do you want to create? Describe the subject and desired outcome."
2. `media_type`: "Should this be an image or a video?"
3. `purpose`: "What will this media be used for?"
4. `platform`: "Where will it be used?"
5. `aspect_ratio`: "Which aspect ratio should I use?"
6. `duration_seconds`: "How many seconds should the video run?" This appears only for video briefs.
7. `audio_mode`: "Should SeedDance generate synchronized audio?" Video only.
8. `video_mode`: "How should SeedDance guide the video?" Video only.
9. `style`: "What visual style or mood should it have?"
10. `first_frame`, then optionally `last_frame`, when a frame-controlled mode is selected.
11. `reference`: image by default; multimodal video accepts `video:` and `audio:` prefixes. This repeats until `skip`, `done`, or `none`.

Agent resume flow:

```bash
kie-pp-cli create --brief <brief_id> --answer <value> --agent
```

The command returns either the next single `next_question` or a qualified plan. Images return `can_submit: true` and `next_action: "review_then_submit"`. Videos first return `can_submit: false` and `next_action: "generate_preview"`.

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
kie-pp-cli create --brief <brief_id> --preview --wait --agent
```

The completed preview generation returns `kind: "preview"` and a `result_urls` image. The terminal prints the URL; an agent or chat host with image rendering must display the image itself to the user. Do not approve from prompt text, metadata, silence, or an agent's private judgment.

If the image is wrong, reject it, revise the brief, and regenerate:

```bash
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --style "revised direction" --agent
kie-pp-cli create --brief <brief_id> --preview --wait --agent
```

If the user explicitly approves the displayed image:

```bash
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

The approved still becomes the final video's visual anchor. SeedDance text/first-frame routes use it as `first_frame_url`; multimodal routes include it in `reference_image_urls`. Editing a creative brief field invalidates the approval and requires a new preview.

After user approval, submit:

```bash
kie-pp-cli create --brief <brief_id> --submit --agent
```

Submit and wait:

```bash
kie-pp-cli create --brief <brief_id> --submit --wait --agent
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

Model override is available with `--model`, but agents should prefer the implemented defaults unless the user requests a specific model or the current CLI/model docs require a different route.

## Approval and MCP

Image generation approval is explicit `--submit`. Video requires separate preview and final actions:

```bash
kie-pp-cli create --brief <brief_id> --preview --wait --agent
# Show result_urls[0] and obtain explicit user approval.
kie-pp-cli create --brief <brief_id> --approve-preview --agent
kie-pp-cli create --brief <brief_id> --submit --agent
```

Do not submit merely because a brief is qualified. The final service rejects every video brief whose current preview has not been returned and explicitly approved. Preview approval is fingerprinted to the current brief, so later creative changes make it stale.

A brief can be submitted only once. After submission it returns `next_action: "check_generation_status"`; create a new brief for another paid generation.

For MCP, brief, reference, identity, and preview approve/reject tools are local state actions. `media_preview_generate` and `media_generate` are distinct live actions and may each consume credits. Poll previews with `media_generation_status`, render the returned image URL in the host, ask for explicit approval, and only then call `media_preview_approve`. Use `media_preview_reject` when the user wants another direction.

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
The twenty-one stable tools are:

- `media_brief_start`
- `media_brief_answer`
- `media_brief_get`
- `media_workflow_list`
- `media_workflow_get`
- `media_reference_add`
- `media_reference_list`
- `media_identity_create`
- `media_identity_get`
- `media_identity_list`
- `media_preview_generate`
- `media_preview_approve`
- `media_preview_reject`
- `media_script_set`
- `media_script_get`
- `media_script_decide`
- `media_storyboard_set`
- `media_storyboard_get`
- `media_storyboard_decide`
- `media_generate`
- `media_generation_status`

## Sources

- Higgsfield CLI: https://higgsfield.ai/cli
- Higgsfield skills repository: https://github.com/higgsfield-ai/skills
- Matt Pocock Grill With Docs skill: https://github.com/mattpocock/skills/tree/main/skills/engineering/grill-with-docs
- Kie API docs: https://docs.kie.ai/
- Kie File Upload API quickstart: https://docs.kie.ai/file-upload-api/quickstart
- Official MCP Go SDK v1.7.0: https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.7.0
- MCP 2026-07-28 release: https://blog.modelcontextprotocol.io/posts/2026-07-28/
