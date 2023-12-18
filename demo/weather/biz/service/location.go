package service

import (
	"context"
	"fmt"
	"labs/service-mesh/helper/cache"
	custcon "labs/service-mesh/helper/concurrent"
	custdb "labs/service-mesh/helper/db"
	custerror "labs/service-mesh/helper/error"
	"labs/service-mesh/helper/logger"
	"labs/service-mesh/weather/helpers/factory"
	"labs/service-mesh/weather/helpers/openweather"
	modelsweather "labs/service-mesh/helper/models/weather"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type LocationService struct {
	db     *custdb.LayeredDb
	cache  *ristretto.Cache
	client *openweather.OpenWeatherClient
	cc     *ants.Pool
}

func NewLocationService() *LocationService {
	return &LocationService{
		db:     custdb.Layered(),
		client: factory.OpenWeather(),
		cache:  cache.Cache(),
		cc:     custcon.New(100),
	}
}

func (s *LocationService) GetLocation(ctx context.Context, req *modelsweather.GetLocationRequest) (openweather.GetLocationResponse, error) {
	cacheData, found := s.cache.Get(fmt.Sprintf("GetLocation-%s", req.City))
	if found {
		result, yes := cacheData.(openweather.GetLocationResponse)
		if yes {
			return result, nil
		}
	}

	resp, err := s.client.GetLocation(ctx, &openweather.GetLocationRequest{
		City: req.City,
	})
	if err != nil {
		logger.SInfo("GetLocation",
			zap.Error(err))
		return nil, err
	}

	if len(resp) == 0 {
		return nil, custerror.ErrorNotFound
	}

	s.cc.Submit(func() {
		success := s.cache.SetWithTTL(fmt.Sprintf("GetLocation-%s", req.City), resp, 100, time.Hour*5)
		if !success {
			logger.SInfo("GetLocation: cache.SetWithTTL not success")
		}
	})

	return resp, nil
}
