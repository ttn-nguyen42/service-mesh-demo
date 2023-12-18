package openweather

import (
	"context"
	"fmt"
	custerror "labs/service-mesh/helper/error"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	"labs/service-mesh/weather/configs"
	"net/url"

	fastshot "github.com/opus-domini/fast-shot"
	"go.uber.org/zap"
)

type OpenWeatherClient struct {
	httpClient *fastshot.Client
}

type GetLocationRequest struct {
	City        string `json:"city"`
	StateCode   string `json:"state_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

type GetLocationResponse []Location

type Location struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names,omitempty"`
	Latitude   float32           `json:"lat"`
	Longitude  float32           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state,omitempty"`
}

type GetCurrentWeatherRequest struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"longitude"`
}

func (c OpenWeatherClient) GetLocation(ctx context.Context, req *GetLocationRequest) (GetLocationResponse, error) {
	uri, err := c.buildGetLocationPath(req)
	if err != nil {
		return nil, err
	}

	reqBuilder := c.httpClient.
		GET(uri).
		Context().Set(ctx)

	resp, err := reqBuilder.Send()
	if err != nil {
		logger.SInfo("GetLocation.Send", zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var result GetLocationResponse
	if err := custhttp.JSONResponse(&resp, &result); err != nil {
		logger.SInfo("GetCurrentTime.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (c OpenWeatherClient) GetCurrentWeather(ctx context.Context, req *GetCurrentWeatherRequest) (*WeatherData, error) {
	uri, _ := c.buildGetCurrentWeatherPath(req)

	reqBuilder := c.httpClient.
		GET(uri).
		Context().Set(ctx)

	resp, err := reqBuilder.Send()
	if err != nil {
		logger.SInfo("GetCurrentWeather.Send", zap.Error(err))
		return nil, err
	}

	if resp.Is4xxClientError() {
		return nil, custerror.ErrorInvalidArgument
	}

	if resp.Is5xxServerError() {
		return nil, custerror.ErrorInternal
	}

	var result WeatherData
	if err := custhttp.JSONResponse(&resp, &result); err != nil {
		logger.SInfo("GetCurrentWeather.JSONResponse",
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (c OpenWeatherClient) buildGetLocationPath(req *GetLocationRequest) (string, error) {
	u, _ := url.ParseRequestURI("/geo/1.0/direct")
	if req.City == "" {
		return "", custerror.FormatInvalidArgument("city name not included")
	}
	q := u.Query()
	q.Add("q", req.City)
	if req.StateCode != "" {
		q.Add("q", req.StateCode)
	}
	if req.CountryCode != "" {
		q.Add("q", req.CountryCode)
	}
	q.Add("limit", "1")
	q.Add("appid",
		configs.Get().
			GetOpenWeatherMap().
			ApiKey)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c OpenWeatherClient) buildGetCurrentWeatherPath(req *GetCurrentWeatherRequest) (string, error) {
	u, _ := url.ParseRequestURI("/data/2.5/weather")

	q := u.Query()
	q.Add("lat", fmt.Sprintf("%f", req.Latitude))
	q.Add("lon", fmt.Sprintf("%f", req.Longitude))
	q.Add("limit", "1")
	q.Add("units", "metric")
	q.Add("appid",
		configs.Get().
			GetOpenWeatherMap().
			ApiKey)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
