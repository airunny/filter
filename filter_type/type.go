package filter_type

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/liyanbing/filter/utils"
	"github.com/mohae/deepcopy"
)

type FilterType int

const (
	STRING FilterType = iota
	NUMBER
	BOOL
	HASH
	ARRAY
	STRUCT
	NULL
	UNKNOWN
)

var FilterTypeValue = map[FilterType]string{
	STRING:  "STRING",
	NUMBER:  "NUMBER",
	BOOL:    "BOOL",
	HASH:    "HASH",
	STRUCT:  "STRUCT",
	ARRAY:   "ARRAY",
	NULL:    "NULL",
	UNKNOWN: "UNKNOWN",
}

func (m FilterType) String() string {
	if filterType, ok := FilterTypeValue[m]; ok {
		return filterType
	}
	return "UNKNOWN"
}

var numberTypes = map[reflect.Kind]bool{
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
}

func GetFilterType(obj interface{}) FilterType {
	if obj == nil {
		return NULL
	}

	kind := reflect.TypeOf(obj).Kind()
	if _, ok := numberTypes[kind]; ok {
		return NUMBER
	}

	if kind == reflect.Ptr {
		kind = reflect.TypeOf(reflect.ValueOf(obj).Elem().Interface()).Kind()
	}

	switch kind {
	case reflect.Slice, reflect.Array:
		return ARRAY
	case reflect.Map:
		return HASH
	case reflect.String:
		return STRING
	case reflect.Bool:
		return BOOL
	case reflect.Struct:
		return STRUCT
	}

	return UNKNOWN
}

func GetBool(obj interface{}) bool {
	if obj == nil {
		return false
	}

	switch v := obj.(type) {
	case bool:
		return v
	case float32:
		return v > 0
	case float64:
		return v > 0
	case int:
		return v > 0
	case int8:
		return v > 0
	case int16:
		return v > 0
	case int32:
		return v > 0
	case int64:
		return v > 0
	case uint:
		return v > 0
	case uint8:
		return v > 0
	case uint16:
		return v > 0
	case uint32:
		return v > 0
	case uint64:
		return v > 0
	case string:
		return v != ""
	default:
		return false
	}
}

func GetFloat(obj interface{}) float64 {
	if obj == nil {
		return 0.0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1.0
		} else {
			return 0.0
		}
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		if v == "" {
			return 0.0
		}

		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0.0
		}
		return n
	default:
		return 0.0
	}
}

func GetInt(obj interface{}) int64 {
	if obj == nil {
		return 0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case string:
		if v == "" {
			return 0
		}
		if n, err := strconv.ParseInt(v, 10, 64); err != nil {
			return 0
		} else {
			return n
		}
	default:
		return 0
	}
}

func GetUint(obj interface{}) uint64 {
	if obj == nil {
		return 0
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	case int:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case uint64:
		return v
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case string:
		if v == "" {
			return 0
		}
		if n, err := strconv.ParseUint(v, 10, 64); err != nil {
			return 0
		} else {
			return n
		}
	default:
		return 0
	}
}

func GetString(obj interface{}) string {
	if obj == nil {
		return ""
	}

	switch v := obj.(type) {
	case bool:
		if v {
			return "1"
		} else {
			return ""
		}
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case string:
		return v
	default:
		return ""
	}
}

func IsHash(obj interface{}) bool {
	return GetFilterType(obj) == HASH
}

func IsArray(obj interface{}) bool {
	return GetFilterType(obj) == ARRAY
}

func IsString(obj interface{}) bool {
	return GetFilterType(obj) == STRING
}

func IsNumber(obj interface{}) bool {
	return GetFilterType(obj) == NUMBER
}

func IsBool(obj interface{}) bool {
	return GetFilterType(obj) == BOOL
}

func IsStruct(obj interface{}) bool {
	return GetFilterType(obj) == STRUCT
}

func IsScalar(obj interface{}) bool {
	filterType := GetFilterType(obj)

	if filterType == NUMBER || filterType == STRING || filterType == BOOL || filterType == NULL {
		return true
	}

	return false
}

func NumberCompare(a, b interface{}) int {
	fa := GetFloat(a)
	fb := GetFloat(b)

	if utils.FloatEquals(fa, fb) {
		return 0
	}

	if fa-fb > 0 {
		return 1
	} else {
		return -1
	}
}

func ObjectCompare(compare, compared interface{}) int {
	if compare == nil && compared == nil {
		return 0
	}

	compareType := GetFilterType(compare)
	comparedType := GetFilterType(compared)

	if compareType == NUMBER || compareType == BOOL || comparedType == NUMBER || comparedType == BOOL {
		return NumberCompare(compare, compared)
	}

	if compareType == STRING || comparedType == STRING {
		return strings.Compare(GetString(compare), GetString(compared))
	}

	if compareType == STRING || comparedType == STRING {
		return strings.Compare(GetString(compare), GetString(compared))
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

	if !IsScalar(obj) {
		return deepcopy.Copy(obj)
	}

	return obj
}
