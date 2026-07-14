#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
# shellcheck source=./_tools.sh
. "$ROOT/scripts/_tools.sh"

MODE="${1:-docker}"

run_docker() {
  echo "==> Backend tests (backend-tools container)"
  backend_tools go test ./...

  echo
  echo "==> Frontend tests (frontend-tools container)"
  frontend_npm_ci
  frontend_tools npm test
}

run_host() {
  echo "==> Backend tests (host)"
  (cd "$ROOT/apps/backend" && go test ./...)

  echo
  echo "==> Frontend tests (host)"
  (cd "$ROOT/apps/frontend" && npm test)
}

case "$MODE" in
  docker|""|--docker)
    run_docker
    ;;
  --host|host)
    run_host
    ;;
  -h|--help)
    cat <<'EOF'
Usage: ./scripts/test.sh [--docker|--host]

  --docker  Default. Run tests in Go/Node tool containers
            (docker-compose.tools.yml). No host Go/Node needed.
  --host    Use Go/Node installed on the machine.
EOF
    exit 0
    ;;
  *)
    echo "Unknown mode: $MODE (use --docker or --host)" >&2
    exit 1
    ;;
esac
