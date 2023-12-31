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
	CitiesApi CitiesApiConfigs `json:"citiesApi" yaml:"citiesApi"`
}

type CitiesApiConfigs struct {
	BaseUrl string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
	ApiKey  string `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
}

func (c ServiceConfigs) GetCitiesApi() *CitiesApiConfigs {
	return &c.CitiesApi
}

func Get() *ServiceConfigs {
	return c
}

func (c ServiceConfigs) String() string {
	cbytes, _ := sonic.Marshal(c)
	return string(cbytes)
}
