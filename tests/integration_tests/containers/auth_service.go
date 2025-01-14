package containers

import (
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"story-pulse/tests/integration_tests/config"
)

func NewAuthServiceRedis(cfg *config.TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:           cfg.AuthService.RedisName,
			Image:          cfg.AuthService.RedisImage,
			ExposedPorts:   []string{fmt.Sprintf("%s/tcp", cfg.AuthService.RedisPort)},
			WaitingFor:     wait.ForLog("Ready to accept connections"),
			Networks:       []string{cfg.Network},
			NetworkAliases: map[string][]string{cfg.Network: {cfg.AuthService.RedisNetworkAlias}},
		},
		Started: true,
	}
}

func BuildRedisURL(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func NewAuthService(cfg *config.TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: cfg.AuthService.Name,
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../../.",
				Dockerfile: "deployments/docker-tests/auth-service.dockerfile",
				KeepImage:  true,
			},
			Env: map[string]string{
				"NAME":               cfg.AuthService.Name,
				"ADDRESS":            cfg.AuthService.Address,
				"PORT":               cfg.AuthService.WebPort,
				"GRPC_PORT":          cfg.AuthService.GrpcPort,
				"GRACEFUL_TIMEOUT":   cfg.AuthService.GracefulTimeout,
				"CONSUL_ADDRESS":     cfg.ConsulConfig.Address,
				"USERS_SERVICE_PATH": cfg.AuthService.UsersServicePath,
				"REDIS_URL":          cfg.AuthService.RedisURL,
			},
			ExposedPorts: []string{cfg.AuthService.WebPort},
			Networks:     []string{cfg.Network},
		},
		Started: true,
	}
}
