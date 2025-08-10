// Package slogdiscard provides a no-op logger and handler for the Go 'slog' package.
// Useful for disabling log output in testing or production when logging is not desired.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger returns a new slog.Logger that discards all log records.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is a no-op implementation of slog.Handler.
type DiscardHandler struct{}

// NewDiscardHandler returns a new instance of DiscardHandler.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle receives a log record and does nothing with it.
// Always returns nil, indicating no error but no action taken.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same DiscardHandler instance.
// Since no attributes are ever used, this method has no effect.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same DiscardHandler instance.
// Grouping is ignored, as all logs are discarded.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled reports whether a record at the given level should be logged.
// Always returns false, meaning no log records are ever enabled.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
