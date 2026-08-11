# Explainer production contract

Use a local manifest with:

```json
{
  "title": "",
  "minutes": 1,
  "language": "en",
  "voice": {"model": "", "id": "", "pace": ""},
  "visual_style": {"description": "", "style_key": ""},
  "aspect_ratio": "16:9",
  "subtitles": {"enabled": true, "format": "srt"},
  "blocks": []
}
```

The master `brief_id`, approved `script_id`, and approved `storyboard_id` are
required. Each block records its returned `shot_brief_id`, index, narration,
visual prompt, exact overlay copy, citations, preview generation/result/user
decision, audio task/result, final video generation/task/result, measured
duration, retry parent, and QA state.

Phase barriers:

1. approve settings;
2. verify research and script;
3. approve voice and style key;
4. finish narration jobs;
5. finish visual jobs;
6. assemble and caption locally;
7. inspect the complete runtime and delivery files.

Do not start downstream phases while required parent artifacts are missing or unapproved.
