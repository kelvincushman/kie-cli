#!/usr/bin/env python3
"""Regression tests for the reproducible Kie documentation merger."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

import build_spec


class BuildSpecTests(unittest.TestCase):
    def test_discovery_keeps_all_english_markdown_and_excludes_cn(self) -> None:
        index = """
        - [English](https://docs.kie.ai/market/vendor/model.md)
        - [Numeric](https://docs.kie.ai/123abc.md)
        - [Chinese](https://docs.kie.ai/cn/market/vendor/model.md)
        - [Chinese market](https://docs.kie.ai/cnmarket/vendor/model.md)
        - [HTML](https://docs.kie.ai/market/quickstart)
        """
        self.assertEqual(
            build_spec.discover_doc_urls(index),
            ["https://docs.kie.ai/123abc.md", "https://docs.kie.ai/market/vendor/model.md"],
        )

    def test_model_extraction_prefers_request_example_over_broken_enum(self) -> None:
        operation = {
            "requestBody": {
                "content": {
                    "application/json": {
                        "example": {"model": "qwen2/text-to-image", "input": {"prompt": "hello"}},
                        "schema": {
                            "type": "object",
                            "properties": {"model": {"type": "string", "enum": ["qwen2/image-edit"]}},
                        },
                    }
                }
            }
        }
        self.assertEqual(build_spec.extract_model_id(operation), ("qwen2/text-to-image", "request_example"))

    def test_shared_operation_unions_model_variants(self) -> None:
        def variant(model: str) -> dict:
            return {
                "url": f"https://docs.kie.ai/{model}.md",
                "operation": {
                    "operationId": model,
                    "summary": model,
                    "requestBody": {
                        "content": {
                            "application/json": {
                                "example": {"model": model},
                                "schema": {
                                    "type": "object",
                                    "required": ["model"],
                                    "properties": {"model": {"type": "string", "enum": [model]}},
                                },
                            }
                        }
                    },
                },
            }

        merged = build_spec.merge_operation_variants(
            "/codex/v1/responses", "post", [variant("gpt-5-5"), variant("gpt-5-6-terra")]
        )
        model_schema = merged["requestBody"]["content"]["application/json"]["schema"]["properties"]["model"]
        self.assertEqual(model_schema["enum"], ["gpt-5-5", "gpt-5-6-terra"])
        self.assertEqual(merged["operationId"], "openai-responses")

    def test_shared_operation_prefers_english_request_descriptions(self) -> None:
        def variant(model: str, description: str, field_description: str) -> dict:
            return {
                "url": f"https://docs.kie.ai/market/chat/{model}.md",
                "operation": {
                    "summary": description,
                    "description": description,
                    "requestBody": {
                        "content": {
                            "application/json": {
                                "example": {"model": model, "input": "hello"},
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "model": {"type": "string"},
                                        "input": {"type": "string", "description": field_description},
                                    },
                                },
                            }
                        }
                    },
                },
            }

        merged = build_spec.merge_operation_variants(
            "/codex/v1/responses",
            "post",
            [
                variant("gpt-cn", "这是一个中文端点说明", "中文输入说明"),
                variant("gpt-en", "English endpoint description", "English input description"),
            ],
        )
        schema = merged["requestBody"]["content"]["application/json"]["schema"]
        self.assertEqual(schema["properties"]["input"]["description"], "English input description")

    def test_duplicate_model_contract_prefers_descriptive_market_page(self) -> None:
        numeric = {"id": "demo/model", "source": "https://docs.kie.ai/123abc.md", "input_schema": {}}
        market = {
            "id": "demo/model",
            "source": "https://docs.kie.ai/market/demo/model.md",
            "input_schema": {"properties": {"prompt": {"type": "string"}}},
        }
        result = build_spec.preferred_model_entry(numeric, market)
        self.assertEqual(result["source"], market["source"])
        self.assertEqual(len(result["source_pages"]), 2)

    def test_default_outputs_are_repo_relative_not_cwd_relative(self) -> None:
        expected_root = Path(build_spec.__file__).resolve().parents[1]
        for path in build_spec.DEFAULT_OUTPUTS.values():
            self.assertTrue(path.is_absolute())
            self.assertTrue(path.is_relative_to(expected_root))

    def test_check_detects_stale_output(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "artifact"
            path.write_text("old", encoding="utf-8")
            with self.assertRaises(build_spec.BuildError):
                build_spec.write_or_check({"spec": "new"}, {"spec": path}, True)

    def test_schema_property_whitespace_is_normalized(self) -> None:
        schema = {
            "type": "object",
            "required": [" image_urls "],
            "properties": {" image_urls ": {"type": "array"}},
        }
        cleaned = build_spec.clean_spec_tree(schema)
        self.assertEqual(cleaned["required"], ["image_urls"])
        self.assertEqual(list(cleaned["properties"]), ["image_urls"])


if __name__ == "__main__":
    unittest.main()
