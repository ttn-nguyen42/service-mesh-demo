package handlers

import (
	"context"
	"sync"
)

var once sync.Once

var locationHandlers *CitiesHandlers

func Init(ctx context.Context) {
	once.Do(func() {
		locationHandlers = NewCitiesHandlers(ctx)
	})
}

func Cities() *CitiesHandlers {
	return locationHandlers
}
