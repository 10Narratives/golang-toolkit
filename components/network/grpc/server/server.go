package grpcsrv

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	cfg  *Config
	log  *slog.Logger
	grpc *grpc.Server
}

func (s *App) New(
	cfg *Config,
	log *slog.Logger,
	unary []grpc.UnaryServerInterceptor,
	stream []grpc.StreamServerInterceptor,
	registers []func(server *grpc.Server),
) (*App, error) {
	server := &App{
		cfg: cfg,
		log: log,
	}

	server.grpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(unary...),
		grpc.ChainStreamInterceptor(stream...),
	)

	for _, register := range registers {
		register(server.grpc)
	}

	return server, nil
}

func (s *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port))
	if err != nil {
		return err
	}

	if err := s.grpc.Serve(l); err != nil {
		return err
	}

	return nil
}

func (s *App) Stop() error {
	s.grpc.GracefulStop()
	return nil
}
