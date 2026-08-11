#!/usr/bin/env python3
"""Rebuild the complete Kie.ai OpenAPI and per-model contract catalog.

The Kie documentation publishes one OpenAPI fragment per documentation page.
Many Market models share ``POST /api/v1/jobs/createTask`` and several chat
models share the same provider-compatible endpoint. A simple path merge loses
all but the last variant, so this builder deliberately preserves two views:

* a valid merged OpenAPI document for CLI Printing Press; and
* a model registry containing every documented request/input schema, setting,
  constraint, default, example, and source page.

Run from any directory:

    python3 research/build_spec.py
    python3 research/build_spec.py --check
"""

from __future__ import annotations

import argparse
import concurrent.futures
import copy
import hashlib
import json
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from collections import defaultdict
from pathlib import Path
from typing import Any, Iterable

import yaml


DOCS_BASE = "https://docs.kie.ai"
DOCS_INDEX = f"{DOCS_BASE}/llms.txt"
HTTP_METHODS = {"get", "put", "post", "delete", "options", "head", "patch", "trace"}
OPENAPI_BLOCK_RE = re.compile(
    r"##\s+OpenAPI Specification\s*\n+```ya?ml\s*\n(.*?)\n```",
    re.IGNORECASE | re.DOTALL,
)
MARKDOWN_URL_RE = re.compile(r"\((https://docs\.kie\.ai/[^)\s]+\.md(?:\?[^)\s]*)?)\)")
MODEL_IN_TEXT_RE = re.compile(
    r"(?:\"model\"\s*:\s*\"|model(?:\s+parameter)?\s+(?:must be|is|to)\s+[`'\"])([A-Za-z0-9_.:/-]+)",
    re.IGNORECASE,
)
REPO_ROOT = Path(__file__).resolve().parents[1]

DEFAULT_OUTPUTS = {
    "spec": REPO_ROOT / "research" / "kie-final-openapi.yaml",
    "models": REPO_ROOT / "docs" / "MODELS.md",
    "inputs": REPO_ROOT / "docs" / "MODEL_INPUTS.md",
    "catalog": REPO_ROOT / "internal" / "kiecatalog" / "catalog.json",
    "coverage": REPO_ROOT / "research" / "kie-api-coverage.json",
    "metadata": REPO_ROOT / ".printing-press.json",
}

STABLE_OPERATION_IDS = {
    ("/api/v1/jobs/createTask", "post"): "market-create-task",
    ("/api/v1/jobs/recordInfo", "get"): "market-query-task",
    ("/codex/v1/responses", "post"): "openai-responses",
    ("/claude/v1/messages", "post"): "claude-messages",
    ("/grok/v1/responses", "post"): "grok-responses",
}

TAG_MAP = {
    "market": "Market",
    "common api": "Common",
    "file upload": "FileUpload",
    "claude": "Chat",
    "codex": "Chat",
    "4o image": "ImageGen4o",
    "flux kontext": "FluxKontext",
    "lyrics generation": "SunoLyrics",
    "music generation": "SunoMusic",
    "music video generation": "SunoMusicVideo",
    "sounds generation": "SunoSounds",
    "vocal removal": "SunoVocal",
    "wav conversion": "SunoWav",
    "suno api/voice": "SunoVoice",
    "veo3": "Veo3",
    "gemini omni": "GeminiOmni",
    "runway api/aleph": "RunwayAleph",
    "runway api": "Runway",
}


class BuildError(RuntimeError):
    """Raised when the live documentation cannot be represented safely."""


def parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--docs-index", default=DOCS_INDEX)
    parser.add_argument("--spec-out", type=Path, default=DEFAULT_OUTPUTS["spec"])
    parser.add_argument("--models-out", type=Path, default=DEFAULT_OUTPUTS["models"])
    parser.add_argument("--model-inputs-out", type=Path, default=DEFAULT_OUTPUTS["inputs"])
    parser.add_argument("--catalog-out", type=Path, default=DEFAULT_OUTPUTS["catalog"])
    parser.add_argument("--coverage-out", type=Path, default=DEFAULT_OUTPUTS["coverage"])
    parser.add_argument("--metadata-out", type=Path, default=DEFAULT_OUTPUTS["metadata"])
    parser.add_argument("--workers", type=int, default=16)
    parser.add_argument("--timeout", type=float, default=30.0)
    parser.add_argument("--retries", type=int, default=3)
    parser.add_argument("--check", action="store_true", help="fail if committed outputs differ from live docs")
    return parser.parse_args(argv)


