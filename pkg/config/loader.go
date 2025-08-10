package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigLoader[T any] struct{}

func (cl *ConfigLoader[T]) Load() (*T, error) {
	var path string
	flag.StringVar(&path, "config", "", "path to configuration file")
	flag.Parse()

	return cl.LoadFromFile(path)
}

func (cl *ConfigLoader[T]) LoadFromFile(path string) (*T, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %v", err)
	}

	var cfg T
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %v", err)
	}

	return &cfg, nil
}
