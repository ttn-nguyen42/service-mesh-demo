package service

import "sync"

var once sync.Once

var (
	weatherService  *WeatherService
	locationService *LocationService
)

func Init() {
	once.Do(func() {
		weatherService = NewWeatherService()
		locationService = NewLocationService()
	})
}

func GetWeatherService() *WeatherService {
	return weatherService
}

func GetLocationService() *LocationService {
	return locationService
}
