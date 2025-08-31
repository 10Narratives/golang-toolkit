package sqlitecomp

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type Component struct {
	cfg *Config
	log *slog.Logger

	DB *sql.DB
}

func New(cfg *Config, log *slog.Logger) (*Component, error) {
	return &Component{
		cfg: cfg,
		log: log.With("component", "sqlite"),
	}, nil
}

func (c *Component) Run() error {
	var err error
	c.DB, err = sql.Open("sqlite3", c.cfg.FilePath)
	if err != nil {
		return err
	}

	c.DB.SetConnMaxLifetime(c.cfg.ConnMaxLifetime)
	c.DB.SetMaxOpenConns(c.cfg.MaxOpenConns)
	c.DB.SetMaxIdleConns(c.cfg.MaxIdleConns)

	if err := c.DB.Ping(); err != nil {
		c.DB.Close()
		return err
	}

	return nil
}

func (c *Component) Stop() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
