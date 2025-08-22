package sqlitecomp

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

type App struct {
	cfg *Config
	log *slog.Logger
	db  *sql.DB
}

func New(cfg *Config, log *slog.Logger) (*App, error) {
	if cfg == nil {
		return nil, errors.New("configuration is required")
	}

	if log == nil {
		return nil, errors.New("logger is required")
	}

	return &App{
		cfg: cfg,
		log: log.With("component", "sqlite"),
	}, nil
}

func (a *App) Run() error {
	a.log.Info("connecting to database", "file", a.cfg.FilePath)

	db, err := sql.Open("sqlite", a.cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if a.cfg.CacheSize > 0 {
		_, err = db.Exec(fmt.Sprintf("PRAGMA cache_size = %d", a.cfg.CacheSize))
		if err != nil {
			return fmt.Errorf("failed to set cache size: %w", err)
		}
	}

	if a.cfg.ForeignKeys {
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	a.db = db
	a.log.Info("successfully connected to database")
	return nil
}

func (a *App) Stop() error {
	if a.db == nil {
		return nil
	}

	a.log.Info("closing database connection")
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	a.db = nil
	return nil
}
