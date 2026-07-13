# Common developer entrypoints (Docker toolchains — no host Go/Node required).
# Style is applied via `make fix`; lint must be green before `make build`.

.PHONY: help fix lint test coverage e2e build up down dev-shell-backend dev-shell-frontend

help:
	@echo "Targets:"
	@echo "  make fix       Apply style (gofmt/goimports, eslint --fix, prettier) then verify lint is green"
	@echo "  make lint      Check style + linters (must pass; does not rewrite files)"
	@echo "  make test      Unit tests (BE + FE tool containers)"
	@echo "  make coverage  Coverage HTML reports under ./coverage"
	@echo "  make e2e       Behave + Selenium against the Compose stack"
	@echo "  make build     Enforce lint (green) then build runtime images"
	@echo "  make up        make build && docker compose up -d"
	@echo "  make down      docker compose down"

fix:
	./scripts/lint.sh --fix
	@echo "Style applied; lint is green."

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

coverage:
	./scripts/coverage.sh

e2e:
	./scripts/e2e.sh

build: lint
	./scripts/build.sh

up: build
	docker compose up -d

down:
	docker compose down

dev-shell-backend:
	./scripts/dev-shell.sh backend

dev-shell-frontend:
	./scripts/dev-shell.sh frontend
