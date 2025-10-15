#!/usr/bin/env bash
set -euo pipefail

OUT_FILE="sample-recommendation.yaml"

echo "Running analyzer sample and writing recommendation to ${OUT_FILE}"

# Run from project root: use the cmd/controller path relative to repo root
go run ./cmd/controller --sample > "${OUT_FILE}"

echo "Wrote ${OUT_FILE}"
