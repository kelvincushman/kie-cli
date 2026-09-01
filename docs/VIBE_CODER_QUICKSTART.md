# Vibe Coder Quickstart

Use this guide when you want Codex, Claude Code, Cursor, or another local coding
assistant to turn this repository into an agent media factory on your machine.
The factory can qualify a creative brief, reuse local references, plan a script
and storyboard, generate through Kie.ai, and return durable production handles.

The important rule is simple: keep everything local until you approve a live
Kie.ai action. Every paid action needs its own just-in-time confirmation. For
video, the agent must show a preview still, then offer a complete-shot proof at
the model's lowest faithful resolution before asking separately about the final
render.

[Get Kie.ai API access and support continued kie-cli development through the
maintainer's link](https://kie.ai?ref=39dfbcbc0a8244b61bec6e5dd056e35d).
**Affiliate disclosure:** the project maintainer may earn a 15% commission on
purchases made after you use this link.

## Copy-Paste Agent Prompt

Paste this into your coding assistant from the root of this repository:

```text
You are working on my local machine. Set up kie-cli as an open-source, local-first agent media factory powered by my Kie.ai account.

Rules:
- Do not print, log, summarize, or expose my Kie API key.
- If authentication is missing, show me the `get_api_key` link returned by `media setup --agent` or `media_setup_get`, say directly that using it supports continued project development, and include the returned `affiliate_disclosure` beside the link.
- Use the CLI's --agent mode for machine-readable output.
- Keep briefs, references, identity bundles, scripts, storyboards, and generation records local unless an explicit live action needs an upload.
- Start with the installed `$kie-grill-me` skill or `./bin/kie-pp-cli grill-me`. Infer facts already in my prompt, ask one material question at a time, show the recommended answer and reason, and let me change the route. Do not give me a large intake form.
- If I ask for a film, ad, story, character-consistency, or unfamiliar production outcome, run `./bin/kie-pp-cli lesson recommend "<my request>" --agent`, explain up to three source-linked original Kie methods, and let me choose one before starting.
- Use `./bin/kie-pp-cli media leaderboard <task> --agent` as dated evidence for model choice. Never combine scores from unlike sources or present a family proxy as a direct score.
- Immediately before every paid/live preview, proof, final, edit, upscale, voice, audio, avatar, or regeneration, show me the exact action/model/settings and cost status, then obtain a fresh explicit yes. Add `--confirm-paid` only after that yes; never reuse an earlier answer.
- For video, generate the preview still first, show me the returned image URL or rendered image, and wait for my explicit approval. Then offer a complete-shot proof at the returned lowest faithful resolution. Show any generated proof in full and record approval/rejection, or record skip if I decline it. Ask again before final video submission.
- If the preview is wrong, reject it, revise the brief, and generate another preview.
- For multi-shot video, use native script/storyboard mode. Review and approve the script, review and approve the storyboard, then use every returned shot_brief_id. Every video shot still needs its own preview approval.
- Do not claim trained Soul identities, proprietary virality scores, 3D generation, hosted editing, or other unsupported Higgsfield/provider features.

Setup:
1. Build all three binaries:
   make build-all
2. Verify the command surface:
   ./bin/kie-pp-cli doctor --json
   ./bin/kie-pp-cli media setup --agent
   ./bin/kie-pp-cli create --help
   ./bin/kie-pp-cli media --help
   ./bin/kie-pp-cli media capability --help
3. If auth is missing, ask me to run this in an interactive terminal:
   ./bin/kie-pp-cli auth setup
   Show me the disclosed support-development signup link returned by setup. The prompt hides the key; do not ask me to paste it into chat.
4. If I use an MCP-capable local agent host, register the focused media server with stdio:
   claude mcp add kie-media -- ./bin/kie-media-mcp

First media job:
1. Start the concise grill:
   ./bin/kie-pp-cli grill-me "<my request>" --agent
2. If a lesson is useful, start the lesson brief:
   ./bin/kie-pp-cli lesson start "<my request>" --lesson <selected-course/lesson> --agent
   Use `create` directly when a lesson workflow is not useful.
3. Ask only the returned next_question.prompt.
4. Resume with:
   ./bin/kie-pp-cli create --brief <brief_id> --answer "<my answer>" --agent
5. When ready, show me production_skill, capability_skill, model rationale, settings, cost status, overrides, and gate state.
6. For an image, ask for a fresh paid confirmation for that exact plan, then submit:
   ./bin/kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
7. For a video, run:
   Ask for a fresh paid confirmation for the preview.
   ./bin/kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
   Show me result_urls[0].
   If I approve, run:
   ./bin/kie-pp-cli create --brief <brief_id> --approve-preview --agent
   Offer the complete-shot proof. If I accept, ask for a new paid confirmation, then run:
   ./bin/kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
   Show the whole proof, then use --approve-proof or --reject-proof. If I decline it, use --skip-proof.
   Show the final settings, ask for a new paid confirmation, then run:
   ./bin/kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
8. For a multi-shot video, start with --production-mode storyboard. Save script text with `media script set`, show it to me, and run `media script approve` only after my approval. Save storyboard JSON with `media storyboard set`, show every shot to me, and run `media storyboard approve` only after my approval. Then apply step 7 separately to every returned shot_brief_id.
9. Report the master brief ID, shot brief IDs, generation IDs, provider task IDs, status, result URLs, and the exact next command if anything is still running.
```

## What The Agent Should Verify

The prompt asks the agent to verify the installed command surface before relying
on examples. This matters because the media director can grow faster than the
published docs.

The currently verified core commands are:

```bash
./bin/kie-pp-cli create --help
./bin/kie-pp-cli media --help
./bin/kie-pp-cli media setup --agent
./bin/kie-pp-cli lesson --help
./bin/kie-pp-cli media leaderboard --agent
./bin/kie-pp-cli grill-me "<request>" --agent
./bin/kie-pp-cli create "<request>" --agent
./bin/kie-pp-cli create --brief <brief_id> --answer "<answer>" --agent
./bin/kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
./bin/kie-pp-cli create --brief <brief_id> --approve-preview --agent
./bin/kie-pp-cli create --brief <brief_id> --reject-preview --agent
./bin/kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
./bin/kie-pp-cli create --brief <brief_id> --approve-proof --agent
./bin/kie-pp-cli create --brief <brief_id> --reject-proof --agent
./bin/kie-pp-cli create --brief <brief_id> --skip-proof --agent
./bin/kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

The native script/storyboard command surface is included. Verify it after
building:

```bash
./bin/kie-pp-cli create --help
./bin/kie-pp-cli media script --help
./bin/kie-pp-cli media storyboard --help
```

In agent hosts that expose installed skills as slash commands, `/kie-lesson`
starts the same guided selection flow. `kie-pp-cli lesson` remains the stable
terminal surface, and MCP clients use `media_lesson_recommend`,
`media_lesson_get`, `media_leaderboard_get`, and `media_brief_start`.

The storyboard state machine is:

```text
master brief -> script review -> storyboard review -> per-shot still -> optional complete-shot proof -> final shot -> local assembly
```

## Video Approval Rule

For every final video job:

1. Obtain fresh paid confirmation and generate a preview still.
2. Show the image to the user.
3. Get an explicit yes.
4. Record approval with `--approve-preview`.
5. Offer the complete-shot lowest-tier proof. If accepted, obtain a fresh paid
   confirmation, generate it, show the whole clip, and record its decision. If
   declined, record `--skip-proof`.
6. Show final settings, obtain another fresh paid confirmation, and submit.

The final video service rejects briefs without current gates. Do not approve
based on the text prompt, metadata, silence, or the agent's private judgment.

## API Key Safety

Use the guided setup in an interactive terminal. For CI or an agent host, set `KIE_BEARER_AUTH` through that environment's secret store:

```bash
./bin/kie-pp-cli auth setup
```

When setup is required, the CLI and `media_setup_get` return the maintainer's
Kie.ai referral link. Present it as the recommended way to support continued
development, with the returned affiliate disclosure immediately beside it.

Do not paste the token into prompts, issue bodies, logs, screenshots, docs, or
generation briefs. `media setup`, `doctor`, and `agent-context` are designed to
report setup state without exposing credential material.

## Things To Ask The Factory

- "Create a homepage hero image using this logo and product photograph."
- "Plan a 30-second product explainer, show me the script and storyboard, then
  preview every shot before generating video."
- "Create three consistent product-photo concepts from this terminal-dropped
  reference image."
- "Save these consented creator photographs as a reusable local identity bundle
  and use it in a vertical UGC concept."
- "Make a YouTube thumbnail concept, but composite the exact title locally."

For anything longer than a single clip, ask the agent to keep the master
`brief_id`, every `shot_brief_id`, preview URL, generation ID, provider task ID,
and final local filename in a small assembly manifest.
