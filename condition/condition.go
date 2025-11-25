package condition

import (
	"context"
	"errors"
	"fmt"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/types"
	"github.com/liyanbing/filter/variables"
)

type Condition interface {
	IsConditionOk(context.Context, interface{}, *cache.Cache) (bool, error)
}

type Logic int

const (
	LogicAnd Logic = iota
	LogicOr
	LogicNot
)

var groupLogicKeys = map[string]Logic{
	"and": LogicAnd,
	"or":  LogicOr,
	"not": LogicNot,
}

type BaseCondition struct {
	variable  variables.Variable
	operation operations.Operation
	value     interface{}
}

func (s *BaseCondition) IsConditionOk(ctx context.Context, data interface{}, cache *cache.Cache) (bool, error) {
	return s.operation.Run(ctx, s.variable, s.value, data, cache)
}

// --------------

func BuildCondition(ctx context.Context, items []interface{}, logic Logic) (Condition, error) {
	if len(items) == 0 {
		return nil, errors.New("condition is empty")
	}

	// group
	if types.IsArray(items[0]) {
		return BuildGroup(ctx, items, logic)
	}

	if len(items) != 3 {
		return nil, errors.New("condition item must contains three element")
	}

	if !types.IsString(items[0]) {
		return nil, fmt.Errorf("condition item 1st element[%v] is not string", items[0])
	}

	key := items[0].(string)
	if logicKey, ok := groupLogicKeys[key]; ok {
		if !types.IsArray(items[2]) {
			return nil, fmt.Errorf("group condition [%s] 3rd element is not array", key)
		}
		return BuildCondition(ctx, items[2].([]interface{}), logicKey)
	}

	variable, ok := variables.Get(key)
	if !ok {
		return nil, fmt.Errorf("condition not exists variable [%s]", key)
	}

	if !types.IsString(items[1]) {
		return nil, fmt.Errorf("condition item 2nd element[%g] is not string", items[1])
	}

	operationName := items[1].(string)
	operation, ok := operations.Get(operationName)
	if !ok {
		return nil, fmt.Errorf("condition not exists operation [%s]", operationName)
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
