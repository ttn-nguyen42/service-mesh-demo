package service

import (
	"context"
	"fmt"
	"labs/service-mesh/helper/cache"
	custcon "labs/service-mesh/helper/concurrent"
	modelslocations "labs/service-mesh/helper/models/locations"
	custdb "labs/service-mesh/helper/db"
	"labs/service-mesh/helper/logger"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"labs/service-mesh/locations/helpers/factory"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type LocationsService struct {
	db    *custdb.LayeredDb
	cache *ristretto.Cache
	lc    *dtclient.LocationApiClient
	cc    *ants.Pool
}

func NewLocationsService() *LocationsService {
	return &LocationsService{
		db:    custdb.Layered(),
		cache: cache.Cache(),
		lc:    factory.LocationsClient(),
		cc:    custcon.New(100),
	}
}

func (s *LocationsService) GetCurrentLocation(ctx context.Context, req *modelslocations.GetCurrentLocationRequest) (*modelslocations.GetCurrentLocationResponse, error) {
	cacheKey := fmt.Sprintf("GetCurrentLocation-%s", req.IP)
	cacheResult, found := s.cache.Get(cacheKey)
	if found {
		resp, yes := cacheResult.(modelslocations.GetCurrentLocationResponse)
		if yes {
			logger.SInfo("GetCurrentLocation cache hit",
				zap.Any("req", req))
			return &resp, nil
		}
	}

	resp, err := s.lc.GetCurrentLocation(ctx, &dtclient.GetCurrentLocationRequest{
		IP: req.IP,
	})
	if err != nil {
		logger.SInfo("GetCurrentLocation: lc.GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	md := modelslocations.GetCurrentLocationResponse{}
	if err := copier.Copy(&md, resp); err != nil {
		logger.SError("GetCurrentLocation: Copy",
			zap.Error(err))
		return nil, err
	}

	s.cc.Submit(func() {
		good := s.cache.SetWithTTL(cacheKey, md, 10, time.Hour*1)
		if !good {
			logger.SInfo("GetCurrentLocation: SetWithTTL failed")
		}
	})

	return &md, nil
}
