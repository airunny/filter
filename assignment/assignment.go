package assignment

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/liyanbing/filter/utils"

	filterType "github.com/liyanbing/filter/filter_type"
)

var innerFactory *factory

func init() {
	innerFactory = &factory{
		assignments: map[string]Assignment{
			"=":      &Set{},
			"+=":     &AddSet{},      // TODO
			"*=":     &MultiplySet{}, // TODO
			"/=":     &DivisionSet{}, // TODO
			"merge":  &Merge{},
			"delete": &Delete{},
		},
	}
}

func Register(name string, ass Assignment) error {
	return innerFactory.Register(name, ass)
}

func Get(name string) Assignment {
	return innerFactory.Get(name)
}

type Assignment interface {
	Run(ctx context.Context, data interface{}, key string, val interface{})
	PrepareValue(ctx context.Context, value interface{}) (interface{}, error)
}

type OriginValue struct{}

func (s *OriginValue) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return value, nil
}

// ----------------
type factory struct {
	assignments map[string]Assignment
}

func (s *factory) Register(name string, ass Assignment) error {
	if _, ok := s.assignments[name]; ok {
		return fmt.Errorf("%v assignment already exists", name)
	}
	s.assignments[name] = ass
	return nil
}

func (s *factory) Get(name string) Assignment {
	if assignment, ok := s.assignments[name]; ok {
		return assignment
	}

	return nil
}

// -------- =
type Set struct{ OriginValue }

func (s *Set) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	if setter, ok := data.(Setter); ok {
		setter.AssignmentSet(key, val)
		return
	}

	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		key = keys[len(keys)-1]

		var ok bool
		data, ok = utils.GetObjectValueByKey(data, strings.Join(keys[:len(keys)-1], "."))
		if !ok {
			return
		}
	}

	if data == nil {
		return
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		dataValue := reflect.ValueOf(data)
		dataValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	case reflect.Slice, reflect.Array:
		dataValue := reflect.ValueOf(data)
		index, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return
		}

		if int(index) >= dataValue.Len() || index < 0 {
			return
		}

		metaData := dataValue.Index(int(index))
		setDataValue(metaData, val)
	case reflect.Ptr:
		dataValue := reflect.ValueOf(data)
		if !dataValue.Elem().CanSet() {
			return
		}

		dataValue = dataValue.Elem()
		f := dataValue.FieldByName(key)
		if !f.IsValid() {
			return
		}

		setDataValue(f, val)
	}
}

func setDataValue(data reflect.Value, value interface{}) {
	switch data.Kind() {
	case reflect.Bool:
		data.SetBool(filterType.GetBool(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		data.SetInt(filterType.GetInt(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		data.SetUint(filterType.GetUint(value))
	case reflect.Float32, reflect.Float64:
		data.SetFloat(filterType.GetFloat(value))
	case reflect.String:
		data.SetString(filterType.GetString(value))
	case reflect.Map:
		data.Set(reflect.ValueOf(filterType.Clone(value)))
	case reflect.Array, reflect.Slice:
		data.Set(reflect.ValueOf(filterType.Clone(value)))
	case reflect.Ptr:
		data.Set(reflect.ValueOf(filterType.Clone(value)))
	default:
		data.Set(reflect.ValueOf(filterType.Clone(value)))
	}
}

// -------- merge
type Merge struct{}

func (s *Merge) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	if merger, ok := data.(Merger); ok {
		merger.AssignmentMerge(key, val)
		return
	}

	originData := data
	var ok bool
	data, ok = utils.GetObjectValueByKey(data, key)
	if !ok {
		return
	}

	if data == nil {
		innerFactory.Get("=").Run(ctx, originData, key, val)
		return
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		dataValue, ok := data.(map[string]interface{})
		valueValue, ok1 := val.(map[string]interface{})
		if ok && ok1 {
			for key, value := range valueValue {
				dataValue[key] = value
			}
		}

	case reflect.Slice:
		if reflect.TypeOf(val).Kind() == reflect.Slice {
			dataValue := reflect.ValueOf(data)
			valueValue := reflect.ValueOf(val)
			for i := 0; i < valueValue.Len(); i++ {
				dataValue = reflect.Append(dataValue, valueValue.Index(i))
			}

			innerFactory.Get("=").Run(ctx, originData, key, dataValue.Interface())
		}
	}
}

func (s *Merge) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, errors.New("assignment[merge] value must be hash or array")
	}

	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Map && t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return nil, errors.New("assignment[merge] value must be hash or array")
	}

	return value, nil
}

// ---------- delete
type Delete struct{}

func (s *Delete) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	if deleter, ok := data.(Deleter); ok {
		deleter.AssignmentDelete(key, val)
		return
	}
	originData := data

	var ok bool
	data, ok = utils.GetObjectValueByKey(data, key)
	if !ok {
		return
	}

	if data == nil {
		return
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		dataValue, ok := data.(map[string]interface{})
		valueValue, ok1 := val.([]string)
		if ok && ok1 {
			for _, key := range valueValue {
				delete(dataValue, key)
			}
		}

	case reflect.Slice:
		delIndexMap := make(map[int]struct{})

		if intArray, ok := val.([]int); ok {
			for _, delIndex := range intArray {
				delIndexMap[delIndex] = struct{}{}
			}
		}

		dataValue := reflect.ValueOf(data)
		newArr := reflect.New(dataValue.Type()).Elem()
		for index := 0; index < dataValue.Len(); index++ {
			if _, ok := delIndexMap[index]; ok {
				continue
			}

			newArr = reflect.Append(newArr, dataValue.Index(index))
		}

		innerFactory.Get("=").Run(ctx, originData, key, newArr.Interface())
	}
}

func (s *Delete) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, errors.New("assignment[delete] value must be hash or array")
	}

	t := reflect.TypeOf(value)
	kind := t.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return nil, errors.New("assignment[delete] value must be int array or string array")
	}

	return value, nil
}

// ------- addSet
type AddSet struct{}

func (s *AddSet) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	return
}

func (s *AddSet) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return nil, nil
}

// ------- multiplySet
type MultiplySet struct{}

func (s *MultiplySet) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	return
}

func (s *MultiplySet) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return nil, nil
}

// --------- divisionSet
type DivisionSet struct{}

func (s *DivisionSet) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	return
}

func (s *DivisionSet) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return nil, nil
}
