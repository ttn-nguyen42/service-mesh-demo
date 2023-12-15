package handlers

import (
	"context"
	"sync"
)

var once sync.Once

var locationHandlers *LocationHandlers

func Init(ctx context.Context) {
	once.Do(func() {
		locationHandlers = NewLocationHandlers(ctx)
	})
}

func Location() *LocationHandlers {
	return locationHandlers
}
