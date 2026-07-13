# webapp-monorepo

Minimal frontend/backend monorepo template: a calculator UI talking to a Go REST API, each service in its own container.

## Layout

```
webapp-monorepo/
  apps/
    frontend/              # static HTML/CSS/JS + nginx
    backend/               # Go REST API + OpenTelemetry
    e2e/                   # Behave + Selenium browser tests
  docker-compose.yml       # runtime: app + Grafana LGTM
  docker-compose.tools.yml # dev tools: Go + Node containers
  docker-compose.e2e.yml   # chrome + behave runner
  scripts/                 # test, lint, coverage, e2e, build helpers
  Makefile                 # make fix | lint | build | up | test | e2e
```

## Run with Docker Compose

```bash
docker compose up --build
```

- UI: http://localhost:8085 (host port mapped in `docker-compose.yml`)
- API (direct): http://localhost:8081
- Swagger UI: http://localhost:8081/swagger/ (also via frontend proxy at `/swagger/`)
- Observability (Grafana LGTM): http://localhost:3000 — traces (Tempo), metrics (Prometheus), logs (Loki)

The frontend nginx proxies `/api/*`, `/swagger/`, and `/openapi.yaml` to the backend.

## Docs

| Topic | Guide |
|-------|--------|
| Swagger / OpenAPI | [apps/backend/api/README.md](apps/backend/api/README.md) |
| Telemetry (OpenTelemetry) | [apps/backend/internal/telemetry/README.md](apps/backend/internal/telemetry/README.md) |
| Testing & coverage | [docs/testing.md](docs/testing.md) |
| E2E (Behave + Selenium) | [apps/e2e/README.md](apps/e2e/README.md) |

## Build & quality gates

Preferred workflow (Docker only on the host):

```bash
make fix      # apply style, then verify linters are green
make lint     # check only (fails if style/lint is dirty)
make build    # requires green lint, then builds images (Dockerfiles re-check too)
make up       # make build && docker compose up -d
make test     # unit tests
make e2e      # Behave + Selenium
```

| Step | What happens |
|------|----------------|
| `make fix` | Auto-format (gofmt/goimports, ESLint `--fix`, Prettier) + lint must end green |
| `make lint` | Style **check** + golangci-lint / ESLint (no rewrite); must be all green |
| `make build` | Runs `make lint`, then `docker compose build` |
| Image build | Backend/frontend Dockerfiles run the same quality gates before shipping the image |

So you cannot ship a bundle with failing lint: both the Makefile and the Dockerfiles enforce it. If lint fails, run `make fix` and re-commit.

## Testing & coverage

Full guide: **[docs/testing.md](docs/testing.md)**.

You do **not** need Go or Node on the host. Tests/coverage run in **tool containers** (`docker-compose.tools.yml`), separate from the app containers started by `docker compose up`.

```bash
./scripts/test.sh                 # Go + Node unit tests in containers
./scripts/lint.sh                 # golangci-lint + ESLint/Prettier
./scripts/lint.sh --fix          # auto-fix style where possible
./scripts/coverage.sh             # → coverage/backend.html, coverage/frontend/index.html
./scripts/dev-shell.sh backend    # interactive Go toolchain
./scripts/dev-shell.sh frontend   # interactive Node toolchain

# optional if you already have Go/Node installed locally
./scripts/test.sh --host
./scripts/lint.sh --host
./scripts/coverage.sh --host
```

### E2E (Behave + Selenium)

Browser tests against the real stack — guide: **[apps/e2e/README.md](apps/e2e/README.md)**.

```bash
./scripts/e2e.sh                  # starts app + Chromium + runs Behave
```

## Observability (OpenTelemetry)

Short summary — full newbie guide with diagrams: **[telemetry README](apps/backend/internal/telemetry/README.md)**.

The backend exports **traces**, **metrics**, and **logs** via OTLP/HTTP when `OTEL_EXPORTER_OTLP_ENDPOINT` is set. Compose includes [`grafana/otel-lgtm`](https://github.com/grafana/docker-otel-lgtm) (Grafana + Tempo + Loki + Prometheus).

In Grafana use **Explore** (not empty Starred dashboards): Tempo / Prometheus / Loki. Without the OTLP endpoint, exporters are off and logs go to stdout as JSON.

## API & Swagger

Short summary — full newbie guide with diagrams: **[Swagger README](apps/backend/api/README.md)**.

All operation endpoints accept `POST` with JSON body `{"a": number, "b": number}` and return `{"result": number}`. Interactive docs: `/swagger/` (spec: `/openapi.yaml`).

| Method | Path | Operation |
|--------|------|-----------|
| GET | `/health` | Health check |
| GET | `/swagger/` | Swagger UI |
| GET | `/openapi.yaml` | OpenAPI 3 spec |
| POST | `/api/add` | a + b |
| POST | `/api/subtract` | a − b |
| POST | `/api/multiply` | a × b |
| POST | `/api/divide` | a ÷ b (400 on divide-by-zero) |

### curl examples

```bash
curl -s http://localhost:8081/health

curl -s -X POST http://localhost:8081/api/add \
  -H 'Content-Type: application/json' \
  -d '{"a":10,"b":2}'

curl -s -X POST http://localhost:8081/api/subtract \
  -H 'Content-Type: application/json' \
  -d '{"a":10,"b":2}'

curl -s -X POST http://localhost:8081/api/multiply \
  -H 'Content-Type: application/json' \
  -d '{"a":10,"b":2}'

curl -s -X POST http://localhost:8081/api/divide \
  -H 'Content-Type: application/json' \
  -d '{"a":10,"b":2}'
```

## Local backend (without Docker)

```bash
cd apps/backend
go test ./...
# optional: point at a local OTLP collector (Compose LGTM on :4318)
export OTEL_SERVICE_NAME=calculator-backend
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
go run ./cmd/server
```

Without `OTEL_EXPORTER_OTLP_ENDPOINT`, exporters are disabled and logs still print as JSON to stdout.

## Frontend without Docker

Serve the static files with any HTTP server (API calls expect `/api` to be reachable, e.g. via a reverse proxy to the backend).
