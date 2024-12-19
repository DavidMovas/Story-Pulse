package config

type Config struct {
	WebPort                 int    `env:"PORT" envDefault:"8010"`
	GRPCPort                int    `env:"GRPC_PORT" envDefault:"8011"`
	GracefulShutdownTimeout int    `env:"GRACEFUL_TIMEOUT" envDefault:"10"`
	DatabaseURL             string `env:"DATABASE_URL"`
}
