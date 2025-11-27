package utils

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/liyanbing/filter/types"
	"github.com/mohae/deepcopy"
)

var EPSILON = 0.00000001

func FloatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func GetObjectValueByKey(ctx context.Context, data interface{}, key string) (interface{}, bool) {
	type Valuer interface {
		Value(ctx context.Context, key string) interface{}
	}
	if valuer, ok := data.(Valuer); ok {
		return valuer.Value(ctx, key), true
	}

	key = strings.TrimSpace(key)
	if key == "." || key == "" {
		return data, true
	}

	segs := strings.Split(key, ".")
	if len(segs) <= 0 {
		return data, true
	}

	for index := 0; index < len(segs); {
		seg := segs[index]
		if data == nil {
			return nil, false
		}

		seg = strings.TrimSpace(seg)
		switch reflect.TypeOf(data).Kind() {
		case reflect.Map:
			if v, ok := data.(map[string]interface{}); ok {
				if data, ok = v[seg]; !ok {
					return nil, false
				}
			} else {
				return nil, false
			}
		case reflect.Array, reflect.Slice:
			value := reflect.ValueOf(data)
			if i, err := strconv.Atoi(seg); err != nil {
				return nil, false
			} else if i < 0 || i >= value.Len() {
				return nil, false
			} else {
				data = value.Index(i).Interface()
			}
		case reflect.Struct:
			var (
				value     = reflect.ValueOf(data)
				valueType = reflect.TypeOf(data)
			)

			f := value.FieldByName(seg)
			if !f.IsValid() {
				existsJson := false
				for i := 0; i < valueType.NumField(); i++ {
					tf := valueType.Field(i)
					if tf.Tag.Get("json") == seg {
						existsJson = true
						f = value.Field(i)
						break
					}
				}

				if !existsJson {
					return nil, false
				}
			}
			data = f.Interface()
		case reflect.Ptr:
			data = reflect.ValueOf(data).Elem().Interface()
			continue
		default:
			return nil, false
		}
		index++
	}
	return data, true
}

func ParseTargetArrayValue(value interface{}) []interface{} {
	var target []interface{}
	switch types.GetFilterType(value) {
	case types.STRING:
		err := json.Unmarshal([]byte(value.(string)), &target)
		if err == nil {
			return target
		}

		targetValue := value.(string)
		values := strings.Split(targetValue, ",")
		for _, v := range values {
			target = append(target, strings.TrimSpace(v))
		}
	case types.ARRAY:
		target = value.([]interface{})
	default:
		target = append(target, value)
	}
	return target
}

func NumberCompare(a, b interface{}) int {
	fa := types.GetFloat(a)
	fb := types.GetFloat(b)

	if FloatEquals(fa, fb) {
		return 0
	}

	if fa-fb > 0 {
		return 1
	} else {
		return -1
	}
}

// ObjectCompare compare
// compare == compared return 0
// compare > compared return 1
// compare < compared return -1
func ObjectCompare(compare, compared interface{}) int {
	if compare == nil && compared == nil {
		return 0
	}

	compareType := types.GetFilterType(compare)
	comparedType := types.GetFilterType(compared)

	if compareType == types.NUMBER || compareType == types.BOOL || comparedType == types.NUMBER || comparedType == types.BOOL {
		return NumberCompare(compare, compared)
	}

	if compareType == types.STRING || comparedType == types.STRING {
		return strings.Compare(types.GetString(compare), types.GetString(compared))
	}

	ok := reflect.DeepEqual(compare, compared)
	if ok {
		return 0
	}
	return 1
}

func Clone(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	if !types.IsScalar(obj) {
		return deepcopy.Copy(obj)
	}
	return obj
}

func SetValue(data reflect.Value, value interface{}) {
	switch data.Kind() {
	case reflect.Bool:
		data.SetBool(types.GetBool(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		data.SetInt(types.GetInt(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		data.SetUint(types.GetUint(value))
	case reflect.Float32, reflect.Float64:
		data.SetFloat(types.GetFloat(value))
	case reflect.String:
		data.SetString(types.GetString(value))
	case reflect.Map:
		data.Set(reflect.ValueOf(Clone(value)))
	case reflect.Array, reflect.Slice:
		data.Set(reflect.ValueOf(Clone(value)))
	case reflect.Ptr:
		data.Set(reflect.ValueOf(Clone(value)))
	default:
		data.Set(reflect.ValueOf(Clone(value)))
	}
}

// VersionCompare
// compare > compared return 1
// compare > compared return 0
// compare < compared return -1
func VersionCompare(compare, compared string) int {
	compare = strings.TrimPrefix(strings.ToLower(compare), "v")
	compared = strings.TrimPrefix(strings.ToLower(compared), "v")

	if compare == compared {
		return 0
	}

	if compare == "" {
		return -1
	}

	if compared == "" {
		return 1
	}

	compareVersions := strings.Split(compare, ".")
	compareVersionCount := len(compareVersions)

	comparedVersions := strings.Split(compared, ".")
	comparedVersionCount := len(comparedVersions)

	maxVersionCount := compareVersionCount
	if comparedVersionCount > compareVersionCount {
		maxVersionCount = comparedVersionCount
	}

	for i := 0; i < maxVersionCount; i++ {
		var (
			compareVersion  string
			comparedVersion string
		)

		if i >= compareVersionCount {
			compareVersion = "0"
		} else {
			compareVersion = compareVersions[i]
		}

		if i >= comparedVersionCount {
			comparedVersion = "0"
		} else {
			comparedVersion = comparedVersions[i]
		}

		if compareVersion == comparedVersion {
			continue
		}

		compareVersionInt, err := strconv.ParseInt(compareVersion, 10, 64)
		if err != nil {
			return -1
		}

		comparedVersionInt, err := strconv.ParseInt(comparedVersion, 10, 64)
		if err != nil {
			return 1
		}

		if compareVersionInt < comparedVersionInt {
			return -1
		} else if compareVersionInt > comparedVersionInt {
			return 1
		}
	}
	return 0
}
