package condition

import (
	"context"
	"errors"
	"fmt"

	"github.com/Liyanbing/filter/cache"
	"github.com/Liyanbing/filter/operations"
	"github.com/Liyanbing/filter/variables"

	filterType "github.com/Liyanbing/filter/type"
)

// ------------- base
type Condition interface {
	IsConditionOk(context.Context, interface{}, *cache.Cache) bool
}

type BaseCondition struct {
	variable  variables.Variable
	operation operations.Operation
	value     interface{}
}

func (s *BaseCondition) IsConditionOk(ctx context.Context, data interface{}, cache *cache.Cache) bool {
	return s.operation.Run(ctx, s.variable, s.value, data, cache)
}

// -------------- group
type LOGIC int

const (
	LOGIC_AND LOGIC = iota // 所有 Operation true
	LOGIC_OR               // 只要有一个 Operation true
	LOGIC_NOT              // 所有 Operation false
)

var groupConditionKeys = map[string]LOGIC{
	"and": LOGIC_AND,
	"or":  LOGIC_OR,
	"not": LOGIC_NOT,
}

type GroupCondition struct {
	logic      LOGIC
	conditions []Condition
}

func NewGroupCondition(logic LOGIC) *GroupCondition {
	return &GroupCondition{
		logic:      logic,
		conditions: make([]Condition, 0),
	}
}

func (s *GroupCondition) add(condition Condition) {
	s.conditions = append(s.conditions, condition)
}

func (s *GroupCondition) IsConditionOk(ctx context.Context, data interface{}, cache *cache.Cache) bool {
	result := true

	for _, condition := range s.conditions {
		if ok := condition.IsConditionOk(ctx, data, cache); ok {
			if s.logic == LOGIC_OR {
				result = true
				break
			}

			if s.logic == LOGIC_NOT {
				result = false
				break
			}
		} else {
			if s.logic == LOGIC_AND {
				result = false
				break
			}

			if s.logic == LOGIC_OR {
				result = false
				break
			}
		}
	}

	return result
}

// --------------
func BuildGroupCondition(items []interface{}, logic LOGIC) (Condition, error) {
	group := NewGroupCondition(logic)

	for _, item := range items {
		if !filterType.IsArray(item) {
			return nil, errors.New("condition item is not array")
		}

		subCondition, err := BuildCondition(item.([]interface{}), LOGIC_AND)
		if err != nil {
			return nil, err
		}
		group.add(subCondition)
	}

	return group, nil
}

func BuildCondition(items []interface{}, logic LOGIC) (Condition, error) {
	if len(items) == 0 {
		return nil, errors.New("condition is empty")
	}

	// group
	if filterType.IsArray(items[0]) {
		return BuildGroupCondition(items, logic)
	}

	if len(items) != 3 {
		return nil, errors.New("condition item must contains 3 elements")
	}

	if !filterType.IsString(items[0]) {
		return nil, fmt.Errorf("condition item 1st element[%g] is not string", items[0])
	}

	key := items[0].(string)
	if logic, ok := groupConditionKeys[key]; ok {
		if !filterType.IsArray(items[2]) {
			return nil, fmt.Errorf("group condition [%s] 3rd element is not array", key)
		}

		return BuildCondition(items[2].([]interface{}), logic)
	}

	variable := variables.Factory.Get(key)
	if variable == nil {
		return nil, fmt.Errorf("condition unknow var [%s]", key)
	}

	if !filterType.IsString(items[1]) {
		return nil, fmt.Errorf("condition item 2nd element[%g] is not string", items[1])
	}

	operationName := items[1].(string)
	operation := operations.Factory.Get(operationName)
	if operation == nil {
		return nil, fmt.Errorf("condition with invalid operation[%s]", operationName)
	}

	prepayValue, err := operation.PrepareValue(items[2])
	if err != nil {
		return nil, err
	}

	return &BaseCondition{
		variable:  variable,
		operation: operation,
		value:     prepayValue,
	}, nil
}
