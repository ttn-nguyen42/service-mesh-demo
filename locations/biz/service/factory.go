package service

import "sync"

var once sync.Once

var (
	datetimeService *DatetimeService
	locationService *LocationsService
)

func Init() {
	once.Do(func() {
		datetimeService = NewDatetimeService()
		locationService = NewLocationsService()
	})
}

func GetDatetimeService() *DatetimeService {
	return datetimeService
}

func GetLocationService() *LocationsService {
	return locationService
}
