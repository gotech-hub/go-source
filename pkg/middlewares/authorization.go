package middlewares

import (
	"context"
	"encoding/json"
	"go-source/pkg/client"
	"go-source/pkg/jwt"
	"go-source/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Authorization(clientAuthorize *client.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// get token
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, Resp{
					ErrorCode:   ErrAuthentication,
					Message:     "Authentication failed",
					Description: "Token is empty",
				})
			}

			// parse token
			var claims map[string]interface{}
			err := jwt.ParseClaims(strings.TrimPrefix(token, "Bearer "), &claims)
			if err == nil {
				// check expired
				exp, ok := claims[utils.JwtExp].(int64)
				if ok && time.Now().Unix() > exp {
					return c.JSON(http.StatusUnauthorized, Resp{
						ErrorCode:   ErrAuthentication,
						Message:     "Authentication failed",
						Description: "Token expired",
					})
				}

				// get sub
				sub, ok := claims[utils.JwtSub].(string)
				if ok {
					// set sub to context
					ctx = context.WithValue(ctx, utils.JwtSub, sub)
				}
			}

			body := map[string]interface{}{
				"type":   "API",
				"method": c.Request().Method,
				"url":    c.Request().URL.String(),
			}

			reqHeader := make(map[string][]string)
			for k, v := range c.Request().Header {
				reqHeader[k] = v
			}
			reqHeader["Content-Type"] = []string{"application/json"}
			reqHeader["Accept-Encoding"] = []string{"gzip"}

			// call authenticate and authorize: request header is c.Request().Header
			res, err := clientAuthorize.R().SetContext(ctx).SetBody(body).SetHeaderMultiValues(reqHeader).Post("")
			if err != nil {
				var respBody Resp
				if len(res.Body()) > 0 {
					if err := json.Unmarshal(res.Body(), &respBody); err != nil {
						return c.JSON(http.StatusUnauthorized, Resp{
							ErrorCode:   ErrAuthentication,
							Message:     "Authentication failed",
							Description: err.Error(),
						})
					}
				}

				if respBody.Message == "" {
					return c.JSON(http.StatusUnauthorized, Resp{
						ErrorCode:   ErrAuthentication,
						Message:     "Authentication failed",
						Description: err.Error(),
					})
				}

				return c.JSON(res.StatusCode(), Resp{
					ErrorCode:   respBody.ErrorCode,
					Message:     respBody.Message,
					Description: respBody.Description,
				})
			}

			me := res.Header().Get(utils.HeaderXMeProfile)
			if me != "" {
				// set me to context
				ctx = context.WithValue(ctx, utils.HeaderXMeProfile, me)
			}

			// set ctx to request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
