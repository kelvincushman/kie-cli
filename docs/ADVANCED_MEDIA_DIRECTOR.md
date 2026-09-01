# Advanced Agent Media Factory

This document is for local users who want to operate `kie-cli` as a
terminal-first media factory, focused MCP server, skill runtime, or automation
target. The repository is `kie-cli`; the installed command is `kie-pp-cli`.

## Current Status

Implemented and verified in the current local command surface:

- `kie-pp-cli create` for guided image and video briefs.
- `kie-pp-cli media` for setup, workflows, briefs, references, identities, and
  generation status.
- Prompt inference, concise grilling, visible route rationale/overrides, video
  still approval, optional complete-shot proof, and final submission blocking.
- `kie-media-mcp` as a focused local MCP server.
- Local reference vault and consented identity-reference bundles.
- Nine routed production workflows delivered through 17 Markdown skills,
  including shared grill and capability-level routes.
- A 16-course, 171-lesson public source map with original Kie-native methods,
  exposed through `/kie-lesson`, terminal commands, and MCP.
- A dated task-specific model evidence ledger refreshed weekly without combining
  unlike benchmark scores.
- Seventy current API operations and a 129-model Market snapshot with complete
  embedded request/input schemas and local validation.
- A fail-closed five-capability classifier over all 129 models, compact
  `media capability`/`media models` views, and a direct `media video` shortcut
  backed by the same canonical registry.
- Local SQLite recall, teach, and playbook state for repeated agent work.

Native storyboard command surface:

- `kie-pp-cli create --production-mode storyboard`
- `kie-pp-cli media script set/show/approve/reject`
- `kie-pp-cli media storyboard set/show/approve/reject`
- MCP tools `media_script_set`, `media_script_get`, `media_script_decide`,
  `media_storyboard_set`, `media_storyboard_get`, and
  `media_storyboard_decide`

Each saved storyboard shot receives a normal `shot_brief_id` and therefore
inherits the existing preview approval gate. The storyboard master cannot be
submitted as one prompt-only video.

## Operating Model

The factory has four layers:

1. **Direction:** a user or agent qualifies intent one question at a time.
2. **Durable local production:** briefs, references, identity bundles, scripts,
   storyboards, approvals, and generation records receive opaque handles.
3. **Explicit live actions:** preview, optional proof, and final jobs are
   separate Kie.ai calls with fresh scoped paid confirmations.
4. **Local finishing:** clip order, exact text, captions, audio mix, publishing,
   and deployment remain inspectable local steps.

The focused MCP transport is stateless, but the application is not ephemeral.
Every tool call carries a durable ID, so any compatible agent can resume a job
without MCP session memory or a long prompt containing the entire production.

This is an open-source, bring-your-own-Kie-key alternative to closed agent media
suites. It does not provide trained cross-model Soul identities, proprietary
virality scoring, 3D/mesh generation, marketplace compliance, or a hosted
all-in-one editor. Identity records are consented local reference bundles, and
unsupported gaps returned by `media workflow show` must remain visible to the
agent.

## Local State

Discover the resolved data directory instead of hardcoding paths:

```bash
kie-pp-cli media setup --agent
kie-pp-cli agent-context --pretty
```

The media director uses opaque local handles:

- `brief_...` for durable media briefs.
- `ref_...` for vaulted references.
- `ref:<id>` when passing a vaulted reference back to `create`.
- `identity_...` for consented local likeness-reference bundles.
- `gen_...` for local generation records.
- `script_...` for the current master script.
- `storyboard_...` for the current ordered shot plan.
- `shot_brief_id` values are ordinary `brief_...` handles linked back to the
  storyboard master.

Local state contains brief answers, selected model, provider input, reference
handles, provider task IDs, generation status, result URLs, preview/proof
approval state, paid-confirmation fingerprints, script/storyboard hashes, and
shot lineage. It must not contain bearer
tokens, cookies, auth headers, or credential file contents.

Use `--home <directory>` to isolate projects or test runs. Do not point two
untrusted users at the same home directory: the local store is application
state, not a multi-tenant authorization boundary.

## CLI Contract

Use `--agent` for agent and automation calls. It enables JSON, compact output,
no color, and no input prompts. It does not authorize paid actions;
`--confirm-paid` is still required for the exact live call.

Core commands:

