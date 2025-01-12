package integration_tests

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
			Name:            "api-gateway",
			Image:           "story-pulse-api-gateway",
			Address:         "http://localhost:8000",
			WebPort:         "8000",
			GrpcPort:        "8001",
			GracefulTimeout: defaultGracefulTimeout,
		},
		ConsulConfig: &ConsulConfig{
			Name:    "consul",
			Image:   "consul:1.15",
			Address: "consul:8500",
			APIPort: "8500",
		},
		UsersServiceCfg: &UsersServiceConfig{
			Name:            "users-service",
			Address:         "users-service",
			Image:           "story-pulse-users-service",
			WebPort:         "8030",
			GrpcPort:        "8031",
			GracefulTimeout: defaultGracefulTimeout,

			PostgresName:     "users-service-postgres",
			PostgresImage:    "postgres:17.2-alpine",
			PostgresUsername: "user",
			PostgresPassword: "pass",
			PostgresDB:       "users",
			PostgresPort:     "5432",
		},
		AuthService: &AuthService{
			Name:            "auth-service",
			Address:         "auth-service",
			Image:           "story-pulse-auth-service",
			WebPort:         "8020",
			GrpcPort:        "8021",
			GracefulTimeout: defaultGracefulTimeout,

			RedisName:  "auth-service-redis",
			RedisImage: "redis:7.4-alpine",
			RedisPort:  "6379",
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
	PostgresName     string
	PostgresImage    string
	PostgresUsername string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
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
	RedisName  string
	RedisImage string
	RedisPort  string
}
