package integration_tests

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func prepareInfrastructure(t *testing.T, cfg *TestConfig, runFunc func(t *testing.T, cfg *TestConfig)) {
	var cleanUpFuncs []func(context.Context) error
	defer cleanUp(t, cleanUpFuncs...)

	// Consul
	consul, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         cfg.ConsulConfig.Name,
			Image:        cfg.ConsulConfig.Image,
			ExposedPorts: []string{cfg.ConsulConfig.APIPort},
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
				fmt.Sprintf("%s/tcp", cfg.UsersServiceCfg.PostgresPort),
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, usersServicePostgres.Terminate)
	postgresMappedPort, err := usersServicePostgres.MappedPort(context.Background(), "5432")
	require.NoError(t, err)
	cfg.UsersServiceCfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.UsersServiceCfg.PostgresUsername, cfg.UsersServiceCfg.PostgresPassword, "localhost", postgresMappedPort.Int(), cfg.UsersServiceCfg.PostgresDB)

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
			ExposedPorts: []string{
				cfg.UsersServiceCfg.WebPort,
				cfg.UsersServiceCfg.GrpcPort,
			},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, usersService.Terminate)

	authServiceRedis, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         cfg.AuthService.RedisName,
			Image:        cfg.AuthService.RedisImage,
			ExposedPorts: []string{fmt.Sprintf("%s/tcp", cfg.AuthService.RedisPort)},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, authServiceRedis.Terminate)
	redisMappedPort, err := authServiceRedis.MappedPort(context.Background(), "6379")
	require.NoError(t, err)

	cfg.AuthService.RedisURL = fmt.Sprintf("localhost:%d", redisMappedPort.Int())

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
			ExposedPorts: []string{
				cfg.AuthService.WebPort,
				cfg.AuthService.GrpcPort,
			},
		},
		Started: true,
	})

	require.NoError(t, err)
	cleanUpFuncs = append(cleanUpFuncs, authService.Terminate)
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
