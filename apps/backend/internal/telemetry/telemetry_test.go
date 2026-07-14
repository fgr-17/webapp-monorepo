package telemetry

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestSetupWithoutEndpoint(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_SERVICE_NAME", "")

	shutdown, err := Setup(context.Background())
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected shutdown func")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestFanoutHandler(t *testing.T) {
	var bufA, bufB bytes.Buffer
	hA := slog.NewJSONHandler(&bufA, nil)
	hB := slog.NewJSONHandler(&bufB, nil)
	logger := slog.New(newFanoutHandler(hA, hB))

	logger.Info("hello", "k", "v")

	if !bytes.Contains(bufA.Bytes(), []byte("hello")) {
		t.Fatalf("handler A missing log: %s", bufA.String())
	}
	if !bytes.Contains(bufB.Bytes(), []byte("hello")) {
		t.Fatalf("handler B missing log: %s", bufB.String())
	}
	if !bytes.Contains(bufA.Bytes(), []byte(`"k":"v"`)) && !bytes.Contains(bufA.Bytes(), []byte(`"k": "v"`)) {
		// JSON handler may format without spaces
		if !bytes.Contains(bufA.Bytes(), []byte(`"k"`)) {
			t.Fatalf("handler A missing attr: %s", bufA.String())
		}
	}
}

func TestFanoutWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	base := newFanoutHandler(slog.NewJSONHandler(&buf, nil))
	h := base.WithAttrs([]slog.Attr{slog.String("service", "calc")})
	slog.New(h).Info("boot")
	if !bytes.Contains(buf.Bytes(), []byte("service")) {
		t.Fatalf("missing attr: %s", buf.String())
	}
}

func TestMain(m *testing.M) {
	// Ensure no leftover endpoint from the environment affects noop tests.
	_ = os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	os.Exit(m.Run())
}
