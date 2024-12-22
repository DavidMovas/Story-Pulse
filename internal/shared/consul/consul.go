package consul

import (
	"github.com/hashicorp/consul/api"
)

func RegisterService(client *api.Client, name, address, tag string, port int, check *api.AgentServiceCheck) error {
	registration := &api.AgentServiceRegistration{
		ID:      name + "-" + tag,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    []string{tag},
	}

	if check != nil {
		registration.Check = check
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	return nil
}
