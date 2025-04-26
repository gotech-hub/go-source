package mongodb

import (
	"go-source/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func readTagEncrypt(data interface{}) map[string]bool {
	result := make(map[string]bool)

	v := reflect.ValueOf(data)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("encrypt")

		if tag == "true" {
			flag := false
			if field.Kind() == reflect.String {
				flag = true
			} else if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.String {
				flag = true
			}

			if flag {
				name := utils.CamelToSnake(t.Field(i).Name)
				tagBson := t.Field(i).Tag.Get("bson")
				if tagBson != "" {
					name = tagBson
				}
				result[name] = true
			}
			continue
		}

		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			subRes := readTagEncrypt(field.Interface())

			name := utils.CamelToSnake(t.Field(i).Name)
			tagBson := t.Field(i).Tag.Get("bson")
			if tagBson != "" {
				name = tagBson
			}

			for key, _ := range subRes {
				result[name+"."+key] = true
			}
		}
	}

	return result
}

func encryptBsonUpdate(input interface{}, mapFieldName map[string]bool, key string) (interface{}, error) {
	data, ok := input.(bson.M)
	if !ok {
		return input, nil
	}

	if setValue, ok1 := data["$set"].(bson.M); ok1 {
		for k, v := range setValue {
			if s, ok2 := v.(string); ok2 {
				if _, ok3 := mapFieldName[k]; ok3 {
					encryptedValue, err := utils.Encrypt(s, key)
					if err != nil {
						return input, err
					}
					setValue[k] = encryptedValue
				}
			}
		}
		data["$set"] = setValue
	}

	if setOnInsertValue, ok1 := data["$setOnInsert"].(bson.M); ok1 {
		for k, v := range setOnInsertValue {
			if s, ok2 := v.(string); ok2 {
				if _, ok3 := mapFieldName[k]; ok3 {
					encryptedValue, err := utils.Encrypt(s, key)
					if err != nil {
						return input, err
					}
					setOnInsertValue[k] = encryptedValue
				}
			}
		}
		data["$setOnInsert"] = setOnInsertValue
	}

	return data, nil
}

func calcIndex(index, key []string) int {
	s := 0
	n := len(key)
	for i, v := range index {
		if i >= n {
			break
		}

		if v == key[i] {
			s++
		}
	}
	return s
}