def fetch(url: str, timeout: float = 30.0, retries: int = 3) -> str:
    request = urllib.request.Request(url, headers={"User-Agent": "kie-cli-spec-builder/2.0"})
    last_error: Exception | None = None
    for attempt in range(retries + 1):
        try:
            with urllib.request.urlopen(request, timeout=timeout) as response:
                return response.read().decode("utf-8", errors="replace")
        except (urllib.error.URLError, TimeoutError, OSError) as exc:
            last_error = exc
            if attempt < retries:
                time.sleep(min(0.4 * (2**attempt), 3.0))
    raise BuildError(f"failed to fetch {url} after {retries + 1} attempts: {last_error}")


def discover_doc_urls(index_text: str) -> list[str]:
    urls: set[str] = set()
    for raw_url in MARKDOWN_URL_RE.findall(index_text):
        parsed = urllib.parse.urlsplit(raw_url)
        path = parsed.path.rstrip("/")
        if path == "/cn.md" or path.startswith("/cn/") or path.startswith("/cnmarket/"):
            continue
        urls.add(urllib.parse.urlunsplit((parsed.scheme, parsed.netloc, parsed.path, "", "")))
    return sorted(urls)


def fetch_docs(urls: Iterable[str], workers: int, timeout: float, retries: int) -> dict[str, str]:
    url_list = list(urls)
    failures: list[str] = []
    pages: dict[str, str] = {}
    with concurrent.futures.ThreadPoolExecutor(max_workers=max(1, workers)) as executor:
        futures = {executor.submit(fetch, url, timeout, retries): url for url in url_list}
        for future in concurrent.futures.as_completed(futures):
            url = futures[future]
            try:
                pages[url] = future.result()
            except Exception as exc:  # report all failures in one actionable error
                failures.append(f"{url}: {exc}")
    if failures:
        raise BuildError("documentation fetch failed:\n  " + "\n  ".join(sorted(failures)))
    return dict(sorted(pages.items()))


def load_openapi_blocks(text: str, url: str = "<memory>") -> list[dict[str, Any]]:
    raw_blocks = OPENAPI_BLOCK_RE.findall(text)
    documents: list[dict[str, Any]] = []
    for number, raw in enumerate(raw_blocks, start=1):
        try:
            document = yaml.safe_load(raw)
        except yaml.YAMLError as exc:
            raise BuildError(f"invalid OpenAPI YAML in {url} block {number}: {exc}") from exc
        if not isinstance(document, dict) or not isinstance(document.get("paths"), dict):
            raise BuildError(f"OpenAPI block {number} in {url} has no paths object")
        documents.append(document)
    if "OpenAPI Specification" in text and not raw_blocks:
        raise BuildError(f"found an OpenAPI heading but could not parse its YAML fence in {url}")
    return documents


def json_pointer(document: Any, pointer: str) -> Any:
    node = document
    for token in pointer.removeprefix("#/").split("/"):
        token = urllib.parse.unquote(token).replace("~1", "/").replace("~0", "~")
        node = node[token]
    return node


def resolve_local_refs(value: Any, document: dict[str, Any], stack: tuple[str, ...] = ()) -> Any:
    if isinstance(value, list):
        return [resolve_local_refs(item, document, stack) for item in value]
    if not isinstance(value, dict):
        return copy.deepcopy(value)
    ref = value.get("$ref")
    if isinstance(ref, str) and ref.startswith("#/"):
        if ref in stack:
            return copy.deepcopy(value)
        try:
            target = json_pointer(document, ref)
        except (KeyError, TypeError):
            raise BuildError(f"unresolved local OpenAPI reference {ref}")
        merged = resolve_local_refs(target, document, stack + (ref,))
        if not isinstance(merged, dict):
            return merged
        for key, child in value.items():
            if key != "$ref":
                merged[key] = resolve_local_refs(child, document, stack)
        return merged
    return {key: resolve_local_refs(child, document, stack) for key, child in value.items()}


def clean_desc(text: Any) -> Any:
    if not isinstance(text, str):
        return text
    output: list[str] = []
    for line in text.splitlines():
        stripped = line.strip()
        if re.match(r"^:::\s*\w*(\[\])?\s*$", stripped):
            continue
        match = re.match(r"^:::\s*\w+\[\]\s*(.*)$", stripped)
        if match and match.group(1):
            output.append(match.group(1))
            continue
        output.append(stripped[2:] if stripped.startswith("> ") else line)
    return re.sub(r"\n{3,}", "\n\n", "\n".join(output).strip())


