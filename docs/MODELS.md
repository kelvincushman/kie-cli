# Kie.ai Market Model Catalog

This reproducible snapshot contains **129 unique Market model contracts** discovered from Kie.ai's official documentation index.
Every model shares the create/query task API, but its full input schema and settings are preserved separately.

- List models: `kie-pp-cli models list`
- List token-light routes: `kie-pp-cli media capability list --capability kie-video --agent`
- Inspect route/proof metadata: `kie-pp-cli media capability show <model-id> --agent`
- Inspect exact settings: `kie-pp-cli models show <model-id>`
- Create a starter payload: `kie-pp-cli models example <model-id>`
- Validate before spending credits: `kie-pp-cli models validate <model-id> --input '{...}'`
- Read all field tables: [MODEL_INPUTS.md](MODEL_INPUTS.md)

Generation: `kie-pp-cli kie-ai-jobs market-create-task --model <model-id> --input '{...}'`

## Capability classifier

Every one of the 129 models is classified locally as `kie-image`, `kie-video`,
`kie-audio`, `kie-avatar`, or `kie-identity`, with secondary capabilities,
production-fit tags, routing notes, and proof-resolution metadata where
documented. The focused MCP exposes the same compact data through
`media_capability_list` and `media_capability_get`.

The embedded classifier is generated from `internal/kiecatalog/catalog.json`
and `research/kie-api-coverage.json`. Source hashes, model count, and all 70
documented operations/219 variants are validated so a refreshed catalog cannot
silently leave an unclassified model. Load the chosen model's full schema only
after capability routing.

The classifier is a technical route map, not a quality leaderboard. Use
[MODEL_LEADERBOARD.md](MODEL_LEADERBOARD.md) for dated external evidence and
[PROOF_AND_PAID_CONFIRMATION.md](PROOF_AND_PAID_CONFIRMATION.md) for the
lowest-faithful-proof contract.

