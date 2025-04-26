package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Validate struct {
	validate *validator.Validate
}

func NewValidate() *Validate {
	validate := Validate{}
	validate.validate = validator.New()
	validate.registerTagJson()
	validate.registerCustomValidations()

	return &validate
}

func (v *Validate) registerCustomValidations() {
	// Register custom validation for 'string' tag
	v.validate.RegisterValidation("string", func(fl validator.FieldLevel) bool {
		// Check if the value is a string
		return fl.Field().Kind() == reflect.String
	})

	// Register custom validation for '[]string' tag
	v.validate.RegisterValidation("array", func(fl validator.FieldLevel) bool {
		// Check if the value is a slice of strings
		return fl.Field().Kind() == reflect.Slice
	})

	// Register custom validation for 'int' tag
	v.validate.RegisterValidation("number", func(fl validator.FieldLevel) bool {
		// Check if the value is an integer
		return fl.Field().Kind() == reflect.Float64
	})

	// Register custom validation for 'date' tag
	v.validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})
}

func (v *Validate) ValidateStruct(obj interface{}) error {
	err := v.validate.Struct(obj)
	if err != nil {
		err = v.customError(err.(validator.ValidationErrors), "")
	}
	return err
}

func (v *Validate) ValidateMap(data, rules map[string]interface{}) error {
	for key, _ := range data {
		if _, ok := rules[key]; !ok {
			return fmt.Errorf("key %v in programData not register", key)
		}
	}

	errs := v.validate.ValidateMap(data, rules)

	if len(errs) > 0 {
		errCustom := ""
		for fieldName, e := range errs {
			errCustom += v.customError(e.(validator.ValidationErrors), fieldName).Error() + " | "
		}

		n := len(errCustom)
		if n > 0 {
			return fmt.Errorf(errCustom[:n-3])
		}
	}

	return nil
}

func (v *Validate) RegisterStructValidation(fn validator.StructLevelFunc, validateStruct interface{}) {
	v.validate.RegisterStructValidation(fn, validateStruct)
}

func (v *Validate) customError(errs validator.ValidationErrors, fieldName string) error {
	if len(errs) > 0 {
		errCustom := ""
		for _, e := range errs {
			errCustom += v.mapValidatorError(e, fieldName) + " | "
		}

		n := len(errCustom)
		if n > 0 {
			return fmt.Errorf(errCustom[:n-3])
		}
	}

	return nil
}

// Define custom error message for each tag
func (v *Validate) mapValidatorError(errMap validator.FieldError, fieldName string) string {
	field := errMap.Field()
	if field == "" {
		field = fieldName
	} else {
		ok, _ := regexp.MatchString(`^\[.*\]$`, field)
		if ok {
			field = fieldName
		}
	}

	switch errMap.Tag() {
	case "required":
		return fmt.Sprintf("%v is required", field)
	case "min":
		return fmt.Sprintf("%v - min=%v", field, errMap.Param())
	case "max":
		return fmt.Sprintf("%v - max=%v", field, errMap.Param())
	case "oneof":
		return fmt.Sprintf("%v - must be one of: [%v]", field, strings.ReplaceAll(errMap.Param(), " ", ", "))
	case "rfe":
		return fmt.Sprintf("%v is required when %v = %v", field, errMap.StructField(), errMap.Value())
	case "email":
		return fmt.Sprintf("%v - invalid email: '%v'", field, errMap.Value())
	case "eth_addr":
		return fmt.Sprintf("%v - invalid ethereum address: %v", field, errMap.Value())
	case "string":
		return fmt.Sprintf("%v - invalid string", field)
	case "array":
		return fmt.Sprintf("%v - invalid array", field)
	case "number":
		return fmt.Sprintf("%v - invalid number", field)
	case "url":
		return fmt.Sprintf("%v - invalid url: '%v'", field, errMap.Value())
	case "base64":
		return fmt.Sprintf("%v - invalid base64", field)
	case "hexadecimal":
		return fmt.Sprintf("%v - invalid hex string", field)
	}

	return errMap.Error()
}

// to use the names which json tag in struct, rather than normal Go field names
func (v *Validate) registerTagJson() {
	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
