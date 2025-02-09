package containers

import (
	. "brain-wave/tests/integration_tests/config"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
)

func NewConsulContainer(cfg *TestConfig) testcontainers.GenericContainerRequest {
	return testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:           cfg.ConsulConfig.Name,
			Image:          cfg.ConsulConfig.Image,
			ExposedPorts:   []string{fmt.Sprintf("%s:%s", cfg.ConsulConfig.APIPort, "8500")},
			Networks:       []string{cfg.Network},
			NetworkAliases: map[string][]string{cfg.Network: {"consul"}},
		},
		Started: true,
	}
}
