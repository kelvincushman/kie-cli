# Script And Storyboard Workflow

This guide explains how to use the `kie-pp-cli` agent media factory for videos
that need scripting, shot-by-shot direction, preview approval, and local
assembly.

## Rule Of Thumb

Use a normal single-shot `create` brief for one continuous clip. Use storyboard
mode for explainers, ads, UGC sequences, product walkthroughs, tutorials, or
anything whose shots must be planned and approved independently.

Every video shot goes through the same human and paid-action gates:

```text
master brief -> approved script -> approved storyboard -> shot brief
shot brief -> paid still -> user approval -> optional paid complete-shot proof -> fresh-confirmed final
```

## Workflow

1. Create a storyboard-mode master brief with the complete production intent.
2. Write, save, show, and explicitly approve `script.md`.
3. Write, save, show, and explicitly approve `storyboard.json`.
4. Read the returned `shot_brief_id` for each ordered shot.
5. Generate and show the preview still for each shot.
6. Approve or reject each displayed shot preview.
7. Offer a complete-shot proof at each model's lowest faithful resolution.
8. Show/approve or reject the proof, or record an explicit skip.
9. Obtain a new paid confirmation for every final shot, then assemble locally.

Start the master production:

```bash
kie-pp-cli create --workflow video-explainer \
  "A 30-second product explainer for founders, ending with a direct call to action" \
  --type video --platform YouTube --aspect-ratio 16:9 --duration 30 \
  --production-mode storyboard --agent

kie-pp-cli media script set <brief_id> --file script.md --agent
kie-pp-cli media script show <brief_id> --agent
kie-pp-cli media script approve <brief_id> --agent

kie-pp-cli media storyboard set <brief_id> --file storyboard.json --agent
kie-pp-cli media storyboard show <brief_id> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

Then use the preview gate for every returned `shot_brief_id`:

```bash
kie-pp-cli create --brief <shot_brief_id> --preview --confirm-paid --wait --agent
# Show result_urls[0] to the user.
kie-pp-cli create --brief <shot_brief_id> --approve-preview --agent
kie-pp-cli create --brief <shot_brief_id> --proof --confirm-paid --wait --agent
# Show the entire proof; approve/reject it, or record --skip-proof if declined.
kie-pp-cli create --brief <shot_brief_id> --approve-proof --agent
kie-pp-cli create --brief <shot_brief_id> --submit --confirm-paid --wait --agent
```

Reject and revise a bad shot:

```bash
kie-pp-cli create --brief <shot_brief_id> --reject-preview --agent
kie-pp-cli create --brief <shot_brief_id> --style "tighter crop, more product detail, less background" --agent
kie-pp-cli create --brief <shot_brief_id> --preview --confirm-paid --wait --agent
```

## Native Command Surface

The current CLI exposes this command surface:

```bash
kie-pp-cli create --production-mode storyboard "<video request>" --agent
kie-pp-cli media script set <brief_id> --file script.md --agent
kie-pp-cli media script show <brief_id> --agent
kie-pp-cli media script approve <brief_id> --agent
kie-pp-cli media script reject <brief_id> --agent
kie-pp-cli media storyboard set <brief_id> --file storyboard.json --agent
kie-pp-cli media storyboard show <brief_id> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
kie-pp-cli media storyboard reject <brief_id> --agent
```

Inspect the installed surface with:

```bash
kie-pp-cli create --help
kie-pp-cli media script --help
kie-pp-cli media storyboard --help
```

The master brief is a production container and cannot be submitted as one
prompt-only video. Each storyboard shot includes a normal `shot_brief_id`, which
uses the ordinary video command surface:

```bash
kie-pp-cli create --brief <shot_brief_id> --preview --confirm-paid --wait --agent
kie-pp-cli create --brief <shot_brief_id> --approve-preview --agent
kie-pp-cli create --brief <shot_brief_id> --proof --confirm-paid --wait --agent
kie-pp-cli create --brief <shot_brief_id> --approve-proof --agent
kie-pp-cli create --brief <shot_brief_id> --submit --confirm-paid --wait --agent
```

The aggregate master actions progress through `draft_script`, `review_script`,
`draft_storyboard`, `review_storyboard`, and `generate_shot_previews`. Rejected
or stale artifacts return revision actions before any child generation can
continue.

## Script Format

Use plain Markdown for scripts:

```markdown
# Product Launch Video

