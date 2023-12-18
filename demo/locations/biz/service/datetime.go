package service

import (
	"context"
	"labs/service-mesh/helper/cache"
	custcon "labs/service-mesh/helper/concurrent"
	custdb "labs/service-mesh/helper/db"
	"labs/service-mesh/helper/logger"
	modelslocations "labs/service-mesh/helper/models/locations"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"labs/service-mesh/locations/helpers/factory"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type DatetimeService struct {
	db    *custdb.LayeredDb
	cache *ristretto.Cache
	dt    *dtclient.DatetimeApiClient
	cc    *ants.Pool
}

func NewDatetimeService() *DatetimeService {
	return &DatetimeService{
		db:    custdb.Layered(),
		cache: cache.Cache(),
		dt:    factory.DatetimeClient(),
		cc:    custcon.New(100),
	}
}

func (s *DatetimeService) GetListAreas(ctx context.Context) (dtclient.Timezones, error) {
	result, cacheFound := s.cache.Get("GetListAreas")
	if cacheFound {
		tzs, ok := result.(dtclient.Timezones)
		if ok {
			return tzs, nil
		}
	}

	tzs, err := s.dt.GetTimezones(ctx, dtclient.GetTimezonesRequest{
		Area: "",
	})
	if err != nil {
		logger.SInfo("GetListAreas: dt.GetTimezones",
			zap.Error(err))
		return nil, err
	}

	s.cc.Submit(func() {
		good := s.cache.SetWithTTL("GetListAreas", tzs, 100, time.Hour*1)
		if !good {
			logger.Error("GetListAreas: cache.SetWithTTL not good")
		}
	})

	return tzs, nil
}

func (s *DatetimeService) GetCurrentTime(ctx context.Context, req *modelslocations.GetCurrentTimeRequest) (*modelslocations.GetCurrentTimeResponse, error) {

	apiResp, err := s.dt.GetCurrentTime(ctx, dtclient.CurrentTimeRequest{
		IP: req.IP,
	})
	if err != nil {
		logger.SInfo("GetCurrentTime: s.dt.GetCurrentTime",
			zap.Error(err))
		return nil, err
	}

	var resp modelslocations.GetCurrentTimeResponse
	if err := copier.Copy(&resp, apiResp); err != nil {
		logger.SError("GetCurrentTime: Copy",
			zap.Error(err))
		return nil, err
	}

	logger.SInfo("GetCurrentTime: Response",
		zap.Any("response", resp))

	return &resp, nil
}
