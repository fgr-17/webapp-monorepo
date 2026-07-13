# Quality gate (lint/style) then compile. `make build` / `docker compose build` both enforce this.
FROM golang:1.25-alpine AS quality

RUN apk add --no-cache git curl \
  && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
    | sh -s -- -b /usr/local/bin v2.12.2 \
  && go install golang.org/x/tools/cmd/goimports@v0.30.0

ENV PATH="/go/bin:${PATH}"
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN unformatted="$(gofmt -l .)" \
  && if [ -n "$unformatted" ]; then \
       echo "gofmt needed on:" >&2; echo "$unformatted" >&2; exit 1; \
     fi
RUN bad="$(goimports -l -local github.com/zigglib/webapp-monorepo/backend .)" \
  && if [ -n "$bad" ]; then \
       echo "goimports needed on:" >&2; echo "$bad" >&2; exit 1; \
     fi
RUN golangci-lint run ./...

FROM quality AS build
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /server /app/server
ENV PORT=8081
EXPOSE 8081
USER nobody
CMD ["/app/server"]
