package containers

import (
	"brain-wave/tests/integration_tests/config"
	"github.com/testcontainers/testcontainers-go"
)

func NewGateway(cfg *config.TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: cfg.GatewayConfig.Name,
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../../.",
				Dockerfile: "deployments/docker-tests/api-gateway.dockerfile",
				KeepImage:  true,
			},
			Env: map[string]string{
				"PORT":               cfg.GatewayConfig.WebPort,
				"GRPC_PORT":          cfg.GatewayConfig.GrpcPort,
				"GRACEFUL_TIMEOUT":   cfg.GatewayConfig.GracefulTimeout,
				"USERS_SERVICE_PATH": cfg.GatewayConfig.UsersServicePath,
				"AUTH_SERVICE_PATH":  cfg.GatewayConfig.AuthServicePath,
			},
			ExposedPorts: []string{cfg.GatewayConfig.WebPort},
			Networks:     []string{cfg.Network},
		},
		Started: true,
	}
}