## Bytedance

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Bytedance Seedance 1.5 Pro | `bytedance/seedance-1.5-pro` | 8 | [docs](https://docs.kie.ai/market/bytedance/seedance-1-5-pro) |
| Bytedance Seedance 2.0 | `bytedance/seedance-2` | 13 | [docs](https://docs.kie.ai/market/bytedance/seedance-2) |
| Bytedance Seedance 2.5 | `bytedance/seedance-2-5` | 14 | [docs](https://docs.kie.ai/market/bytedance/seedance-2-5) |
| Bytedance Seedance 2.0 Fast | `bytedance/seedance-2-fast` | 12 | [docs](https://docs.kie.ai/market/bytedance/seedance-2-fast) |
| Bytedance Seedance 2.0 Mini | `bytedance/seedance-2-mini` | 12 | [docs](https://docs.kie.ai/market/bytedance/seedance-2-mini) |
| Bytedance - V1 Lite Image to Video | `bytedance/v1-lite-image-to-video` | 8 | [docs](https://docs.kie.ai/market/bytedance/v1-lite-image-to-video) |
| Bytedance - V1 Lite Text to Video | `bytedance/v1-lite-text-to-video` | 8 | [docs](https://docs.kie.ai/market/bytedance/v1-lite-text-to-video) |
| Bytedance V1 Pro Fast Image to Video | `bytedance/v1-pro-fast-image-to-video` | 5 | [docs](https://docs.kie.ai/market/bytedance/v1-pro-fast-image-to-video) |
| Bytedance V1 Pro Image to Video | `bytedance/v1-pro-image-to-video` | 7 | [docs](https://docs.kie.ai/market/bytedance/v1-pro-image-to-video) |
| Bytedance - V1 Pro Text to Video | `bytedance/v1-pro-text-to-video` | 8 | [docs](https://docs.kie.ai/market/bytedance/v1-pro-text-to-video) |

## Elevenlabs

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| elevenlabs/audio-isolation | `elevenlabs/audio-isolation` | 1 | [docs](https://docs.kie.ai/market/elevenlabs/audio-isolation) |
| elevenlabs/text-to-dialogue-v3 | `elevenlabs/text-to-dialogue-v3` | 5 | [docs](https://docs.kie.ai/market/elevenlabs/text-to-dialogue-v3) |
| elevenlabs/text-to-speech-multilingual-v2 | `elevenlabs/text-to-speech-multilingual-v2` | 10 | [docs](https://docs.kie.ai/market/elevenlabs/text-to-speech-multilingual-v2) |
| elevenlabs/text-to-speech-turbo-2-5 | `elevenlabs/text-to-speech-turbo-2-5` | 10 | [docs](https://docs.kie.ai/market/elevenlabs/text-to-speech-turbo-2-5) |

## Flux2

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Flux-2 - Image to Image | `flux-2/flex-image-to-image` | 5 | [docs](https://docs.kie.ai/market/flux2/flex-image-to-image) |
| Flux-2 - Text to Image | `flux-2/flex-text-to-image` | 4 | [docs](https://docs.kie.ai/market/flux2/flex-text-to-image) |
| Flux-2 - Pro Image to Image | `flux-2/pro-image-to-image` | 5 | [docs](https://docs.kie.ai/market/flux2/pro-image-to-image) |
| Flux-2 - Pro Text to Image | `flux-2/pro-text-to-image` | 4 | [docs](https://docs.kie.ai/market/flux2/pro-text-to-image) |

## Google

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Gemini 2.5 Pro Text to Speech | `google/gemini-2-5-pro-tts` | 13 | [docs](https://docs.kie.ai/google/gemini-2-5-pro-tts) |
| Gemini 3.1 Flash Text to speech | `google/gemini-3-1-flash-tts` | 13 | [docs](https://docs.kie.ai/market/google/gemini-3-1-flash-tts) |
| Google - imagen4 | `google/imagen4` | 4 | [docs](https://docs.kie.ai/market/google/imagen4) |
| Google - imagen4-fast | `google/imagen4-fast` | 4 | [docs](https://docs.kie.ai/market/google/imagen4-fast) |
| Google - imagen4-ultra | `google/imagen4-ultra` | 4 | [docs](https://docs.kie.ai/market/google/imagen4-ultra) |
| Google - Nano Banana | `google/nano-banana` | 5 | [docs](https://docs.kie.ai/market/google/nano-banana) |
| Google - Nano Banana Edit | `google/nano-banana-edit` | 5 | [docs](https://docs.kie.ai/market/google/nano-banana-edit) |
| Google - Nano Banana 2 | `nano-banana-2` | 5 | [docs](https://docs.kie.ai/market/google/nanobanana2) |
| Google - Nano Banana 2 Lite | `nano-banana-2-lite` | 3 | [docs](https://docs.kie.ai/market/google/nano-banana-2-lite) |
| Google - Nano Banana Pro | `nano-banana-pro` | 5 | [docs](https://docs.kie.ai/market/google/pro-image-to-image) |

## Gpt

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| GPT Image 2 - Image To Image | `gpt-image-2-image-to-image` | 4 | [docs](https://docs.kie.ai/market/gpt/gpt-image-2-image-to-image) |
| GPT Image-2 - Text to Image | `gpt-image-2-text-to-image` | 3 | [docs](https://docs.kie.ai/market/gpt/gpt-image-2-text-to-image) |

## Gpt Image

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| GPT Image-1.5 - Image to Image | `gpt-image/1.5-image-to-image` | 4 | [docs](https://docs.kie.ai/market/gpt-image/1-5-image-to-image) |
| GPT Image-1.5 - Text to Image | `gpt-image/1.5-text-to-image` | 3 | [docs](https://docs.kie.ai/market/gpt-image/1-5-text-to-image) |

## Grok Imagine

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Grok Imagine Video 1.5 Preview | `grok-imagine-video-1-5-preview` | 6 | [docs](https://docs.kie.ai/market/grok-imagine/1-5-preview) |
| Grok Imagine - Video Extend | `grok-imagine/extend` | 4 | [docs](https://docs.kie.ai/market/grok-imagine/extend) |
| Grok Imagine - image to image | `grok-imagine/image-to-image` | 3 | [docs](https://docs.kie.ai/market/grok-imagine/image-to-image) |
| Grok Imagine Image to Video | `grok-imagine/image-to-video` | 9 | [docs](https://docs.kie.ai/market/grok-imagine/image-to-video) |
| Grok Imagine - Text to Image | `grok-imagine/text-to-image` | 4 | [docs](https://docs.kie.ai/market/grok-imagine/text-to-image) |
| Grok Imagine Text to Video | `grok-imagine/text-to-video` | 6 | [docs](https://docs.kie.ai/market/grok-imagine/text-to-video) |
| Grok Imagine - Video Upscale | `grok-imagine/upscale` | 2 | [docs](https://docs.kie.ai/market/grok-imagine/upscale) |

## Hailuo

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Hailuo Pro Image to Video | `hailuo/02-image-to-video-pro` | 5 | [docs](https://docs.kie.ai/market/hailuo/02-image-to-video-pro) |
| Hailuo Standard Image to Video | `hailuo/02-image-to-video-standard` | 7 | [docs](https://docs.kie.ai/market/hailuo/02-image-to-video-standard) |
| Hailuo Pro Text to Video | `hailuo/02-text-to-video-pro` | 3 | [docs](https://docs.kie.ai/market/hailuo/02-text-to-video-pro) |
| Hailuo Standard Text to Video | `hailuo/02-text-to-video-standard` | 4 | [docs](https://docs.kie.ai/market/hailuo/02-text-to-video-standard) |
| Hailuo 2.3 Pro Image to Video | `hailuo/2-3-image-to-video-pro` | 5 | [docs](https://docs.kie.ai/market/hailuo/2-3-image-to-video-pro) |
| Hailuo 2.3 Standard Image to Video | `hailuo/2-3-image-to-video-standard` | 4 | [docs](https://docs.kie.ai/market/hailuo/2-3-image-to-video-standard) |

## Happyhorse

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| HappyHorse - image-to-video | `happyhorse/image-to-video` | 5 | [docs](https://docs.kie.ai/market/happyhorse/image-to-video) |
| HappyHorse - reference-to-video | `happyhorse/reference-to-video` | 6 | [docs](https://docs.kie.ai/market/happyhorse/reference-to-video) |
| HappyHorse - text-to-video | `happyhorse/text-to-video` | 5 | [docs](https://docs.kie.ai/market/happyhorse/text-to-video) |
| HappyHorse - video-edit | `happyhorse/video-edit` | 6 | [docs](https://docs.kie.ai/market/happyhorse/video-edit) |

## Happyhorse 1 1

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| HappyHorse-1-1 image-to-video | `happyhorse-1-1/image-to-video` | 4 | [docs](https://docs.kie.ai/market/happyhorse-1-1/image-to-video) |
| HappyHorse-1-1 reference-to-video | `happyhorse-1-1/reference-to-video` | 5 | [docs](https://docs.kie.ai/market/happyhorse-1-1/reference-to-video) |
| HappyHorse-1-1 text-to-video | `happyhorse-1-1/text-to-video` | 4 | [docs](https://docs.kie.ai/market/happyhorse-1-1/text-to-video) |

## Ideogram

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Ideogram - Character | `ideogram/character` | 9 | [docs](https://docs.kie.ai/market/ideogram/character) |
| Ideogram - Character Edit | `ideogram/character-edit` | 9 | [docs](https://docs.kie.ai/market/ideogram/character-edit) |
| Ideogram - Character Remix | `ideogram/character-remix` | 13 | [docs](https://docs.kie.ai/market/ideogram/character-remix) |
| Ideogram V3 Edit | `ideogram/v3-edit` | 6 | [docs](https://docs.kie.ai/market/ideogram/v3-edit) |
| Ideogram V3 Remix | `ideogram/v3-remix` | 10 | [docs](https://docs.kie.ai/market/ideogram/v3-remix) |
| Ideogram V3 Text to Image | `ideogram/v3-text-to-image` | 7 | [docs](https://docs.kie.ai/market/ideogram/v3-text-to-image) |

## Infinitalk

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Infinitalk - From Audio | `infinitalk/from-audio` | 5 | [docs](https://docs.kie.ai/market/infinitalk/from-audio) |

## Kling

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Kling 2.6 Image to Video | `kling-2.6/image-to-video` | 4 | [docs](https://docs.kie.ai/market/kling/image-to-video) |
| Kling 2.6 motion-control | `kling-2.6/motion-control` | 5 | [docs](https://docs.kie.ai/market/kling/motion-control) |
| Kling 2.6 Text to Video | `kling-2.6/text-to-video` | 4 | [docs](https://docs.kie.ai/market/kling/text-to-video) |
| Kling-3.0 motion-control | `kling-3.0/motion-control` | 6 | [docs](https://docs.kie.ai/market/kling/motion-control-v3) |
| Kling 3.0 | `kling-3.0/video` | 17 | [docs](https://docs.kie.ai/market/kling/kling-3-0) |
| Kling AI Avatar Pro | `kling/ai-avatar-pro` | 3 | [docs](https://docs.kie.ai/market/kling/ai-avatar-pro) |
| Kling AI Avatar Standard | `kling/ai-avatar-standard` | 3 | [docs](https://docs.kie.ai/market/kling/ai-avatar-standard) |
| Kling V2.1 Master Image to Video | `kling/v2-1-master-image-to-video` | 5 | [docs](https://docs.kie.ai/market/kling/v2-1-master-image-to-video) |
| Kling V2.1 Master Text to Video | `kling/v2-1-master-text-to-video` | 5 | [docs](https://docs.kie.ai/market/kling/v2-1-master-text-to-video) |
| Kling V2.1 Pro | `kling/v2-1-pro` | 6 | [docs](https://docs.kie.ai/market/kling/v2-1-pro) |
| Kling V2.1 Standard | `kling/v2-1-standard` | 5 | [docs](https://docs.kie.ai/market/kling/v2-1-standard) |
| Kling - V2.5 Turbo Image to Video Pro | `kling/v2-5-turbo-image-to-video-pro` | 6 | [docs](https://docs.kie.ai/market/kling/v25-turbo-image-to-video-pro) |
| Kling - V2.5 Turbo Text to Video Pro | `kling/v2-5-turbo-text-to-video-pro` | 5 | [docs](https://docs.kie.ai/market/kling/v25-turbo-text-to-video-pro) |
| Kling - V3 Turbo Image to Video | `kling/v3-turbo-image-to-video` | 4 | [docs](https://docs.kie.ai/market/kling/v3-turbo-image-to-video) |
| Kling - V3 Turbo Text to Video | `kling/v3-turbo-text-to-video` | 4 | [docs](https://docs.kie.ai/market/kling/v3-turbo-text-to-video) |

## Market

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Gemini Omni Video | `gemini-omni-video` | 12 | [docs](https://docs.kie.ai/market/gemini-omni-video) |
| Omnihuman 1.5 | `omnihuman-1-5` | 7 | [docs](https://docs.kie.ai/market/omnihuman-1-5) |
| Seedream5.0 Lite - Image to Image | `seedream/5-lite-image-to-image` | 6 | [docs](https://docs.kie.ai/market/seedream-5-lite-image-to-image) |

## Minimax H3

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| MiniMax H3 Image-to-Video | `minimax-h3/image-to-video` | 5 | [docs](https://docs.kie.ai/market/minimax-h3/image-to-video) |
| MiniMax H3 Reference-to-Video | `minimax-h3/reference-to-video` | 7 | [docs](https://docs.kie.ai/market/minimax-h3/reference-to-video) |
| MiniMax H3 Text-to-Video | `minimax-h3/text-to-video` | 4 | [docs](https://docs.kie.ai/market/minimax-h3/text-to-video) |

## Omnihuman 1 5

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Omnihuman 1.5 Human Identification | `omnihuman-1-5/human-identification` | 1 | [docs](https://docs.kie.ai/market/omnihuman-1-5/human-identification) |
| OmniHuman 1.5 Subject Detection | `omnihuman-1-5/subject-detection` | 1 | [docs](https://docs.kie.ai/market/omnihuman-1-5/subject-detection) |

## Pixverse

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| PixVerse V6 Video Extension | `pixverse-v6/extend` | 0 | [docs](https://docs.kie.ai/market/pixverse/extend) |
| PixVerse V6 Image-to-Video | `pixverse-v6/image-to-video` | 8 | [docs](https://docs.kie.ai/market/pixverse/image-to-video) |
| PixVerse V6 Fusion / Reference-to-Video | `pixverse-v6/reference-to-video` | 10 | [docs](https://docs.kie.ai/market/pixverse/reference-to-video) |
| PixVerse V6 Text-to-Video | `pixverse-v6/text-to-video` | 7 | [docs](https://docs.kie.ai/market/pixverse/text-to-video) |
| PixVerse V6 First & Last Frame Transition | `pixverse-v6/transition` | 7 | [docs](https://docs.kie.ai/market/pixverse/transition) |

## Qwen

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Qwen - Image Edit | `qwen/image-edit` | 13 | [docs](https://docs.kie.ai/market/qwen/image-edit) |
| Qwen - Image to Image | `qwen/image-to-image` | 11 | [docs](https://docs.kie.ai/market/qwen/image-to-image) |
| Qwen - Text to Image | `qwen/text-to-image` | 10 | [docs](https://docs.kie.ai/market/qwen/text-to-image) |

## Qwen2

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Qwen2 - Image Edit | `qwen2/image-edit` | 6 | [docs](https://docs.kie.ai/market/qwen2/image-edit) |
| Qwen2 - Text To Image | `qwen2/text-to-image` | 5 | [docs](https://docs.kie.ai/market/qwen2/text-to-image) |

## Qwen3

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Qwen3 Image to Image | `qwen3/image-to-image` | 9 | [docs](https://docs.kie.ai/market/qwen3/image-to-image) |
| Qwen3 Text to Image | `qwen3/text-to-image` | 8 | [docs](https://docs.kie.ai/market/qwen3/text-to-image) |

## Qwen3 Pro

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Qwen3 Pro Image to Image | `qwen3/pro-image-to-image` | 9 | [docs](https://docs.kie.ai/market/qwen3-pro/image-to-image) |
| Qwen3 Pro Text to Image | `qwen3/pro-text-to-image` | 8 | [docs](https://docs.kie.ai/market/qwen3-pro/text-to-image) |

## Recraft

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Recraft - Crisp Upscale | `recraft/crisp-upscale` | 1 | [docs](https://docs.kie.ai/market/recraft/crisp-upscale) |
| Recraft - Remove Background | `recraft/remove-background` | 1 | [docs](https://docs.kie.ai/market/recraft/remove-background) |

## Seedream

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Seedream3.0 - Text to Image | `bytedance/seedream` | 4 | [docs](https://docs.kie.ai/market/seedream/seedream) |
| Seedream4.0 - Edit | `bytedance/seedream-v4-edit` | 7 | [docs](https://docs.kie.ai/market/seedream/seedream-v4-edit) |
| Seedream4.0 - Text to Image | `bytedance/seedream-v4-text-to-image` | 6 | [docs](https://docs.kie.ai/market/seedream/seedream-v4-text-to-image) |
| Seedream4.5 - Edit | `seedream/4.5-edit` | 5 | [docs](https://docs.kie.ai/market/seedream/4-5-edit) |
| Seedream4.5 - Text to Image | `seedream/4.5-text-to-image` | 4 | [docs](https://docs.kie.ai/market/seedream/4-5-text-to-image) |
| Seedream5.0 Lite - Text to Image | `seedream/5-lite-text-to-image` | 5 | [docs](https://docs.kie.ai/market/seedream/5-lite-text-to-image) |
| Seedream5.0 Pro - Image to Image | `seedream/5-pro-image-to-image` | 5 | [docs](https://docs.kie.ai/market/seedream/5-pro-image-to-image) |
| Seedream 5.0 Pro -  Layer Decomposition | `seedream/5-pro-layer-decomposition` | 4 | [docs](https://docs.kie.ai/market/seedream/5-pro-layer-decomposition) |
| Seedream5.0 Pro - Text to Image | `seedream/5-pro-text-to-image` | 5 | [docs](https://docs.kie.ai/market/seedream/5-pro-text-to-image) |

## Topaz

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Topaz - Image Upscale | `topaz/image-upscale` | 2 | [docs](https://docs.kie.ai/market/topaz/image-upscale) |
| Topaz - Video Upscale | `topaz/video-upscale` | 2 | [docs](https://docs.kie.ai/market/topaz/video-upscale) |

## Volcengine

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Volcengine video to video lip sync | `volcengine/video-to-video-lip-sync` | 8 | [docs](https://docs.kie.ai/market/volcengine/video-to-video-lip-sync) |

## Wan

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Wan - Image to Video | `wan/2-2-a14b-image-to-video-turbo` | 7 | [docs](https://docs.kie.ai/market/wan/2-2-a14b-image-to-video-turbo) |
| Wan - 2.2 A14B Speech to Video Turbo | `wan/2-2-a14b-speech-to-video-turbo` | 12 | [docs](https://docs.kie.ai/market/wan/2-2-a14b-speech-to-video-turbo) |
| Wan - Text to Video | `wan/2-2-a14b-text-to-video-turbo` | 6 | [docs](https://docs.kie.ai/market/wan/2-2-a14b-text-to-video-turbo) |
| Wan - Animate Move | `wan/2-2-animate-move` | 4 | [docs](https://docs.kie.ai/market/wan/2-2-animate-move) |
| Wan - Animate Replace | `wan/2-2-animate-replace` | 4 | [docs](https://docs.kie.ai/market/wan/2-2-animate-replace) |
| Wan 2.5 - Image to Video | `wan/2-5-image-to-video` | 8 | [docs](https://docs.kie.ai/market/wan/2-5-image-to-video) |
| Wan 2.5 - Text to Video | `wan/2-5-text-to-video` | 8 | [docs](https://docs.kie.ai/market/wan/2-5-text-to-video) |
| Wan - 2.6-flash-image-to-video | `wan/2-6-flash-image-to-video` | 7 | [docs](https://docs.kie.ai/market/wan/2-6-flash-image-to-video) |
| Wan - 2-6-flash-video-to-video | `wan/2-6-flash-video-to-video` | 7 | [docs](https://docs.kie.ai/market/wan/2-6-flash-video-to-video) |
| Wan 2.6 - Image to Video | `wan/2-6-image-to-video` | 5 | [docs](https://docs.kie.ai/market/wan/2-6-image-to-video) |
| Wan 2.6 - Text to Video | `wan/2-6-text-to-video` | 5 | [docs](https://docs.kie.ai/market/wan/2-6-text-to-video) |
| Wan 2.6 - Video to Video | `wan/2-6-video-to-video` | 6 | [docs](https://docs.kie.ai/market/wan/2-6-video-to-video) |
| Wan 2.7 Image | `wan/2-7-image` | 14 | [docs](https://docs.kie.ai/market/wan/2-7-image) |
| Wan 2.7 Image Pro | `wan/2-7-image-pro` | 14 | [docs](https://docs.kie.ai/market/wan/2-7-image-pro) |
| Wan 2.7 - Image to Video | `wan/2-7-image-to-video` | 12 | [docs](https://docs.kie.ai/market/wan/2-7-image-to-video) |
| Wan 2.7 - Reference to Video | `wan/2-7-r2v` | 13 | [docs](https://docs.kie.ai/market/wan/2-7-r2v) |
| Wan 2.7 - Text to Video | `wan/2-7-text-to-video` | 10 | [docs](https://docs.kie.ai/market/wan/2-7-text-to-video) |
| Wan 2.7 - Video Edit | `wan/2-7-videoedit` | 12 | [docs](https://docs.kie.ai/market/wan/2-7-videoedit) |

## Z Image

| Model | `--model` value | Input fields | Source |
|---|---|---:|---|
| Z-Image | `z-image` | 3 | [docs](https://docs.kie.ai/market/z-image/z-image) |
