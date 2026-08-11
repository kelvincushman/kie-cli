#!/usr/bin/env python3
"""Focused consistency checks for the Market catalog builder."""
import unittest

import build_spec


class CatalogConsistencyTest(unittest.TestCase):
    def test_uses_documented_contract_when_openapi_enum_is_stale(self):
        page = {
            "title": "Kling - V2.5 Turbo Image to Video Pro",
            "url": "https://docs.kie.ai/market/kling/v25-turbo-image-to-video-pro.md",
        }
        text = 'model: kling/v2-5-turbo-image-to-video-pro\n'
        create = {
            "requestBody": {
                "content": {
                    "application/json": {
                        "schema": {
                            "properties": {
                                "model": {"enum": ["kling/v2-1-master-image-to-video"]}
                            }
                        }
                    }
                }
            }
        }

        self.assertEqual(
            build_spec.catalog_model_for_page(page, text, create),
            "kling/v2-5-turbo-image-to-video-pro",
        )

    def test_rejects_title_contract_mismatch(self):
        page = {"title": "Qwen3 Text to Image", "url": "https://docs.kie.ai/market/qwen3/text-to-image.md"}
        create = {"requestBody": {"content": {"application/json": {"schema": {"properties": {}}}}}}

        with self.assertRaisesRegex(ValueError, "conflicts"):
            build_spec.catalog_model_for_page(page, "model: qwen3/image-to-image\n", create)

    def test_rejects_version_conflict(self):
        page = {
            "title": "Kling - V2.5 Turbo Image to Video Pro",
            "url": "https://docs.kie.ai/market/kling/v25-turbo-image-to-video-pro.md",
        }
        create = {"requestBody": {"content": {"application/json": {"schema": {"properties": {}}}}}}

        with self.assertRaisesRegex(ValueError, "conflicts"):
            build_spec.catalog_model_for_page(page, "model: kling/v2-1-image-to-video-pro\n", create)

    def test_rejects_duplicate_model_ids(self):
        catalog = [
            {"model": "qwen3/text-to-image", "source": "https://docs.kie.ai/market/qwen3/text-to-image.md"},
            {"model": "qwen3/text-to-image", "source": "https://docs.kie.ai/market/qwen3-pro/text-to-image.md"},
        ]

        with self.assertRaisesRegex(ValueError, "duplicate Market model IDs"):
            build_spec.require_unique_model_ids(catalog)


if __name__ == "__main__":
    unittest.main()
