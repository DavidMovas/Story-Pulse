package containers

import (
	"brain-wave/tests/integration_tests/config"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewUsersServicePostgres(cfg *config.TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  cfg.UsersServiceCfg.PostgresName,
			Image: cfg.UsersServiceCfg.PostgresImage,
			Env: map[string]string{
				"POSTGRES_USER":     cfg.UsersServiceCfg.PostgresUsername,
				"POSTGRES_PASSWORD": cfg.UsersServiceCfg.PostgresPassword,
				"POSTGRES_DB":       cfg.UsersServiceCfg.PostgresDB,
			},
			ExposedPorts:   []string{fmt.Sprintf("%s/tcp", cfg.UsersServiceCfg.PostgresPort)},
			WaitingFor:     wait.ForLog("database system is ready to accept connections"),
			Networks:       []string{cfg.Network},
			NetworkAliases: map[string][]string{cfg.Network: {cfg.UsersServiceCfg.PostgresNetworkAlias}},
		},
		Started: true,
	}
}

func BuildPostgresURL(username, password, host, port, db string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, db)
}

func NewUsersService(cfg *config.TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: cfg.UsersServiceCfg.Name,
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../../.",
				Dockerfile: "deployments/docker-tests/users-service.dockerfile",
				KeepImage:  true,
			},
			Env: map[string]string{
				"NAME":             cfg.UsersServiceCfg.Name,
				"ADDRESS":          cfg.UsersServiceCfg.Address,
				"PORT":             cfg.UsersServiceCfg.WebPort,
				"GRPC_PORT":        cfg.UsersServiceCfg.GrpcPort,
				"GRACEFUL_TIMEOUT": cfg.UsersServiceCfg.GracefulTimeout,
				"CONSUL_ADDRESS":   cfg.ConsulConfig.Address,
				"DATABASE_URL":     cfg.UsersServiceCfg.DatabaseURL,
			},
			ExposedPorts: []string{cfg.UsersServiceCfg.WebPort},
			Networks:     []string{cfg.Network},
		},
		Started: true,
	}
}
