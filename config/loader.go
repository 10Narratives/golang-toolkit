package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Loader[T any] struct {
}

func (l *Loader[T]) Load() *T {
	var path string
	flag.StringVar(&path, "config", "", "path to configuration file")
	flag.Parse()

	return nil
}

func (l *Loader[T]) LoadFromFile(path string) (*T, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}

	var cfg T
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read configuration: %w", err)
	}

	return &cfg, nil
}