```bash
kie-pp-cli media setup --agent
kie-pp-cli media workflow list --agent
kie-pp-cli media workflow show <workflow> --agent
kie-pp-cli lesson recommend "<production outcome>" --agent
kie-pp-cli media leaderboard character-consistency --agent
kie-pp-cli grill-me "<request>" --agent
kie-pp-cli create "<request>" --agent
kie-pp-cli create --workflow <workflow> "<request>" --agent
kie-pp-cli create --brief <brief_id> --answer "<answer>" --agent
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media brief list --agent
```

Reference and identity commands:

```bash
kie-pp-cli media reference add <path-or-url> --name <name> --agent
kie-pp-cli media reference list --agent
kie-pp-cli media identity create <name> --reference <image> --consent --agent
kie-pp-cli media identity show <identity_id> --agent
kie-pp-cli media identity list --agent
```

Generation commands:

```bash
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
kie-pp-cli create --brief <brief_id> --approve-preview --agent
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
kie-pp-cli create --brief <brief_id> --approve-proof --agent
kie-pp-cli create --brief <brief_id> --reject-proof --agent
kie-pp-cli create --brief <brief_id> --skip-proof --agent
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
kie-pp-cli media generation status <generation_id> --wait --agent
```

Lesson and evidence commands:

```bash
kie-pp-cli lesson list --agent
kie-pp-cli lesson recommend "one character across several worlds" --agent
kie-pp-cli lesson show <course-slug/lesson-slug> --agent
kie-pp-cli lesson start "<request>" --lesson <course/lesson> --agent
kie-pp-cli media leaderboard <task> --agent
```

`/kie-lesson` is available where the host turns installed skills into slash
commands. CLI and MCP remain the portable contracts. The Academy artifact
stores public discovery metadata and source links plus original Kie methods;
it is not a mirror of lesson prose, scripts, prompts, or media.

The advanced direct path omits the director preview gate:

```bash
kie-pp-cli media video "<prompt>" --duration 5 --ratio 16:9 --wait --agent
```

Before constructing custom direct calls, use the local registry rather than
putting whole schemas into an agent prompt:

```bash
kie-pp-cli media models --family video --agent
kie-pp-cli media capability list --capability kie-video --agent
kie-pp-cli media capability show bytedance/seedance-2-5 --agent
kie-pp-cli models show wan/2-6-text-to-video --agent
kie-pp-cli media video --model wan/2-6-text-to-video \
  --input '{"prompt":"<prompt>","duration":"5"}' --agent
```

`media models` is a compact view of `internal/kiecatalog`, not an independently
generated catalog. `media video` validates all custom input fields and settings
against that registry before submission. Its Wan 2.7 flags cover prompt,
negative prompt, audio URL, resolution, ratio, duration, prompt extension,
watermark, seed, and NSFW checking. With `--wait`, it polls the shared Market
task endpoint and returns discovered result URLs. Because this direct command
does not create a brief or preview approval record, use it only when that
explicitly advanced automation contract is appropriate.

## MCP Contract

Build and register the focused MCP server:

```bash
make build-media-mcp
claude mcp add kie-media -- ./bin/kie-media-mcp
```

Stdio is the recommended local host transport. The optional HTTP transport is
for local-only stateless operation:

```bash
kie-media-mcp --transport http --addr 127.0.0.1:7780
```

Do not bind HTTP MCP to a wildcard or non-loopback address. The HTTP transport
implements the finalized MCP `2026-07-28` stateless protocol. Stdio remains the
recommended local child-process transport.

The implemented focused server exposes 40 media tools:

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

The MCP server is stateless at the protocol layer. Durable IDs carry workflow
state between calls, so agents should store and pass `brief_id`, `generation_id`,
and storyboard shot IDs explicitly. `kie-pp-mcp` is a different binary: it
retains the broad generated API tool surface and is not the preferred creative
director.

## Video Approval Gate

Preview generation, optional complete-shot proof, and final video generation are
separate live actions. Each may consume Kie.ai credits and each requires a new
paid confirmation.

Required sequence:

```bash
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
# Show result_urls[0] to the user.
kie-pp-cli create --brief <brief_id> --approve-preview --agent
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
# Show the whole proof. Approve/reject it, or use --skip-proof if declined.
kie-pp-cli create --brief <brief_id> --approve-proof --agent
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

If the preview is wrong:

```bash
kie-pp-cli create --brief <brief_id> --reject-preview --agent
kie-pp-cli create --brief <brief_id> --style "<revised direction>" --agent
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
```

The approved still becomes the SeedDance first frame or multimodal visual
anchor. Changing creative fields invalidates approval and requires a new
preview.

## Script And Storyboard Workflow

Use storyboard mode for multi-shot work:

```bash
kie-pp-cli create --production-mode storyboard "<video request>" --agent
kie-pp-cli media script set <brief_id> --file script.md --agent
kie-pp-cli media script show <brief_id> --agent
kie-pp-cli media script approve <brief_id> --agent
kie-pp-cli media storyboard set <brief_id> --file storyboard.json --agent
kie-pp-cli media storyboard show <brief_id> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

