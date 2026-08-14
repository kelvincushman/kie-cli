#!/usr/bin/env python3
"""Generate the reviewed Kie capability classifier from checked-in API artifacts."""

from __future__ import annotations

import hashlib
import json
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
CATALOG_PATH = ROOT / "internal" / "kiecatalog" / "catalog.json"
COVERAGE_PATH = ROOT / "research" / "kie-api-coverage.json"
OPENAPI_PATH = ROOT / "research" / "kie-final-openapi.yaml"
OUTPUT_PATH = ROOT / "internal" / "kiecatalog" / "capabilities.json"

CAPABILITIES = {"kie-image", "kie-video", "kie-audio", "kie-avatar", "kie-identity"}
RESOLUTION_FIELDS = ("resolution", "quality", "mode", "size", "image_size", "output_resolution")


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def primary_capability(model_id: str, fields: set[str]) -> str:
    lowered = model_id.lower()
    if "human-identification" in lowered or "subject-detection" in lowered:
        return "kie-identity"
    if any(token in lowered for token in ("ai-avatar", "omnihuman", "infinitalk", "speech-to-video")):
        return "kie-avatar"
    if any(token in lowered for token in ("elevenlabs/", "text-to-speech", "text-to-dialogue", "audio-isolation")):
        return "kie-audio"
    if any(token in lowered for token in ("video", "motion-control", "animate-", "/extend", "/transition")):
        return "kie-video"
    if {"video_url", "video_urls", "reference_video_urls", "first_clip_url", "first_frame_url", "last_frame_url"} & fields:
        return "kie-video"
    if {"duration", "resolution"} <= fields:
        return "kie-video"
    if {"audio_url", "dialogue", "dialogue_turns", "voice"} & fields and "image_url" not in fields:
        return "kie-audio"
    return "kie-image"


def secondary_capabilities(primary: str, model_id: str, fields: set[str]) -> list[str]:
    secondary: list[str] = []
    if primary == "kie-video" and ({"audio_url", "audio", "reference_audio_urls"} & fields):
        secondary.append("kie-audio")
    if primary in {"kie-image", "kie-video"} and any(token in model_id.lower() for token in ("character", "reference", "identity")):
        secondary.append("kie-identity")
    if primary == "kie-avatar":
        secondary.extend(["kie-video", "kie-identity"])
    return sorted(set(secondary) - {primary})


def production_fit(primary: str, model_id: str) -> list[str]:
    fits = {
        "kie-image": ["kie-brandkit", "kie-marketplace-cards", "kie-product-photoshoot", "kie-youtube-thumbnail"],
        "kie-video": ["kie-film", "kie-video-explainer"],
        "kie-audio": ["kie-film", "kie-video-explainer"],
        "kie-avatar": ["kie-film", "kie-video-explainer"],
        "kie-identity": ["kie-identity"],
    }[primary]
    if "upscale" in model_id or "remove-background" in model_id:
        return ["kie-generate"]
    return fits


def proof_metadata(primary: str, schema: dict) -> dict:
    properties = schema.get("properties", {}) if isinstance(schema, dict) else {}
    fields = [field for field in RESOLUTION_FIELDS if field in properties]
    values: list[str] = []
    for field in fields:
        raw = properties.get(field, {})
        enum = raw.get("enum", []) if isinstance(raw, dict) else []
        values.extend(str(value) for value in enum if isinstance(value, (str, int, float)))
    return {
        "resolution_fields": fields,
        "lowest_faithful_tier": lowest_tier(values),
        "alternate_allowed": primary == "kie-video",
    }


def lowest_tier(values: list[str]) -> str:
    if not values:
        return ""
    numeric: list[tuple[int, str]] = []
    for value in values:
        match = re.fullmatch(r"(\d+)[pP]", value.strip())
        if match:
            numeric.append((int(match.group(1)), value))
    if numeric:
        return min(numeric)[1]
    order = {"low": 0, "draft": 0, "std": 1, "standard": 1, "medium": 2, "high": 3, "pro": 4}
    return min(values, key=lambda value: (order.get(value.lower(), 100), value.lower()))


