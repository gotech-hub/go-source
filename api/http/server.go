package http

import (
	"errors"
	"fmt"
	"go-source/bootstrap"
	"go-source/config"
	logger "go-source/pkg/log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	healthCheck bool
	mu          sync.RWMutex
)

func SetHealthCheck(status bool) {
	mu.Lock()
	defer mu.Unlock()
	healthCheck = status
}

type ServInterface interface {
	Start(e *echo.Echo)
}

type Server struct {
}

func NewHttpServe(handlers *bootstrap.Handlers) *Server {
	return &Server{}
}

func (app *Server) Start(e *echo.Echo) {
	log := logger.GetLogger()
	httpPort := config.GetInstance().HttpPort
	go func() {
		err := e.Start(fmt.Sprintf(":%d", httpPort))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("can't start echo")
		}
	}()
	log.Info().Msg("all service already")
}
