#!/usr/bin/env python3
"""Refresh task-specific media rankings from public, independent sources.

The output never invents a universal quality score. It preserves each source's
published rank/score and marks character consistency as an image-editing proxy.
Kie availability is resolved against the repository's generated model catalog.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from datetime import UTC, datetime
from html.parser import HTMLParser
from pathlib import Path
from typing import Any
from urllib.request import Request, urlopen


USER_AGENT = "kie-cli-model-leaderboard/1.0 (+https://github.com/kelvincushman/kie-cli)"
SOURCES = {
    "arena-text-to-image": {
        "name": "Arena Text-to-Image",
        "url": "https://arena.ai/leaderboard/text-to-image/",
        "evidence_type": "crowdsourced human preference",
    },
    "arena-image-edit": {
        "name": "Arena Image Edit",
        "url": "https://arena.ai/leaderboard/image-edit/",
        "evidence_type": "crowdsourced human preference",
    },
    "artificial-analysis-text-to-video": {
        "name": "Artificial Analysis Text-to-Video Arena",
        "url": "https://artificialanalysis.ai/embed/text-to-video-leaderboard/leaderboard/text-to-video",
        "evidence_type": "crowdsourced human preference",
    },
    "kie-market": {
        "name": "Kie.ai generated model-contract catalog",
        "url": "https://docs.kie.ai/market/quickstart",
        "evidence_type": "provider availability and input contract",
    },
}


class TableParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.rows: list[list[str]] = []
        self._row_depth = 0
        self._cell_depth = 0
        self._row: list[str] = []
        self._cell: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        del attrs
        if tag == "tr":
            self._row_depth += 1
            self._row = []
        elif self._row_depth and tag in {"td", "th"}:
            self._cell_depth += 1
            self._cell = []

    def handle_data(self, data: str) -> None:
        if self._cell_depth:
            self._cell.append(data)

    def handle_endtag(self, tag: str) -> None:
        if self._cell_depth and tag in {"td", "th"}:
            self._row.append(" ".join(" ".join(self._cell).split()))
            self._cell_depth -= 1
        elif self._row_depth and tag == "tr":
            if self._row:
                self.rows.append(self._row)
            self._row_depth -= 1


IMAGE_ROUTES = [
    {
        "source_model": "gpt-image-2 (medium)",
        "text_model": "gpt-image-2-text-to-image",
        "edit_model": "gpt-image-2-image-to-image",
        "strengths": ["overall image quality", "reference-led generation", "precise image editing"],
    },
    {
        "source_model": "gemini-3.1-flash-image (nano-banana-2) [web-search]",
        "text_model": "nano-banana-2",
        "edit_model": "nano-banana-2",
        "strengths": ["fast iterative generation", "reference edits", "web-grounded variants when supported"],
    },
    {
        "source_model": "gemini-3-pro-image-2k (nano-banana-pro)",
        "text_model": "nano-banana-pro",
        "edit_model": "nano-banana-pro",
        "strengths": ["high-resolution reference editing", "character-sheet variants", "layout control"],
    },
    {
        "source_model": "gemini-3.1-flash-lite-image (nano-banana-2-lite)",
        "text_model": "nano-banana-2-lite",
        "edit_model": "nano-banana-2-lite",
        "strengths": ["lower-cost iteration", "draft variations"],
    },
]

VIDEO_ROUTES = [
    {
        "source_model": "Dreamina Seedance 2.0 720p",
        "kie_model": "bytedance/seedance-2",
        "strengths": ["cinematic motion", "prompt adherence", "audio-capable family"],
    },
    {
        "source_model": "Wan2.7-260612",
        "kie_model": "wan/2-7-text-to-video",
        "strengths": ["text-to-video quality", "current Wan family"],
    },
    {
        "source_model": "Wan 2.7",
        "kie_model": "wan/2-7-text-to-video",
        "strengths": ["text-to-video quality", "current Wan family"],
    },
    {
        "source_model": "Kling 3.0 1080p (Pro)",
        "kie_model": "kling-3.0/video",
        "strengths": ["high-resolution motion", "general video quality"],
    },
]


def fetch(url: str) -> str:
    request = Request(url, headers={"User-Agent": USER_AGENT})
    with urlopen(request, timeout=45) as response:  # noqa: S310 - fixed HTTPS sources
        return response.read().decode("utf-8")


def table_rows(html: str) -> list[list[str]]:
    parser = TableParser()
    parser.feed(html)
    return parser.rows


def parse_number(value: str) -> int:
    match = re.search(r"[\d,]+", value)
    if not match:
        raise RuntimeError(f"could not parse numeric leaderboard value {value!r}")
    return int(match.group(0).replace(",", ""))


def find_row(rows: list[list[str]], model: str) -> list[str]:
    matches = [row for row in rows if any(model.lower() in cell.lower() for cell in row)]
    if not matches:
        raise RuntimeError(f"model {model!r} was not found in the current source table")
    return matches[0]


def image_entry(route: dict[str, Any], rows: list[list[str]], mode: str, available: set[str]) -> dict[str, Any]:
    row = find_row(rows, route["source_model"])
    model = route["text_model"] if mode == "text" else route["edit_model"]
    return {
        "source_rank": parse_number(row[0]),
        "source_score": parse_number(row[3]),
        "score_type": "Arena score",
        "source_model": route["source_model"],
        "kie_model": model,
        "kie_available": model in available,
        "direct_match": True,
        "strengths": route["strengths"],
    }


def video_entry(route: dict[str, Any], rows: list[list[str]], available: set[str]) -> dict[str, Any]:
    row = find_row(rows, route["source_model"])
    return {
        "source_rank": parse_number(row[0]),
        "source_score": parse_number(row[4]),
        "score_type": "Elo",
        "source_model": route["source_model"],
        "kie_model": route["kie_model"],
        "kie_available": route["kie_model"] in available,
        "direct_match": True,
        "strengths": route["strengths"],
    }


def build(kie_catalog_path: Path) -> dict[str, Any]:
    now = datetime.now(UTC).isoformat()
    kie_catalog = json.loads(kie_catalog_path.read_text())
    available = {model["id"] for model in kie_catalog["models"]}
    pages = {key: fetch(source["url"]) for key, source in SOURCES.items() if key != "kie-market"}
    t2i_rows = table_rows(pages["arena-text-to-image"])
    edit_rows = table_rows(pages["arena-image-edit"])
    video_rows = table_rows(pages["artificial-analysis-text-to-video"])

    text_entries = [image_entry(route, t2i_rows, "text", available) for route in IMAGE_ROUTES]
    edit_entries = [image_entry(route, edit_rows, "edit", available) for route in IMAGE_ROUTES]
    text_entries.sort(key=lambda entry: entry["source_rank"])
    edit_entries.sort(key=lambda entry: entry["source_rank"])
    video_entries: list[dict[str, Any]] = []
    seen_video_models: set[str] = set()
    for route in VIDEO_ROUTES:
        entry = video_entry(route, video_rows, available)
        if entry["kie_model"] in seen_video_models:
            continue
        seen_video_models.add(entry["kie_model"])
        video_entries.append(entry)

    # The director targets Seedance 2.5 for its newer Kie contract. No source in
    # this refresh publishes a direct Seedance 2.5 score, so keep the family
    # evidence explicit instead of assigning Seedance 2.0's Elo to it.
    seedance_family = next(entry for entry in video_entries if entry["kie_model"] == "bytedance/seedance-2")
    video_entries.insert(
        0,
        {
            "source_rank": None,
            "source_score": None,
            "score_type": "unscored successor",
            "source_model": "Seedance 2.5 (direct independent score unavailable)",
            "kie_model": "bytedance/seedance-2-5",
            "kie_available": "bytedance/seedance-2-5" in available,
            "direct_match": False,
            "proxy": {
                "source_model": seedance_family["source_model"],
                "source_rank": seedance_family["source_rank"],
                "source_score": seedance_family["source_score"],
                "warning": "Family evidence only; do not transfer the score to Seedance 2.5.",
            },
            "strengths": ["Kie multimodal orchestration", "first/last frames", "image/video/audio references", "generated audio"],
        },
    )

    source_records = []
    for source_id, source in SOURCES.items():
        source_records.append({"id": source_id, **source, "fetched_at": now})
    return {
        "schema_version": 1,
        "generated_at": now,
        "methodology": "Task-specific evidence ledger. Source ranks and scores are not normalized or combined. Kie availability is a separate routing constraint, and character consistency uses image-edit performance only as a disclosed proxy.",
        "sources": source_records,
        "tasks": [
            {
                "id": "text-to-image",
                "label": "Text-to-image quality",
                "source_ids": ["arena-text-to-image", "kie-market"],
                "ranking_kind": "direct independent human-preference evidence",
                "entries": text_entries,
            },
            {
                "id": "image-edit",
                "label": "Reference-led image editing",
                "source_ids": ["arena-image-edit", "kie-market"],
                "ranking_kind": "direct independent human-preference evidence",
                "entries": edit_entries,
            },
            {
                "id": "character-consistency",
                "label": "Character/reference consistency",
                "source_ids": ["arena-image-edit", "kie-market"],
                "ranking_kind": "image-edit proxy plus local workflow constraints",
                "warning": "No source used here publishes a cross-model character-consistency benchmark. Treat this as routing guidance, then inspect every generated reference sheet and shot.",
                "entries": edit_entries[:3],
            },
            {
                "id": "text-to-video",
                "label": "Text-to-video quality",
                "source_ids": ["artificial-analysis-text-to-video", "kie-market"],
                "ranking_kind": "direct evidence where model versions match; explicit family proxy otherwise",
                "entries": video_entries,
            },
        ],
    }


def render_markdown(leaderboard: dict[str, Any]) -> str:
    lines = [
        "# Media Model Leaderboard",
        "",
        f"Refreshed **{leaderboard['generated_at'][:10]}**. {leaderboard['methodology']}",
        "",
        "> This is a routing aid, not a universal declaration of the best model. Human preference, reference fidelity, motion quality, speed, price, and Kie availability are separate dimensions. A paid generation still requires local contract validation and human review.",
        "",
    ]
    source_by_id = {source["id"]: source for source in leaderboard["sources"]}
    for task in leaderboard["tasks"]:
        lines.extend([f"## {task['label']}", "", task["ranking_kind"] + ".", ""])
        if task.get("warning"):
            lines.extend([f"> {task['warning']}", ""])
        lines.extend(
            [
                "| Source rank | Kie model | Published score | Match | Kie route |",
                "| ---: | --- | ---: | --- | --- |",
            ]
        )
        for entry in task["entries"]:
            rank = entry["source_rank"] if entry["source_rank"] is not None else "—"
            score = entry["source_score"] if entry["source_score"] is not None else "—"
            match = "direct" if entry["direct_match"] else "family proxy disclosed"
            available = "available" if entry["kie_available"] else "not in captured Market catalog"
            lines.append(f"| {rank} | `{entry['kie_model']}` | {score} {entry['score_type']} | {match} | {available} |")
        source_links = ", ".join(
            f"[{source_by_id[source_id]['name']}]({source_by_id[source_id]['url']})"
            for source_id in task["source_ids"]
        )
        lines.extend(["", f"Sources: {source_links}.", ""])
    lines.extend(
        [
            "## Director defaults",
            "",
            "- New still: `gpt-image-2-text-to-image` because its current direct text-to-image evidence is strongest among the compared Kie routes.",
            "- Reference/identity still: `gpt-image-2-image-to-image`; use `nano-banana-pro` or `nano-banana-2` when controlled editing, batching, or a user preference makes them a better fit.",
            "- Final video: `bytedance/seedance-2-5` for its Kie multimodal contract. The current independent score shown for Seedance 2.0 is family context only, not a Seedance 2.5 score.",
            "- Every video shot: approved still first, then motion generation, then continuity review.",
            "",
        ]
    )
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--kie-catalog", type=Path, default=Path("internal/kiecatalog/catalog.json"))
    parser.add_argument("--output", type=Path, default=Path("internal/leaderboard/leaderboard.json"))
    parser.add_argument("--docs", type=Path, default=Path("docs/MODEL_LEADERBOARD.md"))
    parser.add_argument(
        "--preserve-generated-at-if-unchanged",
        action="store_true",
        help="retain previous timestamps when source-derived ranking content is unchanged",
    )
    args = parser.parse_args()
    leaderboard = build(args.kie_catalog)
    if args.preserve_generated_at_if_unchanged and args.output.exists():
        previous = json.loads(args.output.read_text())
        comparable = json.loads(json.dumps(leaderboard))
        comparable["generated_at"] = previous.get("generated_at")
        previous_source_times = {source["id"]: source.get("fetched_at") for source in previous.get("sources", [])}
        for source in comparable["sources"]:
            source["fetched_at"] = previous_source_times.get(source["id"], source["fetched_at"])
        if comparable == previous:
            leaderboard["generated_at"] = previous["generated_at"]
            for source in leaderboard["sources"]:
                source["fetched_at"] = previous_source_times[source["id"]]
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.docs.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(leaderboard, indent=2) + "\n")
    args.docs.write_text(render_markdown(leaderboard))
    print(f"Refreshed {len(leaderboard['tasks'])} model ranking tasks", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
