# Kie.ai Market Model Inputs and Settings

Complete field reference for **129 models**, generated from [Kie.ai's official documentation index](https://docs.kie.ai/llms.txt).
The embedded machine-readable source is `internal/kiecatalog/catalog.json`; use `kie-pp-cli models show` or the MCP `media_model_get` tool to avoid loading this whole document into an agent context.

Route compactly before loading a schema with `kie-pp-cli media capability list`
and `media capability show <model-id>`, or MCP `media_capability_list/get`. The
capability record includes the primary Kie skill route, production fit, and any
documented lowest faithful proof tier. This document remains the authoritative
per-field/settings reference; skills do not duplicate it.

Top-level task settings are `model`, `input`, and the optional `callBackUrl`. Tables below describe each model's `input` object.

## `bytedance/seedance-1.5-pro`

**Bytedance Seedance 1.5 Pro** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/seedance-1-5-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 20000 | The text prompt used to generate the video. Required field. (Min length: 3, Max length: 20000 characters) |
| `input_urls` | array | no | max items: 2 | URLs of input images for image-to-video generation. Optional field.  - Accepts 0-2 images - If not provided, the model will perform text-to-video generation - File URLs after upload, not file content - Accepted types: image/jpeg, image/png, image/webp - Max size per image: 10.0MB |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`; default: `"1:1"` | Video aspect ratio configuration. Required field. |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance, 1080p for higher quality |
| `duration` | number | yes |  | Duration of the video in seconds,Optional range 4-12 s |
| `fixed_lens` | boolean | no | default: `false` | Seedance adds dynamic camera movement. Enable this feature to lock the camera for stable, static shots.  - **true**: Lock camera for static shots - **false**: Allow dynamic camera movement |
| `generate_audio` | boolean | no | default: `false` | Whether to generate audio for the video.  - **true**: Generate with audio (higher cost) - **false**: Generate without audio  Note: Enabling audio will increase the generation cost |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A serene beach at sunset with waves gently crashing on the shore, palm trees swaying in the breeze, and seagulls flying across the orange sky",
  "input_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/example1.png"
  ],
  "aspect_ratio": "1:1",
  "resolution": "720p",
  "duration": 8,
  "fixed_lens": false,
  "generate_audio": false,
  "nsfw_checker": false
}
```

## `bytedance/seedance-2`

**Bytedance Seedance 2.0** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/seedance-2)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | min length: 3; max length: 20000 | The text prompt used to generate the video. Required field. (Min length: 3, Max length: 20000 characters) |
| `first_frame_url` | string | no |  | First frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |
| `last_frame_url` | string | no |  | End frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |
| `reference_image_urls` | array | no | max items: 9 | Enter a list of image URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single image requirements: Format: jpeg, png, webp, bmp, tiff, gif. Aspect ratio (width/height): (0.4, 2.5) Width and height (px): (300, 6000) Size: Single image less than 30 MB. Maximum number of files: The sum of the number of frames at the beginning and end must not exceed 9.. |
| `reference_video_urls` | array | no | max items: 3 | Enter a list of video URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj) . Single video requirements: Video format: mp4, mov. Resolution: 480p, 720p Duration: Single video duration [2, 15] s, maximum 3 reference videos, total duration of all videos not exceeding 15 seconds. Dimensions: Aspect ratio (width/height): [0.4, 2.5] Width/height (px): [300, 6000] Total pixels: [640×640=409600, 834×1112=927408], i.e., the product of width and height must meet the range requirement of [409600, 927408]. Size: Single video not exceeding 50 MB. Frame rate (FPS): [24, 60] |
| `reference_audio_urls` | array | no | max items: 3 | Enter a list of audio URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj) . Single audio requirements: Format: wav, mp3 Duration: Single audio duration [2, 15] s, maximum 3 reference audios, total duration of all audios not exceeding 15 s. Size: Single audio file size not exceeding 15 MB. |
| `return_last_frame` | boolean | no | default: `false` | Whether to return the last frame of the video as an image. |
| `generate_audio` | boolean | no | default: `true` | Whether to generate audio for the video.  - **true**: Generate with audio - **false**: Generate without audio |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`, `4k`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance, 1080p for High-quality video, 4K Ultra-High Definition, delivering perfect visual details. |
| `aspect_ratio` | string | no | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`, `adaptive`; default: `"16:9"` | Video aspect ratio configuration. Required field. |
| `duration` | integer | no | default: `5` | Video duration in 4-15 seconds. |
| `web_search` | boolean | no |  | Use online search |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A serene beach at sunset with waves gently crashing on the shore, palm trees swaying in the breeze, and seagulls flying across the orange sky",
  "first_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example2.png",
  "last_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example3.png",
  "reference_image_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.png"
  ],
  "reference_video_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp4"
  ],
  "reference_audio_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp3"
  ],
  "return_last_frame": false,
  "generate_audio": false,
  "resolution": "720p",
  "aspect_ratio": "16:9",
  "duration": 15,
  "web_search": false
}
```

## `bytedance/seedance-2-5`

**Bytedance Seedance 2.5** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/seedance-2-5)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 30000 | The text prompt used to generate the video.. ( Max length: 30000 characters) |
| `first_frame_url` | string | no |  | First frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) Cannot be used simultaneously with reference_image_urls, reference_video_urls, or reference_audio_urls. |
| `last_frame_url` | string | no |  | End frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) last_frame_url cannot be passed alone; first_frame_url must be provided together with it. |
| `reference_image_urls` | array | no | max items: 30 | Enter a list of image URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single image requirements: Format: jpeg, png, webp, bmp, tiff, gif. Aspect ratio (width/height): (0.4, 2.5) Width and height (px): (300, 6000) Size: Single image less than 30 MB. Maximum number of files: The sum of the number of frames at the beginning and end must not exceed 30. Mutually exclusive with the first/last-frame scenario; |
| `reference_video_urls` | array | no | max items: 10 | Enter a list of video URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj) . Single video requirements: Video format: mp4, mov. Resolution: 480p, 720p Dimensions: Aspect ratio (width/height): [0.4, 2.5] Width/height (px): [300, 6000] Total pixels: [640×640=409600, 834×1112=927408], i.e., the product of width and height must meet the range requirement of [409600, 927408]. Size: Single video not exceeding 200 MB. Frame rate (FPS): [24, 60] Single video duration: [2, 30] seconds;  Mutually exclusive with the first/last-frame scenario; Total duration of reference videos must not exceed 30 seconds. |
| `reference_audio_urls` | array | no | max items: 10 | Enter a list of audio URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj) . Single audio requirements: Format: wav, mp3 Size: Single audio file size not exceeding 15 MB. Single audio duration: [2, 30] seconds;  Mutually exclusive with the first/last-frame scenario; Total duration of reference videos must not exceed 30 seconds. |
| `return_last_frame` | boolean | no | default: `false` | Whether to return the last frame of the video. When draft=true, this parameter cannot be set to true. |
| `generate_audio` | boolean | no | default: `true` | Whether to generate audio for the video.  - **true**: Generate with audio (higher cost) - **false**: Generate without audio  Note: Enabling audio will increase the generation cost |
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance. |
| `aspect_ratio` | string | no | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`, `adaptive`; default: `"adaptive"` | Video aspect ratio configuration. |
| `duration` | integer | no | default: `5` | Video duration in 4-30 seconds. Special values -1:  Pass -1 for automatic duration selection — the model picks a suitable duration itself (for video-editing tasks, it matches the input video's length; for other task types, it selects within the valid range), instead of you specifying an exact value. |
| `output_format` | string | no | allowed: `mp4`, `mov`; default: `"mp4"` | Video output format. |
| `web_search` | boolean | no |  | Enable online search. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

## `bytedance/seedance-2-fast`

**Bytedance Seedance 2.0 Fast** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/seedance-2-fast)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance |
| `duration` | integer | no | default: `5` | Video duration in 4-15 seconds. |
| `generate_audio` | boolean | no | default: `true` | Whether to generate audio for the video.  - **true**: Generate with audio  - **false**: Generate without audio |
| `return_last_frame` | boolean | no | default: `false` | Whether to return the last frame of the video as an image. |
| `web_search` | boolean | no |  | Use online search. (Web search only can be used in the scene of t2v) |
| `aspect_ratio` | string | no | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`, `adaptive`; default: `"16:9"` | Video aspect ratio configuration. Required field. |
| `prompt` | string | no | min length: 3; max length: 20000 | The text prompt used to generate the video. Required field. (Min length: 3, Max length: 20000 characters) |
| `reference_image_urls` | array | no | max items: 9 | Enter a list of image URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single image requirements: Format: jpeg, png, webp, bmp, tiff, gif. Aspect ratio (width/height): (0.4, 2.5) Width and height (px): (300, 6000) Size: Single image less than 30 MB. Maximum number of files: The sum of the number of frames at the beginning and end must not exceed 9.. |
| `reference_video_urls` | array | no | max items: 3 | Enter a list of video URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single video requirements: Video format: mp4, mov. Resolution: 480p, 720p Duration: Single video duration [2, 15] s, maximum 3 reference videos, total duration of all videos not exceeding 15 seconds. Dimensions: Aspect ratio (width/height): [0.4, 2.5] Width/height (px): [300, 6000] Total pixels: [640×640=409600, 834×1112=927408], i.e., the product of width and height must meet the range requirement of [409600, 927408]. Size: Single video not exceeding 50 MB. Frame rate (FPS): [24, 60] |
| `reference_audio_urls` | array | no | max items: 3 | Enter a list of audio URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single audio requirements: Format: wav, mp3 Duration: Single audio duration [2, 15] s, maximum 3 reference audios, total duration of all audios not exceeding 15 s. Size: Single audio file size not exceeding 15 MB. |
| `first_frame_url` | string | no |  | First frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |
| `last_frame_url` | string | no |  | End frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |

Example `input`:

```json
{
  "prompt": "A serene beach at sunset with waves gently crashing on the shore, palm trees swaying in the breeze, and seagulls flying across the orange sky",
  "first_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example2.png",
  "last_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example3.png",
  "reference_image_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.png"
  ],
  "reference_video_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp4"
  ],
  "reference_audio_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp3"
  ],
  "return_last_frame": false,
  "generate_audio": false,
  "resolution": "720p",
  "aspect_ratio": "16:9",
  "duration": 15,
  "web_search": false
}
```

## `bytedance/seedance-2-mini`

**Bytedance Seedance 2.0 Mini** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/seedance-2-mini)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | min length: 3; max length: 20000 | The text prompt used to generate the video. Required field. (Min length: 3, Max length: 20000 characters) |
| `first_frame_url` | string | no |  | First frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |
| `last_frame_url` | string | no |  | End frame image url or asset://{assetId}  (for example: asset://asset-20260404242101-76djj) |
| `reference_image_urls` | array | no | max items: 9 | Enter a list of image URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single image requirements: Format: jpeg, png, webp, bmp, tiff, gif. Aspect ratio (width/height): (0.4, 2.5) Width and height (px): (300, 6000) Size: Single image less than 30 MB. Maximum number of files: The sum of the number of frames at the beginning and end must not exceed 9.. |
| `reference_video_urls` | array | no | max items: 3 | Enter a list of video URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single video requirements: Video format: mp4, mov. Resolution: 480p, 720p Duration: Single video duration [2, 15] s, maximum 3 reference videos, total duration of all videos not exceeding 15 seconds. Dimensions: Aspect ratio (width/height): [0.4, 2.5] Width/height (px): [300, 6000] Total pixels: [640×640=409600, 834×1112=927408], i.e., the product of width and height must meet the range requirement of [409600, 927408]. Size: Single video not exceeding 50 MB. Frame rate (FPS): [24, 60] |
| `reference_audio_urls` | array | no | max items: 3 | Enter a list of audio URLs or asset://{assetId} (for example: asset://asset-20260404242101-76djj). Single audio requirements: Format: wav, mp3 Duration: Single audio duration [2, 15] s, maximum 3 reference audios, total duration of all audios not exceeding 15 s. Size: Single audio file size not exceeding 15 MB. |
| `generate_audio` | boolean | no | default: `true` | Whether to generate audio for the video.  - **true**: Generate with audio  - **false**: Generate without audio |
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance |
| `aspect_ratio` | string | no | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`, `adaptive`; default: `"16:9"` | Video aspect ratio configuration. Required field. |
| `duration` | integer | no | default: `5` | Video duration in 4-15 seconds. |
| `web_search` | boolean | no |  | Use online search. (Web search only can be used in the scene of t2v) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A serene beach at sunset with waves gently crashing on the shore, palm trees swaying in the breeze, and seagulls flying across the orange sky",
  "first_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example2.png",
  "last_frame_url": "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example3.png",
  "reference_image_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.png"
  ],
  "reference_video_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp4"
  ],
  "reference_audio_urls": [
    "https://templateb.aiquickdraw.com/custom-page/akr/section-images/example1.mp3"
  ],
  "return_last_frame": false,
  "generate_audio": false,
  "resolution": "720p",
  "aspect_ratio": "16:9",
  "duration": 15,
  "web_search": false
}
```

## `bytedance/seedream`

**Seedream3.0 - Text to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/seedream)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt used to generate the image (Max length: 5000 characters) |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9`; default: `"square_hd"` | Select description |
| `guidance_scale` | number | no | default: `2.5`; min: 1; max: 10 | Controls how closely the output image aligns with the input prompt. Higher values mean stronger prompt correlation. (Min: 1, Max: 10, Step: 0.1) (step: 0.1) |
| `seed` | integer | no |  | Random seed to control the stochasticity of image generation. |

## `bytedance/seedream-v4-edit`

**Seedream4.0 - Edit** · Seedream · [official docs](https://docs.kie.ai/market/seedream/seedream-v4-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt used to edit the image (Max length: 5000 characters) |
| `image_urls` | array | yes | max items: 10 | List of URLs of input images for editing. Presently, up to 10 image inputs are allowed. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_3_2`, `portrait_16_9`, `landscape_4_3`, `landscape_3_2`, `landscape_16_9`, `landscape_21_9`; default: `"square_hd"` | The size of the generated image. |
| `image_resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"1K"` | Final image resolution is determined by combining image_size (aspect ratio) and image_resolution (pixel scale). For example, choosing 4:3 + 4K gives 4096 × 3072px |
| `max_images` | number | no | default: `1`; min: 1; max: 6 | Set this value (1–6) to cap how many images a single generation run can produce in one set—because they’re created in one shot rather than separate requests, you must also state the exact number you want in the prompt so both settings align. (Min: 1, Max: 6, Step: 1) (step: 1) |
| `seed` | integer | no |  | Random seed to control the stochasticity of image generation. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Refer to this logo and create a single visual showcase for an outdoor sports brand named ‘KIE AI’. Display five branded items together in one image: a packaging bag, a hat, a carton box, a wristband, and a lanyard. Use blue as the main visual color, with a fun, simple, and modern style.",
  "image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1757930552966e7f2on7s.png"
  ],
  "image_size": "square_hd",
  "image_resolution": "1K",
  "max_images": 1,
  "seed": 80960659,
  "nsfw_checker": true
}
```

## `bytedance/seedream-v4-text-to-image`

**Seedream4.0 - Text to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/seedream-v4-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt used to generate the image (Max length: 5000 characters) |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_3_2`, `portrait_16_9`, `landscape_4_3`, `landscape_3_2`, `landscape_16_9`, `landscape_21_9`; default: `"square_hd"` | The size of the generated image. |
| `image_resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"1K"` | Final image resolution is determined by combining image_size (aspect ratio) and image_resolution (pixel scale). For example, choosing 4:3 + 4K gives 4096 × 3072px |
| `max_images` | number | no | default: `1`; min: 1; max: 6 | Set this value (1–6) to cap how many images a single generation run can produce in one set—because they’re created in one shot rather than separate requests, you must also state the exact number you want in the prompt so both settings align. (Min: 1, Max: 6, Step: 1) (step: 1) |
| `seed` | integer | no |  | Random seed to control the stochasticity of image generation |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Draw the following system of binary linear equations and the corresponding solution steps on the blackboard: 5x + 2y = 26; 2x -y = 5.",
  "image_size": "square_hd",
  "image_resolution": "1K",
  "max_images": 1,
  "seed": 50331296,
  "nsfw_checker": true
}
```

## `bytedance/v1-lite-image-to-video`

**Bytedance - V1 Lite Image to Video** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/v1-lite-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | The text prompt used to generate the video (Max length: 10000 characters) |
| `image_url` | string | yes |  | The URL of the image used to generate video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for higher quality |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | Duration of the video in seconds |
| `camera_fixed` | boolean | no |  | Whether to fix the camera position (Boolean value (true/false)) |
| `seed` | number | no | default: `-1`; min: -1; max: 2147483647 | Random seed to control video generation. Use -1 for random. (Min: -1, Max: 2147483647, Step: 1) (step: 1) |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |
| `end_image_url` | string | no |  | The URL of the image the video ends with. Defaults to None. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |

Example `input`:

```json
{
  "prompt": "Multiple shots. A traveler crosses an endless desert toward a glowing archway. [Cut to] His cloak whips in the wind as he reaches the massive stone threshold. [Wide shot] He steps through — and vanishes into a burst of light",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17550783375205e9woshz.png",
  "resolution": "720p",
  "duration": "5",
  "camera_fixed": false,
  "seed": -1,
  "enable_safety_checker": true,
  "end_image_url": "",
  "nsfw_checker": false
}
```

## `bytedance/v1-lite-text-to-video`

**Bytedance - V1 Lite Text to Video** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/v1-lite-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | The text prompt used to generate the video (Max length: 10000 characters) |
| `aspect_ratio` | string | no | allowed: `16:9`, `4:3`, `1:1`, `3:4`, `9:16`, `9:21`; default: `"16:9"` | The aspect ratio of the generated video |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for higher quality |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | Duration of the video in seconds |
| `camera_fixed` | boolean | no |  | Whether to fix the camera position (Boolean value (true/false)) |
| `seed` | integer | no |  | Random seed to control video generation. Use -1 for random. |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Wide-angle shot: A serene sailing boat gently sways in the harbor at dawn, surrounded by soft Impressionist hues of pink and orange with ivory accents. The camera slowly pans across the scene, capturing the delicate reflections on the water and the intricate details of the boat's sails as the light gradually brightens.",
  "aspect_ratio": "16:9",
  "resolution": "720p",
  "duration": "5",
  "camera_fixed": false,
  "seed": 91466377,
  "enable_safety_checker": true,
  "nsfw_checker": false
}
```

## `bytedance/v1-pro-fast-image-to-video`

**Bytedance V1 Pro Fast Image to Video** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/v1-pro-fast-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | The text prompt used to generate the video (Max length: 10000 characters) |
| `image_url` | string | yes |  | The URL of the image used to generate video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance, 1080p for higher quality |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | Duration of the video in seconds |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A cinematic close-up sequence of a single elegant ceramic coffee cup with saucer on a rustic wooden table near a sunlit window, hot rich espresso poured in a thin golden stream from above, gradually filling the cup in distinct stages: empty with faint steam, 1/4 filled with dark crema, half-filled with swirling coffee and rising steam, 3/4 filled nearing the rim, perfectly full just below overflow with glossy surface and soft bokeh highlights; ultra-realistic, warm golden-hour light, shallow depth of field, photorealism, detailed textures, subtle steam wisps, serene inviting atmosphere --ar 16:9 --q 2 --style raw",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1762340693669m6sey187.webp",
  "resolution": "720p",
  "duration": "5",
  "nsfw_checker": true
}
```

## `bytedance/v1-pro-image-to-video`

**Bytedance V1 Pro Image to Video** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/v1-pro-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | The text prompt used to generate the video (Max length: 10000 characters) |
| `image_url` | string | yes |  | The URL of the image used to generate video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance, 1080p for higher quality |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | Duration of the video in seconds |
| `camera_fixed` | boolean | no |  | Whether to fix the camera position (Boolean value (true/false)) |
| `seed` | number | no | default: `-1`; min: -1; max: 2147483647 | Random seed to control video generation. Use -1 for random. (Min: -1, Max: 2147483647, Step: 1) (step: 1) |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |

Example `input`:

```json
{
  "prompt": "A golden retriever dashing through shallow surf at the beach, back angle camera low near waterline, splashes frozen in time, blur trails in waves and paws, afternoon sun glinting off wet fur, overcast day, dramatic clouds",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755179021328w1nhip18.webp",
  "resolution": "720p",
  "duration": "5",
  "camera_fixed": false,
  "seed": -1,
  "enable_safety_checker": true,
  "nsfw_checker": true
}
```

## `bytedance/v1-pro-text-to-video`

