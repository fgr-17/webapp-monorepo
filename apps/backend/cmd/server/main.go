package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zigglib/webapp-monorepo/backend/api"
	"github.com/zigglib/webapp-monorepo/backend/internal/calc"
	"github.com/zigglib/webapp-monorepo/backend/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type operands struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var (
	tracer           = otel.Tracer("calculator")
	opsCounter       metric.Int64Counter
	opsDuration      metric.Float64Histogram
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTelemetry, err := telemetry.Setup(ctx)
	if err != nil {
		slog.Error("telemetry setup failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown failed", "error", err)
		}
	}()

	meter := otel.Meter("calculator")
	opsCounter, err = meter.Int64Counter(
		"calculator.operations",
		metric.WithDescription("Number of calculator operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		slog.Error("create operations counter", "error", err)
		os.Exit(1)
	}
	opsDuration, err = meter.Float64Histogram(
		"calculator.operation.duration",
		metric.WithDescription("Calculator operation latency"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		slog.Error("create operations histogram", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	docs := api.Handler()
	mux.Handle("GET /openapi.yaml", docs)
	mux.Handle("GET /swagger", docs)
	mux.Handle("GET /swagger/", docs)
	mux.HandleFunc("GET /health", handleHealth)
	mux.Handle("POST /api/add", handleOp("add", func(a, b float64) (float64, error) {
		return calc.Add(a, b), nil
	}))
	mux.Handle("POST /api/subtract", handleOp("subtract", func(a, b float64) (float64, error) {
		return calc.Subtract(a, b), nil
	}))
	mux.Handle("POST /api/multiply", handleOp("multiply", func(a, b float64) (float64, error) {
		return calc.Multiply(a, b), nil
	}))
	mux.Handle("POST /api/divide", handleOp("divide", calc.Divide))

	handler := withCORS(otelhttp.NewHandler(mux, "calculator-backend",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health"
		}),
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		slog.Info("backend listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleOp(name string, op func(a, b float64) (float64, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "calc."+name)
		defer span.End()
		start := time.Now()

		var in operands
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid JSON body")
			recordOp(ctx, name, "error", start)
			slog.WarnContext(ctx, "invalid request body", "operation", name, "error", err)
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		span.SetAttributes(
			attribute.String("calc.operation", name),
			attribute.Float64("calc.a", in.A),
			attribute.Float64("calc.b", in.B),
		)

		result, err := op(in.A, in.B)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			recordOp(ctx, name, "error", start)
			slog.WarnContext(ctx, "operation failed", "operation", name, "error", err)
			if errors.Is(err, calc.ErrDivideByZero) {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "operation failed")
			return
		}

		span.SetAttributes(attribute.Float64("calc.result", result))
		recordOp(ctx, name, "ok", start)
		slog.InfoContext(ctx, "operation completed",
			"operation", name,
			"a", in.A,
			"b", in.B,
			"result", result,
		)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resultResponse{Result: result})
	})
}

func recordOp(ctx context.Context, name, status string, start time.Time) {
	attrs := metric.WithAttributes(
		attribute.String("operation", name),
		attribute.String("status", status),
	)
	opsCounter.Add(ctx, 1, attrs)
	opsDuration.Record(ctx, float64(time.Since(start).Milliseconds()), attrs)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Traceparent, Tracestate")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
