package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func NewConfig[T any]() (*T, error) {
	_ = godotenv.Load()

	var config = new(T)
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}
