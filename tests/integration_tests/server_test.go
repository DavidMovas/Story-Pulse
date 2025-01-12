package integration_tests

import (
	"context"
	"story-pulse/client"
	"story-pulse/tests/integration_tests/config"
	"testing"
)

func TestEcosystem(t *testing.T) {
	cfg := config.NewTestConfig()
	ctx := context.Background()

	prepareInfrastructure(t, ctx, cfg, runServer)
}

func runServer(t *testing.T, cfg *config.TestConfig) {
	c := client.NewClient(cfg.GatewayConfig.Address)

	usersServiceTest(t, c, cfg)
}
