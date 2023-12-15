package factory

import (
	"context"
	"labs/service-mesh/weather/helpers/openweather"
	"sync"
)

var once sync.Once

var (
	openWeatherClient *openweather.OpenWeatherClient
)

func Init(ctx context.Context) {
	once.Do(func() {
		openWeatherClient = openweather.NewDatetimeApiClient(ctx)
	})
}

func OpenWeather() *openweather.OpenWeatherClient {
	return openWeatherClient
}
