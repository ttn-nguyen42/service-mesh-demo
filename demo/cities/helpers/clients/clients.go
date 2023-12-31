package dtclient

import (
	"context"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/locations/configs"
	"time"
)

func NewCitiesApiClient(ctx context.Context) *CitiesApiClient {
	serviceConfigs := configs.Get()
	return &CitiesApiClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.WithBaseUrl(serviceConfigs.
				GetCitiesApi().
				BaseUrl),
			custhttp.WithHeader("X-CSCAPI-KEY", serviceConfigs.
				GetCitiesApi().
				ApiKey),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}
