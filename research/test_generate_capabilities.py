import pathlib
import sys
import unittest


RESEARCH = pathlib.Path(__file__).resolve().parent
sys.path.insert(0, str(RESEARCH))

import generate_capabilities  # noqa: E402


class CapabilityProofMetadataTests(unittest.TestCase):
    def test_uses_first_field_with_a_rankable_tier(self):
        metadata = generate_capabilities.proof_metadata(
            "kie-video",
            {
                "properties": {
                    "quality": {"enum": ["high", "low"]},
                    "resolution": {"enum": ["balanced", "creative"]},
                }
            },
        )
        self.assertEqual(metadata["resolution_fields"], ["resolution", "quality"])
        self.assertEqual(metadata["lowest_faithful_tier"], "low")

    def test_unrankable_values_do_not_claim_a_lowest_tier(self):
        self.assertEqual(generate_capabilities.lowest_tier(["balanced", "creative"]), "")

    def test_numeric_resolution_without_p_matches_runtime_ranking(self):
        self.assertEqual(generate_capabilities.lowest_tier(["1080", "720"]), "720")

    def test_operations_accept_quoted_paths_and_valid_yaml_indentation(self):
        operations = generate_capabilities.parse_operations(
            """
paths:
    '/v1/images':
        post:
            operationId: image-create-task
""",
            {"shared_operations": {"POST /v1/images": {"documentation_pages": 3}}},
        )
        self.assertEqual(
            operations,
            [{
                "operation_id": "image-create-task",
                "method": "POST",
                "path": "/v1/images",
                "variant_count": 3,
                "primary_capability": "kie-image",
                "creative": True,
                "plumbing": False,
                "reason": "creates or transforms media",
            }],
        )


if __name__ == "__main__":
    unittest.main()
