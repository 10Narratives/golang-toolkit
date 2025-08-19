package slogpretty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type Handler struct {
	attrs []slog.Attr
	out   io.Writer
	opts  *slog.HandlerOptions
}

var _ slog.Handler = &Handler{}

func NewHandler(out io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		out:  out,
		opts: opts,
	}
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	level := rec.Level.String() + ":"
	switch rec.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	timestamp := color.GreenString(rec.Time.Format("15:04:05.000"))
	msg := color.CyanString(rec.Message)

	fields := make(map[string]any)
	rec.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()
		return true
	})

	for _, attr := range h.attrs {
		fields[attr.Key] = attr.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	var logLine string
	logLine = fmt.Sprintf("%s [%s] %s", timestamp, level, msg)
	if len(b) > 0 {
		logLine += "\n" + string(b)
	}
	logLine += "\n"

	_, err = h.out.Write([]byte(logLine))
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &Handler{
		attrs: newAttrs,
		out:   h.out,
		opts:  h.opts,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		attrs: h.attrs,
		out:   h.out,
		opts:  h.opts,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelDebug
	if h.opts != nil && h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}
