package assignment

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	"github.com/Liyanbing/filter"
	"github.com/Liyanbing/filter/utils"

	filterType "github.com/Liyanbing/filter/type"
)

var factory *Factory

func init() {
	factory = &Factory{
		assignments: map[string]Assignment{
			"=": &Equal{},
		},
	}
}

type Assignment interface {
	Run(ctx context.Context, data interface{}, key string, val interface{})
	PrepareValue(ctx context.Context, value interface{}) (interface{}, error)
}

type OriginValue struct{}

func (s *OriginValue) PrepareValue(ctx context.Context, value interface{}) (interface{}, error) {
	return value, nil
}

type Factory struct {
	assignments map[string]Assignment
}

func (s *Factory) Register(name string, ass Assignment) error {
	if _, ok := s.assignments[name]; !ok {
		return filter.ErrAlreadyExists
	}

	return nil
}

func (s *Factory) Get(name string) Assignment {
	if assignment, ok := s.assignments[name]; ok {
		return assignment
	}

	return nil
}

// -------- =
type Equal struct{ OriginValue }

func (s *Equal) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	if setter, ok := data.(Setter); ok {
		setter.AssignmentSet(key, val)
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

// --------
