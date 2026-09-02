#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
# shellcheck source=./_tools.sh
. "$ROOT/scripts/_tools.sh"

MODE="docker"
FIX=0

for arg in "$@"; do
  case "$arg" in
    --host|host) MODE="host" ;;
    --docker|docker) MODE="docker" ;;
    --fix|fix) FIX=1 ;;
    -h|--help)
      cat <<'EOF'
Usage: ./scripts/lint.sh [--docker|--host] [--fix]

  --docker  Default. Run linters in Go/Node tool containers.
  --host    Use tools installed on the machine.
  --fix    Auto-fix where supported (gofmt/goimports, eslint, prettier).

Backend:  gofmt/goimports check + golangci-lint
Frontend: ESLint + Prettier
EOF
      exit 0
      ;;
    *)
      echo "Unknown arg: $arg" >&2
      exit 1
      ;;
  esac
done

lint_backend_docker() {
  echo "==> Backend lint (backend-tools container)"
  if [ "$FIX" -eq 1 ]; then
    backend_tools sh -c 'gofmt -w . && goimports -local github.com/fgr-17/webapp-monorepo/backend -w . && golangci-lint run ./...'
  else
    backend_tools sh -c '
      unformatted=$(gofmt -l .)
      if [ -n "$unformatted" ]; then
        echo "gofmt needed on:" >&2
        echo "$unformatted" >&2
        exit 1
      fi
      golangci-lint run ./...
    '
  fi
}

lint_frontend_docker() {
  echo "==> Frontend lint (frontend-tools container)"
  frontend_npm_ci
  if [ "$FIX" -eq 1 ]; then
    frontend_tools npm run lint:fix
  else
    frontend_tools npm run lint
  fi
}

lint_backend_host() {
  echo "==> Backend lint (host)"
  cd "$ROOT/apps/backend"
  if [ "$FIX" -eq 1 ]; then
    gofmt -w .
    command -v goimports >/dev/null && goimports -local github.com/fgr-17/webapp-monorepo/backend -w .
    golangci-lint run ./...
  else
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
      echo "gofmt needed on:" >&2
      echo "$unformatted" >&2
      exit 1
    fi
    golangci-lint run ./...
  fi
}

lint_frontend_host() {
  echo "==> Frontend lint (host)"
  cd "$ROOT/apps/frontend"
  npm ci
  if [ "$FIX" -eq 1 ]; then
    npm run lint:fix
  else
    npm run lint
  fi
}

case "$MODE" in
  docker)
    lint_backend_docker
    echo
    lint_frontend_docker
    ;;
  host)
    lint_backend_host
    echo
    lint_frontend_host
    ;;
esac

echo
echo "Lint OK"
