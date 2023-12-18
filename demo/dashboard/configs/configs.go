package configs

import (
	"labs/service-mesh/helper/configs"

	"github.com/bytedance/sonic"
)

var (
	c *ServiceConfigs = &ServiceConfigs{}
)

type ServiceConfigs struct {
	custconfigs.Configs

	Services Microservices `json:"services,omitempty" yaml:"services,omitempty"`
}

type Microservices struct {
	Weather   string `json:"weather,omitempty" yaml:"weather,omitempty"`
	Locations string `json:"locations,omitempty" yaml:"locations,omitempty"`
}

func (s ServiceConfigs) GetServices() *Microservices {
	return &c.Services
}

func Get() *ServiceConfigs {
	return c
}

func (c ServiceConfigs) String() string {
	cbytes, _ := sonic.Marshal(c)
	return string(cbytes)
}
