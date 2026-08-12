---
name: kie-video-explainer
description: Plan and produce narrated explainer videos using Kie.ai speech and video jobs plus transparent local assembly. Use for educational explainers, product walkthroughs, faceless or presenter-led videos, block scripts, voice selection, subtitles, recurring visual style, and resumable multi-job production.
---

# Kie Video Explainer

Kie does not document Higgsfield's `explainer_video` assembly endpoint. Generate narration and visual blocks through Kie, then assemble and caption them locally with ffmpeg or the user's editor. Do not imply that Kie performed the final assembly.

Prefer one storyboard master created with `kie-pp-cli create --workflow
video-explainer --production-mode storyboard ... --agent`. Save and obtain
explicit approval for the script and storyboard before using the returned
`shot_brief_id` values. MCP hosts use the same workflow through
`media_workflow_get`, `media_brief_start`, `media_script_*`, and
`media_storyboard_*`.

Read [references/production-contract.md](references/production-contract.md) before scheduling jobs.

## Qualify the production

Ask one question at a time:

1. topic, audience, factual promise, and call to action;
2. target length from one to ten minutes;
3. language, voice qualities, pronunciation, and pace;
4. faceless, mascot, consented presenter, or mixed format;
5. aspect ratio, subtitle treatment, and brand style;
6. factual sources and claims that require verification.

Research factual topics from authoritative sources before scripting. Keep citations in the production manifest.

Default to 16:9 and enabled captions when the user has no preference. Require an explicit voice choice before narration. Support every material factual claim with an authoritative source; use two independent authoritative sources when the claim is consequential or contested.

## Build ten-second blocks

Write exactly six blocks per minute. Each block contains:

- one narration segment designed for about ten seconds;
- one English visual prompt matching that narration;
- any on-screen exact copy to add locally;
- source citations and transition notes.

Create one style key image first when continuity matters:

```bash
kie-pp-cli create --workflow video-explainer "A reusable non-photoreal style key: <direction>" \
  --type image --purpose "explainer style key" --aspect-ratio 16:9 --agent
```

This creates local brief state only until `--submit` is explicitly used. Generate all narration before video so timing controls visual edits. Use `elevenlabs/text-to-speech-turbo-2-5` or another current Kie voice route after inspecting its live schema:

```bash
kie-pp-cli kie-ai-jobs market-create-task \
  --model elevenlabs/text-to-speech-turbo-2-5 --input '<schema-verified narration JSON>' --agent
```

Save task IDs and output URLs in the manifest. Do not guess voice or tuning fields when the current schema is unavailable.

Save the complete narration as the master script, approve it, then save a
storyboard containing one ten-second shot per block. The CLI creates one child
brief per shot. Generate one SeedDance 2.5 clip per returned child brief:

```bash
kie-pp-cli create --workflow video-explainer "<complete explainer request>" --type video \
  --production-mode storyboard --duration <total-seconds> --audio off \
  --aspect-ratio <ratio> --reference <style-key> --agent
kie-pp-cli media script set <brief_id> --file script.md --agent
kie-pp-cli media script approve <brief_id> --agent
kie-pp-cli media storyboard set <brief_id> --file storyboard.json --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

Each block is independently gated. Generate its still with `create --brief <brief_id> --preview --wait --agent`, show the image, and record explicit approval with `--approve-preview` before submitting that block's video. Reject and revise a bad composition with `--reject-preview`; do not spend the final video job hoping motion will repair a visibly wrong anchor.

Use `omnihuman-1-5` or a Kling avatar route only when the user requests a presenter and supplies an authorized portrait plus narration audio. Do not promise a perfect lip-sync or identity match before inspection.

## Assemble locally

Normalize each audio/video pair to the agreed frame rate, resolution, and exactly ten seconds. Center narration, preserve pitch, and use silence or room tone rather than time-stretching beyond intelligibility. Concatenate in manifest order. Burn or attach captions from measured narration timings, not guessed line lengths.

Keep assembly commands and tool versions in the manifest. Never interpolate shell paths from untrusted text; use explicit, quoted local paths.

## Gate and deliver

Review all ready Kie plans and preview stills before paid final-video submission. Preview generation is also a live paid step. Retry only failed or explicitly rejected blocks, preserve lineage, and reassemble after replacements. Deliver the final MP4, optional caption file, clean block outputs, approved anchors, style key, voice/model IDs, task IDs, factual sources, local assembly command/version, and known timing or identity caveats.
