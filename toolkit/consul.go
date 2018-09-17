package toolkit

import (
	"os"

	consul "github.com/hashicorp/consul/api"
)

// GetConsulValueFromKey will return a value of given key
func GetConsulValueFromKey(c *consul.Client, k string) (string, error) {
	pair, _, err := c.KV().Get(k, nil)
	if err != nil {
		return "", err
	}

	return string(pair.Value), nil
}

// GetServiceAddress return ip address, port of given service
func GetServiceAddress(c *consul.Client, svc string) (string, int, error) {
	addr, _, err := c.Catalog().Service(svc, "", nil)
	if err != nil {
		return "", 0, err
	}

	return addr[0].ServiceAddress, addr[0].ServicePort, nil
}

// RegisterService register a service to consul
func RegisterService(c *consul.Client, name string, port int) error {
	address := os.Getenv("SERVICE_ADDRESS")
	if address == "" {
		address = os.Getenv("PRIVATE_IP")
	}
	return c.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Name:    name,
		Port:    port,
		Address: address,
	})
}