Duration: 30 seconds
Audience: founders and growth teams
Platform: YouTube Shorts and website hero

## Voiceover

1. Most product demos waste the first five seconds.
2. This one opens with the finished outcome.
3. Then it shows the three steps that got us there.

## Notes

- Keep the tone direct and practical.
- No unsupported claims.
- No tiny text inside generated video frames.
```

Changing an approved script should mark the storyboard and downstream shots
stale. Regenerate or re-approve affected shots before spending final video jobs.

## Storyboard Format

Use JSON when an agent or script will process the storyboard:

```json
{
  "title": "Product Launch Video",
  "shots": [
    {
      "id": "shot_01",
      "duration_seconds": 5,
      "narration": "Most product demos waste the first five seconds.",
      "visual": "Finished product outcome on a laptop screen, shallow depth of field.",
      "camera": "Slow push-in from medium-wide to close-up.",
      "references": ["ref:ref_product", "ref:ref_logo"]
    }
  ]
}
```

Recommended shot fields:

- `id`
- `duration_seconds`
- `narration`
- `visual`
- `camera`
- `references`
- `title`
- `dialogue`
- `transition`

`shot_brief_id` is output by the CLI; do not put it in input JSON.

Validation rules:

- A storyboard contains 1-60 ordered shots.
- Every shot has a non-empty `visual` and lasts 4-30 seconds.
- Total duration cannot exceed 600 seconds.
- When the master has a positive duration, shot durations must total it exactly.
- A shot accepts at most 29 explicit image references because the approved
  preview occupies the thirtieth SeedDance image slot.
- Script, storyboard, and creative-brief hashes bind approvals to the exact
  reviewed content. Editing an approved artifact makes downstream state stale.

## MCP Storyboard Tools

The focused MCP server exposes 40 tools, including six local
script/storyboard tools:

- `media_script_set`
- `media_script_get`
- `media_script_decide`
- `media_storyboard_set`
- `media_storyboard_get`
- `media_storyboard_decide`

Expected MCP flow:

```text
media_grill_start
media_script_set
media_script_decide approve
media_storyboard_set
media_storyboard_decide approve
for each shot_brief_id:
  media_paid_confirm for preview
  media_preview_generate
  media_generation_status
  display preview image to user
  media_preview_approve or media_preview_reject
  offer a complete-shot proof
  if accepted: media_paid_confirm, media_proof_generate, display full proof,
    then media_proof_approve or media_proof_reject
  otherwise: media_proof_skip
  media_paid_confirm for final
  media_generate only after the current gates
  media_generation_status
```

Use `media_preview_reject` and revise the shot if the still is wrong.

## Local Assembly Boundary

Kie generates assets and clips. The local workflow owns final assembly.

Use `ffmpeg`, Remotion, or a video editor for:

- clip ordering
- trimming
- transitions
- captions
- exact text overlays
- audio mix
- final export
- publishing or deployment

Keep an assembly manifest with source clip URLs, local filenames, approved
preview URLs, task IDs, edit decisions, command versions, and known caveats.

Example manifest fields:

```json
{
  "project": "Product Launch Video",
  "script_status": "approved",
  "storyboard_status": "approved",
  "shots": [],
  "assembly": {
    "tool": "ffmpeg",
    "command": "",
    "output": ""
  }
}
```

## Approval Rules

- Do not approve a script unless the user accepts the message, pacing, claims,
  and audience fit.
- Do not approve a storyboard unless the user accepts shot order, framing,
  identity/reference use, and continuity.
- Do not approve a preview still unless the user has actually seen it.
- Do not approve a complete-shot proof unless the user has seen the whole clip.
- Do not submit final video for any shot whose current still/proof decision is
  incomplete.
- Still and proof approval are not paid confirmation for the final render.
- Obtain a fresh confirmation for every paid still, proof, final, retry, or edit.
- Do not assume motion will fix a bad still.
- Do not submit the storyboard master through `create --submit` or
  `media_generate`; only its child shot briefs are generation targets.