**Bytedance - V1 Pro Text to Video** · Bytedance · [official docs](https://docs.kie.ai/market/bytedance/v1-pro-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | The text prompt used to generate the video (Max length: 10000 characters) |
| `aspect_ratio` | string | no | allowed: `21:9`, `16:9`, `4:3`, `1:1`, `3:4`, `9:16`; default: `"16:9"` | The aspect ratio of the generated video |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"720p"` | Video resolution - 480p for faster generation, 720p for balance, 1080p for higher quality |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | Duration of the video in seconds |
| `camera_fixed` | boolean | no |  | Whether to fix the camera position (Boolean value (true/false)) |
| `seed` | number | no | default: `-1`; min: -1; max: 2147483647 | Random seed to control video generation. Use -1 for random. (Min: -1, Max: 2147483647, Step: 1) (step: 1) |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A boy with curly hair and a backpack rides a bike down a golden-lit rural road at sunset.\n[Cut to] He slows down and looks toward a field of tall grass.\n[Wide shot] His silhouette halts in the orange haze.",
  "aspect_ratio": "16:9",
  "resolution": "720p",
  "duration": "5",
  "camera_fixed": false,
  "seed": -1,
  "enable_safety_checker": true,
  "nsfw_checker": false
}
```

## `elevenlabs/audio-isolation`

**elevenlabs/audio-isolation** · Elevenlabs · [official docs](https://docs.kie.ai/market/elevenlabs/audio-isolation)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `audio_url` | string | yes |  | URL of the audio file to isolate voice from (File URL after upload, not file content; Accepted types: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/mp4, audio/ogg; Max size: 10.0MB) |

Example `input`:

```json
{
  "audio_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1756964657418ljw1jbzr.mp3"
}
```

## `elevenlabs/text-to-dialogue-v3`

**elevenlabs/text-to-dialogue-v3** · Elevenlabs · [official docs](https://docs.kie.ai/market/elevenlabs/text-to-dialogue-v3)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `dialogue` | array | yes |  | Array of dialogue items. Each item contains text and voice. The total character count of all text fields combined must not exceed 5000 characters. |
| `dialogue[].text` | string | yes |  | The dialogue text content. All text fields in the dialogue array combined must not exceed 5000 characters. |
| `dialogue[].voice` | string | yes | allowed: `EkK5I93UQWFDigLMpZcX`, `Z3R5wn05IrDiVCyEkUrK`, `NNl6r8mD7vthiJatiJt1`, `YOq2y2Up4RgXP2HyXjE5`, `B8gJV1IhpuegLxdpXFOE`, `2zRM7PkgwBPiau2jvVXc`, `1SM7GgM6IMuvQlz2BwM3`, `5l5f8iK3YPeGga21rQIX`, `scOwDtmlUjD3prqpp97I`, `NOpBlnGInO9m6vDvFkFC`, `BZgkqPqms7Kj9ulSkVzn`, `wo6udizrrtpIxWGp2qJk`, `gU0LNdkMOQCOrPrwtbee`, `DGzg6RaUqxGRTHSBjfgF`, `x70vRnQBMBu4FAYhjJbO`, `Sm1seazb4gs7RSlUVw7c`, `P1bg08DkjqiVEzOn76yG`, `qDuRKMlYmrm8trt5QyBn`, `qXpMhyvQqiRxWQs4qSSB`, `TX3LPaxmHKxFdv7VOQHJ`, `N2lVS1w4EtoT3dr4eOWO`, `FGY2WhTYpPnrIDTdsKH5`, `kPzsL2i3teMYv0FxEYQ6`, `UgBBYS2sOqTuMpoF3BR0`, `hpp4J3VqNfWAUOO0d1Us`, `nPczCjzI2devNBz1zQrb`, `uYXf8XasLslADfZ2MB4u`, `gs0tAILXbY5DNrJrsM6F`, `DTKMou8ccj1ZaWGBiotd`, `vBKc2FfBKJfcZNyEt1n6`, `DYkrAHD8iwork3YSUBbs`, `56AoDkrOh6qfVPDXZ7Pt`, `eR40ATw9ArzDf9h3v7t7`, `g6xIsTj2HwM6VR4iXFCw`, `lcMyyd2HUfFzxdCaC4Ta`, `6aDn1KB0hjpdcocrUkmq`, `Sq93GQT4X1lKDXsQcixO`, `flHkNRp1BlvT73UL6gyz`, `9yzdeviXkFddZ4Oz8Mok`, `pPdl9cQBQq4p6mRkZy2Z`, `zYcjlYFOd3taleS0gkk3`, `nzeAacJi50IvxcyDnMXa`, `ruirxsoakN0GWmGNIo04`, `TC0Zp7WVFzhA8zpTlRqV`, `ljo9gAlSqKOvF6D8sOsX`, `PPzYpIqttlTYA83688JI`, `8JVbfL6oEdmuxKn5DK2C`, `iCrDUkL56s3C8sCRl7wb`, `wJqPPQ618aTW29mptyoc`, `EiNlNiXeDU1pqqOPrYMO`, `4YYIPFl9wE5c4L2eu2Gb`, `6F5Zhi321D3Oq7v1oNT4`, `YXpFCvM1S3JbWEJhoskW`, `LG95yZDEHg6fCZdQjLqj`, `CeNX9CMwmxDxUF5Q2Inm`, `aD6riP1btT197c6dACmy`, `mtrellq69YZsNwzUSyXh`, `dHd5gvgSOzSfduK4CvEg`, `eVItLK1UvXctxuaRV2Oq`, `esy0r39YPLQjOczyOib8`, `Tsns2HvNFKfGiNjllgqo`, `1U02n4nD6AdIZ9CjF053`, `AeRdCCKzvd23BpJoofzx`, `LruHrtVF6PSyGItzMNHS`, `1wGbFxmAM3Fgw63G1zZJ`, `hqfrgApggtO1785R4Fsn`, `MJ0RnG71ty4LH3dvNfSd`; default: `"EkK5I93UQWFDigLMpZcX"` | The voice to use for speech generation. Can be a preset voice name (e.g. Rachel, Adam) or a voice ID. You can preview a voice by opening https://static.aiquickdraw.com/elevenlabs/voice/<voice_id>.mp3 in your browser (replace <voice_id> with the actual voice ID). For example: https://static.aiquickdraw.com/elevenlabs/voice/N2lVS1w4EtoT3dr4eOWO.mp3  Available voices:  EkK5I93UQWFDigLMpZcX - James - Husky, Engaging and Bold Z3R5wn05IrDiVCyEkUrK - Arabella - Mysterious and Emotive NNl6r8mD7vthiJatiJt1 - Bradford - Expressive and Articulate YOq2y2Up4RgXP2HyXjE5 - Xavier - Dominating, Metallic Announcer B8gJV1IhpuegLxdpXFOE - Kuon - Cheerful, Clear and Steady 2zRM7PkgwBPiau2jvVXc - Monika Sogam - Deep and Natural 1SM7GgM6IMuvQlz2BwM3 - Mark - Casual, Relaxed and Light 5l5f8iK3YPeGga21rQIX - Adeline - Feminine and Conversational scOwDtmlUjD3prqpp97I - Sam - Support Agent NOpBlnGInO9m6vDvFkFC - Spuds Oxley - Wise and Approachable BZgkqPqms7Kj9ulSkVzn - Eve - Authentic, Energetic and Happy wo6udizrrtpIxWGp2qJk - Northern Terry gU0LNdkMOQCOrPrwtbee - British Football Announcer DGzg6RaUqxGRTHSBjfgF - Brock - Commanding and Loud Sergeant x70vRnQBMBu4FAYhjJbO - Nathan - Virtual Radio Host Sm1seazb4gs7RSlUVw7c - Anika - Animated, Friendly and Engaging P1bg08DkjqiVEzOn76yG - Viraj - Rich and Soft qDuRKMlYmrm8trt5QyBn - Taksh - Calm, Serious and Smooth qXpMhyvQqiRxWQs4qSSB - Horatius - Energetic Character Voice TX3LPaxmHKxFdv7VOQHJ - Liam - Energetic, Social Media Creator N2lVS1w4EtoT3dr4eOWO - Callum - Husky Trickster FGY2WhTYpPnrIDTdsKH5 - Laura - Enthusiast, Quirky Attitude kPzsL2i3teMYv0FxEYQ6 - Brittney - Social Media Voice - Fun, Youthful & Informative UgBBYS2sOqTuMpoF3BR0 - Mark - Natural Conversations hpp4J3VqNfWAUOO0d1Us - Bella - Professional, Bright, Warm nPczCjzI2devNBz1zQrb - Brian - Deep, Resonant and Comforting uYXf8XasLslADfZ2MB4u - Hope - Bubbly, Gossipy and Girly gs0tAILXbY5DNrJrsM6F - Jeff - Classy, Resonating and Strong DTKMou8ccj1ZaWGBiotd - Jamahal - Young, Vibrant, and Natural vBKc2FfBKJfcZNyEt1n6 - Finn - Youthful, Eager and Energetic DYkrAHD8iwork3YSUBbs - Tom - Conversations & Books 56AoDkrOh6qfVPDXZ7Pt - Cassidy - Crisp, Direct and Clear eR40ATw9ArzDf9h3v7t7 - Addison 2.0 - Australian Audiobook & Podcast g6xIsTj2HwM6VR4iXFCw - Jessica Anne Bogart - Chatty and Friendly lcMyyd2HUfFzxdCaC4Ta - Lucy - Fresh & Casual 6aDn1KB0hjpdcocrUkmq - Tiffany - Natural and Welcoming Sq93GQT4X1lKDXsQcixO - Felix - Warm, Positive & Contemporary RP flHkNRp1BlvT73UL6gyz - Jessica Anne Bogart - Eloquent Villain 9yzdeviXkFddZ4Oz8Mok - Lutz - Chuckling, Giggly and Cheerful pPdl9cQBQq4p6mRkZy2Z - Emma - Adorable and Upbeat zYcjlYFOd3taleS0gkk3 - Edward - Loud, Confident and Cocky nzeAacJi50IvxcyDnMXa - Marshal - Friendly, Funny Professor ruirxsoakN0GWmGNIo04 - John Morgan - Gritty, Rugged Cowboy TC0Zp7WVFzhA8zpTlRqV - Aria - Sultry Villain ljo9gAlSqKOvF6D8sOsX - Viking Bjorn - Epic Medieval Raider PPzYpIqttlTYA83688JI - Pirate Marshal 8JVbfL6oEdmuxKn5DK2C - Johnny Kid - Serious and Calm Narrator iCrDUkL56s3C8sCRl7wb - Hope - Poetic, Romantic and Captivating wJqPPQ618aTW29mptyoc - Ana Rita - Smooth, Expressive and Bright EiNlNiXeDU1pqqOPrYMO - John Doe - Deep 4YYIPFl9wE5c4L2eu2Gb - Burt Reynolds™ - Deep, Smooth and Clear 6F5Zhi321D3Oq7v1oNT4 - Hank - Deep and Engaging Narrator YXpFCvM1S3JbWEJhoskW - Wyatt - Wise Rustic Cowboy LG95yZDEHg6fCZdQjLqj - Phil - Explosive, Passionate Announcer CeNX9CMwmxDxUF5Q2Inm - Johnny Dynamite - Vintage Radio DJ aD6riP1btT197c6dACmy - Rachel M - Pro British Radio Presenter mtrellq69YZsNwzUSyXh - Rex Thunder - Deep N Tough dHd5gvgSOzSfduK4CvEg - Ed - Late Night Announcer eVItLK1UvXctxuaRV2Oq - Jean - Alluring and Playful Femme Fatale esy0r39YPLQjOczyOib8 - Britney - Calm and Calculative Villain Tsns2HvNFKfGiNjllgqo - Sven - Emotional and Nice 1U02n4nD6AdIZ9CjF053 - Viraj - Smooth and Gentle AeRdCCKzvd23BpJoofzx - Nathaniel - Engaging, British and Calm LruHrtVF6PSyGItzMNHS - Benjamin - Deep, Warm, Calming 1wGbFxmAM3Fgw63G1zZJ - Allison - Calm, Soothing and Meditative hqfrgApggtO1785R4Fsn - Theodore HQ - Serene and Grounded MJ0RnG71ty4LH3dvNfSd - Leon - Soothing and Grounded |
| `stability` | number | no | allowed: `0`, `0.5`, `1`; default: `0.5` | Voice stability parameter. Must be one of the following values: 0.0, 0.5, or 1.0 |
| `language_code` | string | no | allowed: `af`, `ar`, `hy`, `as`, `az`, `be`, `bn`, `bs`, `bg`, `ca`, `ceb`, `ny`, `hr`, `cs`, `da`, `nl`, `en`, `et`, `fil`, `fi`, `fr`, `gl`, `ka`, `de`, `el`, `gu`, `ha`, `he`, `hi`, `hu`, `is`, `id`, `ga`, `it`, `ja`, `jv`, `kn`, `kk`, `ky`, `ko`, `lv`, `ln`, `lt`, `lb`, `mk`, `ms`, `ml`, `zh`, `mr`, `ne`, `no`, `ps`, `fa`, `pl`, `pt`, `pa`, `ro`, `ru`, `sr`, `sd`, `sk`, `sl`, `so`, `es`, `sw`, `sv`, `ta`, `te`, `th`, `tr`, `uk`, `ur`, `vi`, `cy` | Language code for the speech. Default is empty string or omit the parameter for automatic language detection. |

Example `input`:

```json
{
  "dialogue": [
    {
      "text": "I have a pen, I have an apple, ah, Apple pen~",
      "voice": "EkK5I93UQWFDigLMpZcX"
    },
    {
      "text": "a happy dog",
      "voice": "Z3R5wn05IrDiVCyEkUrK"
    },
    {
      "text": "a happy cat",
      "voice": "NNl6r8mD7vthiJatiJt1"
    }
  ],
  "stability": 0.5
}
```

## `elevenlabs/text-to-speech-multilingual-v2`

**elevenlabs/text-to-speech-multilingual-v2** · Elevenlabs · [official docs](https://docs.kie.ai/market/elevenlabs/text-to-speech-multilingual-v2)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `text` | string | yes | max length: 5000 | The text to convert to speech (Max length: 5000 characters) |
| `voice` | string | yes | allowed: `EkK5I93UQWFDigLMpZcX`, `Z3R5wn05IrDiVCyEkUrK`, `NNl6r8mD7vthiJatiJt1`, `YOq2y2Up4RgXP2HyXjE5`, `B8gJV1IhpuegLxdpXFOE`, `2zRM7PkgwBPiau2jvVXc`, `1SM7GgM6IMuvQlz2BwM3`, `5l5f8iK3YPeGga21rQIX`, `scOwDtmlUjD3prqpp97I`, `NOpBlnGInO9m6vDvFkFC`, `BZgkqPqms7Kj9ulSkVzn`, `wo6udizrrtpIxWGp2qJk`, `gU0LNdkMOQCOrPrwtbee`, `DGzg6RaUqxGRTHSBjfgF`, `x70vRnQBMBu4FAYhjJbO`, `Sm1seazb4gs7RSlUVw7c`, `P1bg08DkjqiVEzOn76yG`, `qDuRKMlYmrm8trt5QyBn`, `qXpMhyvQqiRxWQs4qSSB`, `TX3LPaxmHKxFdv7VOQHJ`, `N2lVS1w4EtoT3dr4eOWO`, `FGY2WhTYpPnrIDTdsKH5`, `kPzsL2i3teMYv0FxEYQ6`, `UgBBYS2sOqTuMpoF3BR0`, `hpp4J3VqNfWAUOO0d1Us`, `nPczCjzI2devNBz1zQrb`, `uYXf8XasLslADfZ2MB4u`, `gs0tAILXbY5DNrJrsM6F`, `DTKMou8ccj1ZaWGBiotd`, `vBKc2FfBKJfcZNyEt1n6`, `DYkrAHD8iwork3YSUBbs`, `56AoDkrOh6qfVPDXZ7Pt`, `eR40ATw9ArzDf9h3v7t7`, `g6xIsTj2HwM6VR4iXFCw`, `lcMyyd2HUfFzxdCaC4Ta`, `6aDn1KB0hjpdcocrUkmq`, `Sq93GQT4X1lKDXsQcixO`, `flHkNRp1BlvT73UL6gyz`, `9yzdeviXkFddZ4Oz8Mok`, `pPdl9cQBQq4p6mRkZy2Z`, `zYcjlYFOd3taleS0gkk3`, `nzeAacJi50IvxcyDnMXa`, `ruirxsoakN0GWmGNIo04`, `TC0Zp7WVFzhA8zpTlRqV`, `ljo9gAlSqKOvF6D8sOsX`, `PPzYpIqttlTYA83688JI`, `8JVbfL6oEdmuxKn5DK2C`, `iCrDUkL56s3C8sCRl7wb`, `wJqPPQ618aTW29mptyoc`, `EiNlNiXeDU1pqqOPrYMO`, `4YYIPFl9wE5c4L2eu2Gb`, `6F5Zhi321D3Oq7v1oNT4`, `YXpFCvM1S3JbWEJhoskW`, `LG95yZDEHg6fCZdQjLqj`, `CeNX9CMwmxDxUF5Q2Inm`, `aD6riP1btT197c6dACmy`, `mtrellq69YZsNwzUSyXh`, `dHd5gvgSOzSfduK4CvEg`, `eVItLK1UvXctxuaRV2Oq`, `esy0r39YPLQjOczyOib8`, `Tsns2HvNFKfGiNjllgqo`, `1U02n4nD6AdIZ9CjF053`, `AeRdCCKzvd23BpJoofzx`, `LruHrtVF6PSyGItzMNHS`, `1wGbFxmAM3Fgw63G1zZJ`, `hqfrgApggtO1785R4Fsn`, `MJ0RnG71ty4LH3dvNfSd`; default: `"EkK5I93UQWFDigLMpZcX"` | The voice to use for speech generation. Can be a preset voice name (e.g. Rachel, Adam) or a voice ID. You can preview a voice by opening https://static.aiquickdraw.com/elevenlabs/voice/<voice_id>.mp3 in your browser (replace <voice_id> with the actual voice ID). For example: https://static.aiquickdraw.com/elevenlabs/voice/N2lVS1w4EtoT3dr4eOWO.mp3  Available voices:  EkK5I93UQWFDigLMpZcX - James - Husky, Engaging and Bold Z3R5wn05IrDiVCyEkUrK - Arabella - Mysterious and Emotive NNl6r8mD7vthiJatiJt1 - Bradford - Expressive and Articulate YOq2y2Up4RgXP2HyXjE5 - Xavier - Dominating, Metallic Announcer B8gJV1IhpuegLxdpXFOE - Kuon - Cheerful, Clear and Steady 2zRM7PkgwBPiau2jvVXc - Monika Sogam - Deep and Natural 1SM7GgM6IMuvQlz2BwM3 - Mark - Casual, Relaxed and Light 5l5f8iK3YPeGga21rQIX - Adeline - Feminine and Conversational scOwDtmlUjD3prqpp97I - Sam - Support Agent NOpBlnGInO9m6vDvFkFC - Spuds Oxley - Wise and Approachable BZgkqPqms7Kj9ulSkVzn - Eve - Authentic, Energetic and Happy wo6udizrrtpIxWGp2qJk - Northern Terry gU0LNdkMOQCOrPrwtbee - British Football Announcer DGzg6RaUqxGRTHSBjfgF - Brock - Commanding and Loud Sergeant x70vRnQBMBu4FAYhjJbO - Nathan - Virtual Radio Host Sm1seazb4gs7RSlUVw7c - Anika - Animated, Friendly and Engaging P1bg08DkjqiVEzOn76yG - Viraj - Rich and Soft qDuRKMlYmrm8trt5QyBn - Taksh - Calm, Serious and Smooth qXpMhyvQqiRxWQs4qSSB - Horatius - Energetic Character Voice TX3LPaxmHKxFdv7VOQHJ - Liam - Energetic, Social Media Creator N2lVS1w4EtoT3dr4eOWO - Callum - Husky Trickster FGY2WhTYpPnrIDTdsKH5 - Laura - Enthusiast, Quirky Attitude kPzsL2i3teMYv0FxEYQ6 - Brittney - Social Media Voice - Fun, Youthful & Informative UgBBYS2sOqTuMpoF3BR0 - Mark - Natural Conversations hpp4J3VqNfWAUOO0d1Us - Bella - Professional, Bright, Warm nPczCjzI2devNBz1zQrb - Brian - Deep, Resonant and Comforting uYXf8XasLslADfZ2MB4u - Hope - Bubbly, Gossipy and Girly gs0tAILXbY5DNrJrsM6F - Jeff - Classy, Resonating and Strong DTKMou8ccj1ZaWGBiotd - Jamahal - Young, Vibrant, and Natural vBKc2FfBKJfcZNyEt1n6 - Finn - Youthful, Eager and Energetic DYkrAHD8iwork3YSUBbs - Tom - Conversations & Books 56AoDkrOh6qfVPDXZ7Pt - Cassidy - Crisp, Direct and Clear eR40ATw9ArzDf9h3v7t7 - Addison 2.0 - Australian Audiobook & Podcast g6xIsTj2HwM6VR4iXFCw - Jessica Anne Bogart - Chatty and Friendly lcMyyd2HUfFzxdCaC4Ta - Lucy - Fresh & Casual 6aDn1KB0hjpdcocrUkmq - Tiffany - Natural and Welcoming Sq93GQT4X1lKDXsQcixO - Felix - Warm, Positive & Contemporary RP flHkNRp1BlvT73UL6gyz - Jessica Anne Bogart - Eloquent Villain 9yzdeviXkFddZ4Oz8Mok - Lutz - Chuckling, Giggly and Cheerful pPdl9cQBQq4p6mRkZy2Z - Emma - Adorable and Upbeat zYcjlYFOd3taleS0gkk3 - Edward - Loud, Confident and Cocky nzeAacJi50IvxcyDnMXa - Marshal - Friendly, Funny Professor ruirxsoakN0GWmGNIo04 - John Morgan - Gritty, Rugged Cowboy TC0Zp7WVFzhA8zpTlRqV - Aria - Sultry Villain ljo9gAlSqKOvF6D8sOsX - Viking Bjorn - Epic Medieval Raider PPzYpIqttlTYA83688JI - Pirate Marshal 8JVbfL6oEdmuxKn5DK2C - Johnny Kid - Serious and Calm Narrator iCrDUkL56s3C8sCRl7wb - Hope - Poetic, Romantic and Captivating wJqPPQ618aTW29mptyoc - Ana Rita - Smooth, Expressive and Bright EiNlNiXeDU1pqqOPrYMO - John Doe - Deep 4YYIPFl9wE5c4L2eu2Gb - Burt Reynolds™ - Deep, Smooth and Clear 6F5Zhi321D3Oq7v1oNT4 - Hank - Deep and Engaging Narrator YXpFCvM1S3JbWEJhoskW - Wyatt - Wise Rustic Cowboy LG95yZDEHg6fCZdQjLqj - Phil - Explosive, Passionate Announcer CeNX9CMwmxDxUF5Q2Inm - Johnny Dynamite - Vintage Radio DJ aD6riP1btT197c6dACmy - Rachel M - Pro British Radio Presenter mtrellq69YZsNwzUSyXh - Rex Thunder - Deep N Tough dHd5gvgSOzSfduK4CvEg - Ed - Late Night Announcer eVItLK1UvXctxuaRV2Oq - Jean - Alluring and Playful Femme Fatale esy0r39YPLQjOczyOib8 - Britney - Calm and Calculative Villain Tsns2HvNFKfGiNjllgqo - Sven - Emotional and Nice 1U02n4nD6AdIZ9CjF053 - Viraj - Smooth and Gentle AeRdCCKzvd23BpJoofzx - Nathaniel - Engaging, British and Calm LruHrtVF6PSyGItzMNHS - Benjamin - Deep, Warm, Calming 1wGbFxmAM3Fgw63G1zZJ - Allison - Calm, Soothing and Meditative hqfrgApggtO1785R4Fsn - Theodore HQ - Serene and Grounded MJ0RnG71ty4LH3dvNfSd - Leon - Soothing and Grounded |
| `stability` | number | no | default: `0.5`; min: 0; max: 1 | Voice stability (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `similarity_boost` | number | no | default: `0.75`; min: 0; max: 1 | Similarity boost (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `style` | number | no | default: `0`; min: 0; max: 1 | Style exaggeration (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `speed` | number | no | default: `1`; min: 0.7; max: 1.2 | Speech speed (0.7-1.2). Values below 1.0 slow down the speech, above 1.0 speed it up. Extreme values may affect quality. (Min: 0.7, Max: 1.2, Step: 0.01) (step: 0.01) |
| `timestamps` | boolean | no |  | Whether to return timestamps for each word in the generated speech (Boolean value (true/false)) |
| `previous_text` | string | no | max length: 5000 | The text that came before the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation. (Max length: 5000 characters) |
| `next_text` | string | no | max length: 5000 | The text that comes after the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation. (Max length: 5000 characters) |
| `language_code` | string | no | max length: 500 | Language code (ISO 639-1) used to enforce a language for the model. Currently only Turbo v2.5 and Flash v2.5 support language enforcement. For other models, an error will be returned if language code is provided. (Max length: 500 characters) |

Example `input`:

```json
{
  "text": "Unlock powerful API with Kie.ai! Affordable, scalable APl integration, free trial playground, and secure, reliable performance.",
  "voice": "Rachel",
  "stability": 0.5,
  "similarity_boost": 0.75,
  "style": 0,
  "speed": 1,
  "timestamps": false,
  "previous_text": "",
  "next_text": "",
  "language_code": ""
}
```

## `elevenlabs/text-to-speech-turbo-2-5`

**elevenlabs/text-to-speech-turbo-2-5** · Elevenlabs · [official docs](https://docs.kie.ai/market/elevenlabs/text-to-speech-turbo-2-5)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `text` | string | yes | max length: 5000 | The text to convert to speech (Max length: 5000 characters) |
| `voice` | string | no | allowed: `EkK5I93UQWFDigLMpZcX`, `Z3R5wn05IrDiVCyEkUrK`, `NNl6r8mD7vthiJatiJt1`, `YOq2y2Up4RgXP2HyXjE5`, `B8gJV1IhpuegLxdpXFOE`, `2zRM7PkgwBPiau2jvVXc`, `1SM7GgM6IMuvQlz2BwM3`, `5l5f8iK3YPeGga21rQIX`, `scOwDtmlUjD3prqpp97I`, `NOpBlnGInO9m6vDvFkFC`, `BZgkqPqms7Kj9ulSkVzn`, `wo6udizrrtpIxWGp2qJk`, `gU0LNdkMOQCOrPrwtbee`, `DGzg6RaUqxGRTHSBjfgF`, `x70vRnQBMBu4FAYhjJbO`, `Sm1seazb4gs7RSlUVw7c`, `P1bg08DkjqiVEzOn76yG`, `qDuRKMlYmrm8trt5QyBn`, `qXpMhyvQqiRxWQs4qSSB`, `TX3LPaxmHKxFdv7VOQHJ`, `N2lVS1w4EtoT3dr4eOWO`, `FGY2WhTYpPnrIDTdsKH5`, `kPzsL2i3teMYv0FxEYQ6`, `UgBBYS2sOqTuMpoF3BR0`, `hpp4J3VqNfWAUOO0d1Us`, `nPczCjzI2devNBz1zQrb`, `uYXf8XasLslADfZ2MB4u`, `gs0tAILXbY5DNrJrsM6F`, `DTKMou8ccj1ZaWGBiotd`, `vBKc2FfBKJfcZNyEt1n6`, `DYkrAHD8iwork3YSUBbs`, `56AoDkrOh6qfVPDXZ7Pt`, `eR40ATw9ArzDf9h3v7t7`, `g6xIsTj2HwM6VR4iXFCw`, `lcMyyd2HUfFzxdCaC4Ta`, `6aDn1KB0hjpdcocrUkmq`, `Sq93GQT4X1lKDXsQcixO`, `flHkNRp1BlvT73UL6gyz`, `9yzdeviXkFddZ4Oz8Mok`, `pPdl9cQBQq4p6mRkZy2Z`, `zYcjlYFOd3taleS0gkk3`, `nzeAacJi50IvxcyDnMXa`, `ruirxsoakN0GWmGNIo04`, `TC0Zp7WVFzhA8zpTlRqV`, `ljo9gAlSqKOvF6D8sOsX`, `PPzYpIqttlTYA83688JI`, `8JVbfL6oEdmuxKn5DK2C`, `iCrDUkL56s3C8sCRl7wb`, `wJqPPQ618aTW29mptyoc`, `EiNlNiXeDU1pqqOPrYMO`, `4YYIPFl9wE5c4L2eu2Gb`, `6F5Zhi321D3Oq7v1oNT4`, `YXpFCvM1S3JbWEJhoskW`, `LG95yZDEHg6fCZdQjLqj`, `CeNX9CMwmxDxUF5Q2Inm`, `aD6riP1btT197c6dACmy`, `mtrellq69YZsNwzUSyXh`, `dHd5gvgSOzSfduK4CvEg`, `eVItLK1UvXctxuaRV2Oq`, `esy0r39YPLQjOczyOib8`, `Tsns2HvNFKfGiNjllgqo`, `1U02n4nD6AdIZ9CjF053`, `AeRdCCKzvd23BpJoofzx`, `LruHrtVF6PSyGItzMNHS`, `1wGbFxmAM3Fgw63G1zZJ`, `hqfrgApggtO1785R4Fsn`, `MJ0RnG71ty4LH3dvNfSd`; default: `"EkK5I93UQWFDigLMpZcX"` | The voice to use for speech generation. Can be a preset voice name (e.g. Rachel, Adam) or a voice ID. You can preview a voice by opening https://static.aiquickdraw.com/elevenlabs/voice/<voice_id>.mp3 in your browser (replace <voice_id> with the actual voice ID). For example: https://static.aiquickdraw.com/elevenlabs/voice/N2lVS1w4EtoT3dr4eOWO.mp3  Available voices:  EkK5I93UQWFDigLMpZcX - James - Husky, Engaging and Bold Z3R5wn05IrDiVCyEkUrK - Arabella - Mysterious and Emotive NNl6r8mD7vthiJatiJt1 - Bradford - Expressive and Articulate YOq2y2Up4RgXP2HyXjE5 - Xavier - Dominating, Metallic Announcer B8gJV1IhpuegLxdpXFOE - Kuon - Cheerful, Clear and Steady 2zRM7PkgwBPiau2jvVXc - Monika Sogam - Deep and Natural 1SM7GgM6IMuvQlz2BwM3 - Mark - Casual, Relaxed and Light 5l5f8iK3YPeGga21rQIX - Adeline - Feminine and Conversational scOwDtmlUjD3prqpp97I - Sam - Support Agent NOpBlnGInO9m6vDvFkFC - Spuds Oxley - Wise and Approachable BZgkqPqms7Kj9ulSkVzn - Eve - Authentic, Energetic and Happy wo6udizrrtpIxWGp2qJk - Northern Terry gU0LNdkMOQCOrPrwtbee - British Football Announcer DGzg6RaUqxGRTHSBjfgF - Brock - Commanding and Loud Sergeant x70vRnQBMBu4FAYhjJbO - Nathan - Virtual Radio Host Sm1seazb4gs7RSlUVw7c - Anika - Animated, Friendly and Engaging P1bg08DkjqiVEzOn76yG - Viraj - Rich and Soft qDuRKMlYmrm8trt5QyBn - Taksh - Calm, Serious and Smooth qXpMhyvQqiRxWQs4qSSB - Horatius - Energetic Character Voice TX3LPaxmHKxFdv7VOQHJ - Liam - Energetic, Social Media Creator N2lVS1w4EtoT3dr4eOWO - Callum - Husky Trickster FGY2WhTYpPnrIDTdsKH5 - Laura - Enthusiast, Quirky Attitude kPzsL2i3teMYv0FxEYQ6 - Brittney - Social Media Voice - Fun, Youthful & Informative UgBBYS2sOqTuMpoF3BR0 - Mark - Natural Conversations hpp4J3VqNfWAUOO0d1Us - Bella - Professional, Bright, Warm nPczCjzI2devNBz1zQrb - Brian - Deep, Resonant and Comforting uYXf8XasLslADfZ2MB4u - Hope - Bubbly, Gossipy and Girly gs0tAILXbY5DNrJrsM6F - Jeff - Classy, Resonating and Strong DTKMou8ccj1ZaWGBiotd - Jamahal - Young, Vibrant, and Natural vBKc2FfBKJfcZNyEt1n6 - Finn - Youthful, Eager and Energetic DYkrAHD8iwork3YSUBbs - Tom - Conversations & Books 56AoDkrOh6qfVPDXZ7Pt - Cassidy - Crisp, Direct and Clear eR40ATw9ArzDf9h3v7t7 - Addison 2.0 - Australian Audiobook & Podcast g6xIsTj2HwM6VR4iXFCw - Jessica Anne Bogart - Chatty and Friendly lcMyyd2HUfFzxdCaC4Ta - Lucy - Fresh & Casual 6aDn1KB0hjpdcocrUkmq - Tiffany - Natural and Welcoming Sq93GQT4X1lKDXsQcixO - Felix - Warm, Positive & Contemporary RP flHkNRp1BlvT73UL6gyz - Jessica Anne Bogart - Eloquent Villain 9yzdeviXkFddZ4Oz8Mok - Lutz - Chuckling, Giggly and Cheerful pPdl9cQBQq4p6mRkZy2Z - Emma - Adorable and Upbeat zYcjlYFOd3taleS0gkk3 - Edward - Loud, Confident and Cocky nzeAacJi50IvxcyDnMXa - Marshal - Friendly, Funny Professor ruirxsoakN0GWmGNIo04 - John Morgan - Gritty, Rugged Cowboy TC0Zp7WVFzhA8zpTlRqV - Aria - Sultry Villain ljo9gAlSqKOvF6D8sOsX - Viking Bjorn - Epic Medieval Raider PPzYpIqttlTYA83688JI - Pirate Marshal 8JVbfL6oEdmuxKn5DK2C - Johnny Kid - Serious and Calm Narrator iCrDUkL56s3C8sCRl7wb - Hope - Poetic, Romantic and Captivating wJqPPQ618aTW29mptyoc - Ana Rita - Smooth, Expressive and Bright EiNlNiXeDU1pqqOPrYMO - John Doe - Deep 4YYIPFl9wE5c4L2eu2Gb - Burt Reynolds™ - Deep, Smooth and Clear 6F5Zhi321D3Oq7v1oNT4 - Hank - Deep and Engaging Narrator YXpFCvM1S3JbWEJhoskW - Wyatt - Wise Rustic Cowboy LG95yZDEHg6fCZdQjLqj - Phil - Explosive, Passionate Announcer CeNX9CMwmxDxUF5Q2Inm - Johnny Dynamite - Vintage Radio DJ aD6riP1btT197c6dACmy - Rachel M - Pro British Radio Presenter mtrellq69YZsNwzUSyXh - Rex Thunder - Deep N Tough dHd5gvgSOzSfduK4CvEg - Ed - Late Night Announcer eVItLK1UvXctxuaRV2Oq - Jean - Alluring and Playful Femme Fatale esy0r39YPLQjOczyOib8 - Britney - Calm and Calculative Villain Tsns2HvNFKfGiNjllgqo - Sven - Emotional and Nice 1U02n4nD6AdIZ9CjF053 - Viraj - Smooth and Gentle AeRdCCKzvd23BpJoofzx - Nathaniel - Engaging, British and Calm LruHrtVF6PSyGItzMNHS - Benjamin - Deep, Warm, Calming 1wGbFxmAM3Fgw63G1zZJ - Allison - Calm, Soothing and Meditative hqfrgApggtO1785R4Fsn - Theodore HQ - Serene and Grounded MJ0RnG71ty4LH3dvNfSd - Leon - Soothing and Grounded |
| `stability` | number | no | default: `0.5`; min: 0; max: 1 | Voice stability (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `similarity_boost` | number | no | default: `0.75`; min: 0; max: 1 | Similarity boost (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `style` | number | no | default: `0`; min: 0; max: 1 | Style exaggeration (0-1) (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `speed` | number | no | default: `1`; min: 0.7; max: 1.2 | Speech speed (0.7-1.2). Values below 1.0 slow down the speech, above 1.0 speed it up. Extreme values may affect quality. (Min: 0.7, Max: 1.2, Step: 0.01) (step: 0.01) |
| `timestamps` | boolean | no |  | Whether to return timestamps for each word in the generated speech (Boolean value (true/false)) |
| `previous_text` | string | no | max length: 5000 | The text that came before the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation. (Max length: 5000 characters) |
| `next_text` | string | no | max length: 5000 | The text that comes after the text of the current request. Can be used to improve the speech's continuity when concatenating together multiple generations or to influence the speech's continuity in the current generation. (Max length: 5000 characters) |
| `language_code` | string | no | max length: 500 | Language code (ISO 639-1) used to enforce a language for the model. Currently only Turbo v2.5 and Flash v2.5 support language enforcement. For other models, an error will be returned if language code is provided. (Max length: 500 characters) |

Example `input`:

```json
{
  "text": "Unlock powerful API with Kie.ai! Affordable, scalable APl integration, free trial playground, and secure, reliable performance.",
  "voice": "Rachel",
  "stability": 0.5,
  "similarity_boost": 0.75,
  "style": 0,
  "speed": 1,
  "timestamps": false,
  "previous_text": "",
  "next_text": "",
  "language_code": ""
}
```

## `flux-2/flex-image-to-image`

**Flux-2 - Image to Image** · Flux2 · [official docs](https://docs.kie.ai/market/flux2/flex-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `input_urls` | array | yes | max items: 8 | Input reference images (1-8 images). (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `prompt` | string | yes | max length: 5000 | Must be between 3 and 5000 characters. (Max length: 5000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `3:2`, `2:3`, `auto`; default: `"1:1"` | Aspect ratio for the generated image. Select 'auto' to match the first input image ratio (requires input image). |
| `resolution` | string | yes | allowed: `1K`, `2K`; default: `"1K"` | Output image resolution. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1764235158281_tABmx723.png",
    "https://static.aiquickdraw.com/tools/example/1764235165079_8fIR5MEF.png"
  ],
  "prompt": "Replace the can in image 2 with the can from image 1",
  "aspect_ratio": "1:1",
  "resolution": "1K",
  "nsfw_checker": false
}
```

## `flux-2/flex-text-to-image`

**Flux-2 - Text to Image** · Flux2 · [official docs](https://docs.kie.ai/market/flux2/flex-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Generation prompt, length must be between 3-5000 characters. (Maximum length: 5000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `3:2`, `2:3`; default: `"1:1"` | Aspect ratio of the generated image. When `auto` is selected, it will match the ratio of the first input image (requires input image to be provided). |
| `resolution` | string | yes | allowed: `1K`, `2K`; default: `"1K"` | Output image resolution. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A humanoid figure with a vintage television set for a head, featuring a green-tinted screen displaying a `Hello FLUX.2` writing in ASCII font. The figure is wearing a yellow raincoat, and there are various wires and components attached to the television. The background is cloudy and indistinct, suggesting an outdoor setting",
  "aspect_ratio": "1:1",
  "resolution": "1K",
  "nsfw_checker": false
}
```

## `flux-2/pro-image-to-image`

**Flux-2 - Pro Image to Image** · Flux2 · [official docs](https://docs.kie.ai/market/flux2/pro-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `input_urls` | array | yes | max items: 8 | Input reference images (1-8 images). (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `prompt` | string | yes | max length: 5000 | Must be between 3 and 5000 characters. (Max length: 5000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `3:2`, `2:3`, `auto`; default: `"1:1"` | Aspect ratio for the generated image. Select 'auto' to match the first input image ratio (requires input image). |
| `resolution` | string | yes | allowed: `1K`, `2K`; default: `"1K"` | Output image resolution. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1764235041265_kjJ2sTMR.png",
    "https://static.aiquickdraw.com/tools/example/1764235045490_9SjAUr4Z.png"
  ],
  "prompt": "The jar in image 1 is filled with capsules exactly same as image 2 with the exact logo",
  "aspect_ratio": "1:1",
  "resolution": "1K",
  "nsfw_checker": false
}
```

## `flux-2/pro-text-to-image`

**Flux-2 - Pro Text to Image** · Flux2 · [official docs](https://docs.kie.ai/market/flux2/pro-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Must be between 3 and 5000 characters. (Max length: 5000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `3:2`, `2:3`; default: `"1:1"` | Aspect ratio for the generated image. Select 'auto' to match the first input image ratio (requires input image). |
| `resolution` | string | yes | allowed: `1K`, `2K`; default: `"1K"` | Output image resolution. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Hyperrealistic supermarket blister pack on clean olive green surface. No shadows. Inside: bright pink 3D letters spelling \"FLUX.2\" pressing against stretched plastic film, creating realistic deformation and reflective highlights. Bottom left corner: barcode sticker with text \"GENERATE NOW\" and \"PLAYGROUND\". Plastic shows tension wrinkles and realistic shine where stretched by the volumetric letters.",
  "aspect_ratio": "1:1",
  "resolution": "1K",
  "nsfw_checker": false
}
```

## `gemini-omni-video`

**Gemini Omni Video** · Market · [official docs](https://docs.kie.ai/market/gemini-omni-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 20000 | Video prompt used to describe the target content, style, camera language, or character actions in the generated video. |
| `image_urls` | array | no |  | Array of image URLs. You can provide one or more reference images for characters, scenes, styles, or storyboard guidance.  Image limits: - Each file must be no larger than `20MB` - Use publicly accessible image URLs - Max 7 images |
| `audio_ids` | array | no |  | Array of audio IDs generated by the `gemini-omni-audio` endpoint. Useful for narration, dialogue, music, or audio guidance in the generated video. Max 3 items. |
| `video_list` | array | no |  | Array of video clips. Each item defines a source video and the trim range to use during generation.  Video limits: - Each file must be no larger than `100MB` - Video duration must not exceed `30s` - `ends` should be greater than `start` - The difference between the end time and the start time must not exceed `10s`. - Max 1 items. Equal 2 images |
| `video_list[].url` | string | no | format: uri | Video URL. Each source video file must be no larger than `100MB` and no longer than `30s`. |
| `video_list[].start` | number | no | min: 0 | Start time in seconds. |
| `video_list[].ends` | number | no | min: 0 | End time in seconds. It should be greater than `start`.The difference between the end time and the start time must not exceed 10s |
| `character_ids` | array | no |  | An array of character IDs generated by the `gemini-omni-character` API. Used to provide character appearance, identity, or person references for the video. Each character_id uses 1 image slot. The base limit is 7 image slots; if video_list is also provided, video_list uses 2 image slots, so character_ids can contain up to 3 IDs. |
| `duration` | string | yes | allowed: `4`, `6`, `8`, `10` | The duration of the generated video in seconds. Available values are 4, 6, 8, and 10. When video input is provided, the output duration is determined by the model automatically. This duration parameter will not take effect.Note: when video input is provided, the output duration is determined by the model automatically. This duration parameter will not take effect. |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16` | The aspect ratio of the generated video. `16:9` is landscape, and `9:16` is portrait. |
| `seed` | integer | no |  | Random seed. Range: [0, 2147483647]. If not specified, the system generates a seed automatically. Fixing the seed can improve reproducibility, but results may still vary due to the model’s stochasticity. |
| `resolution` | string | no | allowed: `720p`, `1080p`, `4k`; default: `"720p"` | The resolution of the generated video. Available values are 720p, 1080p, and 4k. |

Example `input`:

```json
{
  "prompt": "Create a futuristic night city short film with a slow push-in shot as the character walks out from a neon-lit street.",
  "image_urls": [
    "https://example.com/assets/scene-1.png",
    "https://example.com/assets/scene-2.png"
  ],
  "audio_ids": [
    "audio_01hx8p0demo"
  ],
  "video_list": [
    {
      "url": "https://example.com/assets/source-video.mp4",
      "start": 0,
      "ends": 10
    }
  ],
  "duration": "4"
}
```

## `google/gemini-2-5-pro-tts`

**Gemini 2.5 Pro Text to Speech** · Google · [official docs](https://docs.kie.ai/google/gemini-2-5-pro-tts)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `temperature` | number | no | default: `1`; min: 0; max: 2 | Sampling temperature, e.g., 1 |
| `scene` | string | no | default: `""` | Scene description, e.g., "A quiet, warm room with a fireplace crackling softly." |
| `sample_context` | string | no | default: `""` | Sample context/overall tone, e.g., "Audiobook style narration. Tone is gentle and inviting." |
| `speakers` | array | yes |  | List of speaker configurations |
| `speakers[].speaker_id` | string | yes |  | Speaker identifier, e.g., "Speaker 1" or "Speaker 2";Must be 「Speaker N」format |
| `speakers[].voice_name` | string | yes | allowed: `Achernar`, `Achird`, `Algenib`, `Algieba`, `Alnilam`, `Aoede`, `Autonoe`, `Callirrhoe`, `Charon`, `Despina`, `Enceladus`, `Erinome`, `Fenrir`, `Gacrux`, `Iapetus`, `Kore`, `Laomedeia`, `Leda`, `Orus`, `Puck`, `Pulcherrima`, `Rasalgethi`, `Sadachbia`, `Sadaltager`, `Schedar`, `Sulafat`, `Umbriel`, `Vindemiatrix`, `Zephyr`, `Zubenelgenubi` | Voice name, e.g., "Zephyr", "Fenrir", "Puck" |
| `speakers[].audio_profile` | string | no |  | Audio profile description, e.g., "A warm and soothing narrator" or "A stern and weary gatekeeper" |
| `speakers[].accent` | string | yes | allowed: `Neutral`, `American (Gen)`, `American (Valley)`, `American (South)`, `British (RP)`, `British (Brixton)`, `Transatlantic`, `Australian` | Accent, e.g., "British (RP)" or "American (Gen)" |
| `speakers[].style` | string | no | allowed: `Vocal Smile`, `Newscaster`, `Whisper`, `Empathetic`, `Promo/Hype`, `Deadpan` | Emotional style, e.g., "Gentle" or "Deadpan" |
| `speakers[].pace` | string | no | allowed: `Natural`, `Rapid Fire`, `The Drift`, `Staccato` | Pace, e.g., "Slow" or "Natural" |
| `dialogue_turns` | array | yes |  | List of dialogue turns, output in sequential order |
| `dialogue_turns[].speaker_id` | string | yes |  | Corresponding speaker identifier, e.g., "Speaker 1" |
| `dialogue_turns[].text` | string | yes | max length: 10000 | The text spoken by the speaker, which may contain tone tags, e.g., "Once upon a time, in a quiet valley hidden away..." |

Example `input`:

```json
{
  "temperature": 1,
  "scene": "A dark, crumbling dungeon...",
  "sample_context": "Fantasy RPG style...",
  "speakers": [
    {
      "speaker_id": "Speaker 1",
      "voice_name": "Fenrir",
      "audio_profile": "A stern and weary gatekeeper",
      "accent": "British (RP)",
      "style": "Deadpan",
      "pace": "Natural"
    },
    {
      "speaker_id": "Speaker 2",
      "voice_name": "Puck",
      "audio_profile": "A determined and courageous traveler seeking answers.",
      "accent": "American (Gen)",
      "style": "Empathetic",
      "pace": "Staccato"
    }
  ],
  "dialogue_turns": [
    {
      "speaker_id": "Speaker 1",
      "text": "[shouting] Halt, traveler! The northern pass is sealed by order of the council."
    },
    {
      "speaker_id": "Speaker 2",
      "text": "[determination] I carry a message for the elder. Step aside, or I will force my way through."
    },
    {
      "speaker_id": "Speaker 1",
      "text": "[caution] No one passes. [pensive] The elder is... he's no longer receiving visitors."
    },
    {
      "speaker_id": "Speaker 2",
      "text": "It's too late. [whispers] The shadow... it reached him first. [urgency] You need to leave. [shouting] Now."
    }
  ]
}
```

## `google/gemini-3-1-flash-tts`

**Gemini 3.1 Flash Text to speech** · Google · [official docs](https://docs.kie.ai/market/google/gemini-3-1-flash-tts)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `temperature` | number | no | default: `1`; min: 0; max: 2 | Sampling temperature, e.g., 1 |
| `scene` | string | no | default: `""` | Scene description, e.g., "A quiet, warm room with a fireplace crackling softly." |
| `sample_context` | string | no | default: `""` | Sample context/overall tone, e.g., "Audiobook style narration. Tone is gentle and inviting." |
| `speakers` | array | yes |  | List of speaker configurations |
| `speakers[].speaker_id` | string | yes |  | Speaker identifier, e.g., "Speaker 1" or "Speaker 2";Must be 「Speaker N」format |
| `speakers[].voice_name` | string | yes | allowed: `Achernar`, `Achird`, `Algenib`, `Algieba`, `Alnilam`, `Aoede`, `Autonoe`, `Callirrhoe`, `Charon`, `Despina`, `Enceladus`, `Erinome`, `Fenrir`, `Gacrux`, `Iapetus`, `Kore`, `Laomedeia`, `Leda`, `Orus`, `Puck`, `Pulcherrima`, `Rasalgethi`, `Sadachbia`, `Sadaltager`, `Schedar`, `Sulafat`, `Umbriel`, `Vindemiatrix`, `Zephyr`, `Zubenelgenubi` | Voice name, e.g., "Zephyr", "Fenrir", "Puck" |
| `speakers[].audio_profile` | string | no |  | Audio profile description, e.g., "A warm and soothing narrator" or "A stern and weary gatekeeper" |
| `speakers[].accent` | string | yes | allowed: `Neutral`, `American (Gen)`, `American (Valley)`, `American (South)`, `British (RP)`, `British (Brixton)`, `Transatlantic`, `Australian` | Accent, e.g., "British (RP)" or "American (Gen)" |
| `speakers[].style` | string | no | allowed: `Vocal Smile`, `Newscaster`, `Whisper`, `Empathetic`, `Promo/Hype`, `Deadpan` | Emotional style, e.g., "Gentle" or "Deadpan" |
| `speakers[].pace` | string | no | allowed: `Natural`, `Rapid Fire`, `The Drift`, `Staccato` | Pace, e.g., "Slow" or "Natural" |
| `dialogue_turns` | array | yes |  | List of dialogue turns, output in sequential order |
| `dialogue_turns[].speaker_id` | string | yes |  | Corresponding speaker identifier, e.g., "Speaker 1" |
| `dialogue_turns[].text` | string | yes | max length: 10000 | The text spoken by the speaker, which may contain tone tags, e.g., "Once upon a time, in a quiet valley hidden away..." |

Example `input`:

```json
{
  "temperature": 1,
  "scene": "A dark, crumbling dungeon...",
  "sample_context": "Fantasy RPG style...",
  "speakers": [
    {
      "speaker_id": "Speaker 1",
      "voice_name": "Fenrir",
      "audio_profile": "A stern and weary gatekeeper",
      "accent": "British (RP)",
      "style": "Deadpan",
      "pace": "Natural"
    },
    {
      "speaker_id": "Speaker 2",
      "voice_name": "Puck",
      "audio_profile": "A determined and courageous traveler seeking answers.",
      "accent": "American (Gen)",
      "style": "Empathetic",
      "pace": "Staccato"
    }
  ],
  "dialogue_turns": [
    {
      "speaker_id": "Speaker 1",
      "text": "[shouting] Halt, traveler! The northern pass is sealed by order of the council."
    },
    {
      "speaker_id": "Speaker 2",
      "text": "[determination] I carry a message for the elder. Step aside, or I will force my way through."
    },
    {
      "speaker_id": "Speaker 1",
      "text": "[caution] No one passes. [pensive] The elder is... he's no longer receiving visitors."
    },
    {
      "speaker_id": "Speaker 2",
      "text": "It's too late. [whispers] The shadow... it reached him first. [urgency] You need to leave. [shouting] Now."
    }
  ]
}
```

## `google/imagen4`

**Google - imagen4** · Google · [official docs](https://docs.kie.ai/market/google/imagen4)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt describing what you want to see (Max length: 5000 characters) |
| `negative_prompt` | string | no | max length: 5000 | A description of what to discourage in the generated images (Max length: 5000 characters) |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `9:16`, `3:4`, `4:3`, `auto`; default: `"1:1"` | The aspect ratio of the generated image |
| `seed` | string | no | max length: 500 | Random seed for reproducible generation (Max length: 500 characters) |

Example `input`:

```json
{
  "prompt": "A lively comic scene where two colleagues are in an office. The first person says, 'Have you heard about Google Imagen 4 Ultra?' The second person responds with excitement, 'It’s the best text-to-image tool out there!' The first person asks again, 'Do you know where to get the API?' The second person smiles and says, 'Kie.ai has it!' In the final panel, the two look at a screen showing Kie.ai’s interface with an API option, with bright and colorful comic-style illustrations.",
  "negative_prompt": "",
  "aspect_ratio": "1:1",
  "seed": ""
}
```

## `google/imagen4-fast`

**Google - imagen4-fast** · Google · [official docs](https://docs.kie.ai/market/google/imagen4-fast)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt describing what you want to see (Max length: 5000 characters) |
| `negative_prompt` | string | no | max length: 5000 | A description of what to discourage in the generated images (Max length: 5000 characters) |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `9:16`, `3:4`, `4:3`, `auto`; default: `"16:9"` | The aspect ratio of the generated image |
| `seed` | integer | no |  | Random seed for reproducible generation |

Example `input`:

```json
{
  "prompt": "Create a cinematic, photorealistic medium shot capturing the nostalgic warmth of a late 90s indie film. The focus is a young woman with brightly dyed pink hair (slightly faded) and freckled skin, looking directly and intently into the camera lens with a hopeful yet slightly uncertain smile. She wears an oversized, vintage band t-shirt (slightly worn, with the faintly cracked white text “KIE AI” across the chest) layered over a long-sleeved striped top, along with simple silver stud earrings. The lighting is soft, golden hour sunlight streaming through a slightly dusty window, creating lens flare and illuminating dust motes in the air. The background shows a blurred, cluttered bedroom with posters on the wall and fairy lights, rendered with a shallow depth of field. Natural film grain, a warm, slightly muted color palette, and sharp focus on her expressive eyes enhance the intimate, authentic feel.",
  "negative_prompt": "",
  "aspect_ratio": "16:9"
}
```

## `google/imagen4-ultra`

**Google - imagen4-ultra** · Google · [official docs](https://docs.kie.ai/market/google/imagen4-ultra)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt describing what you want to see (Max length: 5000 characters) |
| `negative_prompt` | string | no | max length: 5000 | A description of what to discourage in the generated images (Max length: 5000 characters) |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `9:16`, `3:4`, `4:3`, `auto`; default: `"1:1"` | The aspect ratio of the generated image |
| `seed` | string | no | max length: 500 | Random seed for reproducible generation (Max length: 500 characters) |

Example `input`:

```json
{
  "prompt": "A lively comic scene where two colleagues are in an office. The first person says, 'Have you heard about Google Imagen 4 Ultra?' The second person responds with excitement, 'It’s the best text-to-image tool out there!' The first person asks again, 'Do you know where to get the API?' The second person smiles and says, 'Kie.ai has it!' In the final panel, the two look at a screen showing Kie.ai’s interface with an API option, with bright and colorful comic-style illustrations.",
  "negative_prompt": "",
  "aspect_ratio": "1:1",
  "seed": ""
}
```

## `google/nano-banana`

**Google - Nano Banana** · Google · [official docs](https://docs.kie.ai/market/google/nano-banana)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt for image generation (Max length: 5000 characters) |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Output format for the images |
| `aspect_ratio` | string | no | allowed: `1:1`, `9:16`, `16:9`, `3:4`, `4:3`, `3:2`, `2:3`, `5:4`, `4:5`, `21:9`, `auto`; default: `"1:1"` | Radio description |
| `image_size` | string | no | allowed: `1:1`, `9:16`, `16:9`, `3:4`, `4:3`, `3:2`, `2:3`, `5:4`, `4:5`, `21:9`, `auto`; default: `"1:1"` | The aspect ratio of the generated image (this parameter has been replaced by aspect_ratio; please use the latest aspect_ratio parameter). |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A surreal painting of a giant banana floating in space, stars and galaxies in the background, vibrant colors, digital art",
  "output_format": "png",
  "aspect_ratio": "1:1"
}
```

## `google/nano-banana-edit`

**Google - Nano Banana Edit** · Google · [official docs](https://docs.kie.ai/market/google/nano-banana-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt for image editing (Max length: 5000 characters) |
| `image_urls` | array | yes | max items: 10 | List of URLs of input images for editing,up to 10 images. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Output format for the images |
| `aspect_ratio` | string | no | allowed: `1:1`, `9:16`, `16:9`, `3:4`, `4:3`, `3:2`, `2:3`, `5:4`, `4:5`, `21:9`, `auto`; default: `"1:1"` | Radio description |
| `image_size` | string | no | allowed: `1:1`, `9:16`, `16:9`, `3:4`, `4:3`, `3:2`, `2:3`, `5:4`, `4:5`, `21:9`, `auto`; default: `"1:1"` | The aspect ratio of the generated image (this parameter has been replaced by aspect_ratio; please use the latest aspect_ratio parameter). |

Example `input`:

```json
{
  "prompt": "turn this photo into a character figure. Behind it, place a box with the character’s image printed on it, and a computer showing the Blender modeling process on its screen. In front of the box, add a round plastic base with the character figure standing on it. set the scene indoors if possible",
  "image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1756223420389w8xa2jfe.png"
  ],
  "output_format": "png",
  "aspect_ratio": "1:1"
}
```

## `gpt-image-2-image-to-image`

**GPT Image 2 - Image To Image** · Gpt · [official docs](https://docs.kie.ai/market/gpt/gpt-image-2-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes |  | Text prompts, up to 20,000 characters. |
| `input_urls` | array | yes | max items: 16 | Array of input image URLs. |
| `aspect_ratio` | string | no | allowed: `auto`, `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `5:4`, `4:5`, `16:9`, `9:16`, `2:1`, `1:2`, `3:1`, `1:3`, `21:9`, `9:21` | The aspect ratio of the generated image is set to auto by default. Note: 5:4 and 4:5 aspect ratios only support 1K images. |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K` | Image resolution: Note: Images with a 1:1 aspect ratio cannot be converted to 4K images. Images with the aspect ratio set to "auto" or without a specified aspect ratio parameter will only be converted to 1K images; otherwise, the task will fail to create. |

Example `input`:

```json
{
  "prompt": "take a photo with Sam Altman in the conference room",
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1776782793756_wrogXTdd.png"
  ],
  "aspect_ratio": "auto"
}
```

## `gpt-image-2-text-to-image`

**GPT Image-2 - Text to Image** · Gpt · [official docs](https://docs.kie.ai/market/gpt/gpt-image-2-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 1; max length: 20000 | Text prompt. Required, maximum 20,000 characters. |
| `aspect_ratio` | string | no | allowed: `auto`, `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `5:4`, `4:5`, `16:9`, `9:16`, `2:1`, `1:2`, `3:1`, `1:3`, `21:9`, `9:21` | The aspect ratio of the generated image is set to auto by default. Note: for 2K and 4K resolution, the following aspect ratios are not supported: 5:4, 4:5, 3:1, 1:3, and 9:21. |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K` | Image resolution: Note: Images with a 1:1 aspect ratio cannot be converted to 4K images. Images with the aspect ratio set to "auto" or without a specified aspect ratio parameter will only be converted to 1K images; otherwise, the task will fail to create. |

Example `input`:

```json
{
  "prompt": "A cinematic night city poster with neon reflections on a rainy street.",
  "aspect_ratio": "auto"
}
```

## `gpt-image/1.5-image-to-image`

**GPT Image-1.5 - Image to Image** · Gpt Image · [official docs](https://docs.kie.ai/market/gpt-image/1-5-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `input_urls` | array | yes | max items: 16 | Upload an image file to use as input for the API (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `prompt` | string | yes |  | A text description of the image you want to generate |
| `aspect_ratio` | string | yes | allowed: `1:1`, `2:3`, `3:2`; default: `"3:2"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `medium`, `high`; default: `"medium"` | Quality: medium=balanced, high=slow/detailed. |

Example `input`:

```json
{
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1765962794374_GhtqB9oX.webp"
  ],
  "prompt": "Edit the image to dress the woman using the provided clothing images. Preserve her exact likeness, expression, hairstyle, and proportions. Replace only the clothing, fitting the garments naturally to her existing pose and body geometry with realistic fabric behavior. Match lighting, shadows, and color temperature to the original photo so the outfit integrates photorealistically, without looking pasted on.",
  "aspect_ratio": "3:2",
  "quality": "medium"
}
```

## `gpt-image/1.5-text-to-image`

**GPT Image-1.5 - Text to Image** · Gpt Image · [official docs](https://docs.kie.ai/market/gpt-image/1-5-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes |  | A text description of the image you want to generate |
| `aspect_ratio` | string | yes | allowed: `1:1`, `2:3`, `3:2`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `medium`, `high`; default: `"medium"` | Quality: medium=balanced, high=slow/detailed. |

Example `input`:

```json
{
  "prompt": "Create a photorealistic candid photograph of an elderly sailor standing on a small fishing boat.  He has weathered skin with visible wrinkles, pores, and sun texture, and a few faded traditional sailor tattoos on his arms. He is calmly adjusting a net while his dog sits nearby on the deck. Shot like a 35mm film photograph, medium close-up at eye level, using a 50mm lens. The image should feel honest and unposed, with real skin texture, worn materials, and everyday detail. No glamorization, no heavy retouching. ",
  "aspect_ratio": "1:1",
  "quality": "medium"
}
```

## `grok-imagine-video-1-5-preview`

**Grok Imagine Video 1.5 Preview** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/1-5-preview)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 4096 | Prompt for video generation. Maximum length: 4096 characters. |
| `image_urls` | array | no | max items: 7 | Upload image files to be used as API input. Supported file types: image/jpeg, image/png, image/webp, image/jpg. Maximum file size: 20MB. Supports multi-file upload, up to 7 file.Only one is supported when the resolution is 1080p. |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `9:16`, `3:2`, `2:3`, `auto`; default: `"auto"` | The aspect ratio of the video. This parameter is invalid if it is a single image. |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"480p"` | Resolution for video generation. |
| `duration` | integer | no | default: `8`; min: 1; max: 15 | Video duration in seconds. Range: [1, 15]. Default: 8. Minimum: 1. Maximum: 15. Step: 1. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Describe the scene you want to generate.",
  "image_urls": [
    "https://your-domain.com/image/example.png"
  ],
  "aspect_ratio": "16:9",
  "resolution": "480p",
  "duration": 8
}
```

## `grok-imagine/extend`

**Grok Imagine - Video Extend** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/extend)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `task_id` | string | yes | max length: 100 | Task ID from a previously successful video generation task. Required field.  - Must be from a Kie AI video generation model (e.g., grok-imagine/text-to-video) - The original video generation must have completed successfully - Only Kie AI–generated task IDs are supported |
| `prompt` | string | yes |  | Text instructions describing the required movement of the video. Required field.  - Provide a detailed description of how you would like the video to expand and continue. - You can specify camera movements, scene changes, object actions, etc. - The more specific the prompt words are, the more likely the generated effect will match your expectations. - Supports input in both Chinese and English. |
| `extend_at` | number | yes | default: `2`; min: 2 | The starting position of the video extension. Optional field. |
| `extend_times` | number | yes |  | Duration of video extension (in seconds). Required field.  - `6`: Expand 6 seconds of video content - `10`: Expand 10 seconds of video content - The longer the extension duration, the longer the time required for generation - Select the appropriate duration based on the complexity of the scene |

Example `input`:

```json
{
  "task_id": "task_grok_12345678",
  "prompt": "",
  "extend_at": 2,
  "extend_times": "6"
}
```

## `grok-imagine/image-to-image`

**Grok Imagine - image to image** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 390000 | A text description specifying the desired content or style of the generated image. (Max length: 390000 characters) |
| `image_urls` | array | yes | max items: 5 | An array containing up to 1 URL string pointing to reference images. (Use file URLs after upload, not raw file content. Accepted types: image/jpeg, image/png, image/webp. Max size: 10.0MB per image.) In your prompt, reference the uploaded image by typing @image(n) followed by a space (for example: @image1 a sunset over the ocean). |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Recreate the Titanic movie poster with two adorable anthropomorphic cats in the same romantic pose at the bow of the ship. The male cat is an orange tabby wearing a vest, standing behind a white long-haired female cat in a lace dress, holding her paws as they stretch forward in the wind. Both cats are photorealistic with detailed fur, wind-swept hair, and dramatic sunset lighting (warm golden highlights, cool blue shadows). Background: the Titanic ship at dusk with four smokestacks, glowing deck lights, calm ocean, and orange-pink sunset sky. Center title: “CATANIC” in the same gold metallic serif style as Titanic, same size and position.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1767602105243_0MmMCrwq.png"
  ]
}
```

## `grok-imagine/image-to-video`

**Grok Imagine Image to Video** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_urls` | array | no | max items: 7 | Provide an external image URL as a reference for video generation. Up to 7 images are supported. Do not use it simultaneously with task_id. In your prompt, reference an uploaded image by typing @image(n) followed by a space (for example: @image1 a sunset over the ocean). - Supports JPEG, PNG, and WEBP formats - Maximum file size for each image: 10MB - The Spicy mode is not available when using external images - The array can contain a maximum of seven URLs - Only one is supported when the resolution is 1080p. |
| `task_id` | string | no | max length: 100 | Task ID from a previously generated Grok image. Use with index to select a specific image. Do not use with image_urls.  - Use task ID from grok-imagine/text-to-image generations - Supports all modes including Spicy - Maximum length: 100 characters |
| `index` | integer | no | default: `0`; min: 0; max: 5 | When using task_id, specify which image to use (Grok generates 6 images per task). Only works with task_id.  - 0-based index (0-5) - Ignored if image_urls is provided - Default: 0 |
| `prompt` | string | no |  | Text prompt describing the desired video motion. Optional field.  - Should be detailed and specific about the desired visual motion - Describe movement, action sequences, camera work, and timing - Include details about subjects, environments, and motion dynamics - Maximum length: 5000 characters - Supports English language prompts |
| `mode` | string | no | allowed: `fun`, `normal`, `spicy`; default: `"normal"` | Specifies the generation mode affecting the style and intensity of motion. Note: Spicy mode is not available for external image inputs.  - **fun**: More creative and playful interpretation - **normal**: Balanced approach with good motion quality - **spicy**: More dynamic and intense motion effects (not available for external images)  Default: normal |
| `duration` | string | no |  | The duration of the generated video (in seconds) (6-30). (Minimum: 6, Maximum: 30, Step: 1) |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"480p"` | The resolution of the generated video. |
| `aspect_ratio` | string | no | allowed: `2:3`, `3:2`, `1:1`, `16:9`, `9:16`; default: `"16:9"` | Image ratio selection only applies to multi-image generation mode. In single-image mode, the video width and height are referenced to the image width and height. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "task_id": "task_grok_12345678",
  "image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1762247692373tw5di116.png"
  ],
  "prompt": "POV hand comes into frame handing the girl a cup of take away coffee, the girl steps out of the screen looking tired, then takes it and she says happily: \"thanks! Back to work\" she exits the frame and walks right to a different part of the office.",
  "mode": "normal",
  "duration": "6",
  "resolution": "480p",
  "aspect_ratio": "16:9"
}
```

## `grok-imagine/text-to-image`

**Grok Imagine - Text to Image** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes |  | Text prompt describing the desired image. Required field.  - Should be detailed and specific about the desired visual elements - Describe composition, style, lighting, mood, and other visual details - Maximum length: 5000 characters - Supports English language prompts |
| `aspect_ratio` | string | no | allowed: `2:3`, `3:2`, `1:1`, `16:9`, `9:16` | Specifies the width-to-height ratio of the generated image. Controls the aspect ratio of the output.  - **2:3**: Portrait orientation (vertical) - **3:2**: Landscape orientation (horizontal)  - **1:1**: Square format - **16:9**: Wide screen format - **9:16**: Tall screen format  Default: 1:1 |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |
| `enable_pro` | boolean | no |  | Controls the request processing strategy.     - `false`: Corresponds to **speed mode**. The system prioritizes response time and throughput, suitable for latency-sensitive scenarios.     - `true`: Corresponds to **quality mode**. The system prioritizes processing quality and precision, suitable for scenarios requiring higher accuracy. |

Example `input`:

```json
{
  "prompt": "Cinematic portrait of a woman sitting by a vinyl record player, retro living room background, soft ambient lighting, warm earthy tones, nostalgic 1970s wardrobe, reflective mood, gentle film grain texture, shallow depth of field, vintage editorial photography style.",
  "aspect_ratio": "3:2"
}
```

## `grok-imagine/text-to-video`

**Grok Imagine Text to Video** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes |  | Text prompt describing the desired video motion. Required field.  - Should be detailed and specific about the desired visual motion - Describe movement, action sequences, camera work, and timing - Include details about subjects, environments, and motion dynamics - Maximum length: 5000 characters - Supports English language prompts |
| `aspect_ratio` | string | no | allowed: `2:3`, `3:2`, `1:1`, `16:9`, `9:16`; default: `"2:3"` | Specifies the width-to-height ratio of the generated video. Controls the aspect ratio of the output.  - **2:3**: Portrait orientation (vertical) - **3:2**: Landscape orientation (horizontal) - **1:1**: Square format - **16:9**: Wide screen format - **9:16**: Tall screen format  Default: 2:3 |
| `mode` | string | no | allowed: `fun`, `normal`, `spicy`; default: `"normal"` | Specifies the generation mode affecting the style and intensity of motion.  - **fun**: More creative and playful interpretation - **normal**: Balanced approach with good motion quality - **spicy**: More dynamic and intense motion effects  Default: normal |
| `duration` | number | no |  | The duration of the generated video (in seconds) (6-30). (Minimum: 6, Maximum: 30, Step: 1) |
| `resolution` | string | no | allowed: `480p`, `720p`, `1080p`; default: `"480p"` | The resolution of the generated video. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A couple of doors open to the right one by one randomly and stay open, to show the inside, each is either a living room, or a kitchen, or a bedroom or an office, with little people living inside.",
  "aspect_ratio": "2:3",
  "mode": "normal",
  "duration": "6",
  "resolution": "480p"
}
```

## `grok-imagine/upscale`

**Grok Imagine - Video Upscale** · Grok Imagine · [official docs](https://docs.kie.ai/market/grok-imagine/upscale)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `task_id` | string | yes | max length: 100 | Task ID from a previously successful video generation task. Required field.  - Must be from a Kie AI video generation model (e.g., grok-imagine/text-to-video) - The original video generation must have completed successfully - Only Kie AI–generated task IDs are supported |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"720p"` | Video generation resolution. |

Example `input`:

```json
{
  "task_id": "task_grok_12345678"
}
```

## `hailuo/02-image-to-video-pro`

**Hailuo Pro Image to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/02-image-to-video-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | Text prompt describing the desired video animation (Max length: 1500 characters) |
| `image_url` | string | yes |  | Input image to animate (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `end_image_url` | string | no |  | Optional URL of the image to use as the last frame of the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `prompt_optimizer` | boolean | no |  | Whether to use the model's prompt optimizer (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Cinematic wide shot: A colossal starship drifts silently above the rings of Saturn, its metallic hull reflecting streaks of cosmic light. The camera pushes closer, revealing thousands of illuminated windows like a floating city. Smaller fighter crafts dart across the frame, leaving neon trails as they maneuver through the vastness of space. A sudden burst of thrusters scatters asteroid fragments in slow motion, glowing faintly as they collide and drift apart.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17585210783150ispzfo7.png",
  "end_image_url": "",
  "prompt_optimizer": true
}
```

## `hailuo/02-image-to-video-standard`

**Hailuo Standard Image to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/02-image-to-video-standard)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | The text prompt describing the video to generate (Max length: 1500 characters) |
| `image_url` | string | yes |  | The URL of the image to use as the first frame of the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `end_image_url` | string | no |  | Optional URL of the image to use as the last frame of the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `duration` | string | no | allowed: `6`, `10`; default: `"10"` | The duration of the video in seconds. 10 seconds videos are not supported for 1080p resolution. |
| `resolution` | string | no | allowed: `512P`, `768P`; default: `"768P"` | The resolution of the generated video. |
| `prompt_optimizer` | boolean | no |  | Whether to use the model's prompt optimizer (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Epic aerial shot: A lone samurai stands atop a jagged mountain peak as a storm of sakura petals is swept across the wind. Behind him, the sky is split in two — half daylight, half night. The shot pulls back to reveal that the mountain is actually the curved back of a sleeping dragon that spans across the horizon. Lightning crackles in the distance as the dragon's eye slowly opens, glowing with ancient magic. The samurai doesn’t flinch; he lowers his straw hat and places his hand on the hilt of his blade.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17585207681646umf3lz8.png",
  "end_image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1758521423357w8586uq8.png",
  "duration": "10",
  "resolution": "768P",
  "prompt_optimizer": true
}
```

## `hailuo/02-text-to-video-pro`

**Hailuo Pro Text to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/02-text-to-video-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | The text prompt for video generation (Max length: 1500 characters) |
| `prompt_optimizer` | boolean | no |  | Whether to use the model's prompt optimizer (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "High top angle wide mid close-up tracking shot, flying very fast two meters high over prehistoric ferns and moss-covered ground, dominated by a real young boy (pink t-shirt, pink shorts, white shoes, white long socks), with his back to the camera, body stretched out, gliding smoothly forward, flying in the air, casting a clear shadow on the terrain below. His legs and body are high above the surface, his feet not touching the ground, soaring in a Superman pose. The background is a vast Jurassic valley, filled with dense, ancient jungle vegetation and towering cycads. In the distance, rugged volcanic mountains rise with a winding path cutting through. Massive, slow-moving sauropods graze far off on the horizon. Large, fluffy white clouds float in the vibrant blue sky. Strong dynamic motion blur adds a vivid sense of high-speed flight and deep cinematic perspective. Realistic image –raw.",
  "prompt_optimizer": true
}
```

## `hailuo/02-text-to-video-standard`

**Hailuo Standard Text to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/02-text-to-video-standard)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | Text description for video generation (Max length: 1500 characters) |
| `duration` | string | no | allowed: `6`, `10`; default: `"6"` | The duration of the video in seconds. 10 seconds videos are not supported for 1080p resolution. |
| `prompt_optimizer` | boolean | no |  | Whether to use the model's prompt optimizer (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A llama and a raccoon battle it out in an intense table tennis match, inside a roaring Olympic stadium. Slow-mo, wild angles, full comedy mode.",
  "duration": "6",
  "prompt_optimizer": true
}
```

## `hailuo/2-3-image-to-video-pro`

**Hailuo 2.3 Pro Image to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/2-3-image-to-video-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the desired video animation (Max length: 5000 characters) |
| `image_url` | string | yes |  | Input image to animate (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `duration` | string | no | allowed: `6`, `10`; default: `"6"` | The duration of the video in seconds. 10 seconds videos are not supported for 1080p resolution. |
| `resolution` | string | no | allowed: `768P`, `1080P`; default: `"768P"` | The resolution of the generated video. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A graceful geisha performs a traditional Japanese dance indoors. She wears a luxurious red kimono with golden floral embroidery, white obi belt, and white tabi socks. Soft and elegant hand movements, expressive pose, sleeves flowing naturally. Scene set in a Japanese tatami room with warm ambient lighting, shoji paper sliding doors, and cherry blossom branches hanging in the foreground. Cinematic, soft depth of field, high detail fabric texture, hyper-realistic, smooth motion.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1761736831884xl56xfiw.webp",
  "duration": "6",
  "resolution": "768P"
}
```

## `hailuo/2-3-image-to-video-standard`

**Hailuo 2.3 Standard Image to Video** · Hailuo · [official docs](https://docs.kie.ai/market/hailuo/2-3-image-to-video-standard)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the desired video animation (Max length: 5000 characters) |
| `image_url` | string | yes |  | Input image to animate (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `duration` | string | no | allowed: `6`, `10`; default: `"6"` | The duration of the video in seconds. 10 seconds videos are not supported for 1080p resolution. |
| `resolution` | string | no | allowed: `768P`, `1080P`; default: `"768P"` | The resolution of the generated video. |

Example `input`:

```json
{
  "prompt": "Two armored medieval knights clash in an intense duel at sunset, cinematic lighting.  Metal armor reflects warm golden light from the sun and the glowing swords. Sparks explode as the swords collide. Dynamic camera movement, shallow depth of field, dramatic slow motion. The scene takes place in an open desert battlefield, dust in the air, warm orange sun behind them, epic atmosphere.  Highly detailed armor textures, realistic reflections, volumetric lighting, cinematic quality.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1761736401898mpm67du5.webp",
  "duration": "6",
  "resolution": "768P"
}
```

## `happyhorse-1-1/image-to-video`

**HappyHorse-1-1 image-to-video** · Happyhorse 1 1 · [official docs](https://docs.kie.ai/market/happyhorse-1-1/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | default: `""`; max length: 5000 | Describes the video content to generate. Supports any language. Maximum: 5,000 non-Chinese characters or 2,500 Chinese |
| `image_urls` | array | yes | default: `[""]`; max items: 1 | The URL of the first frame image.  Image constraints:  1. Formats: JPEG, JPG, PNG, WEBP.  2. Resolution: Width and height must both be at least 300 pixels.  3. Aspect ratio: Between 1:2.5 and 2.5:1.  4. File size: Up to 20 MB. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | The resolution of the generated video. Options: 720p; 1080p |
| `duration` | number | no | default: `5`; min: 3; max: 15 | The duration of the generated video, in seconds.  The value must be an integer in the range [3, 15]. Default: 5 |

Example `input`:

```json
{
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1782114387854_IufKnPxR.png"
  ],
  "prompt": "A cat running on the grass",
  "resolution": "1080p",
  "duration": 5
}
```

## `happyhorse-1-1/reference-to-video`

**HappyHorse-1-1 reference-to-video** · Happyhorse 1 1 · [official docs](https://docs.kie.ai/market/happyhorse-1-1/reference-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | default: `""`; max length: 5000 | A description of the desired elements and visual style for the generated video.  Input in any language is supported. The length is limited to 5,000 non-Chinese characters or 2,500 Chinese characters. Content exceeding this limit is automatically truncated.  Image referencing: In the prompt, use "[Image 1]" and "[Image 2]" to refer to the corresponding reference image in the media array. The order must be consistent with the order in the media array. When using a reference, specify the object in the image, such as "the woman in a red qipao in [Image 1]". |
| `reference_image` | array | yes | default: `[]`; max items: 9 | The URL  of a reference image.  Image requirements:  1. Formats: JPEG, JPG, PNG, WEBP.  2. Resolution: The shortest side must be at least 400 pixels. A clear image with a resolution of 720P or higher is recommended. Avoid using images that are too small, blurry, or overly compressed, as this can degrade the output quality.  3. Maximum file size: 20 MB. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | The resolution tier of the generated video. Options: 720p; 1080p |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `3:4`, `4:3`, `4:5`, `5:4`, `1:1`, `9:21`, `21:9`; default: `"16:9"` | The aspect ratio of the generated video. Options: 16:9; 9:16; 3:4; 4:3; 4:5; 5:4; 1:1; 9:21; 21:9 |
| `duration` | number | no | default: `5`; min: 3; max: 15 | The duration of the generated video, in seconds. Value range: An integer from 3 to 15. |

Example `input`:

```json
{
  "reference_image": [
    "https://static.aiquickdraw.com/tools/example/1782114387854_IufKnPxR.png"
  ],
  "prompt": "A cat running on the grass",
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 5
}
```

## `happyhorse-1-1/text-to-video`

**HappyHorse-1-1 text-to-video** · Happyhorse 1 1 · [official docs](https://docs.kie.ai/market/happyhorse-1-1/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | default: `""`; max length: 4999 | Text description of the video to generate.  Supports any language. Maximum 5,000 non-Chinese characters or 2,500 Chinese |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution. Options: 720p; 1080p |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `4:5`, `5:4`, `9:21`, `21:9`; default: `"16:9"` | Output video aspect ratio. Options: 16:9; 9:16; 1:1; 4:3; 3:4; 4:5; 5:4; 9:21; 21:9 |
| `duration` | number | no | default: `5`; min: 3; max: 15 | Output video duration in seconds. |

Example `input`:

```json
{
  "prompt": "A dog running on the earth",
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 5
}
```

## `happyhorse/image-to-video`

**HappyHorse - image-to-video** · Happyhorse · [official docs](https://docs.kie.ai/market/happyhorse/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 5000 | Text prompt describing the video to generate (any language). Max 5,000 non‑Chinese characters or 2,500 Chinese characters; extra content is truncated. |
| `image_urls` | array | yes | max items: 1 | First-frame image URL list. Exactly one image is required.  Image constraints: Format: JPEG, JPG, PNG, WEBP. Resolution: Width and height must be at least 300 pixels. Aspect ratio: 1:2.5 to 2.5:1. File size: Up to 10 MB. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution. Valid values: 720P, 1080P (default). |
| `duration` | integer | no | default: `5`; min: 3; max: 15 | Output duration in seconds. Must be between 3 and 15. Defaults to 5. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility (if supported). |

Example `input`:

```json
{
  "prompt": "A cat running on the grass",
  "image_urls": [
    "https://loremflickr.com/400/400?lock=4153750340434616"
  ],
  "resolution": "1080p",
  "duration": 5,
  "seed": 1546095068
}
```

## `happyhorse/reference-to-video`

**HappyHorse - reference-to-video** · Happyhorse · [official docs](https://docs.kie.ai/market/happyhorse/reference-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the video to generate (any language). Max 5,000 non‑Chinese characters or 2,500 Chinese characters; extra content is truncated. |
| `reference_image` | array | yes | min items: 1; max items: 9 | Reference image URL list. Provide 1–9 images. The order defines which image is character1, character2, etc.  Image limits: Format: JPEG, JPG, PNG, and WEBP. Resolution: shortest side at least 400 px. 720P or higher recommended. Avoid small, blurry, or heavily compressed images, as they degrade output quality. File size: 10 MB maximum. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution. Valid values: 720P, 1080P (default). |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`; default: `"16:9"` | Output aspect ratio. Valid values: 16:9 (default), 9:16, 1:1, 4:3, 3:4. |
| `duration` | integer | no | default: `5`; min: 3; max: 15 | Output duration in seconds (integer). Must be between 3 and 15. Defaults to 5. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility (if supported). |

Example `input`:

```json
{
  "prompt": "A woman in a red qipao character1. The shot opens with a side medium view outlining the tailored fit of the qipao and S-curve silhouette, then cuts to a low-angle shot capturing her gracefully unfolding a folding fan character2, with tassel earrings character3 swaying lightly as she turns her head. Finally, the camera pushes into a facial close-up, freezing on her fingertips lightly touching the fan ribs and the subtle, reserved charm in her expressive gaze. Through multiple angles, it comprehensively showcases an aura of Eastern elegance.",
  "reference_image": [
    "https://loremflickr.com/400/400?lock=8132663902229376",
    "https://loremflickr.com/400/400?lock=1716016437146867"
  ],
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 5,
  "seed": 1308038620
}
```

## `happyhorse/text-to-video`

**HappyHorse - text-to-video** · Happyhorse · [official docs](https://docs.kie.ai/market/happyhorse/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the video to generate (any language). Max 5,000 non‑Chinese characters or 2,500 Chinese characters; extra content is truncated. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution. Valid values: 720P, 1080P (default). |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`; default: `"16:9"` | Output aspect ratio. Valid values: 16:9 (default), 9:16, 1:1, 4:3, 3:4. |
| `duration` | integer | no | default: `5`; min: 3; max: 15 | Output duration in seconds (integer). Must be between 3 and 15. Defaults to 5. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility (if supported). |

Example `input`:

```json
{
  "prompt": "A miniature city built from cardboard and bottle caps comes to life at night. A cardboard train slowly passes through, with small lights dotting the scene and illuminating the way ahead.",
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 5,
  "seed": 1622429582
}
```

## `happyhorse/video-edit`

**HappyHorse - video-edit** · Happyhorse · [official docs](https://docs.kie.ai/market/happyhorse/video-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Required edit instruction describing the intended change (e.g., style transfer / local replacement). Max 5,000 non‑Chinese characters or 2,500 Chinese characters; extra content is truncated. |
| `video_url` | string | yes |  | Input video URL list. Exactly one video is required. Video requirements: Format: MP4, MOV (H.264 encoding recommended). Duration: 3–60 seconds. Resolution: the longer side must not exceed 2,160 px; the shorter side must be at least 320 px. Aspect ratio: 1:2.5–2.5:1. File size: up to 100 MB. Frame rate: greater than 8 fps. |
| `reference_image` | array | no | min items: 0; max items: 5 | Optional reference image URL list (0–5). Image requirements: Format: JPEG, JPG, PNG, WEBP. Resolution: both width and height must be at least 300 px. Aspect ratio: 1:2.5–2.5:1. File size: up to 10 MB. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution. Valid values: 720P, 1080P (default). |
| `audio_setting` | string | no | allowed: `auto`, `origin`; default: `"auto"` | Audio handling strategy for the output video. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility (if supported). |

Example `input`:

```json
{
  "prompt": "Make the horse-headed humanoid character in the video wear the striped sweater from the image",
  "video_url": "https://hollow-joy.info/",
  "reference_image ": [
    "https://loremflickr.com/400/400?lock=3320229742640740",
    "https://loremflickr.com/400/400?lock=390084853871038",
    "https://loremflickr.com/400/400?lock=4205160298467577",
    "https://loremflickr.com/400/400?lock=7626507781317900",
    "https://loremflickr.com/400/400?lock=2804855355708229"
  ],
  "resolution": "1080p",
  "audio_setting": "auto",
  "seed": 1764574909
}
```

## `ideogram/character`

**Ideogram - Character** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/character)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to fill the masked part of the image. (Max length: 5000 characters) |
| `reference_image_urls` | array | yes |  | A set of images to use as character references. Currently only 1 image is supported, rest will be ignored. (maximum total size 10MB across all character references). The images should be in JPEG, PNG or WebP format (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY`; default: `"BALANCED"` | The rendering speed to use. Default value: "BALANCED" |
| `style` | string | no | allowed: `AUTO`, `REALISTIC`, `FICTION`; default: `"AUTO"` | The style type to generate with. Cannot be used with style_codes. Default value: "AUTO" |
| `expand_prompt` | boolean | no |  | Determine if MagicPrompt should be used in generating the request or not. Default value: true (Boolean value (true/false)) |
| `num_images` | string | no | allowed: `1`, `2`, `3`, `4`; default: `"1"` | Select description |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9`; default: `"square_hd"` | The resolution of the generated image Default value: square_hd |
| `seed` | integer | no |  | Seed for the random number generator |
| `negative_prompt` | string | no | max length: 5000 | Description of what to exclude from an image. Descriptions in the prompt take precedence to descriptions in the negative prompt. Default value: "" (Max length: 5000 characters) |

Example `input`:

```json
{
  "prompt": "Place the woman from the uploaded portrait, wearing a casual white blouse, in a peaceful garden setting. The scene should feature vibrant green plants and colorful flowers, with soft sunlight filtering through the leaves. She should be sitting on a wooden bench, holding a book and smiling gently. The background should be filled with lush greenery, with a serene, tranquil atmosphere. Golden afternoon light should highlight the woman’s face and create soft shadows on the ground, adding a peaceful, reflective mood to the scene",
  "reference_image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1755767145415pvz49dpi.webp"
  ],
  "rendering_speed": "BALANCED",
  "style": "AUTO",
  "expand_prompt": true,
  "num_images": "1",
  "image_size": "square_hd",
  "negative_prompt": ""
}
```

## `ideogram/character-edit`

**Ideogram - Character Edit** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/character-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to fill the masked part of the image. (Max length: 5000 characters) |
| `image_url` | string | yes |  | The image URL to generate an image from. Needs to match the dimensions of the mask. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `mask_url` | string | yes |  | The mask URL to inpaint the image. Needs to match the dimensions of the input image. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `reference_image_urls` | array | yes |  | A set of images to use as character references. Currently only 1 image is supported, rest will be ignored. (maximum total size 10MB across all character references). The images should be in JPEG, PNG or WebP format (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY`; default: `"BALANCED"` | The rendering speed to use. Default value: "BALANCED" |
| `style` | string | no | allowed: `AUTO`, `REALISTIC`, `FICTION`; default: `"AUTO"` | The style type to generate with. Cannot be used with style_codes. Default value: "AUTO" |
| `expand_prompt` | boolean | no |  | Determine if MagicPrompt should be used in generating the request or not. Default value: true (Boolean value (true/false)) |
| `num_images` | string | no | allowed: `1`, `2`, `3`, `4`; default: `"1"` | Select description |
| `seed` | integer | no |  | Seed for the random number generator |

Example `input`:

```json
{
  "prompt": "A fabulous look head tilted down, looking forward with a smile\n",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17557680349256sa0lk53.webp",
  "mask_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755768046014ftgvma28.webp",
  "reference_image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1755768064644jodsmfhq.webp"
  ],
  "rendering_speed": "BALANCED",
  "style": "AUTO",
  "expand_prompt": true,
  "num_images": "1"
}
```

## `ideogram/character-remix`

**Ideogram - Character Remix** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/character-remix)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to remix the image with (Max length: 5000 characters) |
| `image_url` | string | yes |  | The image URL to remix (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `reference_image_urls` | array | yes |  | A set of images to use as character references. Currently only 1 image is supported, rest will be ignored. (maximum total size 10MB across all character references). The images should be in JPEG, PNG or WebP format (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY`; default: `"BALANCED"` | The rendering speed to use. Default value: "BALANCED" |
| `style` | string | no | allowed: `AUTO`, `REALISTIC`, `FICTION`; default: `"AUTO"` | The style type to generate with. Cannot be used with style_codes. Default value: "AUTO" |
| `expand_prompt` | boolean | no |  | Determine if MagicPrompt should be used in generating the request or not. Default value: true (Boolean value (true/false)) |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9`; default: `"square_hd"` | Select description |
| `num_images` | string | no | allowed: `1`, `2`, `3`, `4`; default: `"1"` | Select description |
| `seed` | integer | no |  | Seed for the random number generator |
| `strength` | number | no | default: `0.8`; min: 0.1; max: 1 | Strength of the input image in the remix Default value: 0.8 (Min: 0.1, Max: 1, Step: 0.1) (step: 0.1) |
| `negative_prompt` | string | no | max length: 500 | Description of what to exclude from an image. Descriptions in the prompt take precedence to descriptions in the negative prompt. Default value: "" (Max length: 500 characters) |
| `image_urls` | array | no |  | A set of images to use as style references (maximum total size 10MB across all style references). The images should be in JPEG, PNG or WebP format (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `reference_mask_urls` | string | no |  | A set of masks to apply to the character references. Currently only 1 mask is supported, rest will be ignored. (maximum total size 10MB across all character references). The masks should be in JPEG, PNG or WebP format (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |

Example `input`:

```json
{
  "prompt": "A fisheye lens selfie photograph taken at night on an urban street. The image is circular with a black border and shows a person wearing dark sunglasses and a black jacket, holding a silver digital camera up to capture the reflection. The background shows a row of shuttered storefronts with red neon lighting visible in the upper portion. The street is empty and dark, with street lights creating a warm glow along the sidewalk. The fisheye effect creates a curved, distorted perspective that bends the straight lines of the street and buildings. The lighting is predominantly red and dark, creating a moody urban atmosphere. The person's reflection shows long dark hair and is positioned in the center of the circular frame. Multiple storefront shutters are visible in the background, creating a repeating pattern of horizontal lines. The overall composition has a cinematic quality with strong contrast between the dark street and the illuminated storefronts above.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755768466167d0tiuc6e.webp",
  "reference_image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1755768479029sugx0g6f.webp"
  ],
  "rendering_speed": "BALANCED",
  "style": "AUTO",
  "expand_prompt": true,
  "image_size": "square_hd",
  "num_images": "1",
  "strength": 0.8,
  "negative_prompt": "",
  "image_urls": [],
  "reference_mask_urls": ""
}
```

## `ideogram/v3-edit`

**Ideogram V3 Edit** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/v3-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to fill the masked part of the image. Maximum length: 5000 characters. |
| `image_url` | string | yes | format: uri | The image URL to generate an image from. Needs to match the dimensions of the mask.  - Please provide the URL of the uploaded file, not raw file content - Accepted types: `image/jpeg`, `image/png`, `image/webp` - Max size: 10.0MB |
| `mask_url` | string | yes | format: uri | The mask URL to inpaint the image. Needs to match the dimensions of the input image.  - Please provide the URL of the uploaded file, not raw file content - Accepted types: `image/jpeg`, `image/png`, `image/webp` - Max size: 10.0MB |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY`; default: `"BALANCED"` | The rendering speed to use. Default value: `BALANCED`.  - `TURBO`: Turbo - `BALANCED`: Balanced - `QUALITY`: Quality |
| `expand_prompt` | boolean | no | default: `true` | Determine if MagicPrompt should be used in generating the request or not. Default value: `true`.  - Boolean value: `true` / `false` |
| `seed` | integer | no |  | Seed for the random number generator. |

Example `input`:

```json
{
  "prompt": "A dog wearing a cowboy hat",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755076859801ryyol1du.webp",
  "mask_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755076871089hx9uonhc.webp",
  "rendering_speed": "BALANCED",
  "expand_prompt": true,
  "seed": 123456
}
```

## `ideogram/v3-remix`

**Ideogram V3 Remix** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/v3-remix)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to remix the image with. Maximum length: 5000 characters. |
| `image_url` | string | yes | format: uri | The image URL to remix.  - Please provide the URL of the uploaded file, not raw file content - Accepted types: `image/jpeg`, `image/png`, `image/webp` - Max size: 10.0MB |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY` | The rendering speed to use.  - `TURBO`: Turbo - `BALANCED`: Balanced - `QUALITY`: Quality |
| `style` | string | no | allowed: `AUTO`, `GENERAL`, `REALISTIC`, `DESIGN` | The style type to generate with. Cannot be used together with `style_codes`.  - `AUTO`: Auto - `GENERAL`: General - `REALISTIC`: Realistic - `DESIGN`: Design |
| `expand_prompt` | boolean | no |  | Determine if MagicPrompt should be used in generating the request or not.  - Boolean value: `true` / `false` |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9` | The resolution of the generated image.  - `square`: Square - `square_hd`: Square HD - `portrait_4_3`: Portrait 3:4 - `portrait_16_9`: Portrait 9:16 - `landscape_4_3`: Landscape 4:3 - `landscape_16_9`: Landscape 16:9 |
| `num_images` | string | no | allowed: `1`, `2`, `3`, `4` | Number of images to generate.  - `1`: 1 image - `2`: 2 images - `3`: 3 images - `4`: 4 images |
| `seed` | integer | no |  | Seed for the random number generator. |
| `strength` | number | no | min: 0.01; max: 1 | Strength of the input image in the remix.  - Minimum: `0.01` - Maximum: `1` - Step: `0.01` |
| `negative_prompt` | string | no | max length: 5000 | Description of what to exclude from the generated image. If the positive prompt conflicts with the negative prompt, the positive prompt takes precedence. Maximum length: 5000 characters. |

Example `input`:

```json
{
  "prompt": "Change the cube into a sphere",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17550782013854ykfihxv.webp",
  "rendering_speed": "BALANCED",
  "style": "AUTO",
  "expand_prompt": true,
  "image_size": "square_hd",
  "num_images": "1",
  "seed": 123456,
  "strength": 0.8,
  "negative_prompt": "blurry, low quality, distorted, watermark"
}
```

## `ideogram/v3-text-to-image`

**Ideogram V3 Text to Image** · Ideogram · [official docs](https://docs.kie.ai/market/ideogram/v3-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Description of the image to generate. Maximum length: 5000 characters. |
| `rendering_speed` | string | no | allowed: `TURBO`, `BALANCED`, `QUALITY` | The rendering speed to use.  - `TURBO`: Turbo - `BALANCED`: Balanced - `QUALITY`: Quality |
| `style` | string | no | allowed: `AUTO`, `GENERAL`, `REALISTIC`, `DESIGN` | The style type to generate with. Cannot be used together with `style_codes`.  - `AUTO`: Auto - `GENERAL`: General - `REALISTIC`: Realistic - `DESIGN`: Design |
| `expand_prompt` | boolean | no |  | Determines whether MagicPrompt should be used to enhance the generation request.  - Boolean value: `true` / `false` |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9` | The resolution of the generated image.  - `square`: Square - `square_hd`: Square HD - `portrait_4_3`: Portrait 3:4 - `portrait_16_9`: Portrait 9:16 - `landscape_4_3`: Landscape 4:3 - `landscape_16_9`: Landscape 16:9 |
| `seed` | integer | no |  | Seed for the random number generator. |
| `negative_prompt` | string | no | max length: 5000 | Description of what to exclude from the generated image. If the positive prompt conflicts with the negative prompt, the positive prompt takes precedence. Maximum length: 5000 characters. |

Example `input`:

```json
{
  "prompt": "A cinematic photograph of a tranquil lakeside at twilight, viewed from a slight elevation. In the center, a cluster of softly glowing reeds and water lilies emit a gentle golden light, their reflections shimmering on the calm surface. The elegant neon-style white text 'Kie.ai' hovers just above the water, subtly illuminated and harmonizing with the natural glow. Surrounding willows and drifting mist frame the scene, creating a serene yet magical atmosphere, with warm highlights contrasting against the cool blues of the evening sky.",
  "rendering_speed": "BALANCED",
  "style": "AUTO",
  "expand_prompt": true,
  "image_size": "square_hd",
  "seed": 123456,
  "negative_prompt": "blurry, low detail, distorted anatomy, extra limbs, malformed text, watermark"
}
```

## `infinitalk/from-audio`

**Infinitalk - From Audio** · Infinitalk · [official docs](https://docs.kie.ai/market/infinitalk/from-audio)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes |  | URL of the input image. If the input image does not match the chosen aspect ratio, it is resized and center cropped. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `audio_url` | string | yes |  | The URL of the audio file. (File URL after upload, not file content; Accepted types: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/mp4, audio/ogg; Max size: 10.0MB) |
| `prompt` | string | yes | max length: 5000 | The text prompt to guide video generation. (Max length: 5000 characters) |
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"480p"` | Resolution of the video to generate. Must be either 480p or 720p. |
| `seed` | number | no |  | Random seed for reproducibility. Valid range is 10000 to 1000000. |

Example `input`:

```json
{
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1757329269873ggqj2hz3.png",
  "audio_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1757329255705mmqwrnri.mp3",
  "prompt": "A young woman with long dark hair talking on a podcast.",
  "resolution": "480p"
}
```

## `kling-2.6/image-to-video`

**Kling 2.6 Image to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1000 | Text prompt for video generation (maximum length: 1000 characters) |
| `image_urls` | array | yes | max items: 1 | Image URLs for video generation. (Uploaded file URLs, not file content; supported types: image/jpeg, image/png; maximum file size: 10.0MB) |
| `sound` | boolean | yes |  | This parameter specifies whether the generated video contains sound (boolean: true/false) |
| `duration` | string | yes | allowed: `5`, `10`; default: `"5"` | Video duration (unit: seconds) |

Example `input`:

```json
{
  "prompt": "In a bright rehearsal room, sunlight streams through the windows, and a standing microphone is placed in the center of the room. [Campus band female lead singer] stands in front of the microphone with her eyes closed, and other members stand around her. [Campus band female lead singer, singing loudly] Lead vocal: \"I will do my best to heal you, with all my heart and soul...\" The background is a cappella harmonies, and the camera slowly pans around the band members.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1764851002741_i0lEiI8I.png"
  ],
  "sound": false,
  "duration": "5"
}
```

## `kling-2.6/motion-control`

**Kling 2.6 motion-control** · Kling · [official docs](https://docs.kie.ai/market/kling/motion-control)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 2500 | A text description of the desired output. Maximum length is 2500 characters. (Max length: 2500 characters) |
| `input_urls` | array | yes | max items: 1 | An array containing a single image URL. The photo must clearly show the subject's head, shoulders, and torso. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/jpg; Max size: 10.0MB,size needs to be greater than 300px, aspect ratio 2:5 to 5:2.) |
| `video_urls` | array | yes | max items: 1 | An array containing a single video URL. The duration must be between 3 to 30 seconds, and the video must clearly show the subject's head, shoulders, and torso. (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime; Max size: 100.0MB) |
| `character_orientation` | string | yes | allowed: `image`, `video`; default: `"video"` | Generate the orientation of the characters in the video. 'image': same orientation as the person in the picture (max 10s video). 'video': consistent with the orientation of the characters in the video (max 30s video). |
| `mode` | string | yes | allowed: `720p`, `1080p`; default: `"720p"` | Output resolution mode. Use 'std' for 720p or 'pro' for 1080p. |

Example `input`:

```json
{
  "prompt": "The cartoon character is dancing.",
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1767694885407_pObJoMcy.png"
  ],
  "video_urls": [
    "https://static.aiquickdraw.com/tools/example/1767525918769_QyvTNib2.mp4"
  ],
  "mode": "720p",
  "character_orientation": "image"
}
```

## `kling-2.6/text-to-video`

**Kling 2.6 Text to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1000 | Text prompt for video generation (maximum length: 1000 characters) |
| `sound` | boolean | yes |  | This parameter specifies whether the generated video contains sound (boolean: true/false) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `16:9`, `9:16`; default: `"1:1"` | This parameter defines the video aspect ratio |
| `duration` | string | yes | allowed: `5`, `10`; default: `"5"` | Video duration (unit: seconds) |

Example `input`:

```json
{
  "prompt": "Scene: A fashion live-streaming sales setting, with clothes hanging on racks and the host's figure reflected in a full-length mirror. Lines: [African female host] turns around to showcase the hoodie's cut. [African female host, in a cheerful tone] says: \"360-degree flawless tailoring, slimming and versatile.\" She then [African female host] leans closer to the camera. [African female host, in a lively tone] says: \"Double-sided fleece fabric, $30 off immediately when you order now.\"",
  "sound": false,
  "aspect_ratio": "1:1",
  "duration": "5"
}
```

## `kling-3.0/motion-control`

**Kling-3.0 motion-control** · Kling · [official docs](https://docs.kie.ai/market/kling/motion-control-v3)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no |  | (Optional) Text prompt words, used to guide the generation of animation content. Can be empty or 0 - 2500 characters long. |
| `input_urls` | array | yes |  | (Required) Include a URL of an image |
| `video_urls` | array | yes |  | (Required) Include a video URL |
| `mode` | string | no |  | (Optional) Video Quality Mode. std: Standard Mode (720p). pro: Professional Mode (1080p) |
| `character_orientation` | string | no |  | (Optional) Reference source for character orientation. video: Refer to video (recommended); image: Refer to image. Default value: video |
| `background_source` | string | no |  | (Optional) Background source. input_video: Use video background; input_image: Use image background. Default value: input_video |

Example `input`:

```json
{
  "prompt": "The cartoon character is dancing.",
  "input_urls": [
    "https://static.aiquickdraw.com/tools/example/1767694885407_pObJoMcy.png"
  ],
  "video_urls": [
    "https://static.aiquickdraw.com/tools/example/1767525918769_QyvTNib2.mp4"
  ],
  "mode": "720p",
  "character_orientation": "image",
  "background_source": "input_video"
}
```

## `kling-3.0/video`

**Kling 3.0** · Kling · [official docs](https://docs.kie.ai/market/kling/kling-3-0)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes |  | Video generation prompt. Takes effect when multi_shots is false. |
| `image_urls` | array | no |  | First and last frame image URLs. Required when elements are referenced in the prompt (using @element_name syntax). When multi_shots is false: if length is 2, index 0 is the first frame and index 1 is the last frame; if length is 1, the array item serves as the first frame. When multi_shots is true: only the first frame is supported. |
| `sound` | boolean | yes | default: `false` | Whether to enable sound effects. true enables sound effects, false disables them. When multi_shots is true, this field defaults to true. |
| `duration` | string | yes | allowed: `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`; default: `"5"` | Total video duration in seconds. Integer value, range: 3 to 15. |
| `aspect_ratio` | string | yes | allowed: `16:9`, `9:16`, `1:1`; default: `"16:9"` | Video aspect ratio. Options: 16:9, 9:16, 1:1. When image_urls(first and last frame images) is provided, this parameter is optional and the aspect ratio will be automatically adapted based on the uploaded images. |
| `mode` | string | yes | allowed: `std`, `pro`, `4K`; default: `"pro"` | Generation mode. std has standard resolution, pro has higher resolution, 4K has 4K resolution.  Resolution mapping: - **std mode**: 16:9 (1280×720), 9:16 (720×1280), 1:1 (720×720) - **pro mode**: 16:9 (1920×1080), 9:16 (1080×1920), 1:1 (1080×1080) - **4K mode**: 16:9 (3840×2160), 9:16 (2160×3840), 1:1 (2160×2160) |
| `multi_shots` | boolean | yes | default: `false` | Whether to use multi-shot mode. true enables multi-shot mode, false enables single-shot mode. |
| `multi_prompt` | array | yes |  | Shot prompts. Takes effect when multi_shots is true. Used to describe the text and duration of each shot. Supports up to 5 shots. Each shot duration is 1-12 seconds. If you need to use elements, add them after the prompt. |
| `multi_prompt[].prompt` | string | yes | max length: 500 | Prompt text for this shot, a maximum of 500 characters per shot. Each @element will occupy 37 characters. |
| `multi_prompt[].duration` | integer | yes | min: 1; max: 12 | Duration of this shot in seconds. Range: 1-12. |
| `kling_elements` | array | no | max items: 3 | Referenced elements. Detailed information about elements referenced in the prompt. A single task can reference a maximum of three elements. |
| `kling_elements[].name` | string | no |  | Element name, used in prompt with @ prefix (e.g., @element_dog) |
| `kling_elements[].description` | string | no |  | Element description |
| `kling_elements[].element_input_urls` | array | no |  | Image URLs for the element. 2-4 URLs required. Accepted formats: JPG, PNG. Maximum file size: 10MB per image. |
| `kling_elements[].element_input_audio_urls` | array | no |  | Optional. List of audio material URLs for characters. The audio duration must be between 5 and 30 seconds. |
| `kling_elements[].start_time` | integer | no |  | Start time for video character material capture (in milliseconds). Only effective when uploading videos through element_input_urls. If not uploaded, it defaults to 0. |
| `kling_elements[].end_time` | integer | no |  | End time for video character material capture (in milliseconds). Only effective when uploading videos through element_input_urls. Must be greater than start_time, and the difference between end_time and start_time must be within 3000 to 8000 milliseconds. |

Example `input`:

```json
{
  "prompt": "In a bright rehearsal room, sunlight streams through the window @element_dog",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1764851002741_i0lEiI8I.png"
  ],
  "sound": true,
  "duration": "5",
  "aspect_ratio": "16:9",
  "mode": "pro",
  "multi_shots": false,
  "multi_prompt": [
    {
      "prompt": "a happy dog in running @element_cat",
      "duration": 3
    },
    {
      "prompt": "a happy dog play with a cat @element_dog",
      "duration": 2
    }
  ],
  "kling_elements": [
    {
      "name": "element_dog",
      "description": "dog",
      "element_input_urls": [
        "https://tempfileb.aiquickdraw.com/kieai/market/1770361808044_4RfUUJrI.jpeg",
        "https://tempfileb.aiquickdraw.com/kieai/market/1770361848336_ABQqRHBi.png"
      ],
      "element_input_audio_urls": [
        "https://your-cdn.com/wjeoiajfosijfoi.mp3"
      ]
    },
    {
      "name": "element_cat",
      "description": "cat",
      "element_input_urls": [
        "https://your-cdn.com/element_image.mp4"
      ],
      "element_input_audio_urls": [
        "https://your-cdn.com/wjeoiajfosijfoi.mp3"
      ],
      "start_time": 0,
      "end_time": 8000
    }
  ]
}
```

## `kling/ai-avatar-pro`

**Kling AI Avatar Pro** · Kling · [official docs](https://docs.kie.ai/market/kling/ai-avatar-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes |  | The URL of the image to use as your avatar (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `audio_url` | string | yes |  | The URL of the audio file (must be the URL of the uploaded file, not the file content; supported formats: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/mp4, audio/ogg; audio size is limited to 100M, and the duration cannot exceed 5 minutes) |
| `prompt` | string | yes | max length: 5000 | The prompt to use for the video generation (Max length: 5000 characters) |

Example `input`:

```json
{
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/175792685809077e8h8k3.png",
  "audio_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1757925802302srqfkcqh.mp3",
  "prompt": ""
}
```

## `kling/ai-avatar-standard`

**Kling AI Avatar Standard** · Kling · [official docs](https://docs.kie.ai/market/kling/ai-avatar-standard)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes |  | The URL of the image to use as your avatar (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `audio_url` | string | yes |  | The URL of the audio file (must be the URL of the uploaded file, not the file content; supported formats: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/mp4, audio/ogg; audio size is limited to 100M, and the duration cannot exceed 5 minutes) |
| `prompt` | string | yes | max length: 5000 | The prompt to use for the video generation (Max length: 5000 characters) |

Example `input`:

```json
{
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17579268936223zs9l3dt.png",
  "audio_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17579258340109gghun47.mp3",
  "prompt": ""
}
```

## `kling/v2-1-master-image-to-video`

**Kling V2.1 Master Image to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/v2-1-master-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt describing the video to generate (Max length: 5000 characters) |
| `image_url` | string | yes |  | URL of the image to be used for the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `negative_prompt` | string | no | max length: 500 | Negative prompt to exclude certain elements from the video (Max length: 500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |

Example `input`:

```json
{
  "prompt": "A team of paratroopers descends into enemy territory, as they pass through clouds, the camera switches to a slow pan above the battlefield lighting up with",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755256297923kmjpynul.png",
  "duration": "5",
  "negative_prompt": "blur, distort, and low quality",
  "cfg_scale": 0.5
}
```

## `kling/v2-1-master-text-to-video`

**Kling V2.1 Master Text to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/v2-1-master-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt describing the video you want to generate (Max length: 5000 characters) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`; default: `"16:9"` | The aspect ratio of the generated video frame |
| `negative_prompt` | string | no | max length: 500 | Elements to avoid in the generated video (Max length: 500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |

Example `input`:

```json
{
  "prompt": "First-person view from a soldier jumping from a transport plane — the camera shakes with turbulence, oxygen mask reflections flicker — as the clouds part, the battlefield below pulses with anti-air fire and missile trails.",
  "duration": "5",
  "aspect_ratio": "16:9",
  "negative_prompt": "blur, distort, and low quality",
  "cfg_scale": 0.5
}
```

## `kling/v2-1-pro`

**Kling V2.1 Pro** · Kling · [official docs](https://docs.kie.ai/market/kling/v2-1-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the video to generate (Max length: 5000 characters) |
| `image_url` | string | yes |  | URL of the image to be used for the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `negative_prompt` | string | no | max length: 500 | Terms to avoid in the generated video (Max length: 500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |
| `tail_image_url` | string | no |  | URL of the image to be used for the end of the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |

Example `input`:

```json
{
  "prompt": "POV shot of a gravity surfer diving between ancient ruins suspended midair, glowing moss lights the path, the board hisses as it carves through thin mist, echoes rise with speed ",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1754892534386c8wt0qfs.png",
  "duration": "5",
  "negative_prompt": "blur, distort, and low quality",
  "cfg_scale": 0.5,
  "tail_image_url": ""
}
```

## `kling/v2-1-standard`

**Kling V2.1 Standard** · Kling · [official docs](https://docs.kie.ai/market/kling/v2-1-standard)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt describing the desired video content (Max length: 5000 characters) |
| `image_url` | string | yes |  | URL of the image to be used for the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `negative_prompt` | string | no | max length: 500 | Description of elements to avoid in the generated video (Max length: 500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |

Example `input`:

```json
{
  "prompt": "Begin with the uploaded image as the first frame. Gradually animate the scene: steam rises and drifts upward from the train; lantern lights flicker subtly; cloaked figures begin to move slowly — walking, turning, adjusting their belongings. Floating dust or magical particles catch the light. The text “KLING 2.1 STANDARD API — Now on Kie.ai” softly pulses with a golden glow. The camera pushes forward slightly, then slowly fades to black.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755256596169mkkwr2ag.png",
  "duration": "5",
  "negative_prompt": "blur, distort, and low quality",
  "cfg_scale": 0.5
}
```

## `kling/v2-5-turbo-image-to-video-pro`

**Kling - V2.5 Turbo Image to Video Pro** · Kling · [official docs](https://docs.kie.ai/market/kling/v25-turbo-image-to-video-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 2500 | The text prompt describing the video to generate (Max length: 2500 characters) |
| `image_url` | string | yes |  | URL of the image to be used for the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `tail_image_url` | string | no |  | Tail frame image of video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `negative_prompt` | string | no | max length: 500 | Negative prompt to exclude certain elements from the video (Max length: 500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |

## `kling/v2-5-turbo-text-to-video-pro`

**Kling - V2.5 Turbo Text to Video Pro** · Kling · [official docs](https://docs.kie.ai/market/kling/v25-turbo-text-to-video-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 2500 | The text description of the video you want to generate (Max length: 2500 characters) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`; default: `"16:9"` | The aspect ratio of the generated video frame |
| `negative_prompt` | string | no | max length: 2500 | Things to avoid in the generated video (Max length: 2500 characters) |
| `cfg_scale` | number | no | default: `0.5`; min: 0; max: 1 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt (Min: 0, Max: 1, Step: 0.1) (step: 0.1) |

Example `input`:

```json
{
  "prompt": "Real-time playback. Wide shot of a ruined city: collapsed towers, fires blazing, storm clouds with lightning. Camera drops fast from the sky over burning streets and tilted buildings. Smoke and dust fill the air. A lone hero walks out of the ruins, silhouetted by fire. Camera shifts front: his face is dirty with dust and sweat, eyes firm, a faint smile. Wind blows, debris rises. Extreme close-up: his eyes reflect the approaching enemy. Music and drums hit. Final wide shot: fire forms a blazing halo behind him — reborn in flames with epic cinematic vibe.",
  "duration": "5",
  "aspect_ratio": "16:9",
  "negative_prompt": "blur, distort, and low quality",
  "cfg_scale": 0.5
}
```

## `kling/v3-turbo-image-to-video`

**Kling - V3 Turbo Image to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/v3-turbo-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 2500 | The text prompt describing the video to generate (Max length: 2500 characters) |
| `image_urls` | array | yes |  | URL of the image to be used for the video (File URL after upload, not file content; Accepted types: image/jpeg, image/png; Max size: 10.0MB) |
| `duration` | string | yes | default: `"5"` | The duration of the generated video in seconds Optional duration: 3s - 15s |
| `resolution` | string | yes | allowed: `720p`, `1080p`; default: `"720p"` | Resolution of the generated video (720p or 1080p) |

Example `input`:

```json
{
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1770688028208_jxcvxCQm.png"
  ],
  "prompt": "Outdoor terrace of a European villa, by a dining table with a blue and white checkered tablecloth, a young white woman in a blue and white striped short-sleeve shirt and khaki shorts, with a brown belt, sits barefoot, opposite a young white man in a white T-shirt.\n\nThe camera zooms in, the woman swirls the juice in a glass, her eyes looking at the distant woods, and says, \"These trees will turn yellow in a month, won't they?\"\n\nClose-up of the man, he lowers his head and says, \"But they'll be green again next summer.\"\n\nThen the woman turns her head, smiles at the man opposite, and says, \"Are you always this optimistic? Or just about summer?\"\n\nThen the man lifts his head, looks at the woman and says, \"Only about summers with you.\"",
  "duration": "5",
  "resolution": "720p"
}
```

## `kling/v3-turbo-text-to-video`

**Kling - V3 Turbo Text to Video** · Kling · [official docs](https://docs.kie.ai/market/kling/v3-turbo-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 2500 | The text description of the video you want to generate (Max length: 2500 characters) |
| `duration` | string | yes | default: `"5"` | The duration of the generated video in seconds Optional duration: 3s - 15s |
| `aspect_ratio` | string | yes | allowed: `1:1`, `9:16`, `16:9`; default: `"16:9"` | The aspect ratio of the generated video frame |
| `resolution` | string | yes | allowed: `720p`, `1080p`; default: `"720p"` | Resolution of the generated video (720p or 1080p) |

Example `input`:

```json
{
  "prompt": "Outdoor terrace of a European villa, by a dining table with a blue and white checkered tablecloth, a young white woman in a blue and white striped short-sleeve shirt and khaki shorts, with a brown belt, sits barefoot, opposite a young white man in a white T-shirt.\n\nThe camera zooms in, the woman swirls the juice in a glass, her eyes looking at the distant woods, and says, \"These trees will turn yellow in a month, won't they?\"\n\nClose-up of the man, he lowers his head and says, \"But they'll be green again next summer.\"\n\nThen the woman turns her head, smiles at the man opposite, and says, \"Are you always this optimistic? Or just about summer?\"\n\nThen the man lifts his head, looks at the woman and says, \"Only about summers with you.\"",
  "duration": "5",
  "aspect_ratio": "16:9",
  "resolution": "720p"
}
```

## `minimax-h3/image-to-video`

**MiniMax H3 Image-to-Video** · Minimax H3 · [official docs](https://docs.kie.ai/market/minimax-h3/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 1; max length: 7000 | Video generation prompt, length is between 1 and 7000 characters. |
| `first_frame_url` | string | no | format: uri | First frame image URL. Note: Either first_frame_url or last_frame_url must be provided. Supports HTTP, HTTPS, and OSS addresses; supports JPG, JPEG, PNG, WEBP, HEIC, HEIF formats; single image size not exceeding 30 MB, image side length between 256 and 5760 pixels, aspect ratio between 0.4 and 2.5. |
| `last_frame_url` | string | no | format: uri | Last frame image URL. Note: Either first_frame_url or last_frame_url must be provided. Image restrictions are the same as first_frame_url. |
| `duration` | integer | yes | allowed: `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`; default: `6`; min: 4; max: 15 | Generated video duration, supporting integer values from 4 to 15 seconds. |
| `resolution` | string | no | allowed: `768P`, `2K`; default: `"2K"` | Video resolution. |

Example `input`:

```json
{
  "prompt": "Let the character in the scene turn around naturally and smile, camera slowly pushing forward",
  "first_frame_url": "https://example.com/first-frame.jpg",
  "last_frame_url": "https://example.com/last-frame.jpg",
  "duration": 6
}
```

## `minimax-h3/reference-to-video`

**MiniMax H3 Reference-to-Video** · Minimax H3 · [official docs](https://docs.kie.ai/market/minimax-h3/reference-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 1; max length: 7000 | Video generation prompt, length is between 1 and 7000 characters. |
| `reference_image_urls` | array | no | max items: 9 | Array of reference image URLs. Supports up to 9 images. Supports HTTP, HTTPS, and OSS addresses; supported image formats include JPG, JPEG, PNG, WEBP, HEIC, and HEIF; a single image size cannot exceed 30 MB; the side length of the image must be between 256 and 5760 pixels, and the aspect ratio must be between 0.4 and 2.5. |
| `reference_video_urls` | array | no | max items: 3 | Array of reference video URLs. Supports up to 3 videos. Supports MP4, MOV; video encoding supports H.264, H.265, audio encoding supports AAC, MP3; a single file size cannot exceed 50 MB; the duration of a single segment is 2 to 15 seconds, and the total duration of all reference videos cannot exceed 15 seconds; the video side length must be between 256 and 5760 pixels, the aspect ratio must be between 0.4 and 2.5, and the frame rate must be between 23.976 and 60. |
| `reference_audio_urls` | array | no | max items: 3 | Array of reference audio URLs. Supports up to 3 audios. Supports WAV, MP3; a single file size cannot exceed 15 MB; the duration of a single segment is 2 to 15 seconds, and the total duration of all reference audios cannot exceed 15 seconds. reference_audio cannot be used alone, it must be accompanied by reference_image or reference_video. |
| `aspect_ratio` | string | no | allowed: `adaptive`, `21:9`, `16:9`, `4:3`, `1:1`, `3:4`, `9:16`; default: `"adaptive"` | Video aspect ratio. The default is adaptive, or a specific ratio can be specified. |
| `duration` | integer | yes | allowed: `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`; default: `6`; min: 4; max: 15 | The duration of the generated video, supporting integer values from 4 to 15 seconds. |
| `resolution` | string | no | allowed: `768P`, `2K`; default: `"2K"` | Video resolution. |

Example `input`:

```json
{
  "prompt": "Generate a continuous cinematic video referencing the characters, actions, and scenes in the input material",
  "reference_image_urls": [
    "https://example.com/reference-image.jpg"
  ],
  "reference_video_urls": [
    "https://example.com/reference-video.mp4"
  ],
  "reference_audio_urls": [
    "https://example.com/reference-audio.mp3"
  ],
  "aspect_ratio": "adaptive",
  "duration": 6
}
```

## `minimax-h3/text-to-video`

**MiniMax H3 Text-to-Video** · Minimax H3 · [official docs](https://docs.kie.ai/market/minimax-h3/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 1; max length: 7000 | Video generation prompt, length is between 1 and 7000 characters. |
| `aspect_ratio` | string | yes | allowed: `21:9`, `16:9`, `4:3`, `1:1`, `3:4`, `9:16` | Video aspect ratio. Required for text-to-video, adaptive is not supported. |
| `duration` | integer | yes | allowed: `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `13`, `14`, `15`; default: `6`; min: 4; max: 15 | Generated video duration, supports integer values from 4 to 15 seconds. |
| `resolution` | string | no | allowed: `768P`, `2K`; default: `"2K"` | Video resolution. |

Example `input`:

```json
{
  "prompt": "A cat walking slowly on the beach at sunset, cinematic shot",
  "aspect_ratio": "16:9",
  "duration": 6
}
```

## `nano-banana-2`

**Google - Nano Banana 2** · Google · [official docs](https://docs.kie.ai/market/google/nanobanana2)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 20000 | A text description of the image you want to generate (Max length: 20000 characters) |
| `image_input` | array | no | max items: 14 | Input images to transform or use as reference (supports up to 14 images) (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 30.0MB) |
| `aspect_ratio` | string | no | allowed: `1:1`, `2:3`, `3:2`, `1:4`, `4:1`, `3:4`, `4:3`, `4:5`, `5:4`, `1:8`, `8:1`, `9:16`, `16:9`, `21:9`, `auto`; default: `"auto"` | Aspect ratio of the generated image |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"1K"` | Resolution of the generated image |
| `output_format` | string | no | allowed: `png`, `jpg`; default: `"jpg"` | Format of the output image |

Example `input`:

```json
{
  "prompt": "Comic poster: cool banana hero in shades leaps from sci-fi pad. Six panels: 1) 4K mountain landscape, 2) banana holds page of long multilingual text with auto translation, 3) Gemini 3 hologram for search/knowledge/reasoning, 4) camera UI sliders for angle focus color, 5) frame trio 1:1-9:16, 6) consistent banana poses. Footer shows Google icons. Tagline: Nano Banana Pro now on Kie AI.",
  "image_input": [],
  "aspect_ratio": "auto",
  "resolution": "1K",
  "output_format": "png"
}
```

## `nano-banana-2-lite`

**Google - Nano Banana 2 Lite** · Google · [official docs](https://docs.kie.ai/market/google/nano-banana-2-lite)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_urls` | array | no | default: `[]`; max items: 10 | Input image URL array. Optional parameter. Set to an empty array or omit it for pure text-to-image generation. Supports up to 10 images.  - Please provide uploaded file URLs, not raw file content - Accepted types: `image/jpeg`, `image/png`, `image/webp` - Max size: 30.0MB |
| `prompt` | string | yes | max length: 20000 | Text prompt used to generate the image. Required. Maximum length: 20000 characters. |
| `aspect_ratio` | string | yes | allowed: `1:1`, `1:4`, `1:8`, `2:3`, `3:2`, `3:4`, `4:1`, `4:3`, `4:5`, `5:4`, `8:1`, `9:16`, `16:9`, `21:9`, `auto`; default: `"auto"` | Generated image aspect ratio. Default value: `auto`. Use `auto` to let the system choose the aspect ratio automatically. |

Example `input`:

```json
{
  "image_urls": [
    "https://file.aiquickdraw.com/custom-page/akr/section-images/1756223420389w8xa2jfe.png"
  ],
  "prompt": "Generate a pig on the grass, cinematic light",
  "aspect_ratio": "auto"
}
```

## `nano-banana-pro`

**Google - Nano Banana Pro** · Google · [official docs](https://docs.kie.ai/market/google/pro-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 10000 | A text description of the image you want to generate (Max length: 10000 characters) |
| `image_input` | array | no | max items: 8 | Input images to transform or use as reference (supports up to 8 images) (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 30.0MB) |
| `aspect_ratio` | string | no | allowed: `1:1`, `2:3`, `3:2`, `3:4`, `4:3`, `4:5`, `5:4`, `9:16`, `16:9`, `21:9`, `auto`; default: `"1:1"` | Aspect ratio of the generated image |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"1K"` | Resolution of the generated image |
| `output_format` | string | no | allowed: `png`, `jpg`; default: `"png"` | Format of the output image |

Example `input`:

```json
{
  "prompt": "Comic poster: cool banana hero in shades leaps from sci-fi pad. Six panels: 1) 4K mountain landscape, 2) banana holds page of long multilingual text with auto translation, 3) Gemini 3 hologram for search/knowledge/reasoning, 4) camera UI sliders for angle focus color, 5) frame trio 1:1-9:16, 6) consistent banana poses. Footer shows Google icons. Tagline: Nano Banana Pro now on Kie AI.",
  "image_input": [],
  "aspect_ratio": "1:1",
  "resolution": "1K",
  "output_format": "png"
}
```

## `omnihuman-1-5`

**Omnihuman 1.5** · Market · [official docs](https://docs.kie.ai/market/omnihuman-1-5)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes | format: uri | Portrait image URL. Supports any aspect ratio with subjects including people, pets, anime, etc. Accepted file types: image/jpeg, image/png, image/webp. Max file size: 10MB. |
| `mask_url` | array | no | max items: 5 | Optional mask image URL(s). To have a specific subject in the image speak, use 'Subject Detection' to get the corresponding mask image and pass it as input. Accepted file types: image/jpeg, image/png, image/webp. Max file size: 10MB. Multiple file upload is supported, up to 5 files. |
| `audio_url` | string | yes | format: uri | Audio URL. Duration must be less than 60 seconds (recommended 15 seconds or less; exceeding this will cause quality degradation). Accepted file types: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/ogg, audio/mp4. Max file size: 10MB. |
| `prompt` | string | no | max length: 300 | Optional prompt text. Limited to Chinese, English, Japanese, Korean, Spanish, and Indonesian. Recommended 300 characters or less. Maximum length: 1000 characters. |
| `output_resolution` | string | no | allowed: `720`, `1080`; default: `"1080"` | Output video resolution.  - `720`: 720P - `1080`: 1080P (default) |
| `pe_fast_mode` | boolean | no | default: `false` | Fast mode. Sacrifices some quality to speed up generation. Default value: `false`. |
| `seed` | integer | no | default: `-1` | Random seed. Default is `-1` (random). When using the same positive integer and keeping all other parameters identical, the result will be highly consistent. |

Example `input`:

```json
{
  "image_url": "https://your-domain.com/image/portrait.png",
  "mask_url": [
    "https://your-domain.com/image/mask.png"
  ],
  "audio_url": "https://your-domain.com/audio/speech.mp3",
  "prompt": "A person speaking naturally with gentle expressions.",
  "output_resolution": "1080",
  "pe_fast_mode": false,
  "seed": -1
}
```

## `omnihuman-1-5/human-identification`

**Omnihuman 1.5 Human Identification** · Omnihuman 1 5 · [official docs](https://docs.kie.ai/market/omnihuman-1-5/human-identification)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes | format: uri | Portrait image URL. Requirements: JPG/PNG/JPEG format, less than 5MB, resolution under 4096x4096. Recommended: single person facing forward, with the face occupying a large proportion of the image. Accepted file types: image/jpeg, image/png, image/jpg. Max file size: 5MB. |

Example `input`:

```json
{
  "image_url": "https://your-domain.com/image/portrait.png"
}
```

## `omnihuman-1-5/subject-detection`

**OmniHuman 1.5 Subject Detection** · Omnihuman 1 5 · [official docs](https://docs.kie.ai/market/omnihuman-1-5/subject-detection)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes | format: uri | Portrait image URL. Supports detection of up to 5 subjects in the image. Accepted file types: image/jpeg, image/png, image/jpg. Max file size: 5MB. |

Example `input`:

```json
{
  "image_url": "https://your-domain.com/image/portrait.png"
}
```

## `pixverse-v6/extend`

**PixVerse V6 Video Extension** · Pixverse · [official docs](https://docs.kie.ai/market/pixverse/extend)

No structured input fields were published on the source page.

Example `input`:

```json
{
  "prompt": "Continue the same camera motion and extend the scene naturally",
  "taskId": "parent_task_id_from_previous_success_video",
  "duration": 5,
  "quality": "720p",
  "generate_audio_switch": false,
  "seed": 123456
}
```

## `pixverse-v6/image-to-video`

**PixVerse V6 Image-to-Video** · Pixverse · [official docs](https://docs.kie.ai/market/pixverse/image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | Generate prompt, cannot be empty, length is limited to 3-5000 characters. |
| `image_urls` | array | yes |  | Image URLs, supports up to 2 images, single image size not exceeding 20 MB. Supports HTTP, HTTPS, and OSS addresses. Supported image formats include JPG, JPEG, PNG, and WebP. |
| `quality` | string | yes | allowed: `360p`, `540p`, `720p`, `1080p`; default: `"720p"` | Output video resolution. Supports 360p, 540p, 720p, and 1080p. |
| `duration` | integer | yes | default: `5`; min: 1; max: 15 | Generated video duration in seconds, ranging from 1 to 15. Required when template_id is not passed; if template_id is passed, the video duration is fixed by the selected template, please do not pass this parameter at the same time. |
| `generate_audio_switch` | boolean | no | default: `false` | Whether to generate audio synchronized with the video content. |
| `generate_multi_clip_switch` | boolean | no | default: `false` | Whether to generate a multi-clip video. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed, value range is 0-2147483647. Using the same parameters and seed helps improve result reproducibility. |
| `template_id` | string | no | allowed: `412736208886848`, `411563216524736`, `411316927927040`, `411174903569216`, `410999246341952`, `408891141511104`, `410317363057408`, `410285445698304`, `410133101388544`, `409899296377728`, `408897485909952`, `408869061406656`, `408661207662528`, `409767750265728`, `409766559675264`, `409589071559552`, `407804339389760`, `407658863287616`, `407474361215744`, `407473438360320`, `407467702283008`, `385844572217469`, `406428904874432`, `406423682060736`, `406411724685760`, `406372274350528`, `406218913317312`, `406218479198656`, `406064763308480`, `406413607395776`, `405662117814720`, `406014934000064`, `405658369331648`, `405321470423488`, `405175211454656`, `404955806201792`, `404820147974080`, `398980393937856`, `403916646846400`, `403739217192896`, `403560060358144`, `403556618212098`, `398965022284579`, `402201060373270`, `403085292466285`, `402061828030966`, `402888569901524`, `402155676592228`, `402047865383360`, `402046136040202` | Used to select a PixVerse video effect template. Please pass in the corresponding template_id. Once a template_id is provided, the video duration is fixed by the selected template, so duration cannot be set at the same time. effect_type indicates the number of images required by the selected template. Please upload the specified number of images accordingly. To preview the template effect, open https://static.aiquickdraw.com/tools/example/<template_id>.mp4 in a browser and replace <template_id> with the actual template ID. For example:https://static.aiquickdraw.com/tools/example/412736208886848.mp4.  Available template list:  412736208886848 - Dive into the deep blue of love - Let every kiss turn into a dream. - effect_type: 2 411563216524736 - Skyline Track Flag - Unwavering. Next is the city-level entrance effect. - effect_type: 1 411316927927040 - Vibe Copines - Capture this vibe with tungtung. - effect_type: 1 411174903569216 - Crowd Focus - First person: The camera caught you - upload a selfie or group photo to become the focus of the audience - effect_type: 1 410999246341952 - Poke my little cutie - Pinch this little guy and watch it get cutely angry. - effect_type: 1 408891141511104 - Today is my birthday - Yes, I know you don't know me, but today is my birthday - effect_type: 1 410317363057408 - Mini Football Hero - Turn any photo into a mini chibi football hero and compete with giant players in a cinematic hyper-realistic stadium. - effect_type: 1 410285445698304 - Magma Rise - Rise from the ashes. Fight fire with fire. - effect_type: 1 410133101388544 - Cyber Armor: Reborn in the Rain 🌧️⚡ - Shatter the storm. Awaken the body of steel within. - effect_type: 1 409899296377728 - Product Landmark - Upload your product and turn it into a building! - effect_type: 1 408897485909952 - Mini Football Pitch - Upload your product and create a mini football pitch inside it! - effect_type: 1 408869061406656 - A kick through the world - Upload your product and let it shine on the night of the grand annual global competition! - effect_type: 1 408661207662528 - Small Town Footballer - Upload your product and enjoy the moment of victory! - effect_type: 1 409767750265728 - Crowned God in One Battle - You unexpectedly became the MVP of the game. - effect_type: 1 409766559675264 - Thriller Dance Steps - Upload a photo and watch it turn into an iconic viral dance trend! - effect_type: 1 409589071559552 - Dai Dai Dance - Vibe cheering dance template. A front-facing photo to easily join the dance floor clip. - effect_type: 1 407804339389760 - Fluffy Chef - Upload your product and cook in the fluffy kitchen - effect_type: 1 407658863287616 - Summer Postcard - Upload your product and join the summer special! - effect_type: 1 407474361215744 - Inhaled into Product Universe - Upload your product and enter the product universe! - effect_type: 1 407473438360320 - Fluffy Factory - Upload your product to make it in the fluffy factory - effect_type: 1 407467702283008 - Surfing Summer - Upload your product photo and start surfing! - effect_type: 1 385844572217469 - Love Launcher - A Valentine's Day moment in one shot. - effect_type: 1 406428904874432 - Kitty Shop - Upload your product image and let the cute cat transform into your street vendor! - effect_type: 1 406423682060736 - Courtyard Makeover Party - Upload your courtyard and give it a magical makeover! - effect_type: 1 406411724685760 - Ski Joy - Upload a photo of your pet, toy or product and let it embark on a fun skiing adventure. - effect_type: 1 406372274350528 - Nail Lab - Upload your nail photos and let our exclusive manicurists create exquisite nails for you. - effect_type: 1 406218913317312 - Rhythm Dash - Hit the beat, take off the crown. Welcome to the perfect score frenzy. - effect_type: 1 406218479198656 - Screen Killer King - The court needs a hero. So, you break through the screen. - effect_type: 1 406064763308480 - Dynamic Football Poster - From selfie to jersey photo (single and multiplayer) - effect_type: 1 406413607395776 - Football Live King! - The moment you score a wonderful goal, super sports car gifts pour down. - effect_type: 1 405662117814720 - Trophy Breakthrough - Upload your photo, transform into a champion football player, break through the screen and win the trophy! - effect_type: 1 406014934000064 - Post-match Sharp Comment - If the microphone is handed to you after the game, what would you say? - effect_type: 1 405658369331648 - Stadium Legend - Run, celebrate, create a legendary football moment. - effect_type: 1 405321470423488 - Step by Step - Step by step, shining with confidence - effect_type: 1 405175211454656 - Superstar Lobby - Welcome to your championship season. - effect_type: 1 404955806201792 - Post-match Sharp Comment 2 - If it's your turn for a post-match interview, what would you say about this game? - effect_type: 1 404820147974080 - Top of the World - Hold the trophy high and become the well-deserved hero of the football feast. - effect_type: 2 398980393937856 - The Last Hug - The last hug before the tsunami comes. - effect_type: 2 403916646846400 - My Future has Infinite Possibilities - Embrace your infinite potential and shine your future. - effect_type: 1 403739217192896 - Apex Dance - Flat rhythm, decadent chaotic dance steps - effect_type: 1 403560060358144 - Idol Ending Shot - One look up is a million direct shots. - effect_type: 1 403556618212098 - MotoGP Live - Live an unchoreographed moment. - effect_type: 1 398965022284579 - Knee Slide - Celebrate the victory with an iconic and energetic knee slide goal. - effect_type: 1 402201060373270 - Tunnel to Captain - Transform from an ordinary girl to a legendary football captain in a cinematic stadium journey. - effect_type: 1 403085292466285 - Golden Field Breaker - He is not on the list, but he still controls the scene. - effect_type: 1 402061828030966 - Trophy Celebration - Have you always wanted a huge trophy? Here, you can at least get it digitally. - effect_type: 1 402888569901524 - Jump into the Crowd 2 - Jump into the crowd - effect_type: 1 402155676592228 - World Champion Lift - You are the captain, lifting the trophy of victory. Mountains of people, golden ribbons, pure victory. - effect_type: 1 402047865383360 - Sideline Ball Boy - Experience the game from a unique perspective on the sidelines. Feel the excitement of the game as a young talent ready to participate. - effect_type: 1 402046136040202 - Epic Save - Become a legendary goalkeeper. Catch the ball in a stunning flying save. - effect_type: 1 |

Example `input`:

```json
{
  "prompt": "Animate the subject with gentle wind and cinematic lighting",
  "image_urls": [
    "https://example.com/input-image.png"
  ],
  "duration": 5,
  "quality": "720p",
  "generate_audio_switch": false,
  "generate_multi_clip_switch": false,
  "seed": 123456
}
```

## `pixverse-v6/reference-to-video`

**PixVerse V6 Fusion / Reference-to-Video** · Pixverse · [official docs](https://docs.kie.ai/market/pixverse/reference-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | Generate prompt, cannot be empty, length is limited to 3-5000 characters. |
| `image_references` | array | yes | min items: 1; max items: 7 | Fusion multi-reference image list.   `ref_name` must be unique within the same list. You can reference the reference object in the prompt via `@ref_name`. |
| `image_references[].image_url` | string | yes | format: uri | Reference image URL. |
| `image_references[].type` | string | no | allowed: `subject`, `background`; default: `"subject"` | Reference image type. `subject` means subject, `background` means background. |
| `image_references[].ref_name` | string | no | min length: 1; max length: 30 | Reference image name, maximum 30 characters. Can be referenced in the prompt using `@name`. |
| `aspect_ratio` | string | yes | allowed: `16:9`, `4:3`, `1:1`, `3:4`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"16:9"` | Output video aspect ratio. Supported in text-to-video and Fusion multi-reference image to video modes. |
| `quality` | string | yes | allowed: `360p`, `540p`, `720p`, `1080p`; default: `"720p"` | Output video resolution. Supports 360p, 540p, 720p, and 1080p. |
| `duration` | integer | yes | default: `5`; min: 1; max: 15 | Output video duration, in seconds. |
| `generate_audio_switch` | boolean | no | default: `false` | Whether to generate audio synchronized with the video content. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed, value range is 0-2147483647. Using the same parameters and seed helps improve the reproducibility of the results. |

## `pixverse-v6/text-to-video`

**PixVerse V6 Text-to-Video** · Pixverse · [official docs](https://docs.kie.ai/market/pixverse/text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | Generate prompt, cannot be empty, length is limited to 3-5000 characters. |
| `aspect_ratio` | string | yes | allowed: `16:9`, `4:3`, `1:1`, `3:4`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"16:9"` | Output video aspect ratio. This parameter is supported in text-to-video and Fusion multi-reference image-to-video modes. |
| `quality` | string | yes | allowed: `360p`, `540p`, `720p`, `1080p`; default: `"720p"` | Output video resolution. Supports 360p, 540p, 720p, and 1080p. |
| `duration` | integer | yes | default: `5`; min: 1; max: 15 | Output video duration in seconds. PixVerse V6 supports 1–15 seconds. |
| `generate_audio_switch` | boolean | no | default: `false` | Whether to generate audio synchronized with the video content. |
| `generate_multi_clip_switch` | boolean | no | default: `false` | Whether to generate a multi-clip video. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed, value range is 0–2147483647. Using the same parameters and seed helps improve result reproducibility. |

## `pixverse-v6/transition`

**PixVerse V6 First & Last Frame Transition** · Pixverse · [official docs](https://docs.kie.ai/market/pixverse/transition)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | Generate prompt, cannot be empty, length is limited to 3-5000 characters. |
| `first_frame_image_url` | string | yes | format: uri | First frame image URL. Supports HTTP, HTTPS, and OSS addresses; supported image formats include JPG, JPEG, PNG, and WebP; the size of a single image file must not exceed 20 MB |
| `last_frame_image_url` | string | yes | format: uri | Last frame image URL. Supports HTTP, HTTPS, and OSS addresses; supported image formats include JPG, JPEG, PNG, and WebP; the size of a single image file must not exceed 20 MB |
| `quality` | string | yes | allowed: `360p`, `540p`, `720p`, `1080p`; default: `"720p"` | Output video resolution. Supports 360p, 540p, 720p, and 1080p. |
| `duration` | integer | yes | default: `5`; min: 1; max: 15 | Output video duration in seconds. PixVerse V6 supports 1–15 seconds. |
| `generate_audio_switch` | boolean | no | default: `false` | Whether to generate audio synchronized with the video content. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed, value range is 0–2147483647. Using the same parameters and seed helps improve result reproducibility. |

## `qwen/image-edit`

**Qwen - Image Edit** · Qwen · [official docs](https://docs.kie.ai/market/qwen/image-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 2000 | The prompt to generate the image with (Max length: 2000 characters) |
| `image_url` | string | yes |  | The URL of the image to edit. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `acceleration` | string | no | allowed: `none`, `regular`, `high`; default: `"none"` | Acceleration level for image generation. Options: 'none', 'regular'. Higher acceleration increases speed. 'regular' balances speed and quality. Default value: "none" |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9`; default: `"landscape_4_3"` | The size of the generated image. Default value: landscape_4_3 |
| `num_inference_steps` | number | no | default: `25`; min: 2; max: 49 | The number of inference steps to perform. Default value: 30 (Min: 2, Max: 49, Step: 1) (step: 1) |
| `seed` | integer | no |  | The same seed and the same prompt given to the same version of the model will output the same image every time. |
| `guidance_scale` | number | no | default: `4`; min: 0; max: 20 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt when looking for a related image to show you. Default value: 4 (Min: 0, Max: 20, Step: 0.1) (step: 0.1) |
| `sync_mode` | boolean | no |  | If set to true, the function will wait for the image to be generated and uploaded before returning the response. This will increase the latency of the function but it allows you to get the image directly in the response without going through the CDN. (Boolean value (true/false)) |
| `num_images` | string | no | allowed: `1`, `2`, `3`, `4` | num_images |
| `enable_safety_checker` | boolean | no |  | If set to true, the safety checker will be enabled. Default value: true (Boolean value (true/false)) |
| `output_format` | string | no | allowed: `jpeg`, `png`; default: `"png"` | The format of the generated image. Default value: "png" |
| `negative_prompt` | string | no | max length: 500 | The negative prompt for the generation Default value: " " (Max length: 500 characters) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755603225969i6j87xnw.jpg",
  "acceleration": "none",
  "image_size": "landscape_4_3",
  "num_inference_steps": 25,
  "guidance_scale": 4,
  "sync_mode": false,
  "enable_safety_checker": true,
  "output_format": "png",
  "negative_prompt": "blurry, ugly"
}
```

## `qwen/image-to-image`

**Qwen - Image to Image** · Qwen · [official docs](https://docs.kie.ai/market/qwen/image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to generate the image with (Max length: 5000 characters) |
| `image_url` | string | yes |  | The reference image to guide the generation (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `strength` | number | no | default: `0.8`; min: 0; max: 1 | Denoising strength. 1.0 = fully remake; 0.0 = preserve original (Min: 0, Max: 1, Step: 0.01) (step: 0.01) |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The format of the generated image |
| `acceleration` | string | no | allowed: `none`, `regular`, `high`; default: `"none"` | Acceleration level for image generation. Options: 'none', 'regular', 'high'. Higher acceleration increases speed. 'regular' balances speed and quality. 'high' is recommended for images without text |
| `negative_prompt` | string | no | max length: 500 | The negative prompt for the generation (Max length: 500 characters) |
| `seed` | integer | no |  | The same seed and the same prompt given to the same version of the model will output the same image every time |
| `num_inference_steps` | number | no | default: `30`; min: 2; max: 250 | The number of inference steps to perform (Min: 2, Max: 250, Step: 1) (step: 1) |
| `guidance_scale` | number | no | default: `2.5`; min: 0; max: 20 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt when looking for a related image to show you (Min: 0, Max: 20, Step: 0.1) (step: 0.1) |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "",
  "image_url": "",
  "strength": 0.8,
  "output_format": "png",
  "acceleration": "none",
  "negative_prompt": "blurry, ugly",
  "num_inference_steps": 30,
  "guidance_scale": 2.5,
  "enable_safety_checker": true
}
```

## `qwen/text-to-image`

**Qwen - Text to Image** · Qwen · [official docs](https://docs.kie.ai/market/qwen/text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The prompt to generate the image with (Max length: 5000 characters) |
| `image_size` | string | no | allowed: `square`, `square_hd`, `portrait_4_3`, `portrait_16_9`, `landscape_4_3`, `landscape_16_9`; default: `"square_hd"` | The size of the generated image |
| `num_inference_steps` | number | no | default: `30`; min: 2; max: 250 | The number of inference steps to perform (Min: 2, Max: 250, Step: 1) (step: 1) |
| `seed` | integer | no |  | The same seed and the same prompt given to the same version of the model will output the same image every time |
| `guidance_scale` | number | no | default: `2.5`; min: 0; max: 20 | The CFG (Classifier Free Guidance) scale is a measure of how close you want the model to stick to your prompt when looking for a related image to show you (Min: 0, Max: 20, Step: 0.1) (step: 0.1) |
| `enable_safety_checker` | boolean | no |  | The safety checker is always enabled in Playground. It can only be disabled by setting false through the API. (Boolean value (true/false)) |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The format of the generated image |
| `negative_prompt` | string | no | max length: 500 | The negative prompt for the generation (Max length: 500 characters) |
| `acceleration` | string | no | allowed: `none`, `regular`, `high`; default: `"none"` | Acceleration level for image generation. Options: 'none', 'regular', 'high'. Higher acceleration increases speed. 'regular' balances speed and quality. 'high' is recommended for images without text |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "",
  "image_size": "square_hd",
  "num_inference_steps": 30,
  "guidance_scale": 2.5,
  "enable_safety_checker": true,
  "output_format": "png",
  "negative_prompt": " ",
  "acceleration": "none"
}
```

## `qwen2/image-edit`

**Qwen2 - Image Edit** · Qwen2 · [official docs](https://docs.kie.ai/market/qwen2/image-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 800 | The prompt to generate the image with (Max length: 800 characters) |
| `image_url` | string | yes |  | The URL of the image to edit. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `image_size` | string | no | allowed: `1:1`, `2:3`, `3:2`, `3:4`, `4:3`, `9:16`, `16:9`, `21:9`; default: `"16:9"` | The size of the generated image. Default value: 16:9 |
| `seed` | integer | no |  | The same seed and the same prompt given to the same version of the model will output the same image every time. |
| `output_format` | string | no | allowed: `jpeg`, `png`; default: `"png"` | The format of the generated image. Default value: "png" |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755603225969i6j87xnw.jpg",
  "image_size": "16:9",
  "output_format": "png",
  "seed": 0
}
```

## `qwen2/text-to-image`

**Qwen2 - Text To Image** · Qwen2 · [official docs](https://docs.kie.ai/market/qwen2/text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 800 | The prompt to generate the image with (Max length: 800 characters) |
| `image_size` | string | no | allowed: `1:1`, `3:4`, `4:3`, `9:16`, `16:9`; default: `"16:9"` | The size of the generated image. Default value: 16:9 |
| `seed` | integer | no |  | The same seed and the same prompt given to the same version of the model will output the same image every time. |
| `output_format` | string | no | allowed: `jpeg`, `png`; default: `"png"` | The format of the generated image. Default value: "png" |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "",
  "image_size": "16:9",
  "seed": 0,
  "output_format": "png"
}
```

## `qwen3/image-to-image`

**Qwen3 Image to Image** · Qwen3 · [official docs](https://docs.kie.ai/market/qwen3/image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_urls` | array | yes | min items: 1; max items: 3 | Input image URL array. Upload images and provide their URLs rather than the file content. Supports up to 3 images. Supported file types: `image/jpeg`, `image/png`, `image/webp`, `image/bmp`, and `image/gif`,`image/tiff`. Maximum file size: 10 MB per image. |
| `prompt` | string | yes | min length: 0; max length: 5000 | The positive prompt that describes the image content, style, and composition you want to generate or edit. Both Chinese and English are supported. |
| `resolution` | string | no | allowed: `1K`, `2K`; default: `"1K"` | The output image resolution. |
| `image_size` | string | no | allowed: `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`; default: `"16:9"` | The output image size. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The generated image format. |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default: true (recommended). When enabled, the model optimizes the positive prompt, which significantly improves results for simple descriptions. |
| `nsfw_checker` | boolean | no | default: `false` | Whether to enable content filtering. Optional. Defaults to false. When set to false, content filtering is disabled and all results are returned directly by the model. Note: Content filtering cannot guarantee that all inappropriate content will be detected. If the filtering results do not meet your requirements, you must implement additional content moderation measures. |
| `negative_prompt` | string | no | min length: 0; max length: 5000 | The negative prompt that describes content you do not want to appear in the image. |
| `seed` | integer | no | default: `1`; min: 0; max: 2147483647 | The random seed. Value range: [0, 2147483647]. A fixed seed helps maintain relatively stable results. |

Example `input`:

```json
{
  "image_urls": [
    "https://example.com/input.png"
  ],
  "prompt": "Turn the input image into a cinematic watercolor illustration.",
  "resolution": "1K",
  "image_size": "1:1",
  "output_format": "png",
  "prompt_extend": true,
  "nsfw_checker": false,
  "negative_prompt": "blurry, low quality, distorted",
  "seed": 1
}
```

## `qwen3/pro-image-to-image`

**Qwen3 Pro Image to Image** · Qwen3 Pro · [official docs](https://docs.kie.ai/market/qwen3-pro/image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_urls` | array | yes | min items: 1; max items: 3 | Input image URL array. Upload images and provide their URLs rather than the file content. Supports up to 3 images. Supported file types: `image/jpeg`, `image/png`, `image/webp`, `image/bmp`, and `image/gif`,`image/tiff`. Maximum file size: 10 MB per image. |
| `prompt` | string | yes | min length: 0; max length: 5000 | The positive prompt that describes the image content, style, and composition you want to generate or edit. Both Chinese and English are supported. |
| `resolution` | string | no | allowed: `1K`, `2K`; default: `"1K"` | The output image resolution. |
| `image_size` | string | no | allowed: `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`; default: `"16:9"` | The output image size. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The generated image format. |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default: true (recommended). When enabled, the model optimizes the positive prompt, which significantly improves results for simple descriptions. |
| `nsfw_checker` | boolean | no | default: `false` | Whether to enable content filtering. Optional. Defaults to false. When set to false, content filtering is disabled and all results are returned directly by the model. Note: Content filtering cannot guarantee that all inappropriate content will be detected. If the filtering results do not meet your requirements, you must implement additional content moderation measures. |
| `negative_prompt` | string | no | min length: 0; max length: 5000 | The negative prompt that describes content you do not want to appear in the image. |
| `seed` | integer | no | default: `1`; min: 0; max: 2147483647 | The random seed. Value range: [0, 2147483647]. A fixed seed helps maintain relatively stable results. |

Example `input`:

```json
{
  "image_urls": [
    "https://example.com/input.png"
  ],
  "prompt": "Turn the input image into a cinematic watercolor illustration.",
  "resolution": "1K",
  "image_size": "1:1",
  "output_format": "png",
  "prompt_extend": true,
  "nsfw_checker": false,
  "negative_prompt": "blurry, low quality, distorted",
  "seed": 1
}
```

## `qwen3/pro-text-to-image`

**Qwen3 Pro Text to Image** · Qwen3 Pro · [official docs](https://docs.kie.ai/market/qwen3-pro/text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 0; max length: 5000 | The positive prompt that describes the image content, style, and composition you want to generate or edit. Both Chinese and English are supported. |
| `resolution` | string | no | allowed: `1K`, `2K` | The output image resolution. |
| `image_size` | string | no | allowed: `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`; default: `"16:9"` | The output image size. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The generated image format. |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default: true (recommended). When enabled, the model optimizes the positive prompt, which significantly improves results for simple descriptions. |
| `nsfw_checker` | boolean | no | default: `false` | Whether to enable content filtering. Optional. Defaults to false. When set to false, content filtering is disabled and all results are returned directly by the model. Note: Content filtering cannot guarantee that all inappropriate content will be detected. If the filtering results do not meet your requirements, you must implement additional content moderation measures. |
| `negative_prompt` | string | no | min length: 0; max length: 5000 | The negative prompt that describes content you do not want to appear in the image. |
| `seed` | integer | no | default: `1`; min: 0; max: 2147483647 | The random seed. Value range: [0, 2147483647]. A fixed seed helps maintain relatively stable results. |

Example `input`:

```json
{
  "prompt": "A futuristic city at sunset, cinematic lighting, ultra detailed.",
  "resolution": "1K",
  "image_size": "1:1",
  "output_format": "png",
  "prompt_extend": true,
  "nsfw_checker": false,
  "negative_prompt": "blurry, low quality, distorted",
  "seed": 1
}
```

## `qwen3/text-to-image`

**Qwen3 Text to Image** · Qwen3 · [official docs](https://docs.kie.ai/market/qwen3/text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 0; max length: 5000 | The positive prompt that describes the image content, style, and composition you want to generate or edit. Both Chinese and English are supported. |
| `resolution` | string | no | allowed: `1K`, `2K` | The output image resolution. |
| `image_size` | string | no | allowed: `1:1`, `3:2`, `2:3`, `4:3`, `3:4`, `16:9`, `9:16`, `21:9`; default: `"16:9"` | The output image size. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | The generated image format. |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default: true (recommended). When enabled, the model optimizes the positive prompt, which significantly improves results for simple descriptions. |
| `nsfw_checker` | boolean | no | default: `false` | Whether to enable content filtering. Optional. Defaults to false. When set to false, content filtering is disabled and all results are returned directly by the model. Note: Content filtering cannot guarantee that all inappropriate content will be detected. If the filtering results do not meet your requirements, you must implement additional content moderation measures. |
| `negative_prompt` | string | no | min length: 0; max length: 5000 | The negative prompt that describes content you do not want to appear in the image. |
| `seed` | integer | no | default: `1`; min: 0; max: 2147483647 | The random seed. Value range: [0, 2147483647]. A fixed seed helps maintain relatively stable results. |

Example `input`:

```json
{
  "prompt": "A futuristic city at sunset, cinematic lighting, ultra detailed.",
  "resolution": "1K",
  "image_size": "1:1",
  "output_format": "png",
  "prompt_extend": true,
  "nsfw_checker": false,
  "negative_prompt": "blurry, low quality, distorted",
  "seed": 1
}
```

## `recraft/crisp-upscale`

**Recraft - Crisp Upscale** · Recraft · [official docs](https://docs.kie.ai/market/recraft/crisp-upscale)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image` | string | yes |  | Image to upscale (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |

Example `input`:

```json
{
  "image": "https://file.aiquickdraw.com/custom-page/akr/section-images/1757169577325ijj8vwvt.jpg"
}
```

## `recraft/remove-background`

**Recraft - Remove Background** · Recraft · [official docs](https://docs.kie.ai/market/recraft/remove-background)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image` | string | yes |  | Image to remove background from. Supported formats: PNG, JPG, WEBP. Max 5MB, max 16MP, max dimension 4096px, min dimension 256px. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 5.0MB) |

Example `input`:

```json
{
  "image": "https://file.aiquickdraw.com/custom-page/akr/section-images/1757057285447k9qcbki1.webp"
}
```

## `seedream/4.5-edit`

**Seedream4.5 - Edit** · Seedream · [official docs](https://docs.kie.ai/market/seedream/4-5-edit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 3000 | A text description of the image you want to generate (Max length: 3000 characters) |
| `image_urls` | array | yes | max items: 14 | Upload an image file to use as input for the API (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`; default: `"basic"` | Basic outputs 2K images, while High outputs 4K images. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Keep the model's pose and the flowing shape of the liquid dress unchanged. Change the clothing material from silver metal to completely transparent clear water (or glass). Through the liquid water, the model's skin details are visible. Lighting changes from reflection to refraction.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1764851484363_ScV1s2aq.webp"
  ],
  "aspect_ratio": "1:1",
  "quality": "basic",
  "nsfw_checker": true
}
```

## `seedream/4.5-text-to-image`

**Seedream4.5 - Text to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/4-5-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 3000 | A text description of the image you want to generate (Max length: 3000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`; default: `"basic"` | Basic outputs 2K images, while High outputs 4K images. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A full-process cafe design tool for entrepreneurs and designers. It covers core needs including store layout, functional zoning, decoration style, equipment selection, and customer group adaptation, supporting integrated planning of \"commercial attributes + aesthetic design.\" Suitable as a promotional image for a cafe design SaaS product, with a 16:9 aspect ratio.",
  "aspect_ratio": "1:1",
  "quality": "basic",
  "nsfw_checker": false
}
```

## `seedream/5-lite-image-to-image`

**Seedream5.0 Lite - Image to Image** · Market · [official docs](https://docs.kie.ai/market/seedream-5-lite-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 3000 | A text description of the image you want to generate (Max length: 3-3000 characters) |
| `image_urls` | array | yes | max items: 14 | Upload an image file to use as input for the API (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`, `ultra`; default: `"basic"` | Basic outputs 2K images, while High outputs 3K images, and Ultra outputs 4K images. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Format of the output image |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Keep the model's pose and the flowing shape of the liquid dress unchanged. Change the clothing material from silver metal to completely transparent clear water (or glass). Through the liquid water, the model's skin details are visible. Lighting changes from reflection to refraction.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1764851484363_ScV1s2aq.webp"
  ],
  "aspect_ratio": "1:1",
  "quality": "basic",
  "output_format": "png",
  "nsfw_checker": true
}
```

## `seedream/5-lite-text-to-image`

**Seedream5.0 Lite - Text to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/5-lite-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 3000 | A text description of the image you want to generate (Max length: 3-3000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`, `ultra`; default: `"basic"` | Basic outputs 2K images, while High outputs 3K images, and Ultra outputs 4K images. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Format of the output image |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A full-process cafe design tool for entrepreneurs and designers. It covers core needs including store layout, functional zoning, decoration style, equipment selection, and customer group adaptation, supporting integrated planning of \"commercial attributes + aesthetic design.\" Suitable as a promotional image for a cafe design SaaS product, with a 16:9 aspect ratio.",
  "aspect_ratio": "1:1",
  "quality": "basic",
  "output_format": "png",
  "nsfw_checker": false
}
```

## `seedream/5-pro-image-to-image`

**Seedream5.0 Pro - Image to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/5-pro-image-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | A text description of the image you want to generate (Max length: 3-5000 characters) |
| `image_urls` | array | yes | max items: 10 | Upload an image file to use as input for the API (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`; default: `"basic"` | Basic outputs 1K images, while High outputs 2K images. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Format of the output image |

Example `input`:

```json
{
  "prompt": "Keep the model's pose and the flowing shape of the liquid dress unchanged. Change the clothing material from silver metal to completely transparent clear water (or glass). Through the liquid water, the model's skin details are visible. Lighting changes from reflection to refraction.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1764851484363_ScV1s2aq.webp"
  ],
  "aspect_ratio": "1:1",
  "quality": "basic",
  "output_format": "png",
  "nsfw_checker": true
}
```

## `seedream/5-pro-layer-decomposition`

**Seedream 5.0 Pro -  Layer Decomposition** · Seedream · [official docs](https://docs.kie.ai/market/seedream/5-pro-layer-decomposition)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | min length: 0; max length: 5000 | Optional prompt for layer separation.  - When a prompt is provided, the model separates the specified elements based on the prompt - When no prompt is provided, the model automatically identifies the main elements in the image - Supports using `<bbox>x1 y1 x2 y2</bbox>` to precisely specify element positions, recommending the use of normalized coordinates in the 0-1000 range |
| `image_url` | string | yes | format: uri | The source image URL for layer separation.  - Must upload exactly 1 image - Supports PNG, JPEG, WebP, BMP, TIFF, GIF - Does not support HEIC, HEIF - Image size must not exceed 30 MB - Total pixel range is 262,144-36,000,000 - Aspect ratio range is 1:16-16:1 |
| `size` | string | no | allowed: `auto`, `1K`, `1.5K`, `2K`; default: `"auto"` | Resolution level of the output image. The base image maintains the aspect ratio of the input image, and each separated layer maintains the aspect ratio of the corresponding element in the original image.  - `auto`: Automatically selected based on the input image size - `1K`: 1K resolution level - `1.5K`: 1.5K resolution level - `2K`: 2K resolution level |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"jpeg"` | Output format of the base image. This parameter only controls the base image format; all separated layers are fixed to output as PNG. |

## `seedream/5-pro-text-to-image`

**Seedream5.0 Pro - Text to Image** · Seedream · [official docs](https://docs.kie.ai/market/seedream/5-pro-text-to-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 3; max length: 5000 | A text description of the image you want to generate (Max length: 3-5000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`, `2:3`, `3:2`, `21:9`; default: `"1:1"` | Width-height ratio of the image, determining its visual form. |
| `quality` | string | yes | allowed: `basic`, `high`; default: `"basic"` | Basic outputs 1K images, while High outputs 2K images. |
| `output_format` | string | no | allowed: `png`, `jpeg`; default: `"png"` | Format of the output image |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A full-process cafe design tool for entrepreneurs and designers. It covers core needs including store layout, functional zoning, decoration style, equipment selection, and customer group adaptation, supporting integrated planning of \"commercial attributes + aesthetic design.\" Suitable as a promotional image for a cafe design SaaS product, with a 16:9 aspect ratio.",
  "aspect_ratio": "1:1",
  "quality": "basic",
  "output_format": "png",
  "nsfw_checker": false
}
```

## `topaz/image-upscale`

**Topaz - Image Upscale** · Topaz · [official docs](https://docs.kie.ai/market/topaz/image-upscale)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes |  | Url of the image to be upscaled (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `upscale_factor` | string | yes | allowed: `1`, `2`, `4`; default: `"2"` | Factor to upscale the video by (e.g. 2.0 doubles width and height) |

Example `input`:

```json
{
  "image_url": "https://static.aiquickdraw.com/tools/example/1762752805607_mErUj1KR.png",
  "upscale_factor": "2"
}
```

## `topaz/video-upscale`

**Topaz - Video Upscale** · Topaz · [official docs](https://docs.kie.ai/market/topaz/video-upscale)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `video_url` | string | yes |  | URL of the video to upscale (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime, video/x-matroska; Max size: 50.0MB) |
| `upscale_factor` | string | no | allowed: `1`, `2`, `4`; default: `"2"` | Factor to upscale the video by (e.g. 2.0 doubles width and height) |

Example `input`:

```json
{
  "video_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1758166466095hvbwkrpw.mp4",
  "upscale_factor": "2"
}
```

## `volcengine/video-to-video-lip-sync`

**Volcengine video to video lip sync** · Volcengine · [official docs](https://docs.kie.ai/market/volcengine/video-to-video-lip-sync)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `mode` | string | yes | allowed: `lite`, `basic` | Service mode for lip-sync generation. - `lite`: For single-person frontal videos. Faster processing. - `basic`: For single-person complex scenes. Supports scene segmentation and speaker identification. |
| `video_url` | string | yes | format: uri | Video URL. Supported resolution: 360p–1080p. Videos above 1080p will be compressed to 1080p, while videos below 360p are not supported. Supported formats: MOV, MP4, HDR. Recommended codec: H.264. Other formats/codecs may be transcoded. Max file size: 500 MB. Bitrate: 1–30 Mbps. Frame rate: 24–60 fps. |
| `audio_url` | string | yes | format: uri | Target pure vocal audio URL; used to drive video lip movements. Accepted file types: audio/mpeg, audio/wav, audio/x-wav, audio/aac, audio/mp4, audio/ogg. Max file size: 10MB. |
| `separate_vocal` | boolean | no | default: `false` | Enable vocal separation to suppress background noise. Default value: `false`. |
| `open_scenedet` | boolean | no | default: `false` | Whether to enable scene segmentation and speaker identification. Supported only in Basic mode. Default value: `false`. |
| `align_audio` | boolean | no | default: `true` | Supported in Lite mode. Whether to loop the video when the audio is longer than the video. Default value: `true`. |
| `align_audio_reverse` | boolean | no | default: `false` | Supported in Lite mode. Whether to loop the video in reverse (backward). Requires `align_audio` to be set to `true`. Default value: `false`. |
| `templ_start_seconds` | number | no | default: `0` | Supported in Lite mode. Start time of the template video, in seconds. Default value: `0`. |

Example `input`:

```json
{
  "mode": "lite",
  "video_url": "https://your-domain.com/video/example.mp4",
  "audio_url": "https://your-domain.com/audio/speech.mp3",
  "separate_vocal": false,
  "open_scenedet": false,
  "align_audio": true,
  "align_audio_reverse": false,
  "templ_start_seconds": 0
}
```

## `wan/2-2-a14b-image-to-video-turbo`

**Wan - Image to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-2-a14b-image-to-video-turbo)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `image_url` | string | yes |  | URL of the input image. If the input image does not match the chosen aspect ratio, it is resized and center cropped. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `prompt` | string | yes | max length: 5000 | The text prompt to guide video generation. (Max length: 5000 characters) |
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"720p"` | Resolution of the generated video (480p or 720p). Default value: "720p" |
| `enable_prompt_expansion` | boolean | no |  | Whether to enable prompt expansion. This will use a large language model to expand the prompt with additional details while maintaining the original meaning. (Boolean value (true/false)) |
| `seed` | number | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility. If None, a random seed is chosen. (Min: 0, Max: 2147483647, Step: 1) (step: 1) |
| `acceleration` | string | no | allowed: `none`, `regular`; default: `"none"` | Acceleration level to use. The more acceleration, the faster the generation, but with lower quality. The recommended value is 'none'. Default value: "none" |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1755166042585gtf2mlrk.png",
  "prompt": "Overcast lighting, medium lens, soft lighting, low contrast lighting, edge lighting, low angle shot, desaturated colors, medium close-up shot, clean single shot, cool colors, center composition.The camera captures a low-angle close-up of a Western man outdoors, sharply dressed in a black coat over a gray sweater, white shirt, and black tie. His gaze is fixed on the lens as he advances. In the background, a brown building looms, its windows glowing with warm, yellow light above a dark doorway. As the camera pushes in, a blurred black object on the right side of the frame drifts back and forth, partially obscuring the view against a dark, nighttime background.",
  "resolution": "720p",
  "enable_prompt_expansion": false,
  "seed": 0,
  "acceleration": "none",
  "nsfw_checker": false
}
```

## `wan/2-2-a14b-speech-to-video-turbo`

**Wan - 2.2 A14B Speech to Video Turbo** · Wan · [official docs](https://docs.kie.ai/market/wan/2-2-a14b-speech-to-video-turbo)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt used for video generation (Max length: 5000 characters) |
| `image_url` | string | yes |  | URL of the input image. If the input image does not match the chosen aspect ratio, it is resized and center cropped (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `audio_url` | string | yes |  | The URL of the audio file (File URL after upload, not file content; Accepted types: audio/mp3, audio/wav, audio/ogg, audio/m4a, audio/flac, audio/aac, audio/x-ms-wma, audio/mpeg; Max size: 10.0MB) |
| `num_frames` | number | no | default: `80`; min: 40; max: 120 | Number of frames to generate. Must be between 40 to 120, (must be multiple of 4) (Min: 40, Max: 120, Step: 4) (step: 4) |
| `frames_per_second` | number | no | default: `16`; min: 4; max: 60 | Frames per second of the generated video. Must be between 4 to 60. When using interpolation and adjust_fps_for_interpolation is set to true (default true,) the final FPS will be multiplied by the number of interpolated frames plus one. For example, if the generated frames per second is 16 and the number of interpolated frames is 1, the final frames per second will be 32. If adjust_fps_for_interpolation is set to false, this value will be used as-is (Min: 4, Max: 60, Step: 1) (step: 1) |
| `resolution` | string | no | allowed: `480p`, `580p`, `720p`; default: `"480p"` | Resolution of the generated video (480p, 580p, or 720p) |
| `negative_prompt` | string | no | max length: 500 | Negative prompt for video generation (Max length: 500 characters) |
| `seed` | integer | no |  | Random seed for reproducibility. If None, a random seed is chosen |
| `num_inference_steps` | number | no | default: `27`; min: 2; max: 40 | Number of inference steps for sampling. Higher values give better quality but take longer (Min: 2, Max: 40, Step: 1) (step: 1) |
| `guidance_scale` | number | no | default: `3.5`; min: 1; max: 10 | Classifier-free guidance scale. Higher values give better adherence to the prompt but may decrease quality (Min: 1, Max: 10, Step: 0.1) (step: 0.1) |
| `shift` | number | no | default: `5`; min: 1; max: 10 | Shift value for the video. Must be between 1.0 and 10.0 (Min: 1, Max: 10, Step: 0.1) (step: 0.1) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "The lady is talking",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1756797663082u4pjmcrq.png",
  "audio_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17567977044127d1emlmc.mp3",
  "num_frames": 80,
  "frames_per_second": 16,
  "resolution": "480p",
  "negative_prompt": "",
  "num_inference_steps": 27,
  "guidance_scale": 3.5,
  "shift": 5
}
```

## `wan/2-2-a14b-text-to-video-turbo`

**Wan - Text to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-2-a14b-text-to-video-turbo)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | The text prompt to guide video generation. (Max length: 5000 characters) |
| `resolution` | string | no | allowed: `480p`, `720p`; default: `"720p"` | Resolution of the generated video (480p or 720p). Default value: "720p" |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`; default: `"16:9"` | Aspect ratio of the generated video (16:9 or 9:16). Default value: "16:9" |
| `enable_prompt_expansion` | boolean | no |  | Whether to enable prompt expansion. This will use a large language model to expand the prompt with additional details while maintaining the original meaning. (Boolean value (true/false)) |
| `seed` | number | no | default: `0`; min: 0; max: 2147483647 | Random seed for reproducibility. If None, a random seed is chosen. (Min: 0, Max: 2147483647, Step: 1) (step: 1) |
| `acceleration` | string | no | allowed: `none`, `regular`; default: `"none"` | Acceleration level to use. The more acceleration, the faster the generation, but with lower quality. The recommended value is 'none'. Default value: "none" |

Example `input`:

```json
{
  "prompt": "Drone shot, fast traversal, starting inside a cracked, frosty circular pipe. The camera bursts upward through the pipe to reveal a vast polar landscape bathed in golden sunrise light. Workers in orange suits operate steaming machinery. The camera tilts up, revealing the scene from the perspective of a rising hot air balloon. It continues ascending into a glowing sky, the balloon trailing steam and displaying the letters \"KIE AI\" as it rises into breathtaking polar majesty.",
  "resolution": "720p",
  "aspect_ratio": "16:9",
  "enable_prompt_expansion": false,
  "seed": 0,
  "acceleration": "none",
  "nsfw_checker": false
}
```

## `wan/2-2-animate-move`

**Wan - Animate Move** · Wan · [official docs](https://docs.kie.ai/market/wan/2-2-animate-move)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `video_url` | string | yes |  | URL of the input video. (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime, video/x-matroska; Max size: 10.0MB) |
| `image_url` | string | yes |  | URL of the input image. If the input image does not match the chosen aspect ratio, it is resized and center cropped. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `resolution` | string | no | allowed: `480p`, `580p`, `720p`; default: `"480p"` | Resolution of the generated video (480p, 580p, or 720p). |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "video_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17586254974931y2hottk.mp4",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1758625466310wpehpbnf.png",
  "resolution": "480p",
  "nsfw_checker": false
}
```

## `wan/2-2-animate-replace`

**Wan - Animate Replace** · Wan · [official docs](https://docs.kie.ai/market/wan/2-2-animate-replace)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `video_url` | string | yes |  | URL of the input video. (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime, video/x-matroska; Max size: 10.0MB) |
| `image_url` | string | yes |  | URL of the input image. If the input image does not match the chosen aspect ratio, it is resized and center cropped. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `resolution` | string | no | allowed: `480p`, `580p`, `720p`; default: `"480p"` | Resolution of the generated video (480p, 580p, or 720p). |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "video_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17586199429271xscyd5d.mp4",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/17586199255323tks43kq.png",
  "resolution": "480p",
  "nsfw_checker": false
}
```

## `wan/2-5-image-to-video`

**Wan 2.5 - Image to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-5-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 800 | The text prompt describing the desired video motion. Maximum length: 800 characters. |
| `image_url` | string | yes | format: uri | URL of the image to use as the first frame. Must be publicly accessible.  - Please provide the URL of the uploaded file, not raw file content - Accepted types: `image/jpeg`, `image/png`, `image/webp` - Max size: 10.0MB |
| `duration` | string | yes | allowed: `5`, `10` | The duration of the generated video in seconds.  - `5`: 5 seconds - `10`: 10 seconds |
| `resolution` | string | no | allowed: `720p`, `1080p` | Video resolution. Valid values: `720p`, `1080p`. |
| `negative_prompt` | string | no | max length: 500 | Negative prompt used to describe content to avoid. Maximum length: 500 characters. |
| `enable_prompt_expansion` | boolean | no |  | Whether to enable prompt rewriting using LLM.  - Boolean value: `true` / `false` |
| `seed` | integer | no |  | Random seed for reproducibility. If omitted, a random seed is chosen. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "The same woman from the reference image looks directly into the camera, takes a breath, then smiles brightly and speaks with enthusiasm: \"Have you heard? Alibaba Wan 2.5 API is now available on Kie.ai!\" Ambient audio: quiet indoor atmosphere, soft natural room tone. Camera: medium close-up, steady framing, natural daylight mood, accurate lip-sync with dialogue.",
  "image_url": "https://file.aiquickdraw.com/custom-page/akr/section-images/1758796480945qb63zxq8.webp",
  "duration": "5",
  "resolution": "1080p",
  "negative_prompt": "blurry, flicker, camera shake, distorted face, low quality",
  "enable_prompt_expansion": true,
  "seed": 123456,
  "nsfw_checker": false
}
```

## `wan/2-5-text-to-video`

**Wan 2.5 - Text to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-5-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 800 | The text prompt for video generation. Supports Chinese and English. Maximum length: 800 characters. |
| `duration` | string | yes | allowed: `5`, `10` | The duration of the generated video in seconds.  - `5`: 5 seconds - `10`: 10 seconds |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1` | The aspect ratio of the generated video.  - `16:9`: Landscape - `9:16`: Portrait - `1:1`: Square |
| `resolution` | string | no | allowed: `720p`, `1080p` | Video resolution tier.  - `720p`: 720p - `1080p`: 1080p |
| `negative_prompt` | string | no | max length: 500 | Negative prompt used to describe content to avoid. Maximum length: 500 characters. |
| `enable_prompt_expansion` | boolean | no |  | Whether to enable prompt rewriting using LLM. Improves results for short prompts but increases processing time.  - Boolean value: `true` / `false` |
| `seed` | integer | no |  | Random seed for reproducibility. If omitted, a random seed is chosen. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A dimly lit jazz bar at night, wooden tables glowing under warm pendant lights. Patrons sip drinks and chat quietly while a three-piece band performs on stage. The saxophone player stands under a spotlight, gleaming instrument reflecting the light. No dialogue. Ambient audio: smooth live jazz music with saxophone and piano, clinking glasses, low murmur of audience conversations, occasional burst of laughter from a nearby table. Camera: slow pan across the crowd, then gentle zoom toward the saxophone player's solo, focusing on expressive hand movements.",
  "duration": "5",
  "aspect_ratio": "16:9",
  "resolution": "1080p",
  "negative_prompt": "blurry, flicker, low quality, distorted people, camera shake",
  "enable_prompt_expansion": true,
  "seed": 123456,
  "nsfw_checker": false
}
```

## `wan/2-6-flash-image-to-video`

**Wan - 2.6-flash-image-to-video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-6-flash-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | Text prompts for video generation. Supports both Chinese and English, with a minimum of 2 characters and a maximum of 5,000 characters. (Max length: 1500 characters) |
| `image_urls` | array | yes | max items: 1 | A list of image URLs. All images must be at least 256x256px. (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`, `15`; default: `"5"` | The duration of the generated video in seconds |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution tier |
| `audio` | boolean | yes |  | Whether to generate video with audio. Audio directly affects the cost, as the pricing differs between videos with sound and silent videos. (Boolean value (true/false)) |
| `multi_shots` | boolean | no |  | The multi shots parameter controls the shot composition style during AI video generation, determining whether the generated video is a single continuous shot or multiple shots with transitions. (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Anthopmopric fox singing a Christmas song at the rubbish dump in the rain.",
  "image_urls": [],
  "duration": "5",
  "resolution": "1080p",
  "audio": false,
  "multi_shots": false,
  "nsfw_checker": false
}
```

## `wan/2-6-flash-video-to-video`

**Wan - 2-6-flash-video-to-video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-6-flash-video-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1500 | Text prompts for video generation. Supports both Chinese and English, with a minimum of 2 characters and a maximum of 5,000 characters. (Max length: 1500 characters) |
| `video_urls` | array | yes | max items: 3 | The URL of the image used to generate video (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime, video/x-matroska; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution tier |
| `audio` | boolean | no |  | Whether to generate video with audio. Audio directly affects the cost, as the pricing differs between videos with sound and silent videos. (Boolean value (true/false)) |
| `multi_shots` | boolean | no |  | The multi shots parameter controls the shot composition style during AI video generation, determining whether the generated video is a single continuous shot or multiple shots with transitions. (Boolean value (true/false)) |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "The video drinks milk tea while doing some improvised dance moves to the music.",
  "video_urls": [],
  "duration": "5",
  "resolution": "1080p",
  "multi_shots": false,
  "nsfw_checker": false
}
```

## `wan/2-6-image-to-video`

**Wan 2.6 - Image to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-6-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompts for video generation. Supports both Chinese and English, with a minimum of 2 characters and a maximum of 5,000 characters. (Max length: 5000 characters) |
| `image_urls` | array | yes | max items: 1 | Upload an image file to use as input for the API (File URL after upload, not file content; Accepted types: image/jpeg, image/png, image/webp; Max size: 10.0MB),All images must be at least 256x256px. |
| `duration` | string | no | allowed: `5`, `10`, `15`; default: `"5"` | The duration of the generated video in seconds |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution tier |
| `multi_shots` | boolean | no | default: `false` | The multi shots parameter controls the shot composition style during AI video generation, determining whether the generated video is a single continuous shot or multiple shots with transitions. |

Example `input`:

```json
{
  "prompt": "Anthopmopric fox singing a Christmas song at the rubbish dump in the rain.",
  "image_urls": [
    "https://static.aiquickdraw.com/tools/example/1765957673717_awiBAidD.webp"
  ],
  "duration": "5",
  "resolution": "1080p",
  "multi_shots": false,
  "nsfw_checker": false
}
```

## `wan/2-6-text-to-video`

**Wan 2.6 - Text to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-6-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompts for video generation. Supports both Chinese and English, with a minimum of 1 characters and a maximum of 5,000 characters. (Max length: 5000 characters) |
| `duration` | string | no | allowed: `5`, `10`, `15`; default: `"5"` | The duration of the generated video in seconds |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution tier |
| `multi_shots` | boolean | no | default: `false` | The multi shots parameter controls the shot composition style during AI video generation, determining whether the generated video is a single continuous shot or multiple shots with transitions. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "In a hyperrealistic ASMR video, a hand uses a knitted knife to slowly slice a burger made entirely of knitted wool. The satisfyingly crisp cut reveals a detailed cross-section of knitted meat, lettuce, and tomato slices. Captured in a close-up with a shallow depth of field, the scene is set against a stark, matte black surface. Cinematic lighting makes the surreal yarn textures shine with clear reflections. The focus is on the deliberate, satisfying motion and the unique, tactile materials.",
  "duration": "5",
  "resolution": "1080p",
  "multi_shots": false,
  "nsfw_checker": false
}
```

## `wan/2-6-video-to-video`

**Wan 2.6 - Video to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-6-video-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompts for video generation. Supports both Chinese and English, with a minimum of 2 characters and a maximum of 5,000 characters. (Max length: 5000 characters) |
| `video_urls` | array | yes | max items: 3 | The URL of the image used to generate video (File URL after upload, not file content; Accepted types: video/mp4, video/quicktime, video/x-matroska; Max size: 10.0MB) |
| `duration` | string | no | allowed: `5`, `10`; default: `"5"` | The duration of the generated video in seconds |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution tier |
| `multi_shots` | boolean | no | default: `false` | The multi shots parameter controls the shot composition style during AI video generation, determining whether the generated video is a single continuous shot or multiple shots with transitions. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "The video drinks milk tea while doing some improvised dance moves to the music.",
  "video_urls": [
    "https://static.aiquickdraw.com/tools/example/1765957777782_cNJpvhRx.mp4"
  ],
  "duration": "5",
  "resolution": "1080p",
  "multi_shots": false,
  "nsfw_checker": false
}
```

## `wan/2-7-image`

**Wan 2.7 Image** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Prompt for image generation or editing. This field supports both Chinese and English, with a maximum length of 5000 characters as per Alibaba Cloud documentation. |
| `input_urls` | array | no | max items: 9 | (Optional) Array of input image URLs. The current project uses `input_urls` as a wrapper field. |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `4:3`, `21:9`, `3:4`, `9:16`, `8:1`, `1:8` | (Optional) Output aspect ratio when no image input is provided. |
| `enable_sequential` | boolean | no | default: `false` | Whether to enable sequential/group image mode. Default is false. |
| `n` | integer | no |  | Number of images to generate. Range is 1-4 when `enable_sequential=false` (default: 4); range is 1-12 when `enable_sequential=true` (default: 12). |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"2K"` | Output resolution. The current project uses `resolution` as a wrapper field corresponding to the underlying resolution parameter. |
| `thinking_mode` | boolean | no | default: `false` | Whether to enable thinking mode. Only available when `enable_sequential=false` and `input_urls` is empty; the frontend will automatically disable it in other cases. |
| `color_palette` | array | no | min items: 3; max items: 10 | (Optional) Custom color theme. Only available when `enable_sequential=false`. Requires 3-10 colors, 8 recommended. |
| `color_palette[].hex` | string | no |  | HEX color value. |
| `color_palette[].ratio` | string | no |  | Color proportion, format must be xx.xx%. |
| `bbox_list` | array | no |  | (Optional) Interactive editing bounding box areas. The outer list length should match `input_urls`; maximum 2 boxes per image; single box format is `[x1, y1, x2, y2]`. |
| `watermark` | boolean | no | default: `false` | Whether to add watermark. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed, range 0-2147483647. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

## `wan/2-7-image-pro`

**Wan 2.7 Image Pro** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-image-pro)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Prompt for image generation or editing. This field supports both Chinese and English, with a maximum length of 5000 characters as per Alibaba Cloud documentation. |
| `input_urls` | array | no | max items: 9 | (Optional) Array of input image URLs. The current project uses `input_urls` as a wrapper field. |
| `aspect_ratio` | string | no | allowed: `1:1`, `16:9`, `4:3`, `21:9`, `3:4`, `9:16`, `8:1`, `1:8` | (Optional) Output aspect ratio when no image input is provided. |
| `enable_sequential` | boolean | no | default: `false` | Whether to enable sequential/group image mode. Default is false. |
| `n` | integer | no |  | Number of images to generate. Range is 1-4 when `enable_sequential=false` (default: 4); range is 1-12 when `enable_sequential=true` (default: 12). |
| `resolution` | string | no | allowed: `1K`, `2K`, `4K`; default: `"2K"` | Output resolution. The current project uses `resolution` as a wrapper field corresponding to the underlying resolution parameter.(4K generation is available only for text-to-image in Standard Mode) |
| `thinking_mode` | boolean | no | default: `false` | Whether to enable thinking mode. Only available when `enable_sequential=false` and `input_urls` is empty; the frontend will automatically disable it in other cases. |
| `color_palette` | array | no | min items: 3; max items: 10 | (Optional) Custom color theme. Only available when `enable_sequential=false`. Requires 3-10 colors, 8 recommended. |
| `color_palette[].hex` | string | no |  | HEX color value. |
| `color_palette[].ratio` | string | no |  | Color proportion, format must be xx.xx%. |
| `bbox_list` | array | no |  | (Optional) Interactive editing bounding box areas. The outer list length should match `input_urls`; maximum 2 boxes per image; single box format is `[x1, y1, x2, y2]`. |
| `watermark` | boolean | no | default: `false` | Whether to add watermark. |
| `seed` | integer | no | default: `0`; min: 0; max: 2147483647 | Random seed, range 0-2147483647. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

## `wan/2-7-image-to-video`

**Wan 2.7 - Image to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-image-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Positive prompt. Maximum length: 5000 characters. |
| `negative_prompt` | string | no | max length: 500 | Negative prompt. Maximum length: 500 characters. |
| `first_frame_url` | string | no | format: uri | First frame image URL. |
| `last_frame_url` | string | no | format: uri | Last frame image URL. |
| `first_clip_url` | string | no | format: uri | First clip video URL, used for video continuation. |
| `driving_audio_url` | string | no | format: uri | Driving audio URL. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution.   - `720p`: 720p - `1080p`: 1080p |
| `duration` | integer | no | default: `5`; min: 2; max: 15 | Total output video duration in seconds.  - Minimum: `2` - Maximum: `15` - Default: `5` |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default value: `true`. |
| `watermark` | boolean | no | default: `false` | Whether to add an AI-generated watermark. Default value: `false`. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed.  - Minimum: `0` - Maximum: `2147483647` |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A white cat stands on a windowsill in warm afternoon light. The camera slowly pushes in as the cat blinks softly and turns to look outside.",
  "negative_prompt": "blurry, flicker, low quality, distorted",
  "first_frame_url": "https://your-domain.com/assets/first-frame.png",
  "last_frame_url": "https://your-domain.com/assets/last-frame.png",
  "resolution": "1080p",
  "duration": 5,
  "prompt_extend": true,
  "watermark": false,
  "seed": 123456
}
```

## `wan/2-7-r2v`

**Wan 2.7 - Reference to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-r2v)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 5000 | Text prompt. Required. Describes the desired elements and visual features in the generated video. Supports Chinese and English. Maximum length: 5000 characters. |
| `negative_prompt` | string | no | max length: 500 | Optional negative prompt describing what should not appear in the video. Supports Chinese and English. Maximum length: 500 characters. |
| `reference_image` | array | no | max items: 5 | Array of reference image URLs. At least one of `reference_image` or `reference_video` must be provided. The total number of images and videos cannot exceed 5. |
| `reference_video` | array | no | max items: 5 | Array of reference video URLs. At least one of `reference_image` or `reference_video` must be provided. The total number of images and videos cannot exceed 5. |
| `first_frame` | string | no | format: uri | First frame image URL. At most one image can be provided. If supplied, `aspect_ratio` is ignored and the output uses a ratio close to the first frame image. |
| `reference_voice` | string | no | format: uri | Audio URL used to specify the voice timbre of the subject in the reference material.  Rules: - If `reference_video` contains audio and `reference_voice` is not provided, the original video audio is used by default - If both `reference_video` and `reference_voice` are provided, `reference_voice` takes priority  Audio limits: - Formats: `wav`, `mp3` - Duration: `1` to `10` seconds - File size: up to `15MB` |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution tier. Available values: `720p`, `1080p`. Default value: `1080p`. |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`; default: `"16:9"` | Output video aspect ratio.  Effective logic: - If `first_frame` is not provided: the video is generated using the specified `aspect_ratio` - If `first_frame` is provided: `aspect_ratio` is ignored and the output uses a ratio close to the first frame image |
| `duration` | integer | no | default: `5`; min: 2; max: 10 | Output video duration in seconds. Valid range is an integer from `2` to `10`. Default value: `5`. |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable prompt rewriting. When enabled, the model expands the input prompt. This usually works better for short prompts but increases processing time. |
| `watermark` | boolean | no | default: `false` | Whether to add a watermark. The watermark is placed in the lower-right corner of the video with the fixed text "AI generated". |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed. Range: `0-2147483647`. If omitted, the system generates one automatically. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Image 1 is eating, while video 1 and image 2 are singing beside it.",
  "negative_prompt": "low resolution, errors, worst quality, low quality, malformed, extra fingers, bad proportions",
  "reference_image": [
    "https://example.com/demo/ref-image-1.png",
    "https://example.com/demo/ref-image-2.png"
  ],
  "reference_video": [
    "https://example.com/demo/ref-video-1.mp4"
  ],
  "first_frame": "https://example.com/demo/first-frame.png",
  "reference_voice": "https://example.com/demo/reference-voice.mp3",
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 5,
  "prompt_extend": true,
  "watermark": false,
  "seed": 0
}
```

## `wan/2-7-text-to-video`

**Wan 2.7 - Text to Video** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-text-to-video)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | min length: 1; max length: 5000 | Positive prompt. Minimum length: 1 character. Maximum length: 5000 characters. |
| `negative_prompt` | string | no | max length: 500 | Negative prompt. Maximum length: 500 characters. |
| `audio_url` | string | no | format: uri | Optional custom audio URL. |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Video resolution.   - `720p`: 720p - `1080p`: 1080p |
| `ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`; default: `"16:9"` | Video aspect ratio.  - `16:9`: Landscape - `9:16`: Portrait - `1:1`: Square - `4:3`: Landscape 4:3 - `3:4`: Portrait 3:4 |
| `duration` | integer | no | default: `5`; min: 2; max: 15 | Video duration in seconds.  - Minimum: `2` - Maximum: `15` - Default: `5` |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable intelligent prompt rewriting. Default value: `true`. |
| `watermark` | boolean | no | default: `false` | Whether to add an AI-generated watermark. Default value: `false`. |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed.  - Minimum: `0` - Maximum: `2147483647` |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "A futuristic city street at night, neon reflections shimmering on the wet ground. The camera slowly pushes forward as a silver hover car glides in from the left. Giant holographic billboards flicker in the distance, creating a cinematic atmosphere.",
  "negative_prompt": "blurry, low quality, flicker, distorted characters",
  "audio_url": "https://your-domain.com/audio/custom-track.mp3",
  "resolution": "1080p",
  "ratio": "16:9",
  "duration": 5,
  "prompt_extend": true,
  "watermark": false,
  "seed": 123456
}
```

## `wan/2-7-videoedit`

**Wan 2.7 - Video Edit** · Wan · [official docs](https://docs.kie.ai/market/wan/2-7-videoedit)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | no | max length: 5000 | Optional text prompt describing the expected elements and visual features in the generated video. Supports Chinese and English. Maximum length: 5000 characters. |
| `negative_prompt` | string | no | max length: 500 | Optional negative prompt describing content that should not appear in the video. Supports Chinese and English. Maximum length: 500 characters. |
| `video_url` | string | yes | format: uri | URL of the source video to edit. Required. Only one video is supported.  - Formats: `mp4`, `mov` - Duration: `2` to `10` seconds - Resolution: width and height range `[240,4096]` pixels - Aspect ratio: `1:8` to `8:1` - File size: up to `100MB` - Supports public `http/https` URLs or temporary `oss` URLs |
| `reference_image` | string | no | format: uri | Optional reference image URL for character, clothing, or style guidance.  - Formats: `JPEG`, `JPG`, `PNG` (no alpha channel), `BMP`, `WEBP` - Resolution: width and height range `[240,8000]` pixels - Aspect ratio: `1:8` to `8:1` - Supports public `http/https` URLs or temporary `oss` URLs |
| `resolution` | string | no | allowed: `720p`, `1080p`; default: `"1080p"` | Output video resolution tier. `1080p` costs more than `720p`. Default value: `1080p`.  - `720p`: 720p - `1080p`: 1080p |
| `aspect_ratio` | string | no | allowed: `16:9`, `9:16`, `1:1`, `4:3`, `3:4` | Output video aspect ratio.  - If omitted: the output uses an aspect ratio close to the input video - If provided: the output uses the specified aspect ratio - Available values: `16:9`, `9:16`, `1:1`, `4:3`, `3:4` |
| `duration` | integer | no | default: `0`; min: 0; max: 10 | Output video duration in seconds.  - Default `0` means using the full input video duration without truncation - If a value is provided, the output is clipped from second `0` to the specified length - Valid values are `0` or any integer in `[2,10]` |
| `audio_setting` | string | no | allowed: `auto`, `origin`; default: `"auto"` | Video audio setting.  - `auto`: default, the model decides whether to regenerate audio based on the `prompt` - `origin`: force keeping the original input video audio |
| `prompt_extend` | boolean | no | default: `true` | Whether to enable prompt rewriting. When enabled, the model expands the input prompt. This usually works better for short prompts but increases processing time. |
| `watermark` | boolean | no | default: `false` | Whether to add a watermark. The watermark is placed in the lower-right corner of the video with the fixed text "AI generated". |
| `seed` | integer | no | min: 0; max: 2147483647 | Random seed. Range: `0-2147483647`. If omitted, the system generates one automatically. |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Change the character's outfit and add the hat shown in the reference image.",
  "negative_prompt": "low resolution, errors, worst quality, low quality, malformed, extra fingers, bad proportions",
  "video_url": "https://example.com/demo/video.mp4",
  "reference_image": "https://example.com/demo/reference.png",
  "resolution": "1080p",
  "aspect_ratio": "16:9",
  "duration": 0,
  "audio_setting": "auto",
  "prompt_extend": true,
  "watermark": false,
  "seed": 0
}
```

## `z-image`

**Z-Image** · Z Image · [official docs](https://docs.kie.ai/market/z-image/z-image)

| Input field | Type | Required | Settings | Description |
|---|---|---|---|---|
| `prompt` | string | yes | max length: 1000 | A text description of the image you want to generate (Max length: 1000 characters) |
| `aspect_ratio` | string | yes | allowed: `1:1`, `4:3`, `3:4`, `16:9`, `9:16`; default: `"1:1"` | Aspect ratio for the generated image. Select 'auto' to match the first input image ratio (requires input image). |
| `nsfw_checker` | boolean | no |  | Defaults to false. You can set it to false based on your needs. If set to false, our content filtering will be disabled, and all results will be returned directly by the model itself. Note: There is no guarantee that everything can be filtered out; if you are not satisfied with the results, you will need to make your own arrangements. |

Example `input`:

```json
{
  "prompt": "Generate a photorealistic image of a cafe terrace in the Marais district of Paris on a Wednesday morning in March 2025. It is a crisp, cool spring morning with clear skies. Locals are drinking coffee. In sharp focus should be a young woman with a pixie cut wearing a scarf, stirring a cappuccino and looking thoughtfully to the side; the waiter and street traffic behind her are blurred. The photo should have the candid, natural morning light feel of an iPhone image.",
  "aspect_ratio": "1:1",
  "nsfw_checker": true
}
```
