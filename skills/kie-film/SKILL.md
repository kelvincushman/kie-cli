---
name: kie-film
description: Develop and produce a multi-shot Kie.ai film, advert, explainer, tutorial, or story through an approved script, storyboard, continuity bible, per-shot briefs, visible still gates, optional low-resolution complete-shot proofs, and fresh paid confirmation for every generation.
---

# Kie Film

Start with `$kie-grilling`. Use one master brief with `--production-mode storyboard`; do not flatten a multi-shot story into one video prompt.

## Pre-production

Agree the purpose, audience, duration, format, story promise, tone, call to action, required facts, brand constraints, characters, locations, props, and available references. Use `$kie-identity` for consented recurring people.

Write a local script with beats, dialogue/narration, visual action, and approximate timing:

```bash
kie-pp-cli media script set <brief_id> --file <script.md|-> --agent
kie-pp-cli media script show <brief_id> --agent
kie-pp-cli media script approve <brief_id> --agent
```

Do not approve it for the user. After explicit approval, create storyboard JSON with shot number, duration, framing, lens/camera, action, lighting, environment, character/wardrobe/prop continuity, dialogue/audio, transition, reference handles, and prompt. Then:

```bash
kie-pp-cli media storyboard set <brief_id> --file <storyboard.json|-> --agent
kie-pp-cli media storyboard show <brief_id> --agent
kie-pp-cli media storyboard approve <brief_id> --agent
```

The storyboard returns one durable child brief per shot. Maintain a continuity bible for appearance, wardrobe, screen direction, geography, lighting, props, palette, and audio. An upstream change marks dependent work stale.

## Produce shot by shot

For each child brief, route through `$kie-video`:

1. Ask for a fresh paid confirmation and generate the mandatory still.
2. Show it and record the user's approval/rejection.
3. Offer the lowest-faithful-resolution complete-shot proof; if accepted, ask for a fresh paid confirmation, generate it, show the whole clip, and record approval/rejection. If declined, record skip.
4. Ask for a new paid confirmation for the final shot.
5. Poll, inspect, and record provenance before moving on.

Review complete shots for composition, motion, identity, continuity, timing, dialogue sync, audio, artifacts, and transition compatibility. Do not approve on the user's behalf. Assemble only approved outputs, and deliver the script, storyboard, continuity bible, shot ledger, generation/task IDs, model/settings, result URLs, and remaining edit/audio work.
