package integration_tests

import "testing"

func TestEcosystem(t *testing.T) {
	cfg := NewTestConfig()
	prepareInfrastructure(t, cfg, runServer)
}
