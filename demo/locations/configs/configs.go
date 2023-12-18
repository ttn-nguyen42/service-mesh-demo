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

	WorldTimeApi WorldTimeApi `json:"worldTimeApi,omitempty" yaml:"worldTimeApi,omitempty"`
	IpApi        IpApi        `json:"ipApi,omitempty" yaml:"ipApi,omitempty"`
}

type WorldTimeApi struct {
	BaseUrl string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
}

type IpApi struct {
	BaseUrl string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
}

func (c ServiceConfigs) GetWorldTimeApi() *WorldTimeApi {
	return &c.WorldTimeApi
}

func Get() *ServiceConfigs {
	return c
}

func (c ServiceConfigs) String() string {
	cbytes, _ := sonic.Marshal(c)
	return string(cbytes)
}