def clean_spec_tree(value: Any) -> Any:
    if isinstance(value, dict):
        for key in list(value):
            if key.startswith("x-apidog") or key.startswith("x-run-in-apidog"):
                del value[key]
        # Apidog exports per-operation security fragments inconsistently. Kie
        # documents Bearer auth globally, so retaining an empty local override
        # would incorrectly make a generated CLI command unauthenticated.
        if "responses" in value and "security" in value:
            del value["security"]
        for field in ("description", "summary"):
            if field in value:
                value[field] = clean_desc(value[field])
        properties = value.get("properties")
        if isinstance(properties, dict):
            normalized: dict[str, Any] = {}
            for raw_name, child in properties.items():
                name = raw_name.strip() if isinstance(raw_name, str) else raw_name
                if name in normalized and normalized[name] != child:
                    raise BuildError(f"schema property names collide after whitespace normalization: {raw_name!r}")
                normalized[name] = child
            value["properties"] = normalized
        required = value.get("required")
        if isinstance(required, list):
            value["required"] = [item.strip() if isinstance(item, str) else item for item in required]
        for child in value.values():
            clean_spec_tree(child)
    elif isinstance(value, list):
        for child in value:
            clean_spec_tree(child)
    return value


def normalize_tag(tag: str) -> str:
    lowered = tag.lower()
    if "gpt" in lowered and "chat" in lowered:
        return "Chat"
    if "gemini" in lowered and "chat" in lowered:
        return "Chat"
    if "grok" in lowered and "chat" in lowered:
        return "Chat"
    for key, value in TAG_MAP.items():
        if key in lowered:
            return value
    return tag or "Kie"


def media_json_content(operation: dict[str, Any]) -> dict[str, Any]:
    content = operation.get("requestBody", {}).get("content", {})
    if not isinstance(content, dict):
        return {}
    return content.get("application/json") or content.get("application/*+json") or {}


def schema_for_operation(operation: dict[str, Any]) -> dict[str, Any]:
    schema = media_json_content(operation).get("schema", {})
    return copy.deepcopy(schema) if isinstance(schema, dict) else {}


def examples_for_operation(operation: dict[str, Any]) -> list[Any]:
    content = media_json_content(operation)
    values: list[Any] = []
    if "example" in content:
        values.append(content["example"])
    examples = content.get("examples", {})
    if isinstance(examples, dict):
        for example in examples.values():
            values.append(example.get("value") if isinstance(example, dict) and "value" in example else example)
    schema = content.get("schema", {})
    if isinstance(schema, dict):
        if "example" in schema:
            values.append(schema["example"])
        if isinstance(schema.get("examples"), list):
            values.extend(schema["examples"])
    return values


def model_from_value(value: Any) -> str:
    if isinstance(value, dict):
        candidate = value.get("model")
        if isinstance(candidate, str) and candidate.strip():
            return candidate.strip()
        for child in value.values():
            nested = model_from_value(child)
            if nested:
                return nested
    elif isinstance(value, list):
        for child in value:
            nested = model_from_value(child)
            if nested:
                return nested
    elif isinstance(value, str):
        match = MODEL_IN_TEXT_RE.search(value)
        if match:
            return match.group(1)
        try:
            parsed = json.loads(value)
        except (json.JSONDecodeError, TypeError):
            return ""
        return model_from_value(parsed)
    return ""


def extract_model_id(operation: dict[str, Any]) -> tuple[str, str]:
    """Return (model ID, provenance), preferring executable examples.

    A small number of Kie pages currently have stale enum/default values while
    their request example and prose show the correct model. The source marker is
    retained in the coverage report so this correction remains auditable.
    """

    for value in examples_for_operation(operation):
        candidate = model_from_value(value)
        if candidate:
            return candidate, "request_example"
    schema = schema_for_operation(operation)
    properties = schema.get("properties", {}) if isinstance(schema, dict) else {}
    model_schema = properties.get("model", {}) if isinstance(properties, dict) else {}
    if isinstance(model_schema, dict):
        for field in ("example", "examples"):
            candidate = model_from_value({"model": model_schema.get(field)})
            if candidate:
                return candidate, f"model_{field}"
        description = model_schema.get("description", "")
        if isinstance(description, str):
            match = re.search(r"(?:must be|use|value)\s+[`'\"]([^`'\"]+)[`'\"]", description, re.I)
            if match:
                return match.group(1).strip(), "model_description"
        enum = model_schema.get("enum", [])
        if isinstance(enum, list) and enum and isinstance(enum[0], str):
            return enum[0].strip(), "model_enum"
        default = model_schema.get("default")
        if isinstance(default, str) and default.strip():
            return default.strip(), "model_default"
    for field in (operation.get("description", ""), operation.get("summary", "")):
        candidate = model_from_value(field)
        if candidate:
            return candidate, "operation_text"
    return "", ""


