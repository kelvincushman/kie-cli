# Kie capability routes

Use current official docs when a schema or price may have changed.

| Need | Preferred Kie model or route | Boundary |
| --- | --- | --- |
| General image | `gpt-image-2-text-to-image` | Use `gpt-image-2-image-to-image` with references. |
| Product/reference image | `seedream/5-pro-image-to-image` | Up to ten documented references; inspect current schema. |
| Character still | `ideogram/character` | Per-generation reference, not a trained cross-model identity. |
| Default video | `bytedance/seedance-2-5` | Text, first-frame, first+last, and multimodal modes are mutually constrained. |
| Narration | `elevenlabs/text-to-speech-turbo-2-5` | Inspect voice and input fields before submission. |
| Dialogue | `elevenlabs/text-to-dialogue-v3` | Combined text and voice limits apply. |
| Presenter | `omnihuman-1-5` or Kling avatar | Requires portrait and audio. |
| Music | dedicated Suno generation API | Use current dedicated CLI help. |
| Sound effects | dedicated Suno sounds API | Use current dedicated CLI help. |
| Background removal | `recraft/remove-background` | Post-processing only. |
| Upscale | `recraft/crisp-upscale` or `topaz/image-upscale` | Post-processing only. |
| 3D mesh | unsupported | Do not promise a mesh or GLB result. |
| Virality score | unsupported | Offer a clearly labeled qualitative critique. |

SeedDance 2.5 currently supports 4–30 seconds or `-1` automatic duration, 480p/720p, optional audio, MP4/MOV, up to 30 image references, up to 10 audio references, and video references totaling at most 30 seconds. Verify the live page before relying on limits: https://docs.kie.ai/market/bytedance/seedance-2-5

The common Market surface is `POST /api/v1/jobs/createTask` plus `GET /api/v1/jobs/recordInfo?taskId=...`. The local director wraps this for image/video approval and reference uploads.
