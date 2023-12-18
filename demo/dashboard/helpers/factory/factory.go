package factory

import (
	"context"
	"labs/service-mesh/dashboard/helpers/clients"
	"sync"
)

var once sync.Once

var (
	locations *clients.LocationsServiceClient
	weather   *clients.WeatherServiceClient
)

func Init(ctx context.Context) {
	once.Do(func() {
		locations = clients.NewLocationsServiceClient(ctx)
		weather = clients.NewWeatherServiceClient(ctx)
	})
}

func Locations() *clients.LocationsServiceClient {
	return locations
}

func Weather() *clients.WeatherServiceClient {
	return weather
}
