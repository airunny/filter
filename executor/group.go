package executor

import (
	"context"
	"fmt"

	"github.com/liyanbing/filter/types"
)

func NewGroup() *Group {
	return &Group{
		executors: make([]Executor, 0),
	}
}

type Group struct {
	executors []Executor
}

func (s *Group) Execute(ctx context.Context, data interface{}) error {
	for _, executor := range s.executors {
		err := executor.Execute(ctx, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Group) Add(executor Executor) {
	s.executors = append(s.executors, executor)
}

func BuildGroup(ctx context.Context, items []interface{}) (Executor, error) {
	group := NewGroup()
	for _, item := range items {
		if !types.IsArray(item) {
			return nil, fmt.Errorf("executor group item must be array")
		}

		subExecutor, err := BuildExecutor(ctx, item.([]interface{}))
		if err != nil {
			return nil, err
		}
		group.Add(subExecutor)
	}
	return group, nil
}
