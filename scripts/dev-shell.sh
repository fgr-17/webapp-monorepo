#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
# shellcheck source=./_tools.sh
. "$ROOT/scripts/_tools.sh"

TARGET="${1:-}"

case "$TARGET" in
  backend|be)
    echo "Opening shell in backend-tools (Go). Working dir: /src (= apps/backend)"
    backend_tools sh
    ;;
  frontend|fe)
    echo "Opening shell in frontend-tools (Node). Working dir: /src (= apps/frontend)"
    frontend_npm_ci
    frontend_tools sh
    ;;
  -h|--help|"" )
    cat <<'EOF'
Usage: ./scripts/dev-shell.sh <backend|frontend>

Drop into a tool container with the matching toolchain (no host install).

  backend   golang:1.25-alpine  → apps/backend mounted at /src
  frontend  node:20-alpine      → apps/frontend mounted at /src

Examples inside the shell:
  go test ./...
  go run ./cmd/server
  npm test
  npm run test:coverage:html
EOF
    [ -n "$TARGET" ] || exit 1
    exit 0
    ;;
  *)
    echo "Unknown target: $TARGET (use backend or frontend)" >&2
    exit 1
    ;;
esac
