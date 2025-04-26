package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	logger "go-source/pkg/log"
	"go-source/pkg/metric"
	"go-source/pkg/utils"
	"net/http"
	neturl "net/url"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Client struct {
	*resty.Client
	keyEncrypt         string
	fnResp             []ResponseMiddlewareFunc
	fnReq              []RequestMiddlewareFunc
	reqBodyEncrypt     []TypeKeyRequestBodyEncrypt
	respBodyEncrypt    []TypeKeyResponseBodyEncrypt
	queryParamsEncrypt []TypeKeyQueryParamsEncrypt
}

type (
	ResponseMiddlewareFunc func(client *resty.Client, response *resty.Response) error
	RequestMiddlewareFunc  func(client *resty.Client, request *resty.Request) error

	TypeKeyRequestBodyEncrypt  string
	TypeKeyResponseBodyEncrypt string
	TypeKeyQueryParamsEncrypt  string
)

const (
	KeyRestRequestBodyEncrypt  = "rest_request_body_encrypt"
	KeyRestResponseBodyEncrypt = "rest_response_body_encrypt"
	KeyRestQueryParamsEncrypt  = "rest_query_params_encrypt"

	KeyOffLogRespBody = "off_log_resp_body"
)

func NewClient(baseURL string, timeout time.Duration, maxRetry int, waitTime time.Duration, fns ...interface{}) *Client {
	var (
		fnResp             []ResponseMiddlewareFunc
		fnReq              []RequestMiddlewareFunc
		reqBodyEncrypt     []TypeKeyRequestBodyEncrypt
		respBodyEncrypt    []TypeKeyResponseBodyEncrypt
		queryParamsEncrypt []TypeKeyQueryParamsEncrypt
	)

	for _, fn := range fns {
		switch fn.(type) {
		case []ResponseMiddlewareFunc:
			fnResp = append(fnResp, fn.([]ResponseMiddlewareFunc)...)
		case []RequestMiddlewareFunc:
			fnReq = append(fnReq, fn.([]RequestMiddlewareFunc)...)
		case []TypeKeyRequestBodyEncrypt:
			reqBodyEncrypt = append(reqBodyEncrypt, fn.([]TypeKeyRequestBodyEncrypt)...)
		case []TypeKeyResponseBodyEncrypt:
			respBodyEncrypt = append(respBodyEncrypt, fn.([]TypeKeyResponseBodyEncrypt)...)
		case []TypeKeyQueryParamsEncrypt:
			queryParamsEncrypt = append(queryParamsEncrypt, fn.([]TypeKeyQueryParamsEncrypt)...)
		}
	}

	// get key encrypt from env
	keyEncrypt := os.Getenv(utils.VGREncryptKey)

	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(timeout)
	client.SetRetryCount(maxRetry)
	client.SetRetryWaitTime(waitTime)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.EnableTrace()
	client.SetDisableWarn(true)
	client.SetLogger(NoOpLogger{})

	result := &Client{
		Client:             client,
		keyEncrypt:         keyEncrypt,
		fnResp:             fnResp,
		fnReq:              fnReq,
		reqBodyEncrypt:     reqBodyEncrypt,
		respBodyEncrypt:    respBodyEncrypt,
		queryParamsEncrypt: queryParamsEncrypt,
	}

	client.OnBeforeRequest(result.beforeRequest)
	client.OnAfterResponse(result.afterResponse)

	return result
}

