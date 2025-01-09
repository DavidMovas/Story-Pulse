package integration_tests

import (
	"story-pulse/client"
	"testing"
)

func TestEcosystem(t *testing.T) {
	cfg := NewTestConfig()
	prepareInfrastructure(t, cfg, runServer)
}

func runServer(t *testing.T, cfg *TestConfig) {
	c := client.NewClient(cfg.GatewayConfig.Address)

	usersServiceTest(t, c, cfg)
}
