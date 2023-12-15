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

	OpenWeatherMap OpenWeatherMapConfigs `json:"owm,omitempty" yaml:"owm,omitempty"`
}

type OpenWeatherMapConfigs struct {
	ApiKey  string `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
	BaseUrl string `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
}

func (c ServiceConfigs) GetOpenWeatherMap() *OpenWeatherMapConfigs {
	return &c.OpenWeatherMap
}

func Get() *ServiceConfigs {
	return c
}

func (c ServiceConfigs) String() string {
	cbytes, _ := sonic.Marshal(c)
	return string(cbytes)
}
