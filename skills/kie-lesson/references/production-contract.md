# Lesson Production Contract

## Script

Write an original script from the user's brief. For every beat record:

- dramatic purpose;
- one visible action;
- subject and continuity state;
- dialogue or narration that fits the duration;
- transition into the next beat;
- objective acceptance criterion.

Do not adapt source-course wording. Use the selected lesson only as a method
label and source link.

## Asset lock

Separate immutable truth from shot-specific direction:

```text
LOCKED: identity, proportions, age, face, hair, product geometry, logos,
materials, wardrobe anchors, location topology, palette.

VARIABLE: pose, expression, action, camera, lens, lighting, weather, background
activity, dialogue, and transition.
```

Generate the master reference before scene work. Reject extra faces, merged
features, inconsistent logos, broken hands, false product details, or ambiguous
silhouettes.

## Storyboard shot

Each shot must include:

```json
{
  "title": "short purpose",
  "duration_seconds": 5,
  "visual": "one observable action",
  "camera": "shot size, lens feel, angle, and movement",
  "narration": "optional",
  "dialogue": "optional",
  "transition": "continuity into the next shot",
  "references": ["ref:<id>"]
}
```

The total shot duration must equal the master duration. Keep one primary action
per shot and define continuity in/out before generating motion.

## Per-shot prompt

Use this order:

1. locked character/product/location references;
2. single visible action and timing;
3. composition and spatial relationships;
4. camera framing, lens, angle, and motion;
5. lighting, material response, and palette;
6. dialogue, narration, ambience, or generated-audio intent;
7. continuity from the previous shot and into the next;
8. negative constraints and acceptance criterion.

## Review gate

Still review: identity, face, hands, product/logo truth, object geometry,
composition, lighting, text, and unwanted elements.

Motion review: all still criteria plus geography, physics, temporal continuity,
camera intent, prop state, dialogue timing, lip sync when required, and start/end
frames.

Reject a failed take and change one controllable variable. Preserve accepted
references and prompt fields so the next iteration is diagnosable.
