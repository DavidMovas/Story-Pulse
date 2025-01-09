package integration_tests

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func prepareInfrastructure(t *testing.T, cfg *TestConfig, runFunc func(t *testing.T, cfg *TestConfig)) {
	var cleanUpFuncs []func(context.Context) error
	defer cleanUp(t, cleanUpFuncs...)

	// Consul container
	consul, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         cfg.ConsulConfig.Name,
			Image:        cfg.ConsulConfig.Image,
			ExposedPorts: []string{fmt.Sprintf("%s:%s", cfg.ConsulConfig.APIPort, cfg.ConsulConfig.APIPort)},
			Networks:     []string{cfg.Network},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, consul.Terminate)

	// Users service Postgres container
	usersServicePostgres, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  cfg.UsersServiceCfg.PostgresName,
			Image: cfg.UsersServiceCfg.PostgresImage,
			Env: map[string]string{
				"POSTGRES_USER":     cfg.UsersServiceCfg.PostgresUsername,
				"POSTGRES_PASSWORD": cfg.UsersServiceCfg.PostgresPassword,
				"POSTGRES_DB":       cfg.UsersServiceCfg.PostgresDB,
			},
			ExposedPorts: []string{
				fmt.Sprintf("%s:%s", cfg.UsersServiceCfg.PostgresPort, cfg.UsersServiceCfg.PostgresPort),
			},
			Networks:   []string{cfg.Network},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, usersServicePostgres.Terminate)
	postgresMappedPort, err := usersServicePostgres.MappedPort(context.Background(), "5432")
	require.NoError(t, err)
	cfg.UsersServiceCfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.UsersServiceCfg.PostgresUsername, cfg.UsersServiceCfg.PostgresPassword, "localhost", postgresMappedPort.Int(), cfg.UsersServiceCfg.PostgresDB)

	time.Sleep(time.Second * 2)
	runMigrations(t, cfg.UsersServiceCfg.DatabaseURL)

	// Users service container
	usersService, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  cfg.UsersServiceCfg.Name,
			Image: cfg.UsersServiceCfg.Image,
			Env: map[string]string{
				"NAME":             cfg.UsersServiceCfg.Name,
				"ADDRESS":          cfg.UsersServiceCfg.Address,
				"GRPC_PORT":        cfg.UsersServiceCfg.GrpcPort,
				"GRACEFUL_TIMEOUT": cfg.UsersServiceCfg.GracefulTimeout,
				"CONSUL_ADDRESS":   cfg.ConsulConfig.Address,
				"DATABASE_URL":     cfg.UsersServiceCfg.DatabaseURL,
			},
			Networks: []string{cfg.Network},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, usersService.Terminate)

	// Auth service Redis container
	authServiceRedis, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         cfg.AuthService.RedisName,
			Image:        cfg.AuthService.RedisImage,
			ExposedPorts: []string{fmt.Sprintf("%s/tcp", cfg.AuthService.RedisPort)},
			Networks:     []string{cfg.Network},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, authServiceRedis.Terminate)
	redisMappedPort, err := authServiceRedis.MappedPort(context.Background(), "6379")
	require.NoError(t, err)

	cfg.AuthService.RedisURL = fmt.Sprintf("localhost:%d", redisMappedPort.Int())

	// Auth service container
	authService, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  cfg.AuthService.Name,
			Image: cfg.AuthService.Image,
			Env: map[string]string{
				"NAME":             cfg.AuthService.Name,
				"ADDRESS":          cfg.AuthService.Address,
				"GRPC_PORT":        cfg.AuthService.GrpcPort,
				"GRACEFUL_TIMEOUT": cfg.AuthService.GracefulTimeout,
				"CONSUL_ADDRESS":   cfg.ConsulConfig.Address,
				"REDIS_URL":        cfg.AuthService.RedisURL,
			},
			Networks: []string{cfg.Network},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, authService.Terminate)

	// API Gateway container
	gateway, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  cfg.GatewayConfig.Name,
			Image: cfg.GatewayConfig.Image,
			Env: map[string]string{
				"PORT":               cfg.GatewayConfig.WebPort,
				"GRPC_PORT":          cfg.GatewayConfig.GrpcPort,
				"GRACEFUL_TIMEOUT":   cfg.GatewayConfig.GracefulTimeout,
				"USERS_SERVICE_PATH": cfg.GatewayConfig.UsersServicePath,
				"AUTH_SERVICE_PATH":  cfg.GatewayConfig.AuthServicePath,
			},
			ExposedPorts: []string{
				fmt.Sprintf("%s:%s", cfg.GatewayConfig.WebPort, cfg.GatewayConfig.WebPort),
			},
			Networks: []string{cfg.Network},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, gateway.Terminate)

	gatewayMappedPort, err := gateway.MappedPort(context.Background(), nat.Port(cfg.GatewayConfig.WebPort))
	require.NoError(t, err)
	cfg.GatewayConfig.WebPort = gatewayMappedPort.Port()

	time.Sleep(time.Second * 5)
	runFunc(t, cfg)
}

func runMigrations(t *testing.T, pgConnString string) {
	conn, err := pgx.Connect(context.Background(), pgConnString)
	require.NoError(t, err)

	migrator, err := migrate.NewMigrator(context.Background(), conn, "migrations")
	require.NoError(t, err)

	err = migrator.LoadMigrations("../../scripts/tern/users_migrations")
	require.NoError(t, err)

	err = migrator.Migrate(context.Background())
	require.NoError(t, err)
}

func cleanUp(t *testing.T, terminateFuncs ...func(ctx context.Context) error) {
	for _, terminate := range terminateFuncs {
		require.NoError(t, terminate(context.Background()))
	}
}
