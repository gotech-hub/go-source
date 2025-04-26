package logger

import (
	"context"
	"go-source/pkg/utils"
	"reflect"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	loggerInstance *Logger
	mu             sync.RWMutex
	keyEncrypt     *string
)

const (
	KeyServiceName = "service_name"
	KeyFileError   = "file_error"
)

func InitLog(serviceName string) {
	mu.Lock()
	defer mu.Unlock()
	if loggerInstance != nil {
		return
	}

	if serviceName == "" {
		log.Fatal().Msg("services name is empty")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	lg := log.With().Str(KeyServiceName, serviceName).Logger()
	loggerInstance = &Logger{lg}
}

func SetKeyEncrypt(key string) {
	if keyEncrypt == nil {
		keyEncrypt = &key
	}
}

func GetLogger() *Logger {
	return loggerInstance
}

func SetEchoReqEncrLog(c echo.Context, req interface{}) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return
	}

	ctx := c.Request().Context()
	if req != nil {
		if newReq, err := utils.StructEncryptTagInterface(req, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt); err == nil {
			if str, err := utils.AnyToString(newReq); err == nil {
				ctx = context.WithValue(ctx, utils.KeyRequestBody, str)
				c.SetRequest(c.Request().WithContext(ctx))
			}
		}
	}
}

func SetEchoRespEncrLog(c echo.Context, resp interface{}) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return
	}

	ctx := c.Request().Context()

	// check response is nil
	if resp == nil {
		return
	}

	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// get value field Data from response
	if v.Kind() == reflect.Struct {
		if data := v.FieldByName("Data"); data.IsValid() {
			if data.Kind() == reflect.Ptr {
				data = data.Elem()
			}

			if newRes, err := utils.InterfaceEncryptTagInterface(data.Interface(), *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt); err == nil {
				if str, err := utils.AnyToString(newRes); err == nil {
					ctx = context.WithValue(ctx, utils.KeyResponseBody, str)
					c.SetRequest(c.Request().WithContext(ctx))
				}
			}
		}
	}
}

func Encrypt[T any](data T) (T, error) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return data, nil
	}

	switch v := interface{}(data).(type) {
	case string:
		res, err := utils.Encrypt(v, *keyEncrypt)
		if err != nil {
			return data, err
		}

		var result interface{} = res
		return result.(T), nil
	case *string:
		res, err := utils.Encrypt(*v, *keyEncrypt)
		if err != nil {
			return data, err
		}

		var result interface{} = &res
		return result.(T), nil
	}

	return utils.InterfaceEncryptTag(data, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
}

func EncryptInterface(data interface{}) (interface{}, error) {
	if keyEncrypt == nil || *keyEncrypt == "" {
		return data, nil
	}

	switch v := data.(type) {
	case string:
		return utils.Encrypt(v, *keyEncrypt)
	case *string:
		return utils.Encrypt(*v, *keyEncrypt)
	}

	return utils.InterfaceEncryptTagInterface(data, *keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
}
