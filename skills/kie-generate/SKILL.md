---
name: kie-generate
description: Create, edit, or direct images, video, narration, dialogue, sound, and music through Kie.ai from a local terminal or MCP host. Use for broad media-generation requests, model routing, reference-driven work, SeedDance 2.5 video, job polling, or honest capability checks for 3D and virality analysis.
---

# Kie Generate

For any multi-shot or narrative video, hand off to the `kie-create` storyboard
protocol: one local master script, one local storyboard, and one independently
previewed/approved child brief per shot. Do not flatten a multi-shot request
into one prompt-only video job.

Use the local CLI or `kie-media-mcp`; keep credentials and source media on the user's machine until an approved live call.

Prefer the compact CLI for repeated steps and token savings: inspect `kie-pp-cli media workflow show generate --agent`, then start with `kie-pp-cli create --workflow generate "<request>" --agent`. MCP hosts use `media_workflow_get` followed by `media_brief_start` with `workflow: "generate"`.

Read [references/capabilities.md](references/capabilities.md) before selecting a model. Inspect current runtime truth when model details matter:

```bash
kie-pp-cli doctor --json
kie-pp-cli agent-context --pretty
kie-pp-cli create --help
kie-pp-cli media --help
```

## Direct the request

Ask one concise question at a time. Establish:

1. the intended outcome and subject;
2. media type and destination;
3. aspect ratio, duration, and audio needs;
4. style, camera, lighting, and motion;
5. truthful copy or factual constraints;
6. image, video, audio, brand, product, or identity references;
7. required variants and delivery format.

Do not batch these as a form. Stop asking once the plan is unambiguous.

For image or video, start the durable director:

```bash
kie-pp-cli create "<request>" --agent
kie-pp-cli create --brief <brief_id> --answer <answer> --agent
kie-pp-cli media brief show <brief_id> --agent
```

Use direct flags when the qualified answer is already known. For SeedDance 2.5, choose exactly one input mode:

```bash
# text-to-video
kie-pp-cli create "<request>" --type video --video-mode text --duration 10 --audio off --agent

# first and last frames
kie-pp-cli create "<request>" --type video --video-mode first-last \
  --first-frame <path-or-url> --last-frame <path-or-url> --duration 10 --audio on --agent

# multimodal references
kie-pp-cli create "<request>" --type video --video-mode multimodal \
  --reference <image> --reference-video <video> --reference-audio <audio> --agent
```

Vault reusable media first when useful:

```bash
kie-pp-cli media reference add <path-or-url> --type image --name "<name>" --agent
kie-pp-cli media reference list --agent
```

Use `identity:<id>` only for a consented local bundle created through `$kie-identity`.

## Route non-image media

Use the Market task surface for narration or dialogue after inspecting the current model page/input schema:

```bash
kie-pp-cli kie-ai-jobs market-create-task --model <model-id> --input '<json>' --agent
kie-pp-cli kie-ai-jobs market-query-task --task-id <task-id> --agent
```

Use the dedicated Suno commands for music and sound when their current help matches the request. Never guess input JSON.

Do not claim Kie can produce a true GLB/OBJ/FBX mesh: no official Kie 3D mesh route is documented. Offer concept renders, turnarounds, texture references, or an explicit handoff to a separate 3D tool.

Treat virality scoring as a qualitative local critique only. Do not present a heuristic review as Higgsfield's proprietary Virality Predictor or as a measured prediction.

## Approve and deliver

Review the ready plan with the user. For video, generate the mandatory review still, show `result_urls[0]`, and wait for an explicit yes:

```bash
kie-pp-cli create --brief <brief_id> --preview --wait --agent
# Render or clearly display result_urls[0].
kie-pp-cli create --brief <brief_id> --approve-preview --agent
```

If rejected, use `--reject-preview`, revise the direction, and generate another still. Do not approve on the user's behalf. The service blocks final video generation until the current preview is approved. Preview and final generation may each consume credits.

Submit only after the relevant approval:

```bash
kie-pp-cli create --brief <brief_id> --submit --agent
kie-pp-cli media generation status <generation_id> --wait --agent
```

Report the brief ID, generation ID, provider task ID, selected model, status, result URLs, warnings, and exact resume command. Never print credentials or claim that an output capability was validated when only the schema was inspected.
