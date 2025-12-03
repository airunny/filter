package executor

import (
	"context"
	"errors"
	"fmt"

	"github.com/airunny/filter/assignment"
	_ "github.com/airunny/filter/assignment/delete"
	_ "github.com/airunny/filter/assignment/set"
	"github.com/airunny/filter/types"
)

type Executor interface {
	Execute(ctx context.Context, data interface{}) error
}

type BaseExecutor struct {
	key        string
	assignment assignment.Assignment
	value      interface{}
}

func (s *BaseExecutor) Execute(ctx context.Context, data interface{}) error {
	return s.assignment.Run(ctx, data, s.key, s.value)
}

func BuildExecutor(ctx context.Context, items []interface{}) (Executor, error) {
	if len(items) == 0 {
		return nil, errors.New("executor item must be array")
	}

	if types.IsArray(items[0]) {
		return BuildGroup(ctx, items)
	}

	if len(items) != 3 {
		return nil, errors.New("executor item must contains 3 elements")
	}

	key, ok := items[0].(string)
	if !ok {
		return nil, fmt.Errorf("executor item 1st item  %v is not string", items[0])
	}

	assignmentName, ok := items[1].(string)
	if !ok {
		return nil, fmt.Errorf("executor item 2nd item  %v is not string", items[1])
	}

	assignInstance, ok := assignment.Get(assignmentName)
	if !ok {
		return nil, fmt.Errorf("executor assignment not exists [%s]", assignmentName)
	}

	prepayValue, err := assignInstance.PrepareValue(ctx, items[2])
	if err != nil {
		return nil, fmt.Errorf("executor assignment [%s] preparevalue err:%s", assignmentName, err)
	}

	return &BaseExecutor{
		key:        key,
		assignment: assignInstance,
		value:      prepayValue,
	}, nil
}
