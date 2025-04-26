package middlewares

import (
	"go-source/pkg/database/redis"
	logger "go-source/pkg/log"
	"go-source/pkg/utils"
	"time"

	"github.com/go-redis/redis_rate/v10"

	"github.com/labstack/echo/v4"
)

func RateLimit(period, rate int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if period == 0 {
				period = 1
			}
			ctx := c.Request().Context()
			log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
			keyRateLimit := c.Request().Context().Value(utils.KeyRateLimit)
			if _, ok := keyRateLimit.(string); !ok {
				log.Info().Msg("key-limit is not string")
				return echo.ErrTooManyRequests
			}
			if keyRateLimit == "" {
				log.Info().Msg("key-limit is empty")
				return echo.ErrTooManyRequests
			}
			limiter := redis_rate.NewLimiter(redis.GetInstance().GetClient())
			result, err := limiter.Allow(ctx, keyRateLimit.(string), redis_rate.Limit{
				Rate:   rate,
				Burst:  rate,
				Period: time.Second * time.Duration(period),
			})
			if err != nil {
				return err
			}

			if result.Allowed == 0 {
				return echo.ErrTooManyRequests
			}

			return next(c)
		}
	}
}
