package operations

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variables"
)

type Operation interface {
	Name() string
	PrepareValue(value interface{}) (interface{}, error)
	Run(ctx context.Context, variable variables.Variable, value interface{}, data interface{}, cache *cache.Cache) (bool, error)
}

type OriginValue struct{}

func (s *OriginValue) PrepareValue(value interface{}) (interface{}, error) {
	return value, nil
}

var defaultFactory = &factory{
	operations: make(map[string]Operation),
}

type factory struct {
	operations map[string]Operation
	sync.Mutex
}

func (s *factory) Register(operation Operation) {
	if operation == nil {
		panic("cannot register a nil Operation")
	}
	if operation.Name() == "" {
		panic("cannot register Operation with empty string result for Name()")
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.operations[operation.Name()]; ok {
		panic(fmt.Sprintf("%v operation already exists", operation.Name()))
	}
	s.operations[operation.Name()] = operation
}

func (s *factory) Get(name string) (Operation, bool) {
	operation, ok := s.operations[name]
	return operation, ok
}

func Register(operation Operation) {
	defaultFactory.Register(operation)
}

func Get(name string) (Operation, bool) {
	return defaultFactory.Get(name)
}

func Print() {
	for name, operation := range defaultFactory.operations {
		fmt.Printf("Operations: \n")
		fmt.Println(name, reflect.TypeOf(operation).Name())
		fmt.Printf("\n\n")
	}
}