func (c *Client) beforeRequest(client *resty.Client, request *resty.Request) error {
	ctx := request.Context()
	ctxNew, traceInfo := utils.NewContextWithRequestId(ctx)
	request.SetContext(ctxNew)
	request.SetHeader(echo.HeaderXRequestID, traceInfo.RequestID)

	// get log
	log := logger.GetLogger().AddTraceInfoContextRequest(request.Context())

	reqBody := request.Body

	reqBaseURL := client.BaseURL + request.URL
	reqQueryParams := copyQueryParams(request.QueryParam)

	if c.keyEncrypt != "" {
		// encrypt request body fields

		var reqBodyEnc []string
		for _, be := range c.reqBodyEncrypt {
			reqBodyEnc = append(reqBodyEnc, string(be))
		}

		if result, err := encryptBodyFields(reqBody, reqBodyEnc, c.keyEncrypt); err != nil {
			log.Error().Err(err).Msg("encrypt request body fields error")
		} else {
			reqBody = result
		}

		if value, ok := ctx.Value(KeyRestRequestBodyEncrypt).([]string); ok {
			if result, err := encryptBodyFields(reqBody, value, c.keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request body fields error")
			} else {
				reqBody = result
			}
		}

		// -----------------------------------------------

		// encrypt request query params

		var queryParamsEnc []string
		for _, qe := range c.queryParamsEncrypt {
			queryParamsEnc = append(queryParamsEnc, string(qe))
		}

		if result, err := encryptQueryParams(reqQueryParams, queryParamsEnc, c.keyEncrypt); err != nil {
			log.Error().Err(err).Msg("encrypt request query params error")
		} else {
			reqQueryParams = result
		}

		if value, ok := ctx.Value(KeyRestQueryParamsEncrypt).([]string); ok {
			if result, err := encryptQueryParams(reqQueryParams, value, c.keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request query params error")
			} else {
				reqQueryParams = result
			}
		}
	}

	// log trace
	log.Info().
		Str("url", reqBaseURL).
		Interface("query_params", reqQueryParams).
		Str("method", request.Method).
		Interface("request_body", reqBody).
		Msg("rest api request")

	for _, fn := range c.fnReq {
		_ = fn(client, request)
	}

	return nil
}

func (c *Client) afterResponse(client *resty.Client, response *resty.Response) error {
	ctx := response.Request.Context()

	// get log
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)

	// get trace request
	trace := response.Request.TraceInfo()

	logTrace := log.Info()
	if response.StatusCode() != http.StatusOK {
		logTrace = log.Error()
	}

	reqBody := response.Request.Body

	var respBody interface{}

	valOffLogRespBody := ctx.Value(KeyOffLogRespBody)
	value, ok := valOffLogRespBody.(bool)
	if ok && value && len(response.Body()) > 0 {
		respBody = "******"
	} else {
		respBody = response.Body()
	}

	reqBaseURL, reqQueryParams, err := parseURL(response.Request.URL)
	if err != nil {
		log.Error().Err(err).Msg("parse request url error")
	}

	if c.keyEncrypt != "" {
		// encrypt request body fields
		var reqBodyEnc []string
		for _, be := range c.reqBodyEncrypt {
			reqBodyEnc = append(reqBodyEnc, string(be))
		}

		if result, err := encryptBodyFields(reqBody, reqBodyEnc, c.keyEncrypt); err != nil {
			log.Error().Err(err).Msg("encrypt request body fields error")
		} else {
			reqBody = result
		}

		if value, ok := ctx.Value(KeyRestRequestBodyEncrypt).([]string); ok {
			if result, err := encryptBodyFields(reqBody, value, c.keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request body fields error")
			} else {
				reqBody = result
			}
		}

		// -----------------------------------------------

		// encrypt response body fields

		var respBodyEnc []string
		for _, be := range c.respBodyEncrypt {
			respBodyEnc = append(respBodyEnc, string(be))
		}

		if result, err := encryptBodyFields(respBody, respBodyEnc, c.keyEncrypt); err != nil {
			log.Error().Err(err).Msg("encrypt response body fields error")
		} else {
			respBody = result
		}

		if value, ok := ctx.Value(KeyRestResponseBodyEncrypt).([]string); ok {
			if result, err := encryptBodyFields(respBody, value, c.keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt response body fields error")
			} else {
				respBody = result
			}
		}

		// -----------------------------------------------

		// encrypt request query params

		var queryParamsEnc []string
		for _, qe := range c.queryParamsEncrypt {
			queryParamsEnc = append(queryParamsEnc, string(qe))
		}

		if result, err := encryptQueryParams(reqQueryParams, queryParamsEnc, c.keyEncrypt); err != nil {
			log.Error().Err(err).Msg("encrypt request query params error")
		} else {
			reqQueryParams = result
		}

		if value, ok := ctx.Value(KeyRestQueryParamsEncrypt).([]string); ok {
			if result, err := encryptQueryParams(reqQueryParams, value, c.keyEncrypt); err != nil {
				log.Error().Err(err).Msg("encrypt request query params error")
			} else {
				reqQueryParams = result
			}
		}
	}

	metric.NewHttpClientHistogramDuration(reqBaseURL, response.Request.Method, response.Status(), response.Time())

	// log trace
	logTrace.
		Str("url", reqBaseURL).
		Interface("query_params", reqQueryParams).
		Str("method", response.Request.Method).
		Str("status", response.Status()).
		Int("status_code", response.StatusCode()).
		Str("latency", response.Time().String()).
		Str("protocol", response.Proto()).
		Time("received_at", response.ReceivedAt()).
		Int64("size", response.Size()).
		Str("conn_time", trace.ConnTime.String()).
		Str("tcp_conn_time", trace.TCPConnTime.String()).
		Str("server_time", trace.ServerTime.String()).
		Str("response_time", trace.ResponseTime.String()).
		Str("total_time", trace.TotalTime.String()).
		Str("conn_idle_time", trace.ConnIdleTime.String()).
		Int("request_attempt", trace.RequestAttempt).
		Str("remote_addr", trace.RemoteAddr.String()).
		Interface("request_body", reqBody).
		Str("response_body", func() string {
			data, err := utils.AnyToString(respBody)
			if err != nil {
				return ""
			}

			return SensitiveDataInstance.FilterSensitiveFields(data)
		}()).
		Msg("rest api response")

	if response.StatusCode() != http.StatusOK {
		if len(response.Body()) > 0 {
			_ = json.Unmarshal(response.Body(), &response.Request.Result)
		}

		return fmt.Errorf("status_code: %d", response.StatusCode())
	}

	return nil
}

