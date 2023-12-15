package main

import (
	"context"
	publicapi "labs/service-mesh/locations/api/public"
	"labs/service-mesh/locations/biz/handlers"
	"labs/service-mesh/locations/biz/service"
	"labs/service-mesh/locations/configs"
	"labs/service-mesh/locations/helpers/factory"
	"labs/service-mesh/helper/app"
	"labs/service-mesh/helper/cache"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	"time"

	"go.uber.org/zap"
)

func main() {
	app.Run(
		time.Second*10,
		configs.Get(),
		func(zl *zap.Logger) []app.Optioner {
			return []app.Optioner{
				app.WithHttpServer(custhttp.New(
					custhttp.WithGlobalConfigs(configs.Get().GetPublic()),
					custhttp.WithErrorHandler(custhttp.GlobalErrorHandler()),
					custhttp.WithRegistration(publicapi.ServiceRegistration()),
					custhttp.WithMiddleware(custhttp.CommonPublicMiddlewares(configs.Get().GetPublic())...),
				)),
				app.WithFactoryHook(func() error {
					ctx := context.Background()
					cache.Init()
					factory.Init(ctx)
					service.Init()
					handlers.Init(ctx)
					return nil
				}),
				app.WithShutdownHook(func(ctx context.Context) {
					cache.Stop()
					logger.Close()
				}),
			}
		},
	)
}
