package executor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Liyanbing/filter/assignment"

	filterType "github.com/Liyanbing/filter/type"
)

type Executor interface {
	Execute(ctx context.Context, data interface{})
}

// ----------
type GenExecutor interface {
	GenExecutor(context.Context, interface{}) (Executor, error)
}

var innerFactory *factory

func init() {
	innerFactory = &factory{
		genExecutors: map[string]GenExecutor{
			"__set": &GroupSetter{Name: "__set"},
		},
	}
}

func Register(name string, genExecutor GenExecutor) error {
	return innerFactory.Register(name, genExecutor)
}

func Get(name string) (GenExecutor, error) {
	return innerFactory.Get(name)
}

// -------------
type factory struct {
	genExecutors map[string]GenExecutor
}

func (s *factory) Register(name string, genExecutor GenExecutor) error {
	if _, ok := s.genExecutors[name]; ok {
		return fmt.Errorf("%v alrealdy exists", name)
	}

	s.genExecutors[name] = genExecutor
	return nil
}

func (s *factory) Get(name string) (GenExecutor, error) {
	gen, ok := s.genExecutors[name]
	if !ok {
		return nil, fmt.Errorf("%v gen executor out found", name)
	}

	return gen, nil
}

// ------------
type GroupSetter struct {
	Name string
}

func (s *GroupSetter) GenExecutor(ctx context.Context, value interface{}) (Executor, error) {
	if !filterType.IsArray(value) {
		return nil, fmt.Errorf("%v GenExecutor value must be array", s.Name)
	}

	valueValue := reflect.ValueOf(value)
	items := make([]interface{}, 0, valueValue.Len())
	for i := 0; i < valueValue.Len(); i++ {
		items = append(items, valueValue.Index(i).Interface())
	}
	return buildGroupExecutor(ctx, items)
}

// ----------
type BaseExecutor struct {
	key        string
	assignment assignment.Assignment
	value      interface{}
}

func (s *BaseExecutor) Execute(ctx context.Context, data interface{}) {
	s.assignment.Run(ctx, data, s.key, s.value)
}

// -----------
type GroupExecutor struct {
	executors []Executor
}

func (s *GroupExecutor) Execute(ctx context.Context, data interface{}) {
	for _, executor := range s.executors {
		executor.Execute(ctx, data)
	}
}

func (s *GroupExecutor) add(executor Executor) {
	s.executors = append(s.executors, executor)
}

func NewExecutorGroup() *GroupExecutor {
	return &GroupExecutor{
		executors: make([]Executor, 0),
	}
}

// ----------
func buildGroupExecutor(ctx context.Context, items []interface{}) (Executor, error) {
	group := NewExecutorGroup()

	for _, item := range items {
		if !filterType.IsArray(item) {
			return nil, fmt.Errorf("buildGroupExecutor items err :%v", item)
		}

		subExecutor, err := BuildExecutor(ctx, item.([]interface{}))
		if err != nil {
			return nil, err
		}

		group.add(subExecutor)
	}

	return group, nil
}

func BuildExecutor(ctx context.Context, items []interface{}) (Executor, error) {
	if len(items) == 0 {
		return nil, errors.New(fmt.Sprintf("BuildExecutor items is empty"))
	}

	if filterType.IsArray(items[0]) {
		return buildGroupExecutor(ctx, items)
	}

	if len(items) != 3 {
		return nil, fmt.Errorf("executor item must contains 3 elements")
	}

	key, ok := items[0].(string)
	if !ok {
		return nil, fmt.Errorf("executor 1st item  %v is not string", items[0])
	}

	assignmentName, ok := items[1].(string)
	if !ok {
		return nil, fmt.Errorf("executor 2nd item  %v is not string", items[1])
	}

	if genExecutor, ok := innerFactory.genExecutors[key]; ok {
		return genExecutor.GenExecutor(ctx, items[2])
	}

	assignFunc := assignment.Get(assignmentName)
	if assignFunc == nil {
		return nil, fmt.Errorf("BuildExecutor with invalid assignment [%s]", assignmentName)
	}

	prepayValue, err := assignFunc.PrepareValue(ctx, items[2])
	if err != nil {
		return nil, fmt.Errorf("BuildExecutor assignment [%s] preparevalue err:%s", assignmentName, err)
	}

	return &BaseExecutor{
		key:        key,
		assignment: assignFunc,
		value:      prepayValue,
	}, nil
}
