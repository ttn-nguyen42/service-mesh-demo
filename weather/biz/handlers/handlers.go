package handlers

import (
	"context"
	"labs/service-mesh/helper/logger"
	modelsweather "labs/service-mesh/helper/models/weather"
	"labs/service-mesh/weather/biz/service"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type Handlers struct {
	locationService *service.LocationService
	weatherService  *service.WeatherService
}

func NewHandlers() *Handlers {
	return &Handlers{
		locationService: service.GetLocationService(),
		weatherService:  service.GetWeatherService(),
	}
}

func (h *Handlers) GetLocation(ctx context.Context, req *modelsweather.GetLocationRequest) (*modelsweather.GetLocationResponse, error) {
	clientResp, err := h.locationService.GetLocation(ctx, req)
	if err != nil {
		logger.SInfo("GetLocation: locationService.GetLocation",
			zap.Error(err))
		return nil, err
	}

	location := clientResp[0]
	resp := &modelsweather.GetLocationResponse{
		Name:      location.Name,
		Country:   location.Country,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}

	return resp, nil
}

func (h *Handlers) GetCurrentWeather(ctx context.Context, req *modelsweather.GetCurrentWeatherRequest) (*modelsweather.GetCurrentWeatherResponse, error) {
	clientResp, err := h.weatherService.GetCurrentWeather(ctx, req)
	if err != nil {
		logger.SInfo("GetCurrentWeather: weatherService.GetCurrentWeather",
			zap.Error(err))
		return nil, err
	}

	resp := &modelsweather.GetCurrentWeatherResponse{}

	if err := copier.Copy(resp, clientResp); err != nil {
		logger.SInfo("GetCurrentWeather: Copy",
			zap.Error(err))
		return nil, err
	}

	return resp, nil
}
