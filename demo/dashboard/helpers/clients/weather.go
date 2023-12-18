package clients

import (
	"context"
	"fmt"
	"labs/service-mesh/dashboard/configs"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	modelsweather "labs/service-mesh/helper/models/weather"
	"net/url"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type WeatherServiceClient struct {
	httpClient *fastshot.Client
}

func NewWeatherServiceClient(ctx context.Context) *WeatherServiceClient {
	cfg := configs.Get()
	return &WeatherServiceClient{
		httpClient: custhttp.NewHttpClient(
			ctx,
			custhttp.
				WithBaseUrl(cfg.GetServices().Weather),
			custhttp.WithTimeout(time.Second*2),
		),
	}
}

func (c *WeatherServiceClient) GetLocation(ctx context.Context, req *modelsweather.GetLocationRequest) (*modelsweather.GetLocationResponse, error) {
	path, err := c.buildGetCurrentLocationPath(req)
	if err != nil {
		logger.SDebug("GetLocation: buildGetCurrentLocationPath",
			zap.Error(err))
		return nil, err
	}

	apiResp, err := c.httpClient.GET(path).
		Context().Set(ctx).
		Send()

	if err != nil {
		logger.SDebug("GetLocation: GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	if apiResp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if apiResp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var resp modelsweather.GetLocationResponse
	if err := custhttp.JSONResponse(&apiResp, &resp); err != nil {
		logger.SDebug("GetCurrentLocation: JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &resp, nil
}

func (c *WeatherServiceClient) buildGetCurrentLocationPath(req *modelsweather.GetLocationRequest) (string, error) {
	if req.City == "" {
		return "", custerror.ErrorInvalidArgument
	}

	p, _ := url.ParseRequestURI("/api/location")
	q := p.Query()
	q.Add("city", req.City)
	p.RawQuery = q.Encode()

	return p.String(), nil
}

func (c *WeatherServiceClient) GetCurrentWeather(ctx context.Context, req *modelsweather.GetCurrentWeatherRequest) (*modelsweather.GetCurrentWeatherResponse, error) {
	path, err := c.buildGetCurrentWeatherPath(req)
	if err != nil {
		logger.SDebug("GetCurrentWeather: buildGetCurrentWeatherPath",
			zap.Error(err))
		return nil, err
	}

	apiResp, err := c.httpClient.GET(path).
		Context().Set(ctx).
		Send()

	if err != nil {
		logger.SDebug("GetCurrentWeather: GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	if apiResp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if apiResp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var resp modelsweather.GetCurrentWeatherResponse
	if err := custhttp.JSONResponse(&apiResp, &resp); err != nil {
		logger.SDebug("GetCurrentWeather: JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &resp, nil
}

func (c *WeatherServiceClient) buildGetCurrentWeatherPath(req *modelsweather.GetCurrentWeatherRequest) (string, error) {
	p, _ := url.ParseRequestURI("/api/weather")
	q := p.Query()
	q.Add("latitude", fmt.Sprintf("%f", req.Latitude))
	q.Add("longitude", fmt.Sprintf("%f", req.Longitude))
	p.RawQuery = q.Encode()

	return p.String(), nil
}
