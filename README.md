# webapp-monorepo

Minimal frontend/backend monorepo template: a calculator UI talking to a Go REST API, each service in its own container.

## Layout

```
webapp-monorepo/
  apps/
    frontend/     # static HTML/CSS/JS + nginx
    backend/      # Go REST API + OpenTelemetry
  docker-compose.yml
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

## Observability (OpenTelemetry)

The backend exports **traces**, **metrics**, and **logs** via OTLP/HTTP when `OTEL_EXPORTER_OTLP_ENDPOINT` is set.

Compose includes [`grafana/otel-lgtm`](https://github.com/grafana/docker9-otel-lgtm): an all-in-one stack (Grafana + Tempo + Loki + Prometheus) that receives OTLP on `:4317`/`:4318`.

| Signal | What is emitted |
|--------|-----------------|
| Traces | HTTP spans (`otelhttp`) + `calc.{add,subtract,multiply,divide}` spans with operands/result |
| Metrics | `calculator.operations` counter, `calculator.operation.duration` histogram (by operation/status); HTTP metrics from instrumentation |
| Logs | Structured `slog` bridged to OTel logs (correlated with trace context) |

In Grafana Explore: use **Tempo** for traces, **Prometheus** for metrics (`calculator_operations_*`), **Loki** for logs. `/health` is excluded from HTTP tracing to reduce noise.

Without `OTEL_EXPORTER_OTLP_ENDPOINT`, telemetry setup is a no-op and logs go to stdout as JSON.

## API

All operation endpoints accept `POST` with JSON body `{"a": number, "b": number}` and return `{"result": number}`. Interactive docs live at `/swagger/` (OpenAPI spec at `/openapi.yaml`).

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
