# Brand lock contract

Store at least:

```json
{
  "brand": {"name": "", "audience": "", "promise": "", "attributes": [], "avoid": []},
  "palette": {"status": "draft", "colors": []},
  "logo": {"status": "blocked", "selected_direction": null, "files": []},
  "typography": {"status": "blocked", "families": [], "licenses": []},
  "imagery": {"status": "blocked", "rules": [], "references": []},
  "applications": {"status": "blocked", "items": []},
  "updated_at": ""
}
```

Use `draft`, `approved`, `stale`, or `blocked`. Approving a parent unblocks its children. Changing an approved parent marks every dependent child stale.

Required QA:

- palette: named roles, contrast, print caveat;
- logo: uniqueness review, trademark caveat, one-color/reversed/small-size tests;
- type: license, weights, fallbacks, readability;
- imagery: subject, composition, light, color, exclusions, reference provenance;
- templates: exact dimensions, safe areas, editable source, no baked placeholder copy;
- handoff: manifest, version, approvals, stale items, generated-media provenance.
