#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"

compose() {
  docker compose -f "$ROOT/docker-compose.yml" -f "$ROOT/docker-compose.e2e.yml" "$@"
}

echo "==> Ensuring app stack is up (recreate backend to clear in-memory history)"
# --force-recreate backend so GET /api/history starts empty for e2e.
compose up -d --build --force-recreate backend
compose up -d --build frontend otel-lgtm chrome

echo "==> Running Behave / Selenium e2e"
compose run --rm --build e2e "$@"

echo "==> E2E finished"
