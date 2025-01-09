package config

type Config struct {
	WebPort                 int                  `env:"PORT" envDefault:"8010"`
	GRPCPort                int                  `env:"GRPC_PORT" envDefault:"8011"`
	GracefulShutdownTimeout int                  `env:"GRACEFUL_TIMEOUT" envDefault:"10"`
	UsersService            UsersServiceConfig   `envPrefix:"USERS_SERVICE_"`
	AuthService             AuthServiceConfig    `envPrefix:"AUTH_SERVICE_"`
	ContentService          ContentServiceConfig `envPrefix:"CONTENT_SERVICE_"`
	CommentService          CommentServiceConfig `envPrefix:"COMMENT_SERVICE_"`
	SearchService           SearchServiceConfig  `envPrefix:"SEARCH_SERVICE_"`
}

type UsersServiceConfig struct {
	ServicePath string `env:"PATH" envDefault:"users-service"`
}

type AuthServiceConfig struct {
	ServicePath string `env:"PATH" envDefault:"auth-service"`
}

type ContentServiceConfig struct {
	ServicePath string `env:"PATH"`
}

type CommentServiceConfig struct {
	ServicePath string `env:"PATH"`
}

type SearchServiceConfig struct {
	ServicePath string `env:"PATH"`
}

type ResolverConfig struct {
	ConsulAddress string `env:"CONSUL_ADDRESS" envDefault:"localhost:8500"`
}
