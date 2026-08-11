---
name: kie-create
description: Run the local-first Kie agent media factory to qualify, script, storyboard, approve, submit, resume, and poll image or video creation through kie-pp-cli or kie-media-mcp. Use for one-question-at-a-time creative intake, multi-shot production, SeedDance 2.5 multimodal video, local media references, consented identity bundles, or durable generation handles.
---

# Kie Create

Use this skill when the user wants to create, resume, revise, or poll an image/video generation brief through Kie. This is a local Codex, Claude Code, Cursor, or MCP-agent workflow. Do not route this workflow through browser-hosted Claude.ai.

When this skill is running from a repository checkout, read
`../../docs/MEDIA_DIRECTOR.md` for the longer product contract. If the skill was
installed by itself, continue with the self-contained workflow below.

## Runtime Discovery First

Inspect the installed CLI and current media command help:

```bash
kie-pp-cli doctor --json
kie-pp-cli agent-context --pretty
kie-pp-cli media setup --agent
kie-pp-cli media workflow list --agent
kie-pp-cli media workflow show <workflow> --agent
kie-pp-cli create --help
kie-pp-cli media --help
kie-pp-cli media brief --help
kie-pp-cli media reference --help
kie-pp-cli media generation --help
kie-pp-cli media script --help
kie-pp-cli media storyboard --help
```

Prefer product commands:

```bash
kie-pp-cli create "what the user wants to make" --agent
kie-pp-cli create --workflow <workflow> "what the user wants to make" --agent
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
kie-pp-cli media script set <brief_id> --file <script.md|-> --agent
kie-pp-cli media script approve <brief_id> --agent
kie-pp-cli media storyboard set <brief_id> --file <storyboard.json|-> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

Raw generated endpoints such as `kie-pp-cli kie-ai-jobs market-create-task` are
advanced/fallback surfaces. Do not use them as the normal media creation path.

Use `--agent` on CLI calls so stdout is parseable JSON and the CLI uses non-interactive agent defaults. Parse the returned `brief`, `next_question`, `ready`, `next_action`, `id`, `task_id`, `status`, and `result_urls` fields according to the command response.

## Credential Rules

Use the guided auth setup in an interactive terminal:

```bash
kie-pp-cli auth setup
kie-pp-cli auth status --agent
```

For agent hosts, set `KIE_BEARER_AUTH` through the host's secret store.

Never print, log, store, summarize, or expose credentials. Do not include bearer tokens, cookies, auth headers, or credential file contents in a brief, response, command transcript, or debugging note.

## Conversation Protocol

Qualify the brief one question at a time. Ask the current `next_question.prompt`, wait for the user's answer, then resume with:

```bash
kie-pp-cli create --brief <brief_id> --answer <value> --agent
```

The implemented v1 question order is exactly:

1. `request`: "What do you want to create? Describe the subject and desired outcome."
2. `media_type`: "Should this be an image or a video?"
3. `purpose`: "What will this media be used for?"
4. `platform`: "Where will it be used?"
5. `aspect_ratio`: "Which aspect ratio should I use?"
6. `duration_seconds`: "How many seconds should the video run?" Video only.
7. `audio_mode`: "Should SeedDance generate synchronized audio?" Video only.
8. `video_mode`: "How should SeedDance guide the video?" Video only.
9. `style`: "What visual style or mood should it have?"
10. `first_frame` and optionally `last_frame`: requested only by the selected frame mode.
11. `reference`: image by default; in multimodal video, prefix answers with `video:` or `audio:`. Repeat until `skip`, `done`, or `none`.

Stay inside the implemented director fields. If the user volunteers extra creative requirements, fold them into the request or style answer only when that preserves the implemented command shape.

## Briefs, References, and Resumption

The director returns durable IDs:

- `brief_...` for briefs
- `ref_...` for vaulted references
- `gen_...` for generations

Inspect and resume:

```bash
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media brief list --agent
kie-pp-cli create --brief <brief_id> --answer <value> --agent
```

Add reusable references to the local vault:

```bash
kie-pp-cli media reference add ./reference.png --name "hero product" --agent
kie-pp-cli media reference add https://example.com/reference.jpg --name "campaign still" --agent
kie-pp-cli media reference list --agent
```

Use references as `ref:<id>` handles:

```bash
kie-pp-cli create --brief <brief_id> --answer ref:<ref_id> --agent
kie-pp-cli create --brief <brief_id> --reference ref:<ref_id> --agent
```

Local files in the vault stay local until explicit live generation. The service uploads local reference files during submission and passes URLs when supported. Use `--type image`, `--type video`, or `--type audio` when URL type cannot be inferred.

Accept supported SeedDance image, video, and audio formats. The local vault supports JPEG, PNG, GIF, WebP, BMP, and TIFF images; MP4, MOV, and MKV video; and MP3, WAV, AAC, M4A, and OGG audio. Local safety limits are 30 MiB for images, 200 MiB for videos, and 15 MiB for audio. The CLI rejects symbolic links and returns reference handles without exposing local filesystem paths.

Create reusable likeness bundles only with explicit consent:

```bash
kie-pp-cli media identity create "Creator" --reference ./front.jpg --reference ./profile.jpg --consent --agent
kie-pp-cli create "Creator presenting the product" --identity identity_<id> --agent
```

This is a local reference bundle, not a trained or cross-model biometric identity.

## Script and Storyboard Protocol

For a narrative, explainer, ad, tutorial, or any multi-shot request, start the
master video brief with `--production-mode storyboard`. Read
[references/storyboard.md](references/storyboard.md), then follow its local
script review, local storyboard review, and per-shot generation gates. Do not
submit the master brief through `media_generate` or `create --submit`; use each
returned `shot_brief_id`.

## Generation Gate

Do not start live generation until the user approves the ready plan. Video has an additional hard visual-anchor gate enforced by the shared service, not merely by these instructions.

For a qualified video brief, `next_action` is `generate_preview` and `can_submit` is false. Generate the separate review image:

```bash
kie-pp-cli create --brief <brief_id> --preview --wait --agent
```

This preview is a paid/live image generation and may consume Kie credits. Read `result_urls[0]` and actually show the image to the user. In Codex, Claude, or another visual host, render the image; in a plain terminal, print the direct URL clearly. Never infer approval from silence, the prompt, metadata, or your own judgment.

If the user rejects it:

```bash
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --style "<revised direction>" --agent
kie-pp-cli create --brief <brief_id> --preview --wait --agent
```

Only after an explicit affirmative response, record approval:

```bash
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

