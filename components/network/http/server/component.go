package servercomp

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

type Component struct {
	cfg    *Config
	log    *slog.Logger
	server *http.Server
}

func New(cfg *Config, log *slog.Logger, handler http.Handler) (*Component, error) {
	return &Component{
		cfg: cfg,
		log: log.With("component", "http_server"),
		server: &http.Server{
			Handler: handler,
		},
	}, nil
}

func (c *Component) Run() error {
	listener, err := net.Listen("tcp", c.cfg.Address)
	if err != nil {
		return err
	}

	if err := c.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (c *Component) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ShutdownTimeout)
	defer cancel()

	if err := c.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
