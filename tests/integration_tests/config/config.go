package config

const (
	defaultGracefulTimeout = "10"
)

type TestConfig struct {
	Network         string
	GatewayConfig   *GatewayConfig
	ConsulConfig    *ConsulConfig
	UsersServiceCfg *UsersServiceConfig
	AuthService     *AuthService
}

func NewTestConfig() *TestConfig {
	return &TestConfig{
		Network: "test-network",
		GatewayConfig: &GatewayConfig{
			Name:            "api-gateway_test",
			Image:           "story-pulse-api-gateway",
			Address:         "http://localhost:8000",
			WebPort:         "9000",
			GrpcPort:        "9001",
			GracefulTimeout: defaultGracefulTimeout,
		},
		ConsulConfig: &ConsulConfig{
			Name:    "consul_test",
			Image:   "consul:1.15",
			Address: "consul:8500",
			APIPort: "9500",
		},
		UsersServiceCfg: &UsersServiceConfig{
			Name:            "users-service_test",
			Address:         "users-service",
			Image:           "story-pulse-users-service",
			WebPort:         "9030",
			GrpcPort:        "9031",
			GracefulTimeout: defaultGracefulTimeout,

			PostgresName:         "users-service-postgres_test",
			PostgresImage:        "postgres:17.2-alpine",
			PostgresUsername:     "user",
			PostgresPassword:     "pass",
			PostgresDB:           "users",
			PostgresPort:         "5432",
			PostgresNetworkAlias: "postgres",
		},
		AuthService: &AuthService{
			Name:            "auth-service_test",
			Address:         "auth-service",
			Image:           "story-pulse-auth-service",
			WebPort:         "9020",
			GrpcPort:        "9021",
			GracefulTimeout: defaultGracefulTimeout,

			RedisName:         "auth-service-redis_test",
			RedisImage:        "redis:7.4-alpine",
			RedisPort:         "6379",
			RedisNetworkAlias: "redis",
		},
	}
}

type GatewayConfig struct {
	Name            string
	Image           string
	Address         string
	WebPort         string
	GrpcPort        string
	GracefulTimeout string

	UsersServicePath string
	AuthServicePath  string
}

type ConsulConfig struct {
	Name    string
	Image   string
	Address string
	APIPort string
}

type UsersServiceConfig struct {
	Name            string
	Address         string
	Image           string
	WebPort         string
	GrpcPort        string
	DatabaseURL     string
	GracefulTimeout string

	// Postgres
	PostgresName         string
	PostgresImage        string
	PostgresUsername     string
	PostgresPassword     string
	PostgresDB           string
	PostgresPort         string
	PostgresNetworkAlias string
}

type AuthService struct {
	Name            string
	Address         string
	Image           string
	WebPort         string
	GrpcPort        string
	RedisURL        string
	GracefulTimeout string

	// Redis
	RedisName         string
	RedisImage        string
	RedisPort         string
	RedisNetworkAlias string
}
