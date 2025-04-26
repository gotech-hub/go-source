package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func GetTypeWithoutPointer(v interface{}) reflect.Type {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Type()
}

func join(elems interface{}, delim string) (string, error) {
	typeElems := GetTypeWithoutPointer(elems)
	if typeElems.Kind() != reflect.Slice {
		return "", fmt.Errorf("elems must be a slice: elems=%v", elems)
	}

	sliceValue := reflect.ValueOf(elems)
	strElems := make([]string, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i).Interface()
		elemType := GetTypeWithoutPointer(elem)

		switch elemType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			// case number
			strElems[i] = fmt.Sprintf("%v", elem)
			break
		case reflect.String:
			// case string
			strElems[i] = fmt.Sprintf("'%v'", elem)
			break
		default:
			return "", fmt.Errorf("elems must be a slice of number or string: elems=%v", elems)
		}
	}

	return strings.Join(strElems, delim), nil
}

func execDefault(query string, key string, value interface{}) (string, error) {
	typeValue := GetTypeWithoutPointer(value)
	switch typeValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// case number
		query = strings.ReplaceAll(query, key, fmt.Sprintf("%v", value))
		break
	case reflect.String:
		// case string
		query = strings.ReplaceAll(query, key, fmt.Sprintf("'%v'", value))
		break
	default:
		return "", fmt.Errorf("value must be a number or string: value=%v", value)
	}

	return query, nil
}

func execJoin(query string, key string, value interface{}) (string, error) {
	strValue, err := join(value, ",")
	if err != nil {
		return "", err
	}

	query = strings.ReplaceAll(query, key, strValue)
	return query, nil
}

func ConvertSqlTemplate(query string, args map[string]interface{}) (string, error) {
	var err error
	for key, value := range args {
		re := regexp.MustCompile(fmt.Sprintf(`{{(?:\w+\.%s|\.%s)}}`, key, key))
		matches := re.FindAllString(query, -1)
		if len(matches) == 0 {
			continue
		}

		firstSubString := matches[0]
		for _, subString := range matches {
			if subString != firstSubString {
				return "", fmt.Errorf("there are 2 keys with 2 different actions: %s, %s", firstSubString, subString)
			}
		}

		actionAndKey := strings.ReplaceAll(firstSubString, "{{", "")
		actionAndKey = strings.ReplaceAll(actionAndKey, "}}", "")

		splitActionAndKey := strings.Split(actionAndKey, ".")

		switch splitActionAndKey[0] {
		case "":
			query, err = execDefault(query, firstSubString, value)
			if err != nil {
				return "", err
			}
		case "join":
			query, err = execJoin(query, firstSubString, value)
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("action %s is not supported", splitActionAndKey[0])
		}
	}

	return query, nil
}
