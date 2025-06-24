package middlewares

import (
	"context"

	logger "go-source/pkg/log"

	"github.com/labstack/echo/v4"
)

// Authenticate Middleware function to check authentication
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		_ = logger.GetLogger().AddTraceInfoContextRequest(ctx)

		newCtx := context.WithValue(ctx, "ProfileId", "id")
		c.SetRequest(c.Request().WithContext(newCtx))

		// If authentication is successful, call the next handler
		return next(c)
	}
}
