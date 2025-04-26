package binding

import (
	"github.com/labstack/echo/v4"
	validate2 "go-source/pkg/validate"
)

var (
	validate     = validate2.NewValidate()
	customBinder = &CustomBinder{}
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := validate.ValidateStruct(i); err != nil {
		return err
	}

	return nil
}

func GetBinding() *CustomBinder {
	return customBinder
}