def model_values(operation: dict[str, Any]) -> list[str]:
    values: list[str] = []
    primary, provenance = extract_model_id(operation)
    if primary:
        values.append(primary)
    # Per-model pages sometimes retain a copied enum/default from the page
    # used as their template. An executable request example is the most
    # reliable declaration for that page and must not be polluted by those
    # stale values (for example Grok pages that still say gpt-5.1-codex).
    if primary and provenance == "request_example":
        return values
    schema = schema_for_operation(operation)
    model_schema = schema.get("properties", {}).get("model", {}) if isinstance(schema, dict) else {}
    if isinstance(model_schema, dict):
        for field in ("enum", "examples"):
            raw = model_schema.get(field, [])
            if isinstance(raw, list):
                values.extend(item for item in raw if isinstance(item, str) and item.strip())
        for field in ("default", "example"):
            raw = model_schema.get(field)
            if isinstance(raw, str) and raw.strip():
                values.append(raw)
    for example in examples_for_operation(operation):
        value = model_from_value(example)
        if value:
            values.append(value)
    return sorted(set(values))


def source_category(url: str, operation: dict[str, Any]) -> str:
    tags = operation.get("tags") or []
    if tags and isinstance(tags[0], str):
        raw = tags[0]
        parts = [p.strip() for p in raw.split("/") if p.strip().lower() not in {"docs", "en", "market"}]
        if parts:
            return " > ".join(parts)
    path = urllib.parse.urlsplit(url).path
    parts = [part for part in path.removesuffix(".md").split("/") if part and part != "market"]
    return " > ".join(parts[:-1]).replace("-", " ").title() or "Market"


def preferred_model_entry(left: dict[str, Any], right: dict[str, Any]) -> dict[str, Any]:
    def score(entry: dict[str, Any]) -> tuple[int, int, int]:
        path = urllib.parse.urlsplit(entry["source"]).path
        descriptive = not re.search(r"/\d+[a-z0-9]*\.md$", path)
        return (int(path.startswith("/market/")), int(descriptive), len(json.dumps(entry.get("input_schema", {}))))

    winner, other = (right, left) if score(right) > score(left) else (left, right)
    pages = sorted(set(winner.get("source_pages", [winner["source"]]) + other.get("source_pages", [other["source"]])))
    winner = copy.deepcopy(winner)
    winner["source_pages"] = pages
    return winner


def market_entry(url: str, operation: dict[str, Any]) -> dict[str, Any]:
    model_id, model_source = extract_model_id(operation)
    if not model_id:
        raise BuildError(f"could not determine Market model identifier from {url}")
    request_schema = schema_for_operation(operation)
    properties = request_schema.get("properties", {}) if isinstance(request_schema, dict) else {}
    input_schema = copy.deepcopy(properties.get("input", {})) if isinstance(properties, dict) else {}
    request_examples = examples_for_operation(operation)
    request_example = request_examples[0] if request_examples else {}
    input_example = request_example.get("input", {}) if isinstance(request_example, dict) else {}
    if not input_example and isinstance(input_schema, dict):
        raw = input_schema.get("example")
        input_example = raw if isinstance(raw, dict) else {}
    return {
        "id": model_id,
        "name": clean_desc(operation.get("summary", "")) or model_id,
        "category": source_category(url, operation),
        "description": clean_desc(operation.get("description", "")),
        "source": url.removesuffix(".md"),
        "source_pages": [url.removesuffix(".md")],
        "model_id_source": model_source,
        "request_schema": request_schema,
        "input_schema": input_schema,
        "request_example": request_example,
        "input_example": input_example,
        "required_input_fields": sorted(input_schema.get("required", [])) if isinstance(input_schema, dict) else [],
    }


def merge_schema(left: Any, right: Any) -> Any:
    if not isinstance(left, dict) or not isinstance(right, dict):
        return copy.deepcopy(left if left not in ({}, None, "") else right)
    merged = copy.deepcopy(left)
    for key, value in right.items():
        if key not in merged:
            merged[key] = copy.deepcopy(value)
        elif key == "properties" and isinstance(value, dict):
            merged[key] = merge_schema(merged[key], value)
        elif key in {"enum", "examples"} and isinstance(value, list) and isinstance(merged[key], list):
            combined = merged[key] + value
            seen: set[str] = set()
            merged[key] = []
            for item in combined:
                marker = json.dumps(item, sort_keys=True, ensure_ascii=False)
                if marker not in seen:
                    seen.add(marker)
                    merged[key].append(copy.deepcopy(item))
        elif key in {"oneOf", "anyOf", "allOf"} and isinstance(value, list) and isinstance(merged[key], list):
            seen = {json.dumps(item, sort_keys=True) for item in merged[key]}
            merged[key].extend(copy.deepcopy(item) for item in value if json.dumps(item, sort_keys=True) not in seen)
    return merged


