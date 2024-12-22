package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/labstack/gommon/log"
)

func RegisterService(client *api.Client, name, address, tag string, port int) error {
	registration := &api.AgentServiceRegistration{
		ID:      name + "-" + address,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    []string{tag},
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	log.Infof("CONSUL: Registered service %s with tag %s on port %d", name, address, port)

	return nil
}