Each storyboard shot carries a `shot_brief_id`. Generate that shot through the
same video gate:

```bash
kie-pp-cli create --brief <shot_brief_id> --preview --confirm-paid --wait --agent
kie-pp-cli create --brief <shot_brief_id> --approve-preview --agent
kie-pp-cli create --brief <shot_brief_id> --proof --confirm-paid --wait --agent
kie-pp-cli create --brief <shot_brief_id> --approve-proof --agent
kie-pp-cli create --brief <shot_brief_id> --submit --confirm-paid --wait --agent
```

The CLI does not claim to assemble a final edited film by itself. Treat local
assembly with `ffmpeg`, Remotion, or an editor as a separate local step. Record
the command, source clips, approved preview URLs, and output file in a manifest.

## Automation Boundaries

Safe to automate:

- setup checks
- workflow discovery
- lesson discovery and task-specific evidence lookup
- brief inspection
- reference listing
- status polling
- local script/storyboard file validation
- local assembly after all inputs are approved
- compact handle passing and local teach/recall

Require explicit user approval:

- storing or changing credentials
- live preview generation
- live complete-shot proof generation
- live final generation
- approving a preview still
- approving a script or storyboard
- publishing or deploying generated media

Do not chain paid calls. The user must see the preview image and any generated
proof first, and every subsequent call needs a fresh confirmation.

## Security

- Run `kie-pp-cli auth setup` in an interactive terminal, or set `KIE_BEARER_AUTH` through an environment secret store.
- Never store credentials in briefs, scripts, storyboards, manifests, issue
  comments, or MCP messages.
- Local reference files remain local until explicit live generation.
- The reference vault rejects symbolic links and omits private filesystem paths
  from public CLI/MCP responses.
- Identity bundles are consented local reference packs, not trained biometric
  models.
- Keep MCP HTTP bound to `127.0.0.1`.
- Preserve unsupported workflow gaps in agent output; do not silently imply
  Higgsfield or provider feature parity.

## Troubleshooting

Auth missing:

```bash
kie-pp-cli auth status --agent
kie-pp-cli auth setup
kie-pp-cli doctor --json
```

Preview is not ready:

```bash
kie-pp-cli media generation status <generation_id> --wait --agent
```

Final video is blocked:

- Confirm the preview generation has a result URL.
- Show the image to the user.
- Run `kie-pp-cli create --brief <brief_id> --approve-preview --agent`.
- Generate/show/decide the offered proof, or record `--skip-proof`.
- Re-run submit after approval.

Preview approval is stale:

- A creative field changed after approval.
- Reject or regenerate the preview.
- Approve the new preview before final generation.

MCP host cannot see tools:

- Rebuild `kie-media-mcp`.
- Register with stdio first.
- Restart the MCP host.
- Use loopback HTTP only when the host supports stateless HTTP MCP.

Storyboard commands missing after an upgrade or custom build:

- Verify with `kie-pp-cli media script --help` and
  `kie-pp-cli media storyboard --help`.
- Rebuild both binaries from the same checkout. CLI, docs, and MCP binaries
  must come from the same revision.

Model behavior differs from the documented plan:

- Run `kie-pp-cli media brief show <brief_id> --agent` and inspect `plan.model`
  plus `plan.input`.
- Run `kie-pp-cli models show <model-id> --agent`, then validate the proposed
  input with `kie-pp-cli models validate <model-id> --input '<json>' --agent`.
- Check [MODELS.md](MODELS.md), [MODEL_INPUTS.md](MODEL_INPUTS.md), and the
  recorded current upstream Kie model page.
- Use `--model` only when the requested model supports the required reference,
  duration, audio, and output constraints.

Evidence looks stale or incomplete:

- Run `kie-pp-cli media leaderboard --agent` and inspect `generated_at`, each
  source URL, direct-versus-proxy status, and current Kie availability.
- Run `scripts/weekly-refresh.sh` to rebuild Kie contracts, verify every mapped
  Academy lesson URL, and refresh the external task tables.
- A missing independent score stays null. In particular, never transfer a
  Seedance 2.0 family score to Seedance 2.5.
