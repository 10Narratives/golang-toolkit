package servercomp

import "time"

type Config struct {
	Address         string        `yaml:"address" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"10s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"15s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"30s"`
}
