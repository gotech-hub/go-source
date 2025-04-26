package middlewares

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/tidwall/sjson"
	logger "go-source/pkg/log"
	"go-source/pkg/utils"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func AddExtraDataForRequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Request().Header.Get(echo.HeaderXRequestID)
		if reqID == "" {
			reqID = c.Response().Header().Get(echo.HeaderXRequestID)
		}

		// set request_id to request
		c.Request().Header.Set(echo.HeaderXRequestID, reqID)

		// set request_id to response
		c.Response().Header().Set(echo.HeaderXRequestID, reqID)

		// set trace_info to context
		traceInfo := utils.TraceInfo{RequestID: reqID}
		ctx := c.Request().Context()
		ctxTraceInfo := context.WithValue(ctx, utils.KeyTraceInfo, traceInfo)
		c.SetRequest(c.Request().WithContext(ctxTraceInfo))

		return next(c)
	}
}

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		start := time.Now()
		req := c.Request()

		// Request
		reqBody := []byte{}
		if c.Request().Body != nil { // Read
			reqBody, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		if !json.Valid(reqBody) && len(reqBody) > 256 {
			reqBody = append(reqBody[:256], []byte("***")...)
		}

		// Response
		resBody := new(bytes.Buffer)
		mw := io.MultiWriter(c.Response().Writer, resBody)
		writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
		c.Response().Writer = writer

		if err := next(c); err != nil {
			c.Error(err)
		}
		res := c.Response()
		latency := time.Since(start)
		latencyInMs := float64(latency.Nanoseconds()) / 1000000.0
		statusCode := res.Status
		method := req.Method
		path := req.URL.Path
		if path == "" {
			path = "/"
		}
		if strings.Contains(path, "/health") {
			return nil
		}

		requestBody := string(reqBody)
		responseBody := resBody.String()

		ctx = c.Request().Context()

		if newReqBody := ctx.Value(utils.KeyRequestBody); newReqBody != nil {
			if str, err := utils.AnyToString(newReqBody); err == nil {
				requestBody = str
			}
		}

		if newResBody := ctx.Value(utils.KeyResponseBody); newResBody != nil {
			if str, err := utils.AnyToString(newResBody); err == nil {
				if v, err := sjson.Set(responseBody, "data", str); err == nil {
					responseBody = v
				}
			}
		}

		log := logger.GetLogger().AddTraceInfoContextRequest(req.Context())

		var newLog logger.Logger
		if region := c.Request().Header.Get(utils.KeyRegion); region != "" {
			newLog = log.With().Str("region", region).Logger()
		} else {
			newLog = *log
		}

		var eventLog *logger.Event
		if statusCode >= 500 {
			eventLog = newLog.Error()
		} else {
			eventLog = newLog.Info()
		}

		eventLog.Str("method", method).
			Str("path", path).
			Str("ip", c.RealIP()).
			Str("user_agent", req.UserAgent()).
			Str("request_id", req.Header.Get(echo.HeaderXRequestID)).
			Int("statusCode", statusCode).
			Float64("latency", latencyInMs).
			Interface("params", c.QueryParams()).
			Str("request_body", requestBody).
			Str("response_body", responseBody).Msg("request income")

		return nil
	}
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
