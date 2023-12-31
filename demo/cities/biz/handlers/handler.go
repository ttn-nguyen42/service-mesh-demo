package handlers

import (
	"context"
	"go.uber.org/zap"
	"labs/service-mesh/helper/logger"
	modelscities "labs/service-mesh/helper/models/cities"
	"labs/service-mesh/locations/biz/service"
)

type CitiesHandlers struct {
	cs *service.CitiesService
}

func NewCitiesHandlers(ctx context.Context) *CitiesHandlers {
	return &CitiesHandlers{
		cs: service.GetCitiesService(),
	}
}

func (h *CitiesHandlers) GetListCountries(ctx context.Context) (modelscities.GetListCountriesResponse, error) {
	resp, err := h.cs.GetCountries(ctx)
	if err != nil {
		logger.SInfo("CitiesHandlers.GetListCountries: error", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (h *CitiesHandlers) GetCities(ctx context.Context, req *modelscities.GetCitiesRequest) (*modelscities.GetListCitiesResponse, error) {
	resp, err := h.cs.GetCities(ctx, req)
	if err != nil {
		logger.SInfo("CitiesHandlers.GetCities: error", zap.Error(err))
		return nil, err
	}
	return &resp, nil
}
