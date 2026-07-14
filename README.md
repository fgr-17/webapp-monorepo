# webapp-monorepo

Minimal frontend/backend monorepo template: a calculator UI talking to a Go REST API, each service in its own container.

## Layout

```
webapp-monorepo/
  apps/
    frontend/              # static HTML/CSS/JS + nginx
    backend/               # Go REST API + OpenTelemetry
  docker-compose.yml       # runtime: app + Grafana LGTM
  docker-compose.tools.yml # dev tools: Go + Node containers
  scripts/                 # test, coverage, dev-shell helpers
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

## Testing & coverage

Full guide: **[docs/testing.md](docs/testing.md)**.

You do **not** need Go or Node on the host. Tests/coverage run in **tool containers** (`docker-compose.tools.yml`), separate from the app containers started by `docker compose up`.

```bash
./scripts/test.sh                 # Go + Node tests in containers
./scripts/coverage.sh             # → coverage/backend.html, coverage/frontend/index.html
./scripts/dev-shell.sh backend    # interactive Go toolchain
./scripts/dev-shell.sh frontend   # interactive Node toolchain

# optional if you already have Go/Node installed locally
./scripts/test.sh --host
./scripts/coverage.sh --host
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
