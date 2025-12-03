package set

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/liyanbing/filter/assignment"
	"github.com/liyanbing/filter/utils"
)

const Name = "="

func init() {
	assignment.Register(&Set{})
}

type Set struct{ assignment.OriginValue }

func (s *Set) Name() string { return Name }
func (s *Set) Run(ctx context.Context, data interface{}, key string, val interface{}) error {
	// 这里优先使用接口
	if setter, ok := data.(assignment.Setter); ok {
		return setter.Set(key, val)
	}

	// 如果没有实现接口就使用反射
	var (
		keys       = strings.Split(key, ".")
		originData = data
	)

	if len(keys) > 0 {
		key = keys[len(keys)-1]
		var (
			ok          bool
			preDataPath = strings.Join(keys[:len(keys)-1], ".")
		)

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
		dataValue := reflect.ValueOf(data)
		dataValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	case reflect.Slice, reflect.Array:
		dataValue := reflect.ValueOf(data)
		index, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return fmt.Errorf("[%s] assignment %v path value %v is a list but key [%s] can not convert to int",
				Name,
				reflect.TypeOf(originData).String(),
				reflect.TypeOf(data).String(),
				key)
		}

		if int(index) >= dataValue.Len() || index < 0 {
			return fmt.Errorf("[%s] assignment %v path value %v length is %d but set index is %s",
				Name,
				reflect.TypeOf(originData).String(),
				reflect.TypeOf(data).String(),
				dataValue.Len(),
				key)
		}
		metaData := dataValue.Index(int(index))
		utils.SetValue(metaData, val)
	case reflect.Ptr:
		var (
			dataValue = reflect.ValueOf(data).Elem()
			dataType  = reflect.TypeOf(data).Elem()
		)

		if !dataValue.CanSet() {
			return fmt.Errorf("[%s] %v path value %v can not set",
				Name,
				reflect.TypeOf(originData).String(),
				reflect.TypeOf(data).String())
		}

		f := dataValue.FieldByName(key)
		if !f.IsValid() {
			for i := 0; i < dataType.NumField(); i++ {
				tf := dataType.Field(i)
				if tf.Tag.Get("json") == key {
					utils.SetValue(dataValue.Field(i), val)
					return nil
				}
			}

			return fmt.Errorf("[%s] %v path value not exists key %s",
				Name,
				reflect.TypeOf(originData).String(),
				key)
		}
		utils.SetValue(f, val)
	default:
		return fmt.Errorf("[%s] assignent not supported %v", Name, reflect.TypeOf(originData).String())
	}
	return nil
}
