package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"go-source/api/http"
	"go-source/pkg/constant"
	logger "go-source/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	logger.InitLog(constant.ServiceName)
	log := logger.GetLogger()
	log.Info().Msgf("Start %s services", constant.ServiceName)

	http.SetHealthCheck(true)
	e := echo.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Handle graceful shutdown.
	<-ctx.Done()
	http.SetHealthCheck(false)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msgf("Force shutdown services")
	}
}
