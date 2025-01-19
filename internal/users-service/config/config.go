package config

import (
	"brain-wave/internal/shared/config"
	"math/rand/v2"
)

var _ config.ServiceConfig = (*Config)(nil)

type Config struct {
	config.DefaultConfig
	DatabaseURL string `env:"DATABASE_URL"`
}

func (c *Config) SetDefaults() {
	c.WebPort = c.WebPort + rand.IntN(100)
	c.GRPCPort = c.GRPCPort + rand.IntN(100)
}
