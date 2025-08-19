package pgapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type App struct {
	cfg  *Config
	log  *slog.Logger
	conn *pgx.Conn
}

func New(cfg *Config, log *slog.Logger) (*App, error) {
	if cfg == nil {
		return nil, errors.New("configuration is required")
	}

	if log != nil {
		return nil, errors.New("logger is required")
	}

	return &App{
		cfg: cfg,
		log: log.With("component", "postgres"),
	}, nil
}

func (a *App) Run() error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		a.cfg.Host, a.cfg.Port, a.cfg.User, a.cfg.Password, a.cfg.DBName, a.cfg.SSLMode,
	)

	a.log.Info("connecting to database")

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	a.conn = conn

	a.log.Info("successfully connected to database")

	return nil
}

func (a *App) Stop() error {
	if a.conn == nil {
		return nil
	}

	a.log.Info("closing database connection")
	err := a.conn.Close(context.Background())
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	a.conn = nil
	return nil
}