def merge_operation_variants(path: str, method: str, variants: list[dict[str, Any]]) -> dict[str, Any]:
    def rank(item: dict[str, Any]) -> tuple[int, float, int]:
        url_path = urllib.parse.urlsplit(item["url"]).path
        text = str(item["operation"].get("description", "")) + str(item["operation"].get("summary", ""))
        ascii_ratio = sum(ord(char) < 128 for char in text) / max(1, len(text))
        return (int(url_path.startswith("/market/")), ascii_ratio, len(json.dumps(item["operation"], sort_keys=True)))

    ranked = sorted(variants, key=rank, reverse=True)
    merged = copy.deepcopy(ranked[0]["operation"])
    schemas = [schema_for_operation(item["operation"]) for item in ranked]
    object_schemas = [schema for schema in schemas if isinstance(schema, dict) and schema]
    if object_schemas:
        request_schema = copy.deepcopy(object_schemas[0])
        for schema in object_schemas[1:]:
            request_schema = merge_schema(request_schema, schema)
        required_sets = [set(schema.get("required", [])) for schema in object_schemas]
        if required_sets:
            common_required = sorted(set.intersection(*required_sets))
            if common_required:
                request_schema["required"] = common_required
            else:
                request_schema.pop("required", None)
        models = sorted({value for item in variants for value in model_values(item["operation"])})
        model_schema = request_schema.get("properties", {}).get("model")
        if models and isinstance(model_schema, dict):
            model_schema["enum"] = models
            model_schema["examples"] = models
            if len(models) > 1:
                model_schema.pop("default", None)
                model_schema["description"] = "Model identifier. Supported values: " + ", ".join(models)
        content = media_json_content(merged)
        if content:
            content["schema"] = request_schema
    stable_id = STABLE_OPERATION_IDS.get((path, method))
    if stable_id:
        merged["operationId"] = stable_id
    models = sorted({value for item in variants for value in model_values(item["operation"])})
    if len(variants) > 1:
        merged["x-kie-variants"] = [
            {
                "source": item["url"].removesuffix(".md"),
                "summary": item["operation"].get("summary", ""),
                "models": model_values(item["operation"]),
            }
            for item in sorted(variants, key=lambda item: item["url"])
        ]
        if models:
            merged["summary"] = f"Create a request using one of {len(models)} documented model variants"
            merged["description"] = "Shared Kie endpoint. Select a documented model value: " + ", ".join(models) + "."
    return merged


def build_market_create(base: dict[str, Any], models: list[dict[str, Any]]) -> dict[str, Any]:
    operation = copy.deepcopy(base)
    operation["summary"] = "Create a Market model generation task"
    operation["operationId"] = "market-create-task"
    operation["tags"] = ["Market"]
    operation["description"] = (
        f"Unified entry point for all {len(models)} currently documented Kie.ai Market models. "
        "Pass a model identifier and its model-specific input payload. Inspect the exact schema with "
        "`kie-pp-cli models show <model>` and validate locally with `kie-pp-cli models validate <model>`."
    )
    schema = schema_for_operation(operation)
    schema["type"] = "object"
    properties = schema.setdefault("properties", {})
    properties["callBackUrl"] = {
        "type": "string",
        "format": "uri",
        "description": "Optional URL that receives the generation task completion callback.",
    }
    properties["model"] = {
        "type": "string",
        "description": "Kie Market model identifier. Use `kie-pp-cli models list` for the current catalog.",
        "enum": [model["id"] for model in models],
        "examples": [model["id"] for model in models[:5]],
    }
    properties["input"] = {
        "type": "object",
        "description": "Model-specific input. The complete per-model JSON Schemas are embedded in the CLI model registry.",
        "additionalProperties": True,
    }
    schema["required"] = sorted(set(schema.get("required", [])) | {"model", "input"})
    media_json_content(operation)["schema"] = schema
    operation["x-kie-model-catalog"] = {
        "count": len(models),
        "embedded_registry": "internal/kiecatalog/catalog.json",
        "documentation": "docs/MODEL_INPUTS.md",
    }
    return operation


def operation_id_fallback(method: str, path: str) -> str:
    tokens = [token for token in re.split(r"[^A-Za-z0-9]+", path) if token and token not in {"api", "v1"}]
    return "-".join([method] + tokens).lower() or method


def unique_operation_ids(paths: dict[str, Any]) -> None:
    seen: dict[str, tuple[str, str]] = {}
    for path in sorted(paths):
        for method in sorted(paths[path]):
            operation = paths[path][method]
            operation_id = str(operation.get("operationId") or operation_id_fallback(method, path))
            base = re.sub(r"[^A-Za-z0-9_-]+", "-", operation_id).strip("-") or operation_id_fallback(method, path)
            candidate = base
            number = 2
            while candidate in seen and seen[candidate] != (path, method):
                candidate = f"{base}-{number}"
                number += 1
            operation["operationId"] = candidate
            seen[candidate] = (path, method)


