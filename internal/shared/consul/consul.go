package consul

import (
	"github.com/hashicorp/consul/api"
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

	return nil
}
