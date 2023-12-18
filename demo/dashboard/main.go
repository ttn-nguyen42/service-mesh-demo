package main

import (
	"context"
	"labs/service-mesh/helper/app"
	"labs/service-mesh/helper/cache"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	publicapi "labs/service-mesh/dashboard/api/public"
	"labs/service-mesh/dashboard/biz/handlers"
	"labs/service-mesh/dashboard/configs"
	"labs/service-mesh/dashboard/helpers/factory"
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
					custhttp.WithTemplatePath("./templates"),
				)),
				app.WithFactoryHook(func() error {
					cache.Init()
					factory.Init(context.Background())
					handlers.Init(context.Background())
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
