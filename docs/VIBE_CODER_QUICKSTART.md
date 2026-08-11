# Vibe Coder Quickstart

Use this guide when you want Codex, Claude Code, Cursor, or another local coding
assistant to turn this repository into an agent media factory on your machine.
The factory can qualify a creative brief, reuse local references, plan a script
and storyboard, generate through Kie.ai, and return durable production handles.

The important rule is simple: keep everything local until you approve a live
Kie.ai action. For video, the agent must generate and show a preview still
first, then wait for your explicit approval before submitting the final video.

## Copy-Paste Agent Prompt

Paste this into your coding assistant from the root of this repository:

```text
You are working on my local machine. Set up kie-cli as an open-source, local-first agent media factory powered by my Kie.ai account.

Rules:
- Do not print, log, summarize, or expose my Kie API key.
- Use the CLI's --agent mode for machine-readable output.
- Keep briefs, references, identity bundles, scripts, storyboards, and generation records local unless an explicit live action needs an upload.
- Ask one useful qualifying question at a time. Do not give me a large intake form.
- Do not submit any paid/live generation until I approve the ready plan.
- For video, generate the preview still first, show me the returned image URL or rendered image, and wait for my explicit approval before final video submission.
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
3. If auth is missing, ask me to run this in an interactive terminal:
   ./bin/kie-pp-cli auth setup
   The prompt hides the key; do not ask me to paste it into chat.
4. If I use an MCP-capable local agent host, register the focused media server with stdio:
   claude mcp add kie-media -- ./bin/kie-media-mcp

First media job:
1. Ask me what I want to create.
2. Start the brief:
   ./bin/kie-pp-cli create "<my request>" --agent
3. Ask only the returned next_question.prompt.
4. Resume with:
   ./bin/kie-pp-cli create --brief <brief_id> --answer "<my answer>" --agent
5. When ready, show me the plan.
6. For an image, submit only after I approve:
   ./bin/kie-pp-cli create --brief <brief_id> --submit --wait --agent
7. For a video, run:
   ./bin/kie-pp-cli create --brief <brief_id> --preview --wait --agent
   Show me result_urls[0].
   If I approve, run:
   ./bin/kie-pp-cli create --brief <brief_id> --approve-preview --agent
   ./bin/kie-pp-cli create --brief <brief_id> --submit --wait --agent
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
./bin/kie-pp-cli create "<request>" --agent
./bin/kie-pp-cli create --brief <brief_id> --answer "<answer>" --agent
./bin/kie-pp-cli create --brief <brief_id> --preview --wait --agent
./bin/kie-pp-cli create --brief <brief_id> --approve-preview --agent
./bin/kie-pp-cli create --brief <brief_id> --reject-preview --agent
./bin/kie-pp-cli create --brief <brief_id> --submit --wait --agent
```

The native script/storyboard command surface is included. Verify it after
building:

```bash
./bin/kie-pp-cli create --help
./bin/kie-pp-cli media script --help
./bin/kie-pp-cli media storyboard --help
```

The storyboard state machine is:

```text
master brief -> script review -> storyboard review -> per-shot preview review -> per-shot video -> local assembly
```

## Video Approval Rule

For every final video job:

1. Generate a preview still.
2. Show the image to the user.
3. Get an explicit yes.
4. Record approval with `--approve-preview`.
5. Submit the final video.

The final video service rejects unapproved video briefs. Do not approve based on
the text prompt, metadata, silence, or the agent's private judgment.

## API Key Safety

Use the guided setup in an interactive terminal. For CI or an agent host, set `KIE_BEARER_AUTH` through that environment's secret store:

```bash
./bin/kie-pp-cli auth setup
```

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
