#!/usr/bin/env sh
# Build runtime images. Lint must already be green (Makefile runs `lint` first).
# Dockerfiles also re-check quality during image build so raw `docker compose build`
# cannot bypass the gate.
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Building runtime images (Dockerfile quality gates enabled)"
docker compose build "$@"

echo "==> Build OK"
