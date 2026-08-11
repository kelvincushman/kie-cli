# Kie.ai Standard Market Task-Model Catalog

All 128 English Market model contracts below use the shared task API:

- Create a task: `kie-pp-cli kie-ai-jobs market-create-task --model <id> --input '{...}'`
- Query a task: `kie-pp-cli kie-ai-jobs market-query-task --task-id <id>`
- Search the embedded CLI snapshot: `kie-pp-cli media models [query]`

This is a point-in-time snapshot of every English docs.kie.ai Market page that documents the standard createTask route (128 model IDs); Kie's common Market page documents the shared recordInfo query. It intentionally excludes Market Chat and Omni pages with other paths; use their generated commands. Each row links to its official source. `research/build_spec.py` rejects a missing page contract, duplicate model ID, or title/contract mismatch during refresh. The CLI snapshot is embedded at build time, not fetched live; run `scripts/weekly-refresh.sh` and rebuild to refresh it.

## Image Models > Flux-2

| Model | `--model` value | Official source |
|---|---|---|
| Flux-2 - Image to Image | `flux-2/flex-image-to-image` | [Kie documentation](https://docs.kie.ai/market/flux2/flex-image-to-image.md) |
| Flux-2 - Text to Image | `flux-2/flex-text-to-image` | [Kie documentation](https://docs.kie.ai/market/flux2/flex-text-to-image.md) |
| Flux-2 - Pro Image to Image | `flux-2/pro-image-to-image` | [Kie documentation](https://docs.kie.ai/market/flux2/pro-image-to-image.md) |
| Flux-2 - Pro Text to Image | `flux-2/pro-text-to-image` | [Kie documentation](https://docs.kie.ai/market/flux2/pro-text-to-image.md) |

## Image Models > GPT Image

| Model | `--model` value | Official source |
|---|---|---|
| GPT Image 2 - Image To Image | `gpt-image-2-image-to-image` | [Kie documentation](https://docs.kie.ai/market/gpt/gpt-image-2-image-to-image.md) |
| GPT Image-2 - Text to Image | `gpt-image-2-text-to-image` | [Kie documentation](https://docs.kie.ai/market/gpt/gpt-image-2-text-to-image.md) |
| GPT Image-1.5 - Image to Image | `gpt-image/1.5-image-to-image` | [Kie documentation](https://docs.kie.ai/market/gpt-image/1-5-image-to-image.md) |
| GPT Image-1.5 - Text to Image | `gpt-image/1.5-text-to-image` | [Kie documentation](https://docs.kie.ai/market/gpt-image/1-5-text-to-image.md) |

## Image Models > Google

| Model | `--model` value | Official source |
|---|---|---|
| Google - imagen4 | `google/imagen4` | [Kie documentation](https://docs.kie.ai/market/google/imagen4.md) |
| Google - imagen4-fast | `google/imagen4-fast` | [Kie documentation](https://docs.kie.ai/market/google/imagen4-fast.md) |
| Google - imagen4-ultra | `google/imagen4-ultra` | [Kie documentation](https://docs.kie.ai/market/google/imagen4-ultra.md) |
| Google - Nano Banana | `google/nano-banana` | [Kie documentation](https://docs.kie.ai/market/google/nano-banana.md) |
| Google - Nano Banana Edit | `google/nano-banana-edit` | [Kie documentation](https://docs.kie.ai/market/google/nano-banana-edit.md) |
| Google - Nano Banana 2 | `nano-banana-2` | [Kie documentation](https://docs.kie.ai/market/google/nanobanana2.md) |
| Google - Nano Banana 2 Lite | `nano-banana-2-lite` | [Kie documentation](https://docs.kie.ai/market/google/nano-banana-2-lite.md) |
| Google - Nano Banana Pro | `nano-banana-pro` | [Kie documentation](https://docs.kie.ai/market/google/pro-image-to-image.md) |

## Image Models > Grok Imagine

| Model | `--model` value | Official source |
|---|---|---|
| Grok Imagine - image to image | `grok-imagine/image-to-image` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/image-to-image.md) |
| Grok Imagine - Text to Image | `grok-imagine/text-to-image` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/text-to-image.md) |

## Image Models > Ideogram

| Model | `--model` value | Official source |
|---|---|---|
| Ideogram - Character | `ideogram/character` | [Kie documentation](https://docs.kie.ai/market/ideogram/character.md) |
| Ideogram - Character Edit | `ideogram/character-edit` | [Kie documentation](https://docs.kie.ai/market/ideogram/character-edit.md) |
| Ideogram - Character Remix | `ideogram/character-remix` | [Kie documentation](https://docs.kie.ai/market/ideogram/character-remix.md) |
| Ideogram V3 Edit | `ideogram/v3-edit` | [Kie documentation](https://docs.kie.ai/market/ideogram/v3-edit.md) |
| Ideogram V3 Remix | `ideogram/v3-remix` | [Kie documentation](https://docs.kie.ai/market/ideogram/v3-remix.md) |
| Ideogram V3 Text to Image | `ideogram/v3-text-to-image` | [Kie documentation](https://docs.kie.ai/market/ideogram/v3-text-to-image.md) |

## Image Models > Qwen

| Model | `--model` value | Official source |
|---|---|---|
| Qwen - Image Edit | `qwen/image-edit` | [Kie documentation](https://docs.kie.ai/market/qwen/image-edit.md) |
| Qwen - Image to Image | `qwen/image-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen/image-to-image.md) |
| Qwen - Text to Image | `qwen/text-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen/text-to-image.md) |
| Qwen2 - Image Edit | `qwen2/image-edit` | [Kie documentation](https://docs.kie.ai/market/qwen2/image-edit.md) |
| Qwen2 - Text To Image | `qwen2/text-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen2/text-to-image.md) |
| Qwen3 Image to Image | `qwen3/image-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen3/image-to-image.md) |
| Qwen3 Pro Image to Image | `qwen3/pro-image-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen3-pro/image-to-image.md) |
| Qwen3 Pro Text to Image | `qwen3/pro-text-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen3-pro/text-to-image.md) |
| Qwen3 Text to Image | `qwen3/text-to-image` | [Kie documentation](https://docs.kie.ai/market/qwen3/text-to-image.md) |

## Image Models > Recraft

| Model | `--model` value | Official source |
|---|---|---|
| Recraft - Crisp Upscale | `recraft/crisp-upscale` | [Kie documentation](https://docs.kie.ai/market/recraft/crisp-upscale.md) |
| Recraft - Remove Background | `recraft/remove-background` | [Kie documentation](https://docs.kie.ai/market/recraft/remove-background.md) |

## Image Models > Seedream

| Model | `--model` value | Official source |
|---|---|---|
| Seedream3.0 - Text to Image | `bytedance/seedream` | [Kie documentation](https://docs.kie.ai/market/seedream/seedream.md) |
| Seedream4.0 - Edit | `bytedance/seedream-v4-edit` | [Kie documentation](https://docs.kie.ai/market/seedream/seedream-v4-edit.md) |
| Seedream4.0 - Text to Image | `bytedance/seedream-v4-text-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream/seedream-v4-text-to-image.md) |
| Seedream4.5 - Edit | `seedream/4.5-edit` | [Kie documentation](https://docs.kie.ai/market/seedream/4-5-edit.md) |
| Seedream4.5 - Text to Image | `seedream/4.5-text-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream/4-5-text-to-image.md) |
| Seedream5.0 Lite - Image to Image | `seedream/5-lite-image-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream-5-lite-image-to-image.md) |
| Seedream5.0 Lite - Text to Image | `seedream/5-lite-text-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream/5-lite-text-to-image.md) |
| Seedream5.0 Pro - Image to Image | `seedream/5-pro-image-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream/5-pro-image-to-image.md) |
| Seedream 5.0 Pro -  Layer Decomposition | `seedream/5-pro-layer-decomposition` | [Kie documentation](https://docs.kie.ai/market/seedream/5-pro-layer-decomposition.md) |
| Seedream5.0 Pro - Text to Image | `seedream/5-pro-text-to-image` | [Kie documentation](https://docs.kie.ai/market/seedream/5-pro-text-to-image.md) |

## Image Models > Topaz

| Model | `--model` value | Official source |
|---|---|---|
| Topaz - Image Upscale | `topaz/image-upscale` | [Kie documentation](https://docs.kie.ai/market/topaz/image-upscale.md) |

## Image Models > Wan

| Model | `--model` value | Official source |
|---|---|---|
| Wan 2.7 Image | `wan/2-7-image` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-image.md) |
| Wan 2.7 Image Pro | `wan/2-7-image-pro` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-image-pro.md) |

## Image Models > Z-image

| Model | `--model` value | Official source |
|---|---|---|
| Z-Image | `z-image` | [Kie documentation](https://docs.kie.ai/market/z-image/z-image.md) |

## Music Models > ElevenLabs

| Model | `--model` value | Official source |
|---|---|---|
| elevenlabs/audio-isolation | `elevenlabs/audio-isolation` | [Kie documentation](https://docs.kie.ai/market/elevenlabs/audio-isolation.md) |
| elevenlabs/text-to-dialogue-v3 | `elevenlabs/text-to-dialogue-v3` | [Kie documentation](https://docs.kie.ai/market/elevenlabs/text-to-dialogue-v3.md) |
| elevenlabs/text-to-speech-multilingual-v2 | `elevenlabs/text-to-speech-multilingual-v2` | [Kie documentation](https://docs.kie.ai/market/elevenlabs/text-to-speech-multilingual-v2.md) |
| elevenlabs/text-to-speech-turbo-2-5 | `elevenlabs/text-to-speech-turbo-2-5` | [Kie documentation](https://docs.kie.ai/market/elevenlabs/text-to-speech-turbo-2-5.md) |

## Music Models > Gemini

| Model | `--model` value | Official source |
|---|---|---|
| Gemini 3.1 Flash Text to speech | `google/gemini-3-1-flash-tts` | [Kie documentation](https://docs.kie.ai/market/google/gemini-3-1-flash-tts.md) |

## Video Models > Bytedance

| Model | `--model` value | Official source |
|---|---|---|
| Bytedance Seedance 1.5 Pro | `bytedance/seedance-1.5-pro` | [Kie documentation](https://docs.kie.ai/market/bytedance/seedance-1-5-pro.md) |
| bytedance-seedance-2 | `bytedance/seedance-2` | [Kie documentation](https://docs.kie.ai/market/bytedance/seedance-2.md) |
| Bytedance Seedance 2.5 | `bytedance/seedance-2-5` | [Kie documentation](https://docs.kie.ai/market/bytedance/seedance-2-5.md) |
| Bytedance Seedance 2.0 Fast | `bytedance/seedance-2-fast` | [Kie documentation](https://docs.kie.ai/market/bytedance/seedance-2-fast.md) |
| Bytedance Seedance 2.0 Mini | `bytedance/seedance-2-mini` | [Kie documentation](https://docs.kie.ai/market/bytedance/seedance-2-mini.md) |
| Bytedance - V1 Lite Image to Video | `bytedance/v1-lite-image-to-video` | [Kie documentation](https://docs.kie.ai/market/bytedance/v1-lite-image-to-video.md) |
| Bytedance - V1 Lite Text to Video | `bytedance/v1-lite-text-to-video` | [Kie documentation](https://docs.kie.ai/market/bytedance/v1-lite-text-to-video.md) |
| Bytedance V1 Pro Fast Image to Video | `bytedance/v1-pro-fast-image-to-video` | [Kie documentation](https://docs.kie.ai/market/bytedance/v1-pro-fast-image-to-video.md) |
| Bytedance V1 Pro Image to Video | `bytedance/v1-pro-image-to-video` | [Kie documentation](https://docs.kie.ai/market/bytedance/v1-pro-image-to-video.md) |
| Bytedance - V1 Pro Text to Video | `bytedance/v1-pro-text-to-video` | [Kie documentation](https://docs.kie.ai/market/bytedance/v1-pro-text-to-video.md) |

## Video Models > Gemini Omni

| Model | `--model` value | Official source |
|---|---|---|
| Gemini Omni Video | `gemini-omni-video` | [Kie documentation](https://docs.kie.ai/market/gemini-omni-video.md) |

## Video Models > Grok Imagine

| Model | `--model` value | Official source |
|---|---|---|
| Grok Imagine Video 1.5 Preview | `grok-imagine-video-1-5-preview` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/1-5-preview.md) |
| Grok Imagine - Video Extend | `grok-imagine/extend` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/extend.md) |
| Grok Imagine Image to Video | `grok-imagine/image-to-video` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/image-to-video.md) |
| Grok Imagine Text to Video | `grok-imagine/text-to-video` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/text-to-video.md) |
| Grok Imagine - Video Upscale | `grok-imagine/upscale` | [Kie documentation](https://docs.kie.ai/market/grok-imagine/upscale.md) |

## Video Models > Hailuo

| Model | `--model` value | Official source |
|---|---|---|
|  Hailuo Pro Image to Video | `hailuo/02-image-to-video-pro` | [Kie documentation](https://docs.kie.ai/market/hailuo/02-image-to-video-pro.md) |
| Hailuo Standard Image to Video | `hailuo/02-image-to-video-standard` | [Kie documentation](https://docs.kie.ai/market/hailuo/02-image-to-video-standard.md) |
| Hailuo Pro Text to Video | `hailuo/02-text-to-video-pro` | [Kie documentation](https://docs.kie.ai/market/hailuo/02-text-to-video-pro.md) |
| Hailuo Standard Text to Video | `hailuo/02-text-to-video-standard` | [Kie documentation](https://docs.kie.ai/market/hailuo/02-text-to-video-standard.md) |
| Hailuo 2.3 Pro Image to Video | `hailuo/2-3-image-to-video-pro` | [Kie documentation](https://docs.kie.ai/market/hailuo/2-3-image-to-video-pro.md) |
| Hailuo 2.3 Standard Image to Video | `hailuo/2-3-image-to-video-standard` | [Kie documentation](https://docs.kie.ai/market/hailuo/2-3-image-to-video-standard.md) |

## Video Models > HappyHorse

| Model | `--model` value | Official source |
|---|---|---|
| HappyHorse-1-1 image-to-video | `happyhorse-1-1/image-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse-1-1/image-to-video.md) |
| HappyHorse-1-1 reference-to-video | `happyhorse-1-1/reference-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse-1-1/reference-to-video.md) |
| HappyHorse-1-1 text-to-video | `happyhorse-1-1/text-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse-1-1/text-to-video.md) |
| happyhorse-image-to-video | `happyhorse/image-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse/image-to-video.md) |
| happyhorse/reference-to-video | `happyhorse/reference-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse/reference-to-video.md) |
| happyhorse-text-to-video | `happyhorse/text-to-video` | [Kie documentation](https://docs.kie.ai/market/happyhorse/text-to-video.md) |
| happyhorse/video-edit | `happyhorse/video-edit` | [Kie documentation](https://docs.kie.ai/market/happyhorse/video-edit.md) |

## Video Models > Infinitalk

| Model | `--model` value | Official source |
|---|---|---|
| Infinitalk - From Audio | `infinitalk/from-audio` | [Kie documentation](https://docs.kie.ai/market/infinitalk/from-audio.md) |

## Video Models > Kling

| Model | `--model` value | Official source |
|---|---|---|
| Kling 2.6 Image to Video | `kling-2.6/image-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/image-to-video.md) |
| Kling 2.6 motion-control | `kling-2.6/motion-control` | [Kie documentation](https://docs.kie.ai/market/kling/motion-control.md) |
| Kling 2.6 Text to Video | `kling-2.6/text-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/text-to-video.md) |
| Kling-3.0 motion-control | `kling-3.0/motion-control` | [Kie documentation](https://docs.kie.ai/market/kling/motion-control-v3.md) |
| Kling 3.0 | `kling-3.0/video` | [Kie documentation](https://docs.kie.ai/market/kling/kling-3-0.md) |
| Kling AI Avatar Pro | `kling/ai-avatar-pro` | [Kie documentation](https://docs.kie.ai/market/kling/ai-avatar-pro.md) |
| Kling AI Avatar Standard | `kling/ai-avatar-standard` | [Kie documentation](https://docs.kie.ai/market/kling/ai-avatar-standard.md) |
| Kling V2.1 Master Image to Video | `kling/v2-1-master-image-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/v2-1-master-image-to-video.md) |
| Kling V2.1 Master Text to Video | `kling/v2-1-master-text-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/v2-1-master-text-to-video.md) |
| Kling V2.1 Pro | `kling/v2-1-pro` | [Kie documentation](https://docs.kie.ai/market/kling/v2-1-pro.md) |
| Kling V2.1 Standard | `kling/v2-1-standard` | [Kie documentation](https://docs.kie.ai/market/kling/v2-1-standard.md) |
| Kling - V2.5 Turbo Image to Video Pro | `kling/v2-5-turbo-image-to-video-pro` | [Kie documentation](https://docs.kie.ai/market/kling/v25-turbo-image-to-video-pro.md) |
| Kling - V2.5 Turbo Text to Video Pro | `kling/v2-5-turbo-text-to-video-pro` | [Kie documentation](https://docs.kie.ai/market/kling/v25-turbo-text-to-video-pro.md) |
| Kling - V3 Turbo Image to Video | `kling/v3-turbo-image-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/v3-turbo-image-to-video.md) |
| Kling - V3 Turbo Text to Video | `kling/v3-turbo-text-to-video` | [Kie documentation](https://docs.kie.ai/market/kling/v3-turbo-text-to-video.md) |

## Video Models > MiniMax H3

| Model | `--model` value | Official source |
|---|---|---|
| MiniMax H3 Image-to-Video | `minimax-h3/image-to-video` | [Kie documentation](https://docs.kie.ai/market/minimax-h3/image-to-video.md) |
| MiniMax H3 Reference-to-Video | `minimax-h3/reference-to-video` | [Kie documentation](https://docs.kie.ai/market/minimax-h3/reference-to-video.md) |
| MiniMax H3 Text-to-Video | `minimax-h3/text-to-video` | [Kie documentation](https://docs.kie.ai/market/minimax-h3/text-to-video.md) |

## Video Models > OmniHuman

| Model | `--model` value | Official source |
|---|---|---|
| Omnihuman 1.5 | `omnihuman-1-5` | [Kie documentation](https://docs.kie.ai/market/omnihuman-1-5.md) |
| Omnihuman 1.5 Human Identification | `omnihuman-1-5/human-identification` | [Kie documentation](https://docs.kie.ai/market/omnihuman-1-5/human-identification.md) |
| OmniHuman 1.5 Subject Detection | `omnihuman-1-5/subject-detection` | [Kie documentation](https://docs.kie.ai/market/omnihuman-1-5/subject-detection.md) |

## Video Models > PixVerse

| Model | `--model` value | Official source |
|---|---|---|
| PixVerse V6 Video Extension | `pixverse-v6/extend` | [Kie documentation](https://docs.kie.ai/market/pixverse/extend.md) |
| PixVerse V6 Image-to-Video | `pixverse-v6/image-to-video` | [Kie documentation](https://docs.kie.ai/market/pixverse/image-to-video.md) |
| PixVerse V6 Fusion / Reference-to-Video | `pixverse-v6/reference-to-video` | [Kie documentation](https://docs.kie.ai/market/pixverse/reference-to-video.md) |
| PixVerse V6 Text-to-Video | `pixverse-v6/text-to-video` | [Kie documentation](https://docs.kie.ai/market/pixverse/text-to-video.md) |
| PixVerse V6 First & Last Frame Transition | `pixverse-v6/transition` | [Kie documentation](https://docs.kie.ai/market/pixverse/transition.md) |

## Video Models > Topaz

| Model | `--model` value | Official source |
|---|---|---|
| Topaz - Video Upscale | `topaz/video-upscale` | [Kie documentation](https://docs.kie.ai/market/topaz/video-upscale.md) |

## Video Models > Volcengine

| Model | `--model` value | Official source |
|---|---|---|
| Volcengine video to video lip sync | `volcengine/video-to-video-lip-sync` | [Kie documentation](https://docs.kie.ai/market/volcengine/video-to-video-lip-sync.md) |

## Video Models > Wan

| Model | `--model` value | Official source |
|---|---|---|
| Wan - Image to Video | `wan/2-2-a14b-image-to-video-turbo` | [Kie documentation](https://docs.kie.ai/market/wan/2-2-a14b-image-to-video-turbo.md) |
| Wan - 2.2 A14B Speech to Video Turbo | `wan/2-2-a14b-speech-to-video-turbo` | [Kie documentation](https://docs.kie.ai/market/wan/2-2-a14b-speech-to-video-turbo.md) |
| Wan - Text to Video | `wan/2-2-a14b-text-to-video-turbo` | [Kie documentation](https://docs.kie.ai/market/wan/2-2-a14b-text-to-video-turbo.md) |
| Wan - Animate Move | `wan/2-2-animate-move` | [Kie documentation](https://docs.kie.ai/market/wan/2-2-animate-move.md) |
| Wan - Animate Replace | `wan/2-2-animate-replace` | [Kie documentation](https://docs.kie.ai/market/wan/2-2-animate-replace.md) |
| Wan 2.5 - Image to Video | `wan/2-5-image-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-5-image-to-video.md) |
| Wan 2.5 - Text to Video | `wan/2-5-text-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-5-text-to-video.md) |
| Wan - 2.6-flash-image-to-video | `wan/2-6-flash-image-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-6-flash-image-to-video.md) |
| Wan - 2-6-flash-video-to-video | `wan/2-6-flash-video-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-6-flash-video-to-video.md) |
| Wan 2.6 - Image to Video | `wan/2-6-image-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-6-image-to-video.md) |
| Wan 2.6 - Text to Video | `wan/2-6-text-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-6-text-to-video.md) |
| Wan 2.6 - Video to Video | `wan/2-6-video-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-6-video-to-video.md) |
| Wan 2.7 - Image to Video | `wan/2-7-image-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-image-to-video.md) |
| Wan 2.7 - Reference to Video | `wan/2-7-r2v` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-r2v.md) |
| Wan 2.7 - Text to Video | `wan/2-7-text-to-video` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-text-to-video.md) |
| Wan 2.7 - Video Edit | `wan/2-7-videoedit` | [Kie documentation](https://docs.kie.ai/market/wan/2-7-videoedit.md) |
