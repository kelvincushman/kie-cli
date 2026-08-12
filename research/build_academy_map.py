#!/usr/bin/env python3
"""Build a source-linked, copyright-safe map of Higgsfield Academy lessons.

The generator visits every public course and lesson page to verify that the
catalog still exists, but stores only factual page metadata (titles, positions,
durations, and URLs) plus original Kie-native workflow guidance. It deliberately
does not archive lesson prose, scripts, prompts, media, or compiled page data.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
import time
from dataclasses import dataclass
from datetime import UTC, datetime
from html.parser import HTMLParser
from pathlib import Path
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen


BASE_URL = "https://higgsfield.ai"
COURSES_URL = f"{BASE_URL}/academy/courses"
USER_AGENT = "kie-cli-academy-index/1.0 (+https://github.com/kelvincushman/kie-cli)"


class PageFacts(HTMLParser):
    """Collect only links, JSON-LD, title, and short metadata from a page."""

    def __init__(self) -> None:
        super().__init__()
        self.links: list[tuple[str, str]] = []
        self.json_ld: list[dict[str, Any]] = []
        self.description = ""
        self._href: str | None = None
        self._link_text: list[str] = []
        self._json = False
        self._json_text: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        values = dict(attrs)
        if tag == "a":
            self._href = values.get("href")
            self._link_text = []
        elif tag == "script" and values.get("type") == "application/ld+json":
            self._json = True
            self._json_text = []
        elif tag == "meta" and values.get("name") == "description":
            self.description = values.get("content") or ""

    def handle_data(self, data: str) -> None:
        if self._href is not None:
            self._link_text.append(data)
        if self._json:
            self._json_text.append(data)

    def handle_endtag(self, tag: str) -> None:
        if tag == "a" and self._href is not None:
            text = " ".join(" ".join(self._link_text).split())
            self.links.append((self._href, text))
            self._href = None
            self._link_text = []
        elif tag == "script" and self._json:
            try:
                value = json.loads("".join(self._json_text))
                if isinstance(value, dict):
                    self.json_ld.append(value)
            except json.JSONDecodeError:
                pass
            self._json = False
            self._json_text = []


@dataclass(frozen=True)
class CourseSummary:
    slug: str
    title: str
    lesson_count: int
    duration_minutes: int
    difficulty: str


STAGE_GUIDANCE = {
    "reference-review": {
        "method": "Study the example for beats, recurring assets, transitions, and visible failure modes; record observations as constraints rather than imitating the source wording.",
        "prompt_focus": "desired outcome, audience reaction, continuity risks, acceptance criteria",
    },
    "creative-brief": {
        "method": "Turn the goal into a bounded production brief with audience, platform, duration, aspect ratio, required truth, references, and a clear definition of done.",
        "prompt_focus": "subject, purpose, format, constraints, non-negotiables",
    },
    "research": {
        "method": "Collect a small evidence set, label what is factual versus inspirational, and translate it into an original visual direction without copying another creator's execution.",
        "prompt_focus": "source evidence, audience language, differentiators, prohibited imitation",
    },
    "script": {
        "method": "Write short producible beats where each beat has one dramatic purpose, one visible action, and only the dialogue or narration that can fit its duration.",
        "prompt_focus": "hook, beat objective, action, dialogue, payoff",
    },
    "asset-lock": {
        "method": "Build reusable character, product, prop, wardrobe, and location references before scene work; define immutable identity traits separately from changeable styling.",
        "prompt_focus": "identity traits, reference views, materials, palette, scale, exclusions",
    },
    "storyboard": {
        "method": "Convert the approved script into ordered shots with duration, composition, action, camera, audio, continuity inputs, and a specific approval criterion for each frame.",
        "prompt_focus": "shot size, camera, blocking, timing, continuity in/out, transition",
    },
    "prompt-design": {
        "method": "Write one model-specific prompt per shot, separating locked references from the requested action, camera behavior, lighting, audio, and negative constraints.",
        "prompt_focus": "locked assets, single action, camera, lighting, audio, negatives",
    },
    "keyframe": {
        "method": "Generate and inspect a still keyframe before motion. Use GPT Image 2 for first-pass quality and precise reference work, or Nano Banana for controlled image edits and continuity variants.",
        "prompt_focus": "approved composition, exact identity, product fidelity, clean geometry, no motion blur",
    },
    "motion": {
        "method": "Animate only an approved keyframe, keep each shot to one primary action, and review geography, physics, identity, hands, props, dialogue timing, and camera intent before accepting the take.",
        "prompt_focus": "first frame, primary action, camera path, temporal continuity, audio timing",
    },
    "post-production": {
        "method": "Assemble only approved takes, then check continuity at every cut before adding sound, captions, grade, reframing, upscale, or export as separate reversible passes.",
        "prompt_focus": "selected take, edit point, transition, mix, captions, delivery format",
    },
    "delivery": {
        "method": "Package the finished media for its destination, verify claims and formatting, publish only with authorization, and feed observed performance back into the next brief.",
        "prompt_focus": "platform requirements, title/thumbnail, claims, export, measurement plan",
    },
    "review": {
        "method": "Compare the result with the brief and storyboard, reject visible defects instead of rationalizing them, and record reusable lessons for the next production.",
        "prompt_focus": "brief match, continuity, defects, strongest take, next iteration",
    },
}


KEYWORDS = [
    ("delivery", ("publish", "upload", "monetiz", "package", "thumbnail", "results", "analytics", "launch", "distribution")),
    ("post-production", ("edit", "vfx", "sound", "music", "upscal", "composit", "assembl", "cutting", "color grade", "export", "refram")),
    ("prompt-design", ("prompt",)),
    ("storyboard", ("storyboard", "shot list", "coverage", "sequence", "beat sheet", "blueprint", "continuity")),
    ("review", ("recap", "conclusion", "review", "diagnos", "failure", "failed", "quality", "final verdict", "check ", "audit ", "test ", "read your", "trade-off")),
    ("reference-review", ("watch", "intro", "challenge", "example", "overview", "why now", "finished ad")),
    ("script", ("script", "story", "narrative", "dialogue", "voiceover", "concept", "hook")),
    ("keyframe", ("image", "keyframe", "first frame", "still", "look development", "background")),
    ("motion", ("video", "motion", "animat", "camera", "scene", "chase", "fight", "action", "transition", "impact", "transformation")),
    ("asset-lock", ("character", "asset", "location", "prop", "wardrobe", "product", "brand", "logo", "reference sheet", "consistent", "outfit")),
    ("research", ("research", "niche", "competitor", "trend", "reference", "inspiration")),
    ("creative-brief", ("plan", "brief", "goal", "choose", "format", "workflow", "setup", "strategy", "evaluate")),
]


def fetch(url: str, attempts: int = 3) -> str:
    error: Exception | None = None
    for attempt in range(attempts):
        try:
            request = Request(url, headers={"User-Agent": USER_AGENT})
            with urlopen(request, timeout=30) as response:  # noqa: S310 - fixed public HTTPS URLs
                return response.read().decode("utf-8")
        except (HTTPError, URLError, TimeoutError) as exc:
            error = exc
            if attempt + 1 < attempts:
                time.sleep(0.5 * (attempt + 1))
    raise RuntimeError(f"could not fetch {url}: {error}")


def parse_page(html: str) -> PageFacts:
    parser = PageFacts()
    parser.feed(html)
    return parser


def discover_courses(html: str) -> list[CourseSummary]:
    slugs = list(dict.fromkeys(re.findall(r"/academy/courses/([a-z0-9-]+)", html)))
    summaries: list[CourseSummary] = []
    for slug in slugs:
        pattern = re.compile(
            rf'slug:"{re.escape(slug)}",title:"([^"]+)",description:.*?,'
            rf"lesson_count:(\d+),reading_minutes:\d+,duration_minutes:(\d+),difficulty:\"([^\"]+)\"",
            re.DOTALL,
        )
        match = pattern.search(html)
        if not match:
            raise RuntimeError(f"could not parse public course metadata for {slug}")
        summaries.append(
            CourseSummary(
                slug=slug,
                title=match.group(1).replace(r"\'", "'"),
                lesson_count=int(match.group(2)),
                duration_minutes=int(match.group(3)),
                difficulty=match.group(4),
            )
        )
    return summaries


def extract_lessons(course: CourseSummary, facts: PageFacts) -> list[dict[str, Any]]:
    prefix = f"/academy/courses/{course.slug}/"
    lessons: dict[str, dict[str, Any]] = {}
    for href, text in facts.links:
        if not href.startswith(prefix):
            continue
        slug = href[len(prefix) :].strip("/")
        if not slug or "/" in slug:
            continue
        match = re.match(r"Lesson\s+(\d+)\s+(.+)", text)
        if not match:
            continue
        lessons[slug] = {
            "slug": slug,
            "title": match.group(2).strip(),
            "position": int(match.group(1)),
            "source_url": f"{COURSES_URL}/{course.slug}/{slug}",
        }
    result = sorted(lessons.values(), key=lambda item: item["position"])
    if len(result) != course.lesson_count:
        raise RuntimeError(
            f"{course.slug}: expected {course.lesson_count} lessons, found {len(result)}"
        )
    return result


def classify(title: str, description: str) -> str:
    # The public lesson title is the authoritative classification signal.
    # A page description often repeats the wider course synopsis and can make
    # an editing lesson look like an asset lesson (or vice versa).
    for text in (title.lower(), description.lower()):
        for stage, words in KEYWORDS:
            if any(word in text for word in words):
                return stage
    return "creative-brief"


def model_routes(stage: str) -> list[str]:
    if stage in {"asset-lock", "keyframe"}:
        return [
            "gpt-image-2-text-to-image",
            "gpt-image-2-image-to-image",
            "nano-banana-2",
            "nano-banana-pro",
        ]
    if stage == "motion":
        return ["bytedance/seedance-2-5"]
    return []


def capability_boundary(course_slug: str, stage: str) -> str:
    if course_slug == "build-3d-games-mcp":
        return "Apply the agent-planning and asset-production method only; Kie does not provide a true 3D mesh/game build-and-publish endpoint."
    if stage == "delivery":
        return "Publishing, analytics, and platform compliance remain outside Kie and require explicit local or platform-specific actions."
    if stage == "post-production":
        return "Final edit, compositing, captions, and assembly remain transparent local post-production steps."
    return ""


def build_catalog(verify_lessons: bool = True) -> dict[str, Any]:
    index_html = fetch(COURSES_URL)
    course_summaries = discover_courses(index_html)
    courses: list[dict[str, Any]] = []
    for summary in course_summaries:
        course_url = f"{COURSES_URL}/{summary.slug}"
        facts = parse_page(fetch(course_url))
        course_schema = next((item for item in facts.json_ld if item.get("@type") == "Course"), None)
        title = str(course_schema.get("name")) if course_schema else summary.title
        lessons = extract_lessons(summary, facts)
        for lesson in lessons:
            lesson_facts = parse_page(fetch(lesson["source_url"])) if verify_lessons else PageFacts()
            stage = classify(lesson["title"], lesson_facts.description)
            guidance = STAGE_GUIDANCE[stage]
            lesson.update(
                {
                    "key": f"{summary.slug}/{lesson['slug']}",
                    "production_stage": stage,
                    "kie_method": guidance["method"],
                    "prompt_focus": guidance["prompt_focus"],
                    "recommended_kie_models": model_routes(stage),
                }
            )
            boundary = capability_boundary(summary.slug, stage)
            if boundary:
                lesson["capability_boundary"] = boundary
        courses.append(
            {
                "slug": summary.slug,
                "title": title,
                "source_url": course_url,
                "lesson_count": summary.lesson_count,
                "duration_minutes": summary.duration_minutes,
                "difficulty": summary.difficulty,
                "lessons": lessons,
            }
        )
    return {
        "schema_version": 1,
        "generated_at": datetime.now(UTC).isoformat(),
        "source": COURSES_URL,
        "copyright_policy": "Public titles, positions, durations, and links only; Kie methods and prompt guidance are original adaptations. No Higgsfield lesson prose, scripts, prompts, media, or compiled page content is stored.",
        "course_count": len(courses),
        "lesson_count": sum(len(course["lessons"]) for course in courses),
        "courses": courses,
    }


def render_markdown(catalog: dict[str, Any]) -> str:
    lines = [
        "# Academy-Inspired Production Methods",
        "",
        f"Generated from the public [Higgsfield Academy course index]({catalog['source']}) on "
        f"{catalog['generated_at'][:10]}. The snapshot maps **{catalog['course_count']} courses** and "
        f"**{catalog['lesson_count']} lessons**.",
        "",
        "> Copyright boundary: this document stores public lesson titles, ordering, durations, and source links, then adds original Kie-native method abstractions. It does not reproduce Higgsfield lesson prose, scripts, prompts, videos, or downloadable resources. Open the linked source to study the original lesson.",
        "",
        "Use `kie-pp-cli lesson recommend \"<what you want to create>\" --agent` to find a pattern, "
        "or invoke the installed `$kie-lesson` skill (`/kie-lesson` in hosts that expose skills as slash commands).",
        "",
        "Every model listed below is only a route candidate. Inspect `kie-pp-cli media leaderboard`, "
        "then validate its live Kie input contract before spending credits.",
        "",
    ]
    for course in catalog["courses"]:
        lines.extend(
            [
                f"## [{course['title']}]({course['source_url']})",
                "",
                f"{course['lesson_count']} lessons · {course['duration_minutes']} minutes · {course['difficulty']}",
                "",
                "| # | Public lesson | Production stage | Original Kie-native adaptation |",
                "| ---: | --- | --- | --- |",
            ]
        )
        for lesson in course["lessons"]:
            method = lesson["kie_method"]
            if lesson.get("capability_boundary"):
                method += " **Boundary:** " + lesson["capability_boundary"]
            lines.append(
                f"| {lesson['position']} | [{lesson['title']}]({lesson['source_url']}) | "
                f"`{lesson['production_stage']}` | {method} |"
            )
        lines.append("")
    return "\n".join(lines).rstrip() + "\n"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, default=Path("internal/academy/catalog.json"))
    parser.add_argument("--docs", type=Path, default=Path("docs/ACADEMY_METHODS.md"))
    parser.add_argument(
        "--skip-lesson-verification",
        action="store_true",
        help="Fetch course pages but do not visit all individual lesson URLs",
    )
    parser.add_argument(
        "--preserve-generated-at-if-unchanged",
        action="store_true",
        help="retain the previous timestamp when all source-derived content is unchanged",
    )
    args = parser.parse_args()
    catalog = build_catalog(verify_lessons=not args.skip_lesson_verification)
    if catalog["course_count"] == 0 or catalog["lesson_count"] == 0:
        raise RuntimeError("academy catalog was unexpectedly empty")
    if args.preserve_generated_at_if_unchanged and args.output.exists():
        previous = json.loads(args.output.read_text())
        comparable = dict(catalog)
        comparable["generated_at"] = previous.get("generated_at")
        if comparable == previous:
            catalog["generated_at"] = previous["generated_at"]
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.docs.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(catalog, indent=2, ensure_ascii=False) + "\n")
    args.docs.write_text(render_markdown(catalog))
    print(
        f"Mapped {catalog['course_count']} public courses and {catalog['lesson_count']} lessons "
        f"to {args.output} and {args.docs}",
        file=sys.stderr,
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
