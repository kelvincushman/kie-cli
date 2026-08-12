# Media Model Leaderboard

Refreshed **2026-08-12**. Task-specific evidence ledger. Source ranks and scores are not normalized or combined. Kie availability is a separate routing constraint, and character consistency uses image-edit performance only as a disclosed proxy.

> This is a routing aid, not a universal declaration of the best model. Human preference, reference fidelity, motion quality, speed, price, and Kie availability are separate dimensions. A paid generation still requires local contract validation and human review.

## Text-to-image quality

direct independent human-preference evidence.

| Source rank | Kie model | Published score | Match | Kie route |
| ---: | --- | ---: | --- | --- |
| 1 | `gpt-image-2-text-to-image` | 1381 Arena score | direct | available |
| 7 | `nano-banana-2` | 1264 Arena score | direct | available |
| 11 | `nano-banana-2-lite` | 1251 Arena score | direct | available |
| 12 | `nano-banana-pro` | 1246 Arena score | direct | available |

Sources: [Arena Text-to-Image](https://arena.ai/leaderboard/text-to-image/), [Kie.ai generated model-contract catalog](https://docs.kie.ai/market/quickstart).

## Reference-led image editing

direct independent human-preference evidence.

| Source rank | Kie model | Published score | Match | Kie route |
| ---: | --- | ---: | --- | --- |
| 1 | `gpt-image-2-image-to-image` | 1463 Arena score | direct | available |
| 8 | `nano-banana-pro` | 1389 Arena score | direct | available |
| 10 | `nano-banana-2` | 1385 Arena score | direct | available |
| 17 | `nano-banana-2-lite` | 1314 Arena score | direct | available |

Sources: [Arena Image Edit](https://arena.ai/leaderboard/image-edit/), [Kie.ai generated model-contract catalog](https://docs.kie.ai/market/quickstart).

## Character/reference consistency

image-edit proxy plus local workflow constraints.

> No source used here publishes a cross-model character-consistency benchmark. Treat this as routing guidance, then inspect every generated reference sheet and shot.

| Source rank | Kie model | Published score | Match | Kie route |
| ---: | --- | ---: | --- | --- |
| 1 | `gpt-image-2-image-to-image` | 1463 Arena score | direct | available |
| 8 | `nano-banana-pro` | 1389 Arena score | direct | available |
| 10 | `nano-banana-2` | 1385 Arena score | direct | available |

Sources: [Arena Image Edit](https://arena.ai/leaderboard/image-edit/), [Kie.ai generated model-contract catalog](https://docs.kie.ai/market/quickstart).

## Text-to-video quality

direct evidence where model versions match; explicit family proxy otherwise.

| Source rank | Kie model | Published score | Match | Kie route |
| ---: | --- | ---: | --- | --- |
| — | `bytedance/seedance-2-5` | — unscored successor | family proxy disclosed | available |
| 3 | `bytedance/seedance-2` | 1223 Elo | direct | available |
| 4 | `wan/2-7-text-to-video` | 1160 Elo | direct | available |
| 8 | `kling-3.0/video` | 1109 Elo | direct | available |

Sources: [Artificial Analysis Text-to-Video Arena](https://artificialanalysis.ai/embed/text-to-video-leaderboard/leaderboard/text-to-video), [Kie.ai generated model-contract catalog](https://docs.kie.ai/market/quickstart).

## Director defaults

- New still: `gpt-image-2-text-to-image` because its current direct text-to-image evidence is strongest among the compared Kie routes.
- Reference/identity still: `gpt-image-2-image-to-image`; use `nano-banana-pro` or `nano-banana-2` when controlled editing, batching, or a user preference makes them a better fit.
- Final video: `bytedance/seedance-2-5` for its Kie multimodal contract. The current independent score shown for Seedance 2.0 is family context only, not a Seedance 2.5 score.
- Every video shot: approved still first, then motion generation, then continuity review.