def collect(index_url: str, index_text: str, pages: dict[str, str]) -> tuple[dict[str, Any], dict[str, Any], dict[str, Any]]:
    operation_variants: dict[tuple[str, str], list[dict[str, Any]]] = defaultdict(list)
    market_by_id: dict[str, dict[str, Any]] = {}
    openapi_pages = 0
    openapi_blocks = 0
    for url, text in pages.items():
        documents = load_openapi_blocks(text, url)
        if documents:
            openapi_pages += 1
        openapi_blocks += len(documents)
        for document in documents:
            for path, path_item in document["paths"].items():
                if not isinstance(path_item, dict):
                    continue
                for method, raw_operation in path_item.items():
                    lowered = method.lower()
                    if lowered not in HTTP_METHODS or not isinstance(raw_operation, dict):
                        continue
                    operation = clean_spec_tree(resolve_local_refs(raw_operation, document))
                    tags = operation.get("tags") or ["Kie"]
                    operation["tags"] = [normalize_tag(str(tags[0]))]
                    variant = {"url": url, "operation": operation}
                    operation_variants[(path, lowered)].append(variant)
                    if path == "/api/v1/jobs/createTask" and lowered == "post":
                        entry = clean_spec_tree(market_entry(url, operation))
                        existing = market_by_id.get(entry["id"])
                        market_by_id[entry["id"]] = preferred_model_entry(existing, entry) if existing else entry

    models = sorted(market_by_id.values(), key=lambda item: item["id"])
    if not models:
        raise BuildError("no Kie Market createTask model contracts were discovered")

    merged_paths: dict[str, dict[str, Any]] = {}
    shared_variants: dict[str, Any] = {}
    for (path, method), variants in sorted(operation_variants.items()):
        if path == "/api/v1/jobs/createTask" and method == "post":
            operation = build_market_create(variants[0]["operation"], models)
        else:
            operation = merge_operation_variants(path, method, variants)
        merged_paths.setdefault(path, {})[method] = operation
        unique_variants = {
            json.dumps({"models": model_values(item["operation"]), "summary": item["operation"].get("summary", "")}, sort_keys=True)
            for item in variants
        }
        if len(variants) > 1:
            shared_variants[f"{method.upper()} {path}"] = {
                "documentation_pages": len(variants),
                "distinct_variants": len(unique_variants),
                "models": sorted({model for item in variants for model in model_values(item["operation"])}),
                "sources": sorted(item["url"].removesuffix(".md") for item in variants),
            }
    unique_operation_ids(merged_paths)

    index_hash = hashlib.sha256(index_text.encode()).hexdigest()
    spec = {
        "openapi": "3.0.1",
        "info": {
            "title": "Kie.ai API",
            "description": (
                "Complete CLI-oriented merge of the OpenAPI fragments indexed by Kie.ai. Shared endpoints retain "
                "all documented model variants; every Market input contract is embedded in the CLI model registry."
            ),
            "version": "1.0.0",
            "x-kie-docs-index-sha256": index_hash,
        },
        "servers": [{"url": "https://api.kie.ai", "description": "Production"}],
        "security": [{"BearerAuth": []}],
        "paths": merged_paths,
        "components": {
            "securitySchemes": {
                "BearerAuth": {
                    "type": "http",
                    "scheme": "bearer",
                    "bearerFormat": "API Key",
                    "description": "Create a key at https://kie.ai/api-key and send it as `Authorization: Bearer YOUR_API_KEY`.",
                }
            }
        },
    }
    registry = {
        "schema_version": 1,
        "source_index": index_url,
        "source_index_sha256": index_hash,
        "model_count": len(models),
        "models": models,
    }
    corrections = [
        {
            "model": model["id"],
            "source": model["source"],
            "reason": "request example was preferred over a conflicting enum/default",
        }
        for model in models
        if model["model_id_source"] == "request_example"
        and _declared_model_value(model["request_schema"])
        and _declared_model_value(model["request_schema"]) != model["id"]
    ]
    coverage = {
        "schema_version": 1,
        "source_index": index_url,
        "source_index_sha256": index_hash,
        "indexed_english_pages": len(pages),
        "pages_with_openapi": openapi_pages,
        "openapi_blocks": openapi_blocks,
        "unique_paths": len(merged_paths),
        "unique_operations": sum(len(methods) for methods in merged_paths.values()),
        "documented_operation_variants": sum(len(items) for items in operation_variants.values()),
        "market_models": len(models),
        "shared_operations": shared_variants,
        "source_corrections": corrections,
        "failures": [],
    }
    return spec, registry, coverage


