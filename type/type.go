package _type

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Liyanbing/filter/utils"
	"github.com/mohae/deepcopy"
)

type MyType int

const (
	STRING MyType = iota
	NUMBER
	BOOL
	HASH
	ARRAY
	STRUCT
	NULL
	UNKNOWN
)

var MyTypeValue = map[MyType]string{
	STRING:  "STRING",
	NUMBER:  "NUMBER",
	BOOL:    "BOOL",
	HASH:    "HASH",
	STRUCT:  "STRUCT",
	ARRAY:   "ARRAY",
	NULL:    "NULL",
	UNKNOWN: "UNKNOWN",
}

func (m MyType) String() string {
	if myType, ok := MyTypeValue[m]; ok {
		return myType
	}
	return ""
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

func GetMyType(obj interface{}) MyType {
	kind := reflect.TypeOf(obj).Kind()
	if _, ok := numberTypes[kind]; ok {
		return NUMBER
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
	return GetMyType(obj) == HASH
}

func IsArray(obj interface{}) bool {
	return GetMyType(obj) == ARRAY
}

func IsString(obj interface{}) bool {
	return GetMyType(obj) == STRING
}

func IsNumber(obj interface{}) bool {
	return GetMyType(obj) == NUMBER
}

func IsBool(obj interface{}) bool {
	return GetMyType(obj) == BOOL
}

func IsStruct(obj interface{}) bool {
	return GetMyType(obj) == STRUCT
}

func IsScalar(obj interface{}) bool {
	myType := GetMyType(obj)

	if myType == NUMBER || myType == STRING || myType == BOOL || myType == NULL {
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

	compareType := GetMyType(compare)
	comparedType := GetMyType(compared)

	if compareType == NUMBER || compareType == BOOL || comparedType == NUMBER || comparedType == BOOL {
		return NumberCompare(compare, compared)
	}

	if compareType == STRING || comparedType == STRING {
		return strings.Compare(GetString(compare), GetString(compared))
	}

	if compareType == ARRAY && comparedType == ARRAY {
		targetCompare, ok := compare.([]string)
		if !ok {
			return 1
		}
		compareCount := len(targetCompare)

		targetCompared, ok := compared.([]string)
		if !ok {
			return 1
		}
		comparedCount := len(targetCompared)

		if compareCount > comparedCount {
			return 1
		} else if compareCount < comparedCount {
			return -1
		}

		for i := 0; i < compareCount; i++ {
			if ret := strings.Compare(targetCompare[i], targetCompared[i]); ret != 0 {
				return ret
			}
		}

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
