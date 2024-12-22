package config

type Config struct {
	Name                    string `env:"NAME"`
	Address                 string `env:"ADDRESS"`
	Tag                     string `env:"TAG"`
	WebPort                 int    `env:"PORT" envDefault:"8010"`
	GRPCPort                int    `env:"GRPC_PORT" envDefault:"8011"`
	GracefulShutdownTimeout int    `env:"GRACEFUL_TIMEOUT" envDefault:"10"`
	ConsulAddr              string `env:"CONSUL_ADDRESS"`
	DatabaseURL             string `env:"DATABASE_URL"`
}
