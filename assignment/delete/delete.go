package delete

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/liyanbing/filter/assignment"
	"github.com/liyanbing/filter/utils"
)

const Name = "del"

func init() {
	assignment.Register(&Delete{})
}

type Delete struct {
	assignment.OriginValue
}

func (s *Delete) Name() string { return Name }

func (s *Delete) Run(ctx context.Context, data interface{}, key string, val interface{}) error {
	if deleter, ok := data.(assignment.Deleter); ok {
		return deleter.Delete(key, val)
	}

	var (
		keys        = strings.Split(key, ".")
		originData  = data
		preDataPath string
	)

	if len(keys) > 0 {
		key = keys[len(keys)-1]
		var ok bool
		preDataPath = strings.Join(keys[:len(keys)-1], ".")
		data, ok = utils.GetObjectValueByKey(data, preDataPath)
		if !ok {
			return fmt.Errorf("[%s] assignment %v not exists key %s", Name, reflect.TypeOf(originData).String(), preDataPath)
		}
	}

	if data == nil {
		return fmt.Errorf("[%s] assignment data is nil", Name)
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		dataValue, ok := data.(map[string]interface{})
		if ok {
			delete(dataValue, key)
		}
	//case reflect.Slice:
	//	delIndexMap := make(map[int]struct{})
	//	switch val.(type) {
	//	case []int:
	//		if intArray, ok := val.([]int); ok {
	//			for _, delIndex := range intArray {
	//				delIndexMap[delIndex] = struct{}{}
	//			}
	//		}
	//	case []uint:
	//		if intArray, ok := val.([]uint); ok {
	//			for _, delIndex := range intArray {
	//				delIndexMap[int(delIndex)] = struct{}{}
	//			}
	//		}
	//	case []int64:
	//		if intArray, ok := val.([]int64); ok {
	//			for _, delIndex := range intArray {
	//				delIndexMap[int(delIndex)] = struct{}{}
	//			}
	//		}
	//	case []uint64:
	//		if intArray, ok := val.([]uint64); ok {
	//			for _, delIndex := range intArray {
	//				delIndexMap[int(delIndex)] = struct{}{}
	//			}
	//		}
	//	case string:
	//		index, err := strconv.ParseInt(val.(string), 10, 64)
	//		if err != nil {
	//			return fmt.Errorf("[%s] assignment %v path value %v is a list but value %v can not convert to int",
	//				Name,
	//				reflect.TypeOf(originData).String(),
	//				reflect.TypeOf(data).String(),
	//				key)
	//		}
	//		delIndexMap[int(index)] = struct{}{}
	//	}
	//
	//	dataValue := reflect.ValueOf(data)
	//	if dataValue.Kind() != reflect.Ptr {
	//
	//	}
	//	newArr := reflect.New(dataValue.Type()).Elem()
	//	for index := 0; index < dataValue.Len(); index++ {
	//		if _, ok := delIndexMap[index]; ok {
	//			continue
	//		}
	//		newArr = reflect.Append(newArr, dataValue.Index(index))
	//	}
	//	dataValue.Set(newArr)
	default:
		return fmt.Errorf("[%s] assignent not supported %v", Name, reflect.TypeOf(originData).String())
	}
	return nil
}
