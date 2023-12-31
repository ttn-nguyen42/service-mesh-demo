package service

import "sync"

var once sync.Once

var (
	citiesService *CitiesService
)

func Init() {
	once.Do(func() {
		citiesService = NewCitiesService()
	})
}

func GetCitiesService() *CitiesService {
	return citiesService
}
