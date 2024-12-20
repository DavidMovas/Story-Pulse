package config

type Config struct {
	WebPort                 int                   `env:"PORT" envDefault:"8010"`
	GRPCPort                int                   `env:"GRPC_PORT" envDefault:"8011"`
	GracefulShutdownTimeout int                   `env:"GRACEFUL_TIMEOUT" envDefault:"10"`
	UsersService            *UsersServiceConfig   `envPrefix:"USERS_SERVICE_"`
	AuthService             *AuthServiceConfig    `envPrefix:"AUTH_SERVICE_"`
	ContentService          *ContentServiceConfig `envPrefix:"CONTENT_SERVICE_"`
	CommentService          *CommentServiceConfig `envPrefix:"COMMENT_SERVICE_"`
	SearchService           *SearchServiceConfig  `envPrefix:"SEARCH_SERVICE_"`
}

type UsersServiceConfig struct {
	ServiceURL  string `env:"URL"`
	ServicePath string `env:"PATH"`
}

type AuthServiceConfig struct {
	ServiceURL  string `env:"URL"`
	ServicePath string `env:"PATH"`
}

type ContentServiceConfig struct {
	ServiceURL  string `env:"URL"`
	ServicePath string `env:"PATH"`
}

type CommentServiceConfig struct {
	ServiceURL  string `env:"URL"`
	ServicePath string `env:"PATH"`
}

type SearchServiceConfig struct {
	ServiceURL  string `env:"URL"`
	ServicePath string `env:"PATH"`
}
