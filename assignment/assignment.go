package assignment

import (
	"context"
	"reflect"
	"strconv"

	"github.com/Liyanbing/filter"
)

var factory *Factory

func init() {
	factory = &Factory{
		assignments: map[string]Assignment{},
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

// --------
type Equal struct{ OriginValue }

func (s *Equal) Run(ctx context.Context, data interface{}, key string, val interface{}) {
	if setter, ok := data.(Setter); ok {
		setter.AssignmentSet(key, val)
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		reflect.ValueOf(data).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	case reflect.Slice, reflect.Array:
		dataValue := reflect.ValueOf(data)
		if index, err := strconv.ParseInt(key, 10, 32); err == nil {
			if int(index) >= dataValue.Len() {
				return
			}
		}

	case reflect.Ptr:
	}
}
