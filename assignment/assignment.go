package assignment

import (
	"context"
	"fmt"
	"sync"
)

type Assignment interface {
	Name() string
	PrepareValue(ctx context.Context, value interface{}) (interface{}, error)
	Run(ctx context.Context, data interface{}, key string, val interface{}) error
}

type OriginValue struct{}

func (s *OriginValue) PrepareValue(_ context.Context, value interface{}) (interface{}, error) {
	return value, nil
}

type Setter interface {
	Set(key string, value interface{}) error
}

type Deleter interface {
	Delete(key string, value interface{}) error
}

// ----------------
var defaultFactory = &factory{
	assignments: make(map[string]Assignment),
}

type factory struct {
	assignments map[string]Assignment
	sync.Mutex
}

func (s *factory) Register(assignment Assignment) {
	if assignment == nil {
		panic("cannot register a nil Assignment")
	}
	if assignment.Name() == "" {
		panic("cannot register Assignment with empty string result for Name()")
	}

	s.Lock()
	defer s.Unlock()
	if _, ok := s.assignments[assignment.Name()]; ok {
		panic(fmt.Sprintf("%v operation already exists", assignment.Name()))
	}
	s.assignments[assignment.Name()] = assignment
}

func (s *factory) Get(name string) (Assignment, bool) {
	operation, ok := s.assignments[name]
	return operation, ok
}

func Register(operation Assignment) {
	defaultFactory.Register(operation)
}

func Get(name string) (Assignment, bool) {
	return defaultFactory.Get(name)
}
