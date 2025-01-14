package integration_tests

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"story-pulse/tests/integration_tests/config"
	"story-pulse/tests/integration_tests/containers"
	"testing"
	"time"
)

func prepareInfrastructure(t *testing.T, ctx context.Context, cfg *config.TestConfig, runFunc func(t *testing.T, cfg *config.TestConfig)) {
	testNetwork, err := network.New(ctx)
	require.NoError(t, err)
	defer cleanUp(t, testNetwork.Remove)
	testcontainers.CleanupNetwork(t, testNetwork)
	cfg.Network = testNetwork.Name

	// Consul container
	consul, err := testcontainers.GenericContainer(ctx, containers.NewConsulContainer(cfg))

	require.NoError(t, err)
	defer cleanUp(t, consul.Terminate)

	// Users service Postgres container
	usersServicePostgres, err := testcontainers.GenericContainer(ctx, containers.NewUsersServicePostgres(cfg))

	require.NoError(t, err)
	defer cleanUp(t, usersServicePostgres.Terminate)
	time.Sleep(time.Second)

	usersServicePostgresPort, err := usersServicePostgres.MappedPort(ctx, nat.Port(cfg.UsersServiceCfg.PostgresPort))
	require.NoError(t, err)

	internalPostgresURL := containers.BuildPostgresURL(
		cfg.UsersServiceCfg.PostgresUsername,
		cfg.UsersServiceCfg.PostgresPassword,
		cfg.UsersServiceCfg.PostgresNetworkAlias,
		cfg.UsersServiceCfg.PostgresPort,
		cfg.UsersServiceCfg.PostgresDB,
	)

	externalPostgresURL := containers.BuildPostgresURL(
		cfg.UsersServiceCfg.PostgresUsername,
		cfg.UsersServiceCfg.PostgresPassword,
		"localhost",
		usersServicePostgresPort.Port(),
		cfg.UsersServiceCfg.PostgresDB,
	)

	cfg.UsersServiceCfg.DatabaseURL = internalPostgresURL
	runMigrations(t, externalPostgresURL)

	// Users service container
	usersService, err := testcontainers.GenericContainer(ctx, containers.NewUsersService(cfg))

	require.NoError(t, err)
	defer cleanUp(t, usersService.Terminate)

	// Auth service Redis container
	authServiceRedis, err := testcontainers.GenericContainer(ctx, containers.NewAuthServiceRedis(cfg))

	require.NoError(t, err)
	defer cleanUp(t, authServiceRedis.Terminate)

	cfg.AuthService.RedisURL = containers.BuildRedisURL(cfg.AuthService.RedisNetworkAlias, cfg.AuthService.RedisPort)

	// Auth service container
	authService, err := testcontainers.GenericContainer(ctx, containers.NewAuthService(cfg))

	require.NoError(t, err)
	defer cleanUp(t, authService.Terminate)

	// API Gateway container
	gateway, err := testcontainers.GenericContainer(ctx, containers.NewGateway(cfg))

	require.NoError(t, err)
	defer cleanUp(t, gateway.Terminate)

	time.Sleep(time.Second)
	gatewayMappedPort, err := gateway.MappedPort(ctx, nat.Port(cfg.GatewayConfig.WebPort))
	require.NoError(t, err)
	cfg.GatewayConfig.WebPort = gatewayMappedPort.Port()
	cfg.GatewayConfig.Address = fmt.Sprintf("http://localhost:%s", gatewayMappedPort.Port())

	time.Sleep(time.Second * 2)
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

func cleanUp(t *testing.T, terminateFunc func(ctx context.Context) error) {
	require.NoError(t, terminateFunc(context.Background()))
}