func OffLogRespBody(client *resty.Client, request *resty.Request) error {
	ctx := request.Context()
	ctxNew := context.WithValue(ctx, KeyOffLogRespBody, true)
	request.SetContext(ctxNew)
	return nil
}

func encryptBodyFields(body interface{}, bodyFields []string, secretKey string) (interface{}, error) {
	if len(bodyFields) == 0 {
		return body, nil
	}

	if secretKey == "" {
		return body, nil
	}

	bodyStr, err := utils.AnyToString(body)
	if err != nil {
		return "", err
	}

	if bodyStr == "" {
		return body, nil
	}

	mappingBodyFields := make(map[string]bool)
	for _, field := range bodyFields {
		if _, ok := mappingBodyFields[field]; ok {
			continue
		}

		mappingBodyFields[field] = true

		value := gjson.Get(bodyStr, field).Value()
		if value == nil {
			continue
		}

		encryptedValue, err := utils.Encrypt(fmt.Sprintf("%v", value), secretKey)
		if err != nil {
			return "", err
		}

		bodyStr, err = sjson.Set(bodyStr, field, encryptedValue)
		if err != nil {
			return "", err
		}
	}

	return bodyStr, nil
}

func parseURL(_url string) (string, neturl.Values, error) {
	if _url == "" {
		return "", nil, nil
	}

	url, err := neturl.Parse(_url)
	if err != nil {
		return "", nil, err
	}

	path := url.Scheme + "://" + url.Host + url.Path
	queryParams := url.Query()
	return path, queryParams, nil
}

func encryptQueryParams(queryParams neturl.Values, fields []string, secretKey string) (neturl.Values, error) {
	if len(queryParams) == 0 {
		return queryParams, nil
	}

	result := make(neturl.Values)

	mappingFields := make(map[string]bool)
	for _, field := range fields {
		mappingFields[field] = true
	}

	for key, value := range queryParams {
		if _, ok := mappingFields[key]; ok {
			encryptedValue, err := utils.Encrypt(value[0], secretKey)
			if err != nil {
				return nil, err
			}

			result[key] = []string{encryptedValue}
		} else {
			result[key] = value
		}
	}

	return result, nil
}

func copyQueryParams(queryParams neturl.Values) neturl.Values {
	result := make(neturl.Values)
	for key, value := range queryParams {
		result[key] = value
	}
	return result
}
