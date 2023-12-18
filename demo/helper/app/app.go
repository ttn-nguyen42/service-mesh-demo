package app

import (
	"context"
	"errors"
	"labs/service-mesh/helper/configs"
	custcron "labs/service-mesh/helper/cron"
	"labs/service-mesh/helper/events"
	custhttp "labs/service-mesh/helper/http"
	"labs/service-mesh/helper/logger"
	custmqtt "labs/service-mesh/helper/mqtt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"
)

func Run(shutdownTimeout time.Duration, destConfig custconfigs.ServiceConfigs, registration RegistrationFunc) {
	ctx := context.Background()
	custconfigs.Init(ctx, destConfig)

	globalConfigs := destConfig

	loggerConfigs := globalConfigs.GetLogger()
	logger.Init(ctx, logger.WithGlobalConfigs(loggerConfigs))

	options := registration(logger.Logger())

	opts := Options{}
	for _, optioner := range options {
		optioner(&opts)
	}

	logger := zap.L().Sugar()

	logger.Infof("Run: configs = %s", globalConfigs.String())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	failures := make(chan bool)
	for _, s := range opts.httpServers {
		s := s
		go func() {
			logger.Infof("Run: start HTTP server name = %s", s.Name())
			if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Infof("Run: start HTTP server err = %s", err)
				failures <- true
			}
		}()
	}

	for _, s := range opts.natsServers {
		s := s
		go func() {
			logger.Infof("Run: start embedded NATS server name = %s", s.Name())
			if err := s.Start(); err != nil {
				logger.Infof("Run: start embedded NATS server err = %s", err)
				failures <- true
			}
		}()
	}

	for _, s := range opts.mqttServers {
		s := s
		go func() {
			logger.Infof("Run: start embedded MQTT server name = %s", s.Name())
			if err := s.Start(); err != nil {
				logger.Infof("Run: start embedded MQTT server err = %s", err)
				failures <- true
			}
		}()
	}

	for _, s := range opts.schedulers {
		s := s
		go func() {
			logger.Infof("Run: start scheduler name = %s", s.Name())
			if err := s.Start(); err != nil {
				logger.Infof("Run: start scheduler err = %s", err)
				failures <- true
			}
		}()
	}

	if opts.factoryHook != nil {
		if err := opts.factoryHook(); err != nil {
			logger.Fatalf("Run: factoryHook err = %s", err)
			failures <- true
			return
		}
	}

	select {
	case <-quit:
		logger.Infof("Run: exit signal received, exiting")
	case <-failures:
		logger.Infof("Run: failure occurred, exiting")
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if opts.shutdownHook != nil {
		opts.shutdownHook(ctx)
	}

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		for _, s := range opts.httpServers {
			s := s
			logger.Infof("Run: stop HTTP server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for _, s := range opts.natsServers {
			s := s
			logger.Infof("Run: stop NATS embedded server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for _, s := range opts.mqttServers {
			s := s
			logger.Infof("Run: stop MQTT embedded server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for _, s := range opts.schedulers {
			s := s
			logger.Infof("Run: stop scheduler name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()

	zap.L().Sync()
	log.Print("Run: shutdown complete")
}

type RegistrationFunc func(logger *zap.Logger) []Optioner
type FactoryHook func() error
type ShutdownHook func(ctx context.Context)

type Options struct {
	httpServers []*custhttp.HttpServer
	natsServers []*events.EmbeddedNats
	mqttServers []*custmqtt.EmbeddedMqtt
	schedulers  []*custcron.Scheduler

	factoryHook  FactoryHook
	shutdownHook ShutdownHook
}

type Optioner func(opts *Options)

func WithHttpServer(server *custhttp.HttpServer) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.httpServers = append(opts.httpServers, server)
		}
	}
}

func WithNatsServer(server *events.EmbeddedNats) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.natsServers = append(opts.natsServers, server)
		}
	}
}

func WithMqttServer(server *custmqtt.EmbeddedMqtt) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.mqttServers = append(opts.mqttServers, server)
		}
	}
}

func WithFactoryHook(cb FactoryHook) Optioner {
	return func(opts *Options) {
		opts.factoryHook = cb
	}
}

func WithShutdownHook(cb ShutdownHook) Optioner {
	return func(opts *Options) {
		opts.shutdownHook = cb
	}
}

func WithScheduling(scheduler *custcron.Scheduler) Optioner {
	return func(opts *Options) {
		opts.schedulers = append(opts.schedulers, scheduler)
	}
}
