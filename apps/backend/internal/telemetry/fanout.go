package telemetry

import (
	"context"
	"log/slog"
)

// fanoutHandler writes each record to every child handler.
type fanoutHandler struct {
	handlers []slog.Handler
}

func newFanoutHandler(handlers ...slog.Handler) slog.Handler {
	return &fanoutHandler{handlers: handlers}
}

func (h *fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, child := range h.handlers {
		if child.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *fanoutHandler) Handle(ctx context.Context, record slog.Record) error {
	var first error
	for _, child := range h.handlers {
		if !child.Enabled(ctx, record.Level) {
			continue
		}
		if err := child.Handle(ctx, record.Clone()); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func (h *fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	children := make([]slog.Handler, len(h.handlers))
	for i, child := range h.handlers {
		children[i] = child.WithAttrs(attrs)
	}
	return newFanoutHandler(children...)
}

func (h *fanoutHandler) WithGroup(name string) slog.Handler {
	children := make([]slog.Handler, len(h.handlers))
	for i, child := range h.handlers {
		children[i] = child.WithGroup(name)
	}
	return newFanoutHandler(children...)
}
