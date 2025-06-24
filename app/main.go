package main

import (
	"context"
	"go-source/api/http"
	"go-source/pkg/constant"
	"go-source/pkg/database/redis"
	logger "go-source/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-source/bootstrap"
	"go-source/config"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		logger.GetLogger().Fatal().Msgf("Failed to load configuration: %v", err)
		return
	}

	// Initialize logger
	logger.InitLog(config.ServiceName)
	log := logger.GetLogger()
	log.Info().Msgf("Start %s services", constant.ServiceName)

	http.SetHealthCheck(true)
	e := echo.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize rs
	redisClient, err := redis.ConnectRedis(ctx, &config.RedisConfig)
	if err != nil {
		log.Fatal().Msgf("Connect rs fail! %s", err)
	}

	// Initialize database connections
	storage := bootstrap.NewDatabaseConnection(ctx)

	// Initialize clients
	clients := bootstrap.NewClients()

	// Initialize repositories
	repo := bootstrap.NewRepositories(storage.Connection)

	// Initialize services
	services := bootstrap.NewServices(repo, redisClient, clients)

	// Initialize handlers
	handlers := bootstrap.NewHandlers(services)

	// Start message broker
	// msgBroker := msg.NewMsgBroker(conf, services.TierService)
	// msgBroker.Start(ctx)

	// Start HTTP server
	srv := http.NewHttpServe(handlers)
	srv.Start(e)

	// Handle graceful shutdown.
	<-ctx.Done()
	http.SetHealthCheck(false)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msgf("Force shutdown services")
	}
}
