#!/usr/bin/env bash
# Weekly refresh of the standard Market task-model catalog (docs/MODELS.md and
# the embedded internal/cli/market_catalog.json snapshot) and OpenAPI spec
# (research/kie-final-openapi.yaml) from docs.kie.ai.
#
# This does NOT regenerate the command tree for a new model that keeps the
# shared /api/v1/jobs/createTask task contract. Rebuild the CLI to embed the
# refreshed catalog. A Market Chat, Omni, or other distinct endpoint shape
# needs a spec update and the full cli-printing-press generate pipeline (see
# README.md).
#
# Intended to run from cron (see README.md "Keeping this up to date").
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_DIR"

python3 research/build_spec.py

if git diff --quiet -- docs/MODELS.md internal/cli/market_catalog.json research/kie-final-openapi.yaml; then
  echo "$(date -Iseconds) weekly-refresh: no changes"
  exit 0
fi

git add docs/MODELS.md internal/cli/market_catalog.json research/kie-final-openapi.yaml
git commit -m "chore: weekly model catalog refresh from docs.kie.ai"
git push

echo "$(date -Iseconds) weekly-refresh: pushed catalog update"
