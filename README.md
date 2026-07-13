# webapp-monorepo

Minimal frontend/backend monorepo template: a calculator UI talking to a Go REST API, each service in its own container.

## Why REST (not gRPC)

REST + JSON is the better fit for this template: browsers call it with `fetch`, curl debugging is trivial, and there is no need for protobuf tooling or a grpc-web proxy. Prefer gRPC when you have service-to-service traffic, strong contracts, and streaming — not a browser calculator.

## Layout

```
webapp-monorepo/
  apps/
    frontend/     # static HTML/CSS/JS + nginx
    backend/      # Go REST API
  docker-compose.yml
```

## Run with Docker Compose

```bash
docker compose up --build
```

- UI: http://localhost:8080
- API (direct): http://localhost:8081

The frontend nginx proxies `/api/*` to the backend, so the browser only needs port 8080.

## API

All operation endpoints accept `POST` with JSON body `{"a": number, "b": number}` and return `{"result": number}`.

| Method | Path | Operation |
|--------|------|-----------|
| GET | `/health` | Health check |
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
go run ./cmd/server
```

## Frontend without Docker

Serve the static files with any HTTP server (API calls expect `/api` to be reachable, e.g. via a reverse proxy to the backend).
