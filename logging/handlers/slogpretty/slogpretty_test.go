package slogpretty_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/10Narratives/golang-toolkit/logging/handlers/slogpretty"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestHandler_WithGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		out     io.Writer
		initial slog.Level
	}{
		{name: "os.Stdout", out: os.Stdout, initial: slog.LevelError},
		{name: "os.Stderr", out: os.Stderr, initial: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := slogpretty.NewHandler(tt.out, &slog.HandlerOptions{
				Level: tt.initial,
			})
			assert.Equal(t, h, h.WithGroup(""))
		})
	}
}

func TestHandler_Enabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initial, level slog.Level
		want           assert.BoolAssertionFunc
	}{
		{name: "Debug enabled when min=Debug", initial: slog.LevelDebug, level: slog.LevelDebug, want: assert.True},
		{name: "Info enabled when min=Debug", initial: slog.LevelDebug, level: slog.LevelInfo, want: assert.True},
		{name: "Debug disabled when min=Info", initial: slog.LevelInfo, level: slog.LevelDebug, want: assert.False},
		{name: "Info enabled when min=Info", initial: slog.LevelInfo, level: slog.LevelInfo, want: assert.True},
		{name: "Warn enabled when min=Info", initial: slog.LevelInfo, level: slog.LevelWarn, want: assert.True},
		{name: "Info disabled when min=Warn", initial: slog.LevelWarn, level: slog.LevelInfo, want: assert.False},
		{name: "Error enabled when min=Warn", initial: slog.LevelWarn, level: slog.LevelError, want: assert.True},
		{name: "Error enabled when min=Error", initial: slog.LevelError, level: slog.LevelError, want: assert.True},
		{name: "Warn disabled when min=Error", initial: slog.LevelError, level: slog.LevelWarn, want: assert.False},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := slogpretty.NewHandler(nil, &slog.HandlerOptions{
				Level: tt.initial,
			})
			enabled := h.Enabled(context.Background(), tt.level)
			tt.want(t, enabled)
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                        string
		initial, additional, wanted []slog.Attr
	}{
		{
			name:       "empty initial and additional",
			initial:    nil,
			additional: nil,
			wanted:     nil,
		},
		{
			name:       "only initial attrs",
			initial:    []slog.Attr{slog.String("env", "test")},
			additional: nil,
			wanted:     []slog.Attr{slog.String("env", "test")},
		},
		{
			name:       "only additional attrs",
			initial:    nil,
			additional: []slog.Attr{slog.Int("count", 42)},
			wanted:     []slog.Attr{slog.Int("count", 42)},
		},
		{
			name:       "merge different keys",
			initial:    []slog.Attr{slog.Bool("active", true)},
			additional: []slog.Attr{slog.String("service", "auth")},
			wanted:     []slog.Attr{slog.Bool("active", true), slog.String("service", "auth")},
		},
		{
			name: "duplicate keys (new appended after old)",
			initial: []slog.Attr{
				slog.String("user", "old"),
				slog.Int("version", 1),
			},
			additional: []slog.Attr{
				slog.String("user", "new"),
				slog.Int("version", 2),
			},
			wanted: []slog.Attr{
				slog.String("user", "old"),
				slog.Int("version", 1),
				slog.String("user", "new"),
				slog.Int("version", 2),
			},
		},
		{
			name:    "multiple sequential additions",
			initial: []slog.Attr{slog.String("base", "value")},
			additional: []slog.Attr{
				slog.Bool("flag", true),
				slog.Float64("pi", 3.14),
			},
			wanted: []slog.Attr{
				slog.String("base", "value"),
				slog.Bool("flag", true),
				slog.Float64("pi", 3.14),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := slogpretty.NewHandler(io.Discard, nil).
				WithAttrs(tt.wanted)

			h := slogpretty.NewHandler(io.Discard, nil).
				WithAttrs(tt.initial).
				WithAttrs(tt.additional)

			assert.Equal(t, expected, h)
		})
	}
}
func TestHandler_Handle(t *testing.T) {
	oldNoColor := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = oldNoColor }()

	fixedTime := time.Date(2023, 1, 1, 12, 30, 45, 123000000, time.UTC)

	tests := []struct {
		name         string
		level        slog.Level
		msg          string
		recordAttrs  []slog.Attr
		handlerAttrs []slog.Attr
		expected     string
	}{
		{
			name:     "basic info log",
			level:    slog.LevelInfo,
			msg:      "Hello",
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[34mINFO:\x1b[0m] \x1b[36mHello\x1b[0m\n",
		},
		{
			name:        "with record attributes",
			level:       slog.LevelInfo,
			msg:         "Request processed",
			recordAttrs: []slog.Attr{slog.String("method", "GET"), slog.Int("status", 200)},
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[34mINFO:\x1b[0m] \x1b[36mRequest processed\x1b[0m\n" +
				"{\n  \"method\": \"GET\",\n  \"status\": 200\n}\n",
		},
		{
			name:         "with handler attributes",
			level:        slog.LevelInfo,
			msg:          "System check",
			handlerAttrs: []slog.Attr{slog.String("service", "api"), slog.Bool("healthy", true)},
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[34mINFO:\x1b[0m] \x1b[36mSystem check\x1b[0m\n" +
				"{\n  \"healthy\": true,\n  \"service\": \"api\"\n}\n", // Исправлен порядок ключей
		},
		{
			name:         "with both attributes",
			level:        slog.LevelInfo,
			msg:          "Combined",
			recordAttrs:  []slog.Attr{slog.String("user", "admin")},
			handlerAttrs: []slog.Attr{slog.Int("id", 42)},
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[34mINFO:\x1b[0m] \x1b[36mCombined\x1b[0m\n" +
				"{\n  \"id\": 42,\n  \"user\": \"admin\"\n}\n", // Исправлен порядок ключей
		},
		{
			name:     "error level",
			level:    slog.LevelError,
			msg:      "Failed",
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[31mERROR:\x1b[0m] \x1b[36mFailed\x1b[0m\n",
		},
		{
			name:     "warn level",
			level:    slog.LevelWarn,
			msg:      "Warning",
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[33mWARN:\x1b[0m] \x1b[36mWarning\x1b[0m\n",
		},
		{
			name:     "debug level",
			level:    slog.LevelDebug,
			msg:      "Debug info",
			expected: "\x1b[32m12:30:45.123\x1b[0m [\x1b[35mDEBUG:\x1b[0m] \x1b[36mDebug info\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := slogpretty.NewHandler(buf, nil).WithAttrs(tt.handlerAttrs)

			rec := slog.NewRecord(fixedTime, tt.level, tt.msg, 0)
			rec.AddAttrs(tt.recordAttrs...)

			err := h.Handle(context.Background(), rec)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}
