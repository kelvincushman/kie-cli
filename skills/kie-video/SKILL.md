---
name: kie-video
description: Direct a Kie.ai single-shot video with text, first-frame, first-and-last-frame, or multimodal references; mandatory visible still approval; optional low-resolution complete-shot proof; and fresh confirmation for each paid action. Use for short video, image-to-video, cinematic shot, motion ad, or SeedDance work.
---

# Kie Video

Start with `$kie-grilling`. For multiple shots, narrative continuity, an advert, explainer, tutorial, or film, hand off to `$kie-film` before generating anything.

Inspect the compact route and chosen schema:

```bash
kie-pp-cli media capability list --capability kie-video --agent
kie-pp-cli media capability show <model-id> --agent
kie-pp-cli models show <model-id> --agent
```

Choose one documented video mode: text, first frame, first plus last frame, or multimodal references. Collect only the inputs that mode needs. State why the route fits and let the user override it.

Follow the shared gates exactly:

```bash
# Each live command below needs a new explicit user yes first.
kie-pp-cli create --brief <brief_id> --preview --confirm-paid --wait --agent
# Show the actual still; then record the user's decision.
kie-pp-cli create --brief <brief_id> --approve-preview --agent

# Optional complete-shot dry run at the returned lowest faithful tier.
kie-pp-cli create --brief <brief_id> --proof --confirm-paid --wait --agent
# Show the actual clip; then approve/reject, or use --skip-proof if declined.
kie-pp-cli create --brief <brief_id> --approve-proof --agent

# Final render needs another new explicit yes.
kie-pp-cli create --brief <brief_id> --submit --confirm-paid --wait --agent
```

Never record approval without showing the artifact. Never treat still/proof approval as final-render payment authorization. A creative change can invalidate prior approvals; obey `gate_state`, `next_action`, and `resume_command`. Report every durable ID, model/settings, result URL, and known limitation.
