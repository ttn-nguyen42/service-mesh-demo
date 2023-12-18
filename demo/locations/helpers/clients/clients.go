package dtclient

import (
	"context"
	"labs/service-mesh/locations/configs"
	custhttp "labs/service-mesh/helper/http"
	"time"
)

func NewDatetimeApiClient(ctx context.Context) *DatetimeApiClient {
	serviceConfigs := configs.Get()
	return &DatetimeApiClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.WithBaseUrl(serviceConfigs.
				GetWorldTimeApi().
				BaseUrl),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}

func NewLocationApiClient(ctx context.Context) *LocationApiClient {
	serviceConfigs := configs.Get()
	return &LocationApiClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.WithBaseUrl(serviceConfigs.
				IpApi.BaseUrl),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}
