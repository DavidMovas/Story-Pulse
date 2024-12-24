package config

import (
	"math/rand/v2"
	"story-pulse/internal/shared/config"
)

var _ config.ServiceConfig = (*Config)(nil)

type Config struct {
	config.DefaultConfig
}

func (c *Config) SetDefaults() {
	c.WebPort = c.WebPort + rand.IntN(100)
	c.GRPCPort = c.GRPCPort + rand.IntN(100)
}
