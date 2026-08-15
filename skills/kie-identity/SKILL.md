---
name: kie-identity
description: Create and use consented local likeness-reference bundles for consistent people in Kie.ai images or videos. Use when a user wants recurring personal, creator, presenter, character, or brand-face media, supplies portrait references, or asks for a Higgsfield Soul ID equivalent without uploading a trained cross-model identity.
---

# Kie Identity

This is a local reference-bundle workflow, not biometric model training. Kie does not document a cross-model Soul-equivalent endpoint. State that limitation before creating the profile.

Use `$kie-grilling` for the surrounding creative brief. This skill owns the
rights and reusable-reference decision, then hands the durable identity handle
back to `$kie-image`, `$kie-video`, `$kie-avatar`, or `$kie-film`.

Prefer `kie-pp-cli media workflow show identity --agent` plus the `media identity` CLI commands. MCP hosts use `media_workflow_get`, `media_identity_create`, `media_identity_get`, and `media_identity_list` so photo bundles stay handle-based and token-light.

Read [references/photo-guide.md](references/photo-guide.md) before accepting a set.

## Confirm authority

Ask one question at a time:

1. Whose likeness is this?
2. Does the user have explicit permission to create and reuse these media references?
3. Is the intended use personal, editorial, commercial, or fictional?
4. Which outputs and time period are authorized?

Do not proceed without explicit consent. Apply extra care to minors, sexual content, political persuasion, fraud, impersonation, or deceptive endorsement. Never imply that a generated result is authentic footage.

## Create the local bundle

Review the photos for diversity and quality. Five to twenty images are useful; eight to twelve varied, clear images are the preferred set. Then save the bundle locally:

```bash
kie-pp-cli media identity create "<name>" \
  --reference <image-1> --reference <image-2> --reference <image-3> \
  --consent --agent
```

The returned `identity_...` record contains only private `ref:...` handles. Local files are copied to the private vault and upload only when an approved generation needs them.

Inspect profiles without exposing source paths:

```bash
kie-pp-cli media identity list --agent
kie-pp-cli media identity show identity_<id> --agent
```

## Use the bundle

Pass the identity to the director:

```bash
kie-pp-cli create "<portrait or scene>" --type image --identity identity_<id> --agent
kie-pp-cli create "<video scene>" --type video --video-mode multimodal \
  --identity identity_<id> --duration 10 --audio off --agent
```

For stills, the director expands the profile into image-to-image references. For video, SeedDance 2.5 receives the images as multimodal references. Inspect current docs before using provider-specific `characterId` features such as Gemini Omni; keep those IDs provider-scoped and do not describe them as portable Soul models.

Before submitting any likeness video, ask for a fresh paid confirmation, run
`create --brief <brief_id> --preview --confirm-paid --wait --agent`, display the
returned still, and let the person/user inspect identity match. Record approval
with `--approve-preview` only after an explicit yes. Use `--reject-preview` for
identity drift or an unwanted portrayal. Offer the complete-shot low-resolution
proof. If accepted, obtain a fresh confirmation, generate and display it, then
record `--approve-proof` or `--reject-proof`. If the user declines the proof,
record `--skip-proof` without generating it. In either branch, obtain a separate
fresh confirmation for the final render.
The final video call is blocked until the current gates are satisfied.

Review identity match, age, skin tone, facial geometry, hair, distinguishing features, and unwanted stereotyping. Do not claim perfect identity preservation. Deliver the identity handle, authorized use summary, generation provenance, and any match caveat.
