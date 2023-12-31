package factory

import (
	"context"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"sync"
)

var once sync.Once

var (
	citiesApiClient *dtclient.CitiesApiClient
)

func Init(ctx context.Context) {
	once.Do(func() {
		citiesApiClient = dtclient.NewCitiesApiClient(ctx)
	})
}

func CitiesApiClient() *dtclient.CitiesApiClient {
	return citiesApiClient
}
