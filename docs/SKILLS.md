# Kie Agent Skill Suite

The repository ships 17 small Markdown skills. They teach coding agents how to
operate the local CLI/MCP safely without embedding Kie.ai's 129 model schemas in
every prompt. Runtime commands and the generated capability catalog remain the
source of truth.

Install them from the repository:

```bash
npx skills add kelvincushman/kie-cli
```

Hosts that expose skills as slash commands can start with `/kie-grill-me`.
Other hosts can say “use `$kie-grill-me`”, or invoke the portable terminal
contract directly:

```bash
kie-pp-cli grill-me "<what I want to make>" --agent
```

## Routing model

| Skill | Role |
| --- | --- |
| `kie-grill-me` | Thin, user-invoked entry that starts the shared media interview |
| `kie-grilling` | Infers known facts, asks one material question, explains the recommendation, accepts overrides, and owns gates |
| `kie-create` | Durable CLI/MCP coordinator for briefs, references, generations, and resumption |
| `kie-image` | Text/reference image creation, edits, removals, and upscales |
| `kie-video` | Single continuous shot with still and complete-shot proof review |
| `kie-audio` | Narration, dialogue, voice, music, and sound effects |
| `kie-avatar` | Consented presenter, talking portrait, avatar, and lip-sync media |
| `kie-film` | Script, storyboard, continuity, per-shot production, and assembly handoff |
| `kie-identity` | Consented local likeness reference bundles; not a trained portable identity |
| `kie-lesson` | Source-linked Academy pattern discovery and original Kie-native methods |
| `kie-generate` | Advanced/manual cross-capability route when a narrower skill is unsuitable |
| `kie-brandkit` | Approval-led brand system and deterministic production-art handoff |
| `kie-marketplace-cards` | Truthful product listing image sets |
| `kie-product-photoshoot` | Consistent reference-led product campaign imagery |
| `kie-video-explainer` | Narrated explainer blocks using the film/audio/video skills |
| `kie-websites` | Kie media assets plus a separate local website build workflow |
| `kie-youtube-thumbnail` | Truthful thumbnail concepts with deterministic final typography |

The route returned in `brief.plan` includes `production_skill`,
`capability_skill`, `rationale`, `cost_status`, and `override_options`. Show this
briefly and let the user choose another supported route.

## Shared conversation contract

`kie-grilling` follows the useful part of Matt Pocock's pattern: one question
at a time, recommendations included, environmental facts discovered rather
than asked, and decisions left with the user. It is deliberately less
relentless for media production: the runtime infers explicit prompt facts and
stops once the plan is unambiguous.

```bash
kie-pp-cli grill-me "a five-second vertical product reveal with no audio" --agent
kie-pp-cli create --brief <brief_id> --answer "<answer>" --agent
kie-pp-cli create --brief <brief_id> --wrap-up --agent  # only when requested
```

MCP equivalents are `media_grill_start`, `media_grill_answer`, and
`media_grill_wrap_up`.

## Runtime introspection

Skills must not guess which fields a current model accepts:

```bash
kie-pp-cli media capability list --capability kie-video --agent
kie-pp-cli media capability show bytedance/seedance-2-5 --agent
kie-pp-cli models show bytedance/seedance-2-5 --agent
kie-pp-cli models example bytedance/seedance-2-5 --agent
kie-pp-cli models validate bytedance/seedance-2-5 --input '<json>' --agent
```

MCP uses `media_capability_list/get` and `media_model_list/get/example/validate`.
This keeps model updates in the generated catalog and prevents skill text from
becoming a second, drifting schema source.

## Approval and credit contract

Every live generation is a separate paid action. Immediately before it, the
agent must show the exact action, brief/model/settings, known cost status, and
ask for a fresh explicit yes. CLI director mode records that yes with
`--confirm-paid`; focused MCP preview, proof, and final calls require a new
`media_paid_confirm` result. That director confirmation is scoped, expires, and
can be consumed once. Broad generated endpoint tools remain an advanced manual
surface: skills must show their complete payload and ask again, but must not
claim those raw tools are protected by the director's plan-hash transaction.
Use the current turn's `paid_action` fields when creating the director
confirmation; agents must not calculate or guess its plan hash.

For every video shot:

```text
ready brief
  -> fresh-confirmed paid still
  -> user sees and approves/rejects still
  -> optional fresh-confirmed lowest-tier complete-shot proof
  -> user sees and approves/rejects proof, or explicitly skips it
  -> fresh-confirmed final render
```

Creative approval is not credit authorization. A changed prompt, setting,
model, reference, script, or storyboard can invalidate earlier artifacts. See
[Proof and Paid Confirmation](PROOF_AND_PAID_CONFIRMATION.md).

## Files and terminal input

When the terminal supports drag-and-drop, the shell normally receives a quoted
local path. Skills may pass a path, HTTP(S) URL, or `ref:<id>` handle. They must
not promise visual display in a text-only host; print the direct artifact URL
clearly when rendering is unavailable.

Vault reusable media and consented identities locally. Never put secrets in a
skill prompt, brief, script, storyboard, log, or issue. Local references upload
only for an explicitly authorized live call.
