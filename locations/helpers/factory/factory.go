package factory

import (
	"context"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"sync"
)

var once sync.Once

var (
	datetimeApiClient *dtclient.DatetimeApiClient
	locationsClient   *dtclient.LocationApiClient
)

func Init(ctx context.Context) {
	once.Do(func() {
		datetimeApiClient = dtclient.NewDatetimeApiClient(ctx)
		locationsClient = dtclient.NewLocationApiClient(ctx)
	})
}

func DatetimeClient() *dtclient.DatetimeApiClient {
	return datetimeApiClient
}

func LocationsClient() *dtclient.LocationApiClient {
	return locationsClient
}
