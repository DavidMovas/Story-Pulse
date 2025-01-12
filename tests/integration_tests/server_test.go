package integration_tests

import (
	"context"
	"story-pulse/client"
	"testing"
)

func TestEcosystem(t *testing.T) {
	cfg := NewTestConfig()
	ctx := context.Background()

	prepareInfrastructure(t, ctx, cfg, runServer)
}

func runServer(t *testing.T, cfg *TestConfig) {
	c := client.NewClient(cfg.GatewayConfig.Address)

	usersServiceTest(t, c, cfg)
}