def _declared_model_value(request_schema: dict[str, Any]) -> str:
    model_schema = request_schema.get("properties", {}).get("model", {}) if isinstance(request_schema, dict) else {}
    if not isinstance(model_schema, dict):
        return ""
    enum = model_schema.get("enum", [])
    if isinstance(enum, list) and enum and isinstance(enum[0], str):
        return enum[0]
    return model_schema.get("default", "") if isinstance(model_schema.get("default"), str) else ""


def markdown_escape(value: Any) -> str:
    text = str(value if value is not None else "")
    return text.replace("|", "\\|").replace("\n", " ").strip()


def schema_type(schema: dict[str, Any]) -> str:
    raw = schema.get("type")
    if isinstance(raw, list):
        return " / ".join(str(item) for item in raw)
    if isinstance(raw, str):
        return raw
    if "oneOf" in schema:
        return "oneOf"
    if "anyOf" in schema:
        return "anyOf"
    return "object" if "properties" in schema else "any"


def constraint_text(schema: dict[str, Any]) -> str:
    values: list[str] = []
    if "enum" in schema:
        values.append("allowed: " + ", ".join(f"`{item}`" for item in schema["enum"]))
    if "const" in schema:
        values.append(f"constant: `{schema['const']}`")
    if "default" in schema:
        values.append(f"default: `{json.dumps(schema['default'], ensure_ascii=False)}`")
    if "minimum" in schema:
        values.append(f"min: {schema['minimum']}")
    if "maximum" in schema:
        values.append(f"max: {schema['maximum']}")
    if "minLength" in schema:
        values.append(f"min length: {schema['minLength']}")
    if "maxLength" in schema:
        values.append(f"max length: {schema['maxLength']}")
    if "minItems" in schema:
        values.append(f"min items: {schema['minItems']}")
    if "maxItems" in schema:
        values.append(f"max items: {schema['maxItems']}")
    if "format" in schema:
        values.append(f"format: {schema['format']}")
    return "; ".join(values)


def flatten_schema(schema: dict[str, Any], prefix: str = "", parent_required: bool = True) -> list[dict[str, Any]]:
    rows: list[dict[str, Any]] = []
    properties = schema.get("properties", {}) if isinstance(schema, dict) else {}
    required = set(schema.get("required", [])) if isinstance(schema, dict) else set()
    if isinstance(properties, dict):
        for name, child in properties.items():
            if not isinstance(child, dict):
                continue
            path = f"{prefix}.{name}" if prefix else name
            is_required = parent_required and name in required
            rows.append(
                {
                    "path": path,
                    "type": schema_type(child),
                    "required": is_required,
                    "constraints": constraint_text(child),
                    "description": clean_desc(child.get("description", "")),
                }
            )
            rows.extend(flatten_schema(child, path, is_required))
            items = child.get("items")
            if isinstance(items, dict) and isinstance(items.get("properties"), dict):
                rows.extend(flatten_schema(items, path + "[]", is_required))
    return rows


def render_models(registry: dict[str, Any]) -> str:
    models = registry["models"]
    categories: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for model in models:
        categories[model["category"]].append(model)
    lines = [
        "# Kie.ai Market Model Catalog",
        "",
        f"This reproducible snapshot contains **{len(models)} unique Market model contracts** discovered from Kie.ai's official documentation index.",
        "Every model shares the create/query task API, but its full input schema and settings are preserved separately.",
        "",
        "- List models: `kie-pp-cli models list`",
        "- Inspect exact settings: `kie-pp-cli models show <model-id>`",
        "- Create a starter payload: `kie-pp-cli models example <model-id>`",
        "- Validate before spending credits: `kie-pp-cli models validate <model-id> --input '{...}'`",
        "- Read all field tables: [MODEL_INPUTS.md](MODEL_INPUTS.md)",
        "",
        "Generation: `kie-pp-cli kie-ai-jobs market-create-task --model <model-id> --input '{...}'`",
        "",
    ]
    for category in sorted(categories):
        lines.extend([f"## {category}", "", "| Model | `--model` value | Input fields | Source |", "|---|---|---:|---|"])
        for model in sorted(categories[category], key=lambda item: item["id"]):
            fields = len(flatten_schema(model["input_schema"]))
            lines.append(
                f"| {markdown_escape(model['name'])} | `{model['id']}` | {fields} | [docs]({model['source']}) |"
            )
        lines.append("")
    return "\n".join(lines).rstrip() + "\n"


