// Package slogdiscard provides a no-op logger and handler for the Go 'slog' package.
// Useful for disabling log output in testing or production when logging is not desired.
package slogdiscard

import (
	"context"
	"io"
	"log/slog"
)

// Handler is a no-op implementation of slog.Handler.DiscardHandler
type Handler struct{}

var _ slog.Handler = &Handler{}

// NewHandler returns a new instance of discard handler.
func NewHandler(_ io.Writer, _ *slog.HandlerOptions) *Handler {
	return &Handler{}
}

// Handle receives a log record and does nothing with it.
// Always returns nil, indicating no error but no action taken.
func (h *Handler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same discard handler instance.
// Since no attributes are ever used, this method has no effect.
func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same discard handler instance.
// Grouping is ignored, as all logs are discarded.
func (h *Handler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled reports whether a record at the given level should be logged.
// Always returns false, meaning no log records are ever enabled.
func (h *Handler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
