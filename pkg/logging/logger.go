package logging

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/10Narratives/golang-toolkit/pkg/logging/handlers/slogdiscard"
	"github.com/10Narratives/golang-toolkit/pkg/logging/handlers/slogpretty"
	"github.com/natefinch/lumberjack"
)

func New(opts ...LoggerOption) (*slog.Logger, error) {
	handler, err := createHandler(opts...)
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

func createOutput(fpath string) (io.Writer, error) {
	if fpath == "stdout" {
		return os.Stdout, nil
	}

	dir := filepath.Dir(fpath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.New("failed to create log directory")
	}

	return &lumberjack.Logger{
		Filename:  fpath,
		LocalTime: true,
	}, nil
}

func createHandler(opts ...LoggerOption) (slog.Handler, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	out, err := createOutput(options.output)
	if err != nil {
		return nil, err
	}

	slogOpts := &slog.HandlerOptions{
		Level: slog.Level(options.level),
	}

	switch options.format {
	case "json":
		return slog.NewJSONHandler(out, slogOpts), nil
	case "pretty":
		return slogpretty.NewHandler(out, slogOpts), nil
	case "plain":
		return slog.NewTextHandler(out, slogOpts), nil
	case "discard":
		return slogdiscard.NewHandler(out, slogOpts), nil
	default:
		return nil, errors.New("unsupported logger format")
	}
}