The returned turn must have `can_submit: true`. The approved still becomes the SeedDance first frame or multimodal visual anchor. Any subsequent creative change makes the fingerprinted approval stale and requires a new preview.

In terminal/agent flow, explicit approval is:

```bash
kie-pp-cli create --brief <brief_id> --submit --agent
```

Use this when the user has approved the plan. Add `--wait` only when the user wants the command to poll:

```bash
kie-pp-cli create --brief <brief_id> --submit --wait --agent
```

For MCP video, call `media_preview_generate`, poll it with `media_generation_status`, display the returned image URL, obtain explicit approval, then call `media_preview_approve`. Use `media_preview_reject` for another direction. Only then call `media_generate`. Preview and final generation are separate live actions; brief/reference/identity and preview approve/reject tools are local state actions.

Register the focused local MCP server with an agent host using stdio:

```bash
claude mcp add kie-media -- kie-media-mcp
```

`kie-media-mcp` exposes twenty-one focused workflow, brief, script, storyboard,
reference, identity, preview, generation, and status tools. Its optional HTTP
mode is local-only, stateless, and negotiates MCP `2026-07-28`:

```bash
kie-media-mcp --transport http --addr 127.0.0.1:7780
```

Never expose that listener on a wildcard or non-loopback address.

Cost estimation is off by default. Do not invent prices or token/credit estimates. If the user asks about cost, inspect current Kie docs/API context first and state uncertainty clearly.

## Create, Poll, and Report

Start:

```bash
kie-pp-cli create "a cinematic product image for a landing page" --agent
```

Answer:

```bash
kie-pp-cli create --brief <brief_id> --answer <value> --agent
```

Review:

```bash
kie-pp-cli media brief show <brief_id> --agent
```

For video, preview and approve:

```bash
kie-pp-cli create --brief <brief_id> --preview --wait --agent
# Show result_urls[0] and wait for an explicit yes.
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

Submit:

```bash
kie-pp-cli create --brief <brief_id> --submit --agent
```

Poll:

```bash
kie-pp-cli media generation status <generation_id> --agent
kie-pp-cli media generation status <generation_id> --wait --agent
```

Report the final or current state with:

- `brief_id`
- `generation_id`
- provider `task_id`
- selected model
- status/job state
- result URL(s)
- any warning or failure reason
- exact next command if the generation is still running

## Defaults

Implemented defaults:

- Image without references: `gpt-image-2-text-to-image`
- Image with references: `gpt-image-2-image-to-image`
- Video: `bytedance/seedance-2-5`
- Video duration recommendation: `5` seconds
- Video resolution defaults to `720p`, audio defaults off, and output defaults to MP4.
- SeedDance supports text, first-frame, first+last-frame, and multimodal image/video/audio reference modes.
- Aspect-ratio recommendations: Instagram image `3:4`, Instagram/TikTok video `9:16`, YouTube/website/LinkedIn `16:9`, general video `16:9`, general image `1:1`

Use `--model` only when the user asks for a specific model or current docs/help indicate the default route is unsuitable.

## Sources

- Higgsfield CLI: https://higgsfield.ai/cli
- Higgsfield agent skills: https://github.com/higgsfield-ai/skills
- Matt Pocock Grill With Docs: https://github.com/mattpocock/skills/tree/main/skills/engineering/grill-with-docs
- Kie API docs: https://docs.kie.ai/
- Kie File Upload API quickstart: https://docs.kie.ai/file-upload-api/quickstart
