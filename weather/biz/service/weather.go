package service

import (
	"context"
	"fmt"
	"labs/service-mesh/helper/cache"
	custcon "labs/service-mesh/helper/concurrent"
	custdb "labs/service-mesh/helper/db"
	"labs/service-mesh/helper/logger"
	modelsweather "labs/service-mesh/helper/models/weather"
	"labs/service-mesh/weather/helpers/factory"
	"labs/service-mesh/weather/helpers/openweather"
	"math"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type WeatherService struct {
	db     *custdb.LayeredDb
	cache  *ristretto.Cache
	client *openweather.OpenWeatherClient
	cc     *ants.Pool
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		db:     custdb.Layered(),
		client: factory.OpenWeather(),
		cache:  cache.Cache(),
		cc:     custcon.New(100),
	}
}

func (s *WeatherService) GetCurrentWeather(ctx context.Context, req *modelsweather.GetCurrentWeatherRequest) (*openweather.WeatherData, error) {
	currentTime := time.Now()
	currentHour := currentTime.Hour()
	year, month, date := currentTime.Date()

	cacheKey := fmt.Sprintf("GetCurrentWeather-(%d,%s,%d)-%d-(%f,%f)",
		year,
		month,
		date,
		currentHour,
		math.Round(float64(req.Latitude)),
		math.Round(float64(req.Longitude)))

	cacheResult, found := s.cache.Get(cacheKey)
	if found {
		resp, yes := cacheResult.(openweather.WeatherData)
		if yes {
			logger.SDebug("GetCurrentWeather: cache.Get success",
				zap.String("key", cacheKey))
			return &resp, nil
		}
	}

	weather, err := s.client.GetCurrentWeather(ctx, &openweather.GetCurrentWeatherRequest{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		logger.SInfo("GetCurrentWeather: client.GetCurrentWeather",
			zap.Error(err))
		return nil, err
	}

	s.cc.Submit(func() {
		good := s.cache.SetWithTTL(cacheKey, *weather, 10, time.Hour*1)
		if !good {
			logger.SInfo("GetCurrentTime: cache.SetWithTTL failed")
		}
	})

	return weather, nil
}
