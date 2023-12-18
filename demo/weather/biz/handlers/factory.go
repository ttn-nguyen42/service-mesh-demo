package handlers

import (
	"context"
	"sync"
)

var once sync.Once

var handlers *Handlers

func Init(ctx context.Context) {
	once.Do(func() {
		handlers = NewHandlers()
	})
}

func GetHandlers() *Handlers {
	return handlers
}
