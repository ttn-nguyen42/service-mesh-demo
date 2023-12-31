package service

import (
	"context"
	"fmt"
	"labs/service-mesh/helper/cache"
	custcon "labs/service-mesh/helper/concurrent"
	custdb "labs/service-mesh/helper/db"
	"labs/service-mesh/helper/logger"
	modelscities "labs/service-mesh/helper/models/cities"
	dtclient "labs/service-mesh/locations/helpers/clients"
	"labs/service-mesh/locations/helpers/factory"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/jinzhu/copier"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type CitiesService struct {
	db    *custdb.LayeredDb
	cache *ristretto.Cache
	lc    *dtclient.CitiesApiClient
	cc    *ants.Pool
}

func NewCitiesService() *CitiesService {
	return &CitiesService{
		db:    custdb.Layered(),
		cache: cache.Cache(),
		lc:    factory.CitiesApiClient(),
		cc:    custcon.New(100),
	}
}

func (s *CitiesService) GetCountries(ctx context.Context) (modelscities.GetListCountriesResponse, error) {
	cacheKey := fmt.Sprintf("GetCountries")
	cacheResult, found := s.cache.Get(cacheKey)
	if found {
		resp, yes := cacheResult.(modelscities.GetListCountriesResponse)
		if yes {
			logger.SInfo("GetCountries cache hit")
			return resp, nil
		}
	}

	resp, err := s.lc.GetCountries(ctx)
	if err != nil {
		logger.SInfo("GetCountries: lc.GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	md := modelscities.GetListCountriesResponse{}
	if err := copier.Copy(&md, resp); err != nil {
		logger.SError("GetCountries: Copy",
			zap.Error(err))
		return nil, err
	}

	s.cc.Submit(func() {
		good := s.cache.SetWithTTL(cacheKey, md, 10, time.Hour*1)
		if !good {
			logger.SInfo("GetCountries: SetWithTTL failed")
		}
	})

	return md, nil
}

func (s *CitiesService) GetCities(ctx context.Context, req *modelscities.GetCitiesRequest) (modelscities.GetListCitiesResponse, error) {
	cacheKey := fmt.Sprintf("GetCities-%s", req.Iso2)
	cacheResult, found := s.cache.Get(cacheKey)
	if found {
		resp, yes := cacheResult.(modelscities.GetListCitiesResponse)
		if yes {
			logger.SInfo("GetCities cache hit")
			return resp, nil
		}
	}

	resp, err := s.lc.GetCities(ctx, &dtclient.GetCitiesRequest{
		Iso2: req.Iso2,
	})
	if err != nil {
		logger.SInfo("GetCities: lc.GetCurrentLocation",
			zap.Error(err))
		return nil, err
	}

	md := modelscities.GetListCitiesResponse{}
	if err := copier.Copy(&md, resp); err != nil {
		logger.SError("GetCities: Copy",
			zap.Error(err))
		return nil, err
	}

	s.cc.Submit(func() {
		good := s.cache.SetWithTTL(cacheKey, md, 10, time.Hour*1)
		if !good {
			logger.SInfo("GetCities: SetWithTTL failed")
		}
	})

	return md, nil
}
