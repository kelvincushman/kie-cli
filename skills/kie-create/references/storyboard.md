# Storyboard Production Contract

Use one master brief for story intent and one generated child brief per shot.
Script/storyboard operations are local-only; preview and final generation are
separate live Kie.ai calls.

## Sequence

1. Start a fully qualified video brief with `--production-mode storyboard`.
2. Write the complete script to a file or stdin and call `media script set`.
3. Show the script. Call `media script approve` only after explicit approval.
4. Write storyboard JSON. Shot durations must be 4-30 seconds and sum to the
   master duration.
5. Call `media storyboard set`, show all shots, then record the user's explicit
   approve or reject decision.
6. For each returned `shot_brief_id`, generate the preview, poll it, display the
   image, and wait for explicit approval. Reject and revise wrong shots.
7. Submit only approved shot briefs. Poll and record every result URL.
8. Assemble clips locally and disclose the assembly tool and command.

## Storyboard JSON

```json
{
  "title": "Launch story",
  "shots": [
    {
      "id": "opening",
      "title": "The problem",
      "duration_seconds": 5,
      "visual": "Founder faces a cluttered desk before sunrise",
      "camera": "slow dolly in",
      "narration": "Work should not begin in chaos.",
      "dialogue": "",
      "transition": "match cut to the clean workspace",
      "references": ["ref:ref_brand"]
    }
  ]
}
```

Do not include output-only fields such as `shot_brief_id`. Script edits
invalidate the storyboard. Master creative edits make script/storyboard
approvals stale. A storyboard master is never a final generation target.

MCP uses `media_script_set/get/decide` and
`media_storyboard_set/get/decide`; each shot then uses the normal
`media_preview_*`, `media_generate`, and `media_generation_status` tools.
