package slogdiscard_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/10Narratives/golang-toolkit/logging/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
)

func TestDiscardHandler_Handle(t *testing.T) {
	t.Run("returns nil error", func(t *testing.T) {
		h := slogdiscard.NewHandler(nil, nil)
		assert.NoError(t, h.Handle(context.Background(), slog.Record{}))
	})
}

func TestDiscardHandler_WithAttrs(t *testing.T) {
	t.Run("returns same handler instance", func(t *testing.T) {
		h := slogdiscard.NewHandler(nil, nil)
		withAttrs := h.WithAttrs([]slog.Attr{})
		assert.Equal(t, h, withAttrs)
	})
}

func TestDiscardHandler_WithGroup(t *testing.T) {
	t.Run("returns same handler instance", func(t *testing.T) {
		h := slogdiscard.NewHandler(nil, nil)
		withGroups := h.WithGroup("some group")
		assert.Equal(t, h, withGroups)
	})
}

func TestDiscardHandler_Enabled(t *testing.T) {
	t.Run("always return false", func(t *testing.T) {
		h := slogdiscard.NewHandler(nil, nil)
		assert.False(t, h.Enabled(context.Background(), slog.LevelInfo))
	})
}