def render_model_inputs(registry: dict[str, Any]) -> str:
    lines = [
        "# Kie.ai Market Model Inputs and Settings",
        "",
        f"Complete field reference for **{registry['model_count']} models**, generated from [Kie.ai's official documentation index]({registry['source_index']}).",
        "The embedded machine-readable source is `internal/kiecatalog/catalog.json`; use `kie-pp-cli models show` or the MCP `media_model_get` tool to avoid loading this whole document into an agent context.",
        "",
        "Top-level task settings are `model`, `input`, and the optional `callBackUrl`. Tables below describe each model's `input` object.",
        "",
    ]
    for model in registry["models"]:
        lines.extend(
            [
                f"## `{model['id']}`",
                "",
                f"**{markdown_escape(model['name'])}** · {markdown_escape(model['category'])} · [official docs]({model['source']})",
                "",
            ]
        )
        rows = flatten_schema(model["input_schema"])
        if rows:
            lines.extend(["| Input field | Type | Required | Settings | Description |", "|---|---|---|---|---|"])
            for row in rows:
                lines.append(
                    f"| `{row['path']}` | {markdown_escape(row['type'])} | {'yes' if row['required'] else 'no'} | "
                    f"{markdown_escape(row['constraints'])} | {markdown_escape(row['description'])} |"
                )
        else:
            lines.append("No structured input fields were published on the source page.")
        if model.get("input_example"):
            lines.extend(["", "Example `input`:", "", "```json", json.dumps(model["input_example"], indent=2, ensure_ascii=False), "```"])
        lines.append("")
    return "\n".join(lines).rstrip() + "\n"


def serialized_outputs(spec: dict[str, Any], registry: dict[str, Any], coverage: dict[str, Any]) -> dict[str, str]:
    return {
        "spec": yaml.safe_dump(spec, sort_keys=False, allow_unicode=True, width=120),
        "models": render_models(registry),
        "inputs": render_model_inputs(registry),
        "catalog": json.dumps(registry, indent=2, ensure_ascii=False, sort_keys=True) + "\n",
        "coverage": json.dumps(coverage, indent=2, ensure_ascii=False, sort_keys=True) + "\n",
    }


def render_printing_press_metadata(path: Path, spec_text: str, registry: dict[str, Any], coverage: dict[str, Any]) -> str | None:
    if not path.exists():
        return None
    metadata = json.loads(path.read_text(encoding="utf-8"))
    description = (
        f"Complete Kie.ai CLI: {coverage['unique_operations']} current API operations and "
        f"{registry['model_count']} Market models with embedded per-model input/settings schemas, "
        "plus a local agent media factory."
    )
    metadata["spec_checksum"] = "sha256:" + hashlib.sha256(spec_text.encode()).hexdigest()
    metadata["description"] = description
    metadata["catalog_description"] = description
    for feature in metadata.get("novel_features", []):
        if feature.get("name") == "Unified Market Catalog Access":
            feature["description"] = (
                f"Single create/query command pair that reaches all {registry['model_count']} Kie.ai Market models, "
                "backed by local model discovery, examples, and formal input validation."
            )
    return json.dumps(metadata, indent=2, ensure_ascii=False) + "\n"


def write_or_check(outputs: dict[str, str], paths: dict[str, Path], check: bool) -> None:
    stale: list[str] = []
    for key, content in outputs.items():
        path = paths[key].resolve()
        if check:
            if not path.exists() or path.read_text(encoding="utf-8") != content:
                stale.append(str(path))
            continue
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")
        print(f"wrote {path}", file=sys.stderr)
    if stale:
        raise BuildError("generated Kie API artifacts are stale:\n  " + "\n  ".join(stale))


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)
    try:
        print(f"fetching documentation index {args.docs_index}", file=sys.stderr)
        index_text = fetch(args.docs_index, args.timeout, args.retries)
        urls = discover_doc_urls(index_text)
        if not urls:
            raise BuildError("documentation index contained no English Markdown pages")
        print(f"fetching {len(urls)} indexed English documentation pages", file=sys.stderr)
        pages = fetch_docs(urls, args.workers, args.timeout, args.retries)
        spec, registry, coverage = collect(args.docs_index, index_text, pages)
        outputs = serialized_outputs(spec, registry, coverage)
        metadata = render_printing_press_metadata(args.metadata_out, outputs["spec"], registry, coverage)
        if metadata is not None:
            outputs["metadata"] = metadata
        paths = {
            "spec": args.spec_out,
            "models": args.models_out,
            "inputs": args.model_inputs_out,
            "catalog": args.catalog_out,
            "coverage": args.coverage_out,
            "metadata": args.metadata_out,
        }
        write_or_check(outputs, paths, args.check)
        verb = "verified" if args.check else "captured"
        print(
            f"{verb} {coverage['unique_operations']} operations across {coverage['unique_paths']} paths and "
            f"{coverage['market_models']} complete Market model contracts",
            file=sys.stderr,
        )
        return 0
    except BuildError as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
