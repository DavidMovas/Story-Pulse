package config

import (
	"math/rand/v2"
	"story-pulse/internal/shared/config"
	"time"
)

var _ config.ServiceConfig = (*Config)(nil)

type Config struct {
	config.DefaultConfig
	Secret                string        `env:"JWT_SECRET" envDefault:"secret"`
	AccessExpirationTime  time.Duration `env:"JWT_ACCESS_EXPIRATION_TIME" envDefault:"15m"`
	RefreshExpirationTime time.Duration `env:"JWT_REFRESH_EXPIRATION_TIME" envDefault:"168h"`
	RedisURL              string        `env:"REDIS_URL" envDefault:"redis:6379"`
}

func (c *Config) SetDefaults() {
	c.WebPort = c.WebPort + rand.IntN(100)
	c.GRPCPort = c.GRPCPort + rand.IntN(100)
}
