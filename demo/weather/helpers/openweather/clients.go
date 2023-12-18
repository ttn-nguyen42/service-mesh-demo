package openweather

import (
	"context"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/weather/configs"
	"time"
)

func NewDatetimeApiClient(ctx context.Context) *OpenWeatherClient {
	serviceConfigs := configs.Get()
	return &OpenWeatherClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.WithBaseUrl(serviceConfigs.
				GetOpenWeatherMap().
				BaseUrl),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}
