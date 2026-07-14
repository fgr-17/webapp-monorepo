#!/usr/bin/env sh
# Shared helpers for tool containers (Go / Node). Sourced by other scripts.
# shellcheck shell=sh

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
COMPOSE_TOOLS="docker compose -p webapp-monorepo-tools -f $ROOT/docker-compose.tools.yml"


backend_tools() {
  $COMPOSE_TOOLS run --rm backend-tools "$@"
}

frontend_tools() {
  $COMPOSE_TOOLS run --rm frontend-tools "$@"
}

frontend_npm_ci() {
  # Install into the named node_modules volume (idempotent).
  frontend_tools npm ci
}
