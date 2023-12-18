package handlers

import (
	"context"
	"labs/service-mesh/dashboard/helpers/clients"
	"labs/service-mesh/dashboard/helpers/factory"
	custcon "labs/service-mesh/helper/concurrent"
	"labs/service-mesh/helper/logger"
	modelsdash "labs/service-mesh/helper/models/dashboard"
	modelslocations "labs/service-mesh/helper/models/locations"
	modelsweather "labs/service-mesh/helper/models/weather"

	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type Handlers struct {
	ls *clients.LocationsServiceClient
	ws *clients.WeatherServiceClient

	cc *ants.Pool
}

func NewHandlers() *Handlers {
	return &Handlers{
		ls: factory.Locations(),
		ws: factory.Weather(),

		cc: custcon.New(100),
	}
}

func (h *Handlers) GetDashboardData(ctx context.Context, req *modelsdash.GetDashboardDataRequest) (*modelsdash.GetDashboardDataResponse, error) {
	resp := modelsdash.GetDashboardDataResponse{
		IP: req.IP,
	}

	locationResp, err := h.ls.GetCurrentLocation(ctx, &modelslocations.GetCurrentLocationRequest{
		IP: req.IP,
	})
	if err != nil {
		logger.SInfo("GetDashboardData.GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}
	resp.Location = locationResp

	city := locationResp.City

	coordinatesResp, err := h.ws.GetLocation(ctx, &modelsweather.GetLocationRequest{
		City: city,
	})
	if err != nil {
		logger.SError("GetDashboardData.GetLocation",
			zap.Error(err))
		return nil, err
	}

	timeRespFunc := func() error {
		timeResp, err := h.ls.GetCurrentTime(ctx, &modelslocations.GetCurrentTimeRequest{
			IP: req.IP,
		})
		if err != nil {
			logger.SError("GetDashboardData.GetCurrentTime",
				zap.Error(err))
			return err
		}
		resp.CurrentTime = timeResp
		return nil
	}

	weatherRespFunc := func() error {
		weatherResp, err := h.ws.GetCurrentWeather(ctx, &modelsweather.GetCurrentWeatherRequest{
			Latitude:  coordinatesResp.Latitude,
			Longitude: coordinatesResp.Longitude,
		})
		if err != nil {
			logger.SError("GetDashboardData.GetCurrentWeather",
				zap.Error(err))
			return err
		}
		resp.Weather = weatherResp
		return nil
	}

	err = custcon.Do(timeRespFunc, weatherRespFunc)
	if err != nil {
		return nil, err
	}

	logger.SInfo("GetDashboardData: Response",
		zap.Any("response", resp))
	return &resp, nil
}
