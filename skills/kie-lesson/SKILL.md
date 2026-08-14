---
name: kie-lesson
description: Select a source-linked public media-production lesson, adapt its method into an original Kie.ai workflow, and guide the user through brief, script, reusable character/assets, storyboard, shot prompts, still approvals, video generation, continuity review, and local assembly. Use when the user invokes /kie-lesson or $kie-lesson, asks what production method to follow, wants a cinematic ad/film/story workflow, or needs character-consistent image-to-video direction through kie-pp-cli or kie-media-mcp.
---

# Kie Lesson Director

Start with `$kie-grilling`, then use the selected lesson as a production pattern,
not as a replacement for the user's creative decisions. Route multi-shot work to
`$kie-film` and each child shot to `$kie-video`.

Use the repository's original Kie-native method map. Do not copy a source
lesson's prose, script, prompt text, video, or downloadable resources. Public
lesson titles and links are discovery metadata; the `kie_method` returned by the
CLI is the reusable production instruction.

In hosts that expose skills as slash commands, treat `/kie-lesson` as the human
invocation. The portable command contract is `kie-pp-cli lesson`.

## Discover the Runtime

```bash
kie-pp-cli lesson --help
kie-pp-cli lesson list --agent
kie-pp-cli lesson recommend "what the user wants to create" --agent
kie-pp-cli media leaderboard --agent
kie-pp-cli media setup --agent
```

With MCP, use `media_lesson_recommend`, `media_lesson_get`,
`media_leaderboard_get`, and `media_brief_start`.

## Select, Explain, Then Start

1. Ask what the user wants the audience to see, feel, or do.
2. Run `lesson recommend` and show at most three results with source links,
   production stage, original Kie method, and capability boundary.
3. Recommend one result and explain why. Let the user override it.
4. Start the durable workflow:

```bash
kie-pp-cli lesson start "the request" \
  --lesson <course-slug/lesson-slug> --agent
```

The command starts workflow `academy`, selects video plus storyboard mode, and
returns exactly one next question. Ask only that question, then resume with
`kie-pp-cli create --brief <id> --answer <value> --agent`.

For an image-only outcome, start directly with:

```bash
kie-pp-cli create --workflow academy --lesson <lesson-key> \
  --type image "the request" --agent
```

## Direct the Production

For narrative or multi-shot media, require this order:

1. Qualify the creative brief one question at a time.
2. Write a producible script with one visible action per beat.
3. Lock reusable character, product, prop, wardrobe, and location references.
4. Save and obtain explicit approval for the script.
5. Build a timed storyboard with camera, action, audio, transition, continuity,
   and approval criteria for every shot.
6. Save and obtain explicit approval for the storyboard.
7. Generate a still keyframe for every returned `shot_brief_id`.
8. Display the still and record approval only after the user says yes.
9. Generate motion from the approved frame, inspect the result, and reject
   broken continuity instead of rationalizing it.
10. Assemble selected clips locally; keep captions, sound mix, exact copy,
    export, publishing, and measurement as explicit downstream steps.

Read [references/production-contract.md](references/production-contract.md)
before drafting the script, storyboard, or per-shot prompts.

## Preserve Character Consistency

Obtain consent before using a real person's likeness. Create one master
character sheet with stable face, body, hair, age, wardrobe anchors, palette,
and front/three-quarter/profile views. Save approved images as local references
or a consented identity bundle.

Use the dated leaderboard as evidence, not doctrine:

- Prefer `gpt-image-2-text-to-image` for a new master sheet.
- Prefer `gpt-image-2-image-to-image` for direct reference-led variations.
- Offer `nano-banana-pro` or `nano-banana-2` for controlled edits and alternate
  variants when their settings or user preference fit better.
- Keep identity traits immutable in every shot prompt; change pose, expression,
  wardrobe, lighting, and setting as separate variables.

Do not describe local identity bundles as trained biometric models or a
portable Soul equivalent.

## Enforce Paid-Action Gates

Lesson search, brief work, scripts, storyboards, references, model inspection,
and approvals are local. Preview and final generation are separate live Kie
actions that may consume credits.

Before each paid action:

1. Inspect the chosen model with `models show`.
2. Validate its complete input with `models validate`.
3. Show the plan and obtain explicit approval.

For every video shot, call preview generation, poll, display the returned image,
and wait. Call preview approval only after an explicit yes; otherwise reject,
revise one variable, and regenerate. Never infer approval from silence.
