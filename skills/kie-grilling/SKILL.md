---
name: kie-grilling
description: Shared Kie media interview protocol for agents. Use when a media request needs concise qualification, reference or likeness intake, visible model rationale, route overrides, approval gates, or a safe handoff to image, video, audio, avatar, film, identity, or outcome skills.
---

# Kie Grilling

Turn an idea into a durable local brief without turning the conversation into a form. The CLI/MCP is the source of truth; do not reproduce model schemas or invent settings in the conversation.

## Start from runtime truth

Prefer the token-light CLI:

```bash
kie-pp-cli doctor --json
kie-pp-cli media setup --agent
kie-pp-cli grill-me "<the user's request>" --agent
```

An MCP host uses `media_grill_start`. Preserve the returned `brief.id` and resume with `media_grill_answer`; use `media_grill_wrap_up` only when the user asks the agent to choose sensible remaining defaults.

Ask exactly the returned `next_question.prompt`, one question at a time. Show its `recommendation`, `recommendation_reason`, and choices briefly. Let the user accept the recommendation, give another answer, add a file/URL/reference handle, or ask to wrap up. Look up discoverable facts locally instead of asking.

## Route visibly

When the brief is ready, show:

- the recommended production skill and capability skill;
- model and settings;
- why they fit the request;
- cost status and known uncertainty;
- meaningful override options;
- references, rights acknowledgement, and missing constraints;
- the next gate and exact resume command.

Use `kie-pp-cli media capability show <model-id> --agent` and `kie-pp-cli models show <model-id> --agent` when route detail matters. MCP hosts use `media_capability_get` and `media_model_get`. The user may change the suggested route; validate the replacement locally before any paid call.

Route by outcome:

- still, edit, product image, thumbnail, brand asset: `$kie-image` or the matching outcome skill;
- motion, single shot, image-to-video: `$kie-video`;
- narration, dialogue, music, sound effect: `$kie-audio`;
- talking presenter or lip-synced person: `$kie-avatar`;
- script, scenes, continuity, multi-shot production: `$kie-film`;
- recurring likeness or character references: `$kie-identity` before production;
- broad or advanced direct-model work: `$kie-generate` only after explaining that it is a manual surface.

## Reference and rights intake

Accept terminal paths, drag-and-drop paths, URLs, or `ref:<id>` handles when the host supports them. Never imply that a terminal can display or receive files when it cannot. Vault reusable inputs with `media reference add`; local paths stay local until an approved generation uploads them.

Ask for only material assets: product artwork, logo, packaging, style frames, props, first/last frames, video/audio references, and consented likeness images. Do not continue with a person's likeness until the rights acknowledgement is explicit. `$kie-identity` creates local reference bundles, not a trained portable identity.

## Gates and paid calls

No live generation is authorized merely because the brief is ready. Immediately before every paid action, name the action, brief, model, settings/plan hash, and cost uncertainty, then obtain a fresh explicit yes. In CLI agent mode, add `--confirm-paid` only after that yes. For focused director MCP preview, proof, and final calls, use the current turn's `paid_action` fields to create a fresh `media_paid_confirm` transaction, then pass its exact `confirmation_id`; there is no boolean shortcut and the agent must not calculate or guess the plan hash. Broad generated Kie endpoint tools are an advanced surface and are not wrapped by that director transaction, so show their complete payload and obtain a fresh explicit yes immediately before calling them. Never reuse confirmation or infer it from an earlier approval.

For video:

1. Generate the mandatory paid still with `--preview --confirm-paid --wait`.
2. Actually render or clearly link the returned image. Wait for explicit approval or rejection.
3. Record `--approve-preview` or `--reject-preview` locally.
4. Offer one optional complete-shot proof at the selected model's lowest faithful resolution. Explain that this is a separate paid call.
5. If accepted, run `--proof --confirm-paid --wait`, show the whole clip, then record `--approve-proof` or `--reject-proof`. If declined, record `--skip-proof`.
6. Ask again before the final paid render; submit with `--submit --confirm-paid`.

Approval of a still or proof is creative approval, not payment authorization. A changed prompt, route, model, settings, script, storyboard, or reference can stale earlier approvals; follow the returned gate state rather than assuming they still apply.

Stop grilling once the shared plan is unambiguous. Report durable IDs, current gate, exact next action, and resume command. Do not make a paid call unless the current user response explicitly authorizes that exact call.
