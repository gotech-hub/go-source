package utils

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

func BuildFromTemplateAndData(template string, data interface{}) (string, error) {
	dataStr, err := AnyToString(data)
	if err != nil {
		return "", err
	}

	keys := FindKeysWithRegex(template, `{{(.*?)}}`)

	for _, key := range keys {
		value := gjson.Get(dataStr, key).Value()
		if valueStr, err := AnyToString(value); err == nil {
			template = strings.ReplaceAll(template, "{{"+key+"}}", valueStr)
		} else {
			template = strings.ReplaceAll(template, "{{"+key+"}}", fmt.Sprintf("%+v", value))
		}
	}

	return template, nil
}
