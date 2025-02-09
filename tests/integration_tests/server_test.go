package integration_tests

import (
	"brain-wave/client"
	"brain-wave/tests/integration_tests/config"
	"brain-wave/tests/integration_tests/modules"
	"context"
	"testing"
)

func TestEcosystem(t *testing.T) {
	cfg := config.NewTestConfig()
	ctx := context.Background()

	prepareInfrastructure(t, ctx, cfg, runServer)
}

func runServer(t *testing.T, cfg *config.TestConfig) {
	c := client.NewClient(cfg.GatewayConfig.Address)

	modules.AuthServiceTest(t, c, cfg)
	modules.UsersServiceTest(t, c, cfg)
}
