package sqlitecomp

import "time"

type Config struct {
	FilePath        string        `yaml:"file_path" env-required:"true"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"1"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"1"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"0"`
}
