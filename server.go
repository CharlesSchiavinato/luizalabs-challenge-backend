// @title           API de Normalização das Informações dos Pedidos do Sistema Legado
// @version         1.0.0
// @description     Está API tem como objetivo carregar os pedidos do sistema legado a partir de arquivos no formato TXT desnormalizado e devolver as informações normalizadas no formato JSON.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:9000
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/route"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/router"
	cache "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache/redis"
	repository "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
	gohandlers "github.com/gorilla/handlers"
	"github.com/hashicorp/go-hclog"
)

func main() {
	// create a new logger
	log := hclog.New(&hclog.LoggerOptions{
		Name:       "luizalabs-order",
		JSONFormat: false,
		Level:      hclog.LevelFromString("DEBUG"),
	})

	// load configs
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Error("Cannot load application configs", "error", err)
		return
	}

	// update logger with new options
	log = hclog.New(&hclog.LoggerOptions{
		Name:       config.ServerAppName,
		JSONFormat: config.ServerLogJSONFormat,
		Level:      hclog.LevelFromString(config.ServerLogLevel),
	})

	log.Info("Configs loaded successfuly")

	// run database migration
	// err = migration.Run(config)

	// if err != nil {
	// 	log.Error("Cannot run db migration", "error", err)
	// 	os.Exit(0)
	// }

	// log.Info("DB migration run successfuly")

	// create a new repository
	repository, err := repository.NewInMemory(config)

	if err != nil {
		log.Error("Cannot connect to database", "error", err)
		os.Exit(0)
	}

	defer repository.Close()

	log.Info("Connected database successfuly")

	// create a new cache
	cache, err := cache.NewRedis(config)

	if err != nil {
		log.Error("Cannot connect to cache", "error", err)
		os.Exit(0)
	}

	defer cache.Close()

	log.Info("Connected cache successfuly")

	// set server address
	serverAddr := config.ServerAddress

	// create a new router
	appRouter := router.NewHttpRouter()
	routerParameters := &route.RouteParameters{
		AppRouter:  appRouter,
		Log:        log,
		Repository: repository,
		Cache:      cache,
	}

	// include the routes
	route.OrderRoute(routerParameters)
	route.SwaggerRoute(appRouter)
	route.HealthzRoute(routerParameters)

	// create HTTP handler
	httpHandler := appRouter.Serve()

	// include the middleware handler CORS
	corsOption := gohandlers.AllowedOrigins(strings.Split(config.ServerCORSAllowedOrigins, ";"))
	corsHandler := gohandlers.CORS(corsOption)
	httpHandler = corsHandler(httpHandler)

	// include the middleware handler logger
	httpHandler = router.HttpLogger(httpHandler, log)

	requestTimeout, err := time.ParseDuration(config.ServerRequestTimeout)

	if err != nil {
		requestTimeout, err = time.ParseDuration("1m")
		log.Error("Cannot load config RequestTimeout", "error", err)
	}

	// create a new server
	httpServer := http.Server{
		Addr:    serverAddr,
		Handler: httpHandler,
		// ErrorLog: log,
		ReadTimeout:  requestTimeout,
		WriteTimeout: requestTimeout,
		IdleTimeout:  requestTimeout,
	}

	// start the server
	go func() {
		log.Info(fmt.Sprintf("HTTP server running on port %v", serverAddr))

		err := httpServer.ListenAndServe()

		if err != nil {
			log.Error("Error running HTTP server", "error", err)
			os.Exit(0)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt, os.Kill)

	// Block until a signal is received.
	sig := <-chanSignal
	log.Info(fmt.Sprintf("HTTP server terminate signal %v", sig))

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)
}