def parse_operations(openapi_text: str, coverage: dict) -> list[dict]:
    shared = coverage.get("shared_operations", {})
    operations: list[dict] = []
    path = ""
    method = ""
    for line in openapi_text.splitlines():
        path_match = re.match(r"^  (/.+):\s*$", line)
        if path_match:
            path = path_match.group(1)
            method = ""
            continue
        method_match = re.match(r"^    (get|post|put|patch|delete|options|head):\s*$", line)
        if method_match:
            method = method_match.group(1).upper()
            continue
        id_match = re.match(r"^      operationId:\s*(.+?)\s*$", line)
        if not (id_match and path and method):
            continue
        operation_id = id_match.group(1).strip("'\"")
        key = f"{method} {path}"
        shared_entry = shared.get(key, {})
        variant_count = int(shared_entry.get("documentation_pages", 1))
        creative = operation_is_creative(operation_id, path, method)
        operations.append({
            "operation_id": operation_id,
            "method": method,
            "path": path,
            "variant_count": variant_count,
            "primary_capability": operation_capability(operation_id, path) if creative else "plumbing",
            "creative": creative,
            "plumbing": not creative,
            "reason": "creates or transforms media" if creative else "account, upload, status, download, or general language-model plumbing",
        })
    return operations


def operation_is_creative(operation_id: str, path: str, method: str) -> bool:
    if method == "GET" or any(token in operation_id for token in ("details", "record-info", "status", "credits", "download-url", "validate-info")):
        return False
    if any(prefix in path for prefix in ("/claude/", "/codex/", "/gemini", "/grok/v1/responses")):
        return False
    if any(token in operation_id for token in ("responses", "chat-completions", "messages", "stream-generate-content")):
        return False
    return not any(token in operation_id for token in ("upload-file", "market-query", "check-voice", "validate"))


def operation_capability(operation_id: str, path: str) -> str:
    value = f"{operation_id} {path}".lower()
    if any(token in value for token in ("avatar", "character")):
        return "kie-avatar"
    if "video" in value or "aleph" in value or "veo" in value or "runway" in value or "mp4" in value:
        return "kie-video"
    if any(token in value for token in ("audio", "music", "lyrics", "midi", "persona", "vocal", "voice", "wav", "sounds", "instrumental", "mashup", "replace-section", "/suno/")):
        return "kie-audio"
    if "image" in value or "background" in value:
        return "kie-image"
    if operation_id == "market-create-task":
        return "creative-router"
    raise RuntimeError(f"unclassified creative operation {operation_id} {path}")


def main() -> None:
    catalog = json.loads(CATALOG_PATH.read_text())
    coverage = json.loads(COVERAGE_PATH.read_text())
    models = []
    for model in catalog["models"]:
        properties = model.get("input_schema", {}).get("properties", {})
        fields = set(properties)
        primary = primary_capability(model["id"], fields)
        if primary not in CAPABILITIES:
            raise RuntimeError(f"unclassified model {model['id']}")
        models.append({
            "model_id": model["id"],
            "primary_capability": primary,
            "secondary_capabilities": secondary_capabilities(primary, model["id"], fields),
            "production_fit": production_fit(primary, model["id"]),
            "creative": True,
            "proof": proof_metadata(primary, model.get("input_schema", {})),
            "routing_notes": f"catalog-derived {primary} route; inspect the model schema before generation",
        })
    operations = parse_operations(OPENAPI_PATH.read_text(), coverage)
    document = {
        "schema_version": 1,
        "catalog_source_sha256": sha256(CATALOG_PATH),
        "coverage_source_sha256": sha256(COVERAGE_PATH),
        "catalog_model_count": catalog["model_count"],
        "documented_operation_variants": coverage["documented_operation_variants"],
        "models": models,
        "operations": operations,
    }
    if len(models) != catalog["model_count"]:
        raise RuntimeError("capability model count drift")
    if sum(item["variant_count"] for item in operations) != coverage["documented_operation_variants"]:
        raise RuntimeError("operation variant coverage drift")
    OUTPUT_PATH.write_text(json.dumps(document, indent=2, sort_keys=True) + "\n")


if __name__ == "__main__":
    main()
