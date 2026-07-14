#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
COMPOSE="docker compose -f $ROOT/docker-compose.yml -f $ROOT/docker-compose.e2e.yml"

echo "==> Ensuring app stack is up"
$COMPOSE up -d --build backend frontend otel-lgtm chrome

echo "==> Running Behave / Selenium e2e"
$COMPOSE run --rm --build e2e "$@"

echo "==> E2E finished"
