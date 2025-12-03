package condition

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/airunny/filter/cache"
	"github.com/airunny/filter/operations"
	"github.com/airunny/filter/types"
	"github.com/airunny/filter/variables"

	// variables
	_ "github.com/airunny/filter/variables/area"
	_ "github.com/airunny/filter/variables/calc"
	_ "github.com/airunny/filter/variables/channel"
	_ "github.com/airunny/filter/variables/ctx"
	_ "github.com/airunny/filter/variables/data"
	_ "github.com/airunny/filter/variables/device"
	_ "github.com/airunny/filter/variables/freq"
	_ "github.com/airunny/filter/variables/ip"
	_ "github.com/airunny/filter/variables/is_login"
	_ "github.com/airunny/filter/variables/platform"
	_ "github.com/airunny/filter/variables/rand"
	_ "github.com/airunny/filter/variables/referer"
	_ "github.com/airunny/filter/variables/success"
	_ "github.com/airunny/filter/variables/time"
	_ "github.com/airunny/filter/variables/ua"
	_ "github.com/airunny/filter/variables/uid"
	_ "github.com/airunny/filter/variables/user_tag"
	_ "github.com/airunny/filter/variables/version"

	// operations
	_ "github.com/airunny/filter/operations/any"
	_ "github.com/airunny/filter/operations/between"
	_ "github.com/airunny/filter/operations/equal"
	_ "github.com/airunny/filter/operations/greater_than"
	_ "github.com/airunny/filter/operations/greater_than_equal"
	_ "github.com/airunny/filter/operations/has"
	_ "github.com/airunny/filter/operations/in"
	_ "github.com/airunny/filter/operations/in_ip_range"
	_ "github.com/airunny/filter/operations/less_than"
	_ "github.com/airunny/filter/operations/less_than_equal"
	_ "github.com/airunny/filter/operations/match"
	_ "github.com/airunny/filter/operations/match_any"
	_ "github.com/airunny/filter/operations/match_none"
	_ "github.com/airunny/filter/operations/not"
	_ "github.com/airunny/filter/operations/not_in"
	_ "github.com/airunny/filter/operations/not_in_ip_range"
	_ "github.com/airunny/filter/operations/not_match"
	_ "github.com/airunny/filter/operations/version_greater_than"
	_ "github.com/airunny/filter/operations/version_greater_than_equal"
	_ "github.com/airunny/filter/operations/version_less_than"
	_ "github.com/airunny/filter/operations/version_less_than_equal"
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
	if logicKey, ok := groupLogicKeys[strings.ToLower(key)]; ok {
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
		return nil, fmt.Errorf("condition operation should be string [%v]", items[1])
	}

	operationName := items[1].(string)
	operation, ok := operations.Get(operationName)
	if !ok {
		return nil, fmt.Errorf("condition not exists operation [%s]", operationName)
	}

	operationValue, err := operation.PrepareValue(items[2])
	if err != nil {
		return nil, err
	}

	return &BaseCondition{
		variable:  variable,
		operation: operation,
		value:     operationValue,
	}, nil
}
