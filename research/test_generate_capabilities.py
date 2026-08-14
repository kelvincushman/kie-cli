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


if __name__ == "__main__":
    unittest.main()
