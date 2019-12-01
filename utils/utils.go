package utils

import (
	"reflect"
	"strconv"
	"strings"
)

var EPSILON = 0.00000001

func FloatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}

	return false
}

func GetObjectValueByKey(data interface{}, key string) (interface{}, bool) {
	for _, seg := range strings.Split(key, ".") {
		if data == nil {
			return nil, false
		}

		switch reflect.TypeOf(data).Kind() {
		case reflect.Map:
			if v, ok := data.(map[string]interface{}); ok {
				if data, ok = v[seg]; !ok {
					return nil, false
				}
			} else {
				return nil, false
			}

		case reflect.Array:
			v := data.([]interface{})
			if i, err := strconv.Atoi(seg); err != nil {
				return nil, false
			} else if i < 0 || i >= len(v) {
				return nil, false
			} else {
				data = v[i]
			}

		case reflect.Struct:
			value := reflect.ValueOf(data)
			f := value.FieldByName(seg)
			if !f.IsValid() {
				return nil, false
			}
			data = f.Interface()

		case reflect.Ptr:
			value := reflect.ValueOf(data).Elem()
			f := value.FieldByName(seg)
			if !f.IsValid() {
				return nil, false
			}
			data = f.Interface()

		default:
			return nil, false
		}
	}
	return data, true
}
