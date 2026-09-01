# Proof, Approval, and Paid Confirmation

The media director separates three decisions that automated media tools often
collapse:

1. Is the creative plan acceptable?
2. Does the generated artifact look acceptable?
3. May this exact paid API call run now?

They are not interchangeable. A ready plan, approved still, or approved proof
never authorizes a later paid generation.

## Paid confirmation transaction

Immediately before every live preview, proof, final render, regeneration, edit,
upscale, voice, audio, avatar, or other generation, show:

- the action and durable brief ID;
- selected model and relevant settings;
- the plan fingerprint/hash when available;
- known cost information, or an explicit “unknown/not estimated” status;
- what reference media will leave the local machine.

After the user explicitly agrees to that exact call, CLI director mode may add
`--confirm-paid`. Interactive director mode asks locally. Focused MCP preview,
proof, and final clients must create a scoped record:

```text
media_paid_confirm(brief_id, kind, expected_model, expected_plan_hash)
media_preview_generate(..., confirmation_id)
```

Use the current turn's `paid_action` object for `kind`, `model`, and
`plan_hash`; do not calculate or guess the hash in the agent.

The record is brief-, action-, model-, and plan-scoped, expires after ten
minutes, and is single-use. Wrong, stale, expired, or reused records fail before
the Kie API is called.

The broad generated CLI/MCP endpoint surface is retained for advanced operations
that do not fit the image/video director. Its raw tools are not wrapped by this
brief-and-plan transaction. Skills must still display the exact payload and ask
for a fresh explicit yes immediately before each raw live call, but must not
claim the director confirmation protects those calls.

## Image flow

```bash
# Ask for an explicit yes for this exact plan first.
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

A retry, variation, edit, or upscale needs another confirmation.

## Video flow

### 1. Mandatory still

Ask for confirmation for the paid preview image, then:

```bash
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
```

Render `result_urls[0]` in a visual host or print its direct URL in a text-only
terminal. Wait for the user's decision:

```bash
kie-pp-cli create --brief <brief_id> --approve-preview --agent
# or
kie-pp-cli create --brief <brief_id> --reject-preview --agent
```

### 2. Optional complete-shot proof

After still approval, the director returns a catalog-derived proof option. It
uses the selected video model and its lowest documented faithful resolution;
for `bytedance/seedance-2-5`, that is 480p. This is a complete shot rather than
a shortened or substituted-model proxy.

Explain that it is another paid call. If accepted, obtain a fresh confirmation:

```bash
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
```

Show the whole proof clip, then:

```bash
kie-pp-cli create --brief <brief_id> --approve-proof --agent
# or
kie-pp-cli create --brief <brief_id> --reject-proof --agent
```

If the user does not want the paid proof:

```bash
kie-pp-cli create --brief <brief_id> --skip-proof --agent
```

Skipping is a recorded creative workflow decision, not final-render
authorization.

### 3. Final render

Show the final settings and ask again. Only after that new yes:

```bash
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

## Invalidation and resumption

Approvals are fingerprinted. A change to creative direction, route, model,
settings, reference media, identity, script, or storyboard can make a preview
or proof stale. Follow `gate_state`, `next_action`, and `resume_command`; never
reason from an old transcript alone.

Inspect durable state:

```bash
kie-pp-cli media brief show <brief_id> --agent
kie-pp-cli media generation status <generation_id> --agent
```

No paid Kie.ai call was used to validate this implementation. Automated tests
prove local gate order and that invalid confirmation paths make zero API calls;
they do not prove the subjective quality of generated media.
