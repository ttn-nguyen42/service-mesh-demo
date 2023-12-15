package handlers

import (
	"context"
	"labs/service-mesh/helper/logger"
	modelslocations "labs/service-mesh/helper/models/locations"
	"labs/service-mesh/locations/biz/service"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"labs/service-mesh/locations/helpers/utils"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type LocationHandlers struct {
	dts *service.DatetimeService
	ls  *service.LocationsService
}

func NewLocationHandlers(ctx context.Context) *LocationHandlers {
	return &LocationHandlers{
		dts: service.GetDatetimeService(),
		ls:  service.GetLocationService(),
	}
}

func (h *LocationHandlers) GetListAreas(ctx context.Context) (*modelslocations.GetListAreasResponse, error) {
	tzs, err := h.dts.GetListAreas(ctx)
	if err != nil {
		return nil, err
	}

	var filteredTimezones []dtclient.Timezone
	for _, tz := range tzs {
		if tz.Area == "" {
			continue
		}
		if tz.Area == "Etc" {
			continue
		}
		newTz := utils.TimezoneToReadable(tz)
		filteredTimezones = append(filteredTimezones, newTz)
	}

	tz := []modelslocations.Timezone{}
	if err := copier.Copy(&tz, &filteredTimezones); err != nil {
		logger.SInfo("GetListAreas: Copy",
			zap.Error(err))
		return nil, err
	}

	return &modelslocations.GetListAreasResponse{
		Areas: tz,
		Total: uint32(len(filteredTimezones)),
	}, nil
}

func (h *LocationHandlers) GetCurrentLocation(ctx context.Context, req *modelslocations.GetCurrentLocationRequest) (*modelslocations.GetCurrentLocationResponse, error) {
	resp, err := h.ls.GetCurrentLocation(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *LocationHandlers) GetCurrentTime(ctx context.Context, req *modelslocations.GetCurrentTimeRequest) (*modelslocations.GetCurrentTimeResponse, error) {
	resp, err := h.dts.GetCurrentTime(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
