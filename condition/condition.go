package condition

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/operations"
	"github.com/liyanbing/filter/types"
	"github.com/liyanbing/filter/variables"

	// variables
	_ "github.com/liyanbing/filter/variables/area"
	_ "github.com/liyanbing/filter/variables/calc"
	_ "github.com/liyanbing/filter/variables/channel"
	_ "github.com/liyanbing/filter/variables/ctx"
	_ "github.com/liyanbing/filter/variables/data"
	_ "github.com/liyanbing/filter/variables/device"
	_ "github.com/liyanbing/filter/variables/freq"
	_ "github.com/liyanbing/filter/variables/ip"
	_ "github.com/liyanbing/filter/variables/is_login"
	_ "github.com/liyanbing/filter/variables/platform"
	_ "github.com/liyanbing/filter/variables/rand"
	_ "github.com/liyanbing/filter/variables/referer"
	_ "github.com/liyanbing/filter/variables/success"
	_ "github.com/liyanbing/filter/variables/time"
	_ "github.com/liyanbing/filter/variables/ua"
	_ "github.com/liyanbing/filter/variables/uid"
	_ "github.com/liyanbing/filter/variables/user_tag"
	_ "github.com/liyanbing/filter/variables/version"

	// operations
	_ "github.com/liyanbing/filter/operations/any"
	_ "github.com/liyanbing/filter/operations/between"
	_ "github.com/liyanbing/filter/operations/equal"
	_ "github.com/liyanbing/filter/operations/greater_than"
	_ "github.com/liyanbing/filter/operations/greater_than_equal"
	_ "github.com/liyanbing/filter/operations/has"
	_ "github.com/liyanbing/filter/operations/in"
	_ "github.com/liyanbing/filter/operations/in_ip_range"
	_ "github.com/liyanbing/filter/operations/less_than"
	_ "github.com/liyanbing/filter/operations/less_than_equal"
	_ "github.com/liyanbing/filter/operations/match"
	_ "github.com/liyanbing/filter/operations/match_any"
	_ "github.com/liyanbing/filter/operations/match_none"
	_ "github.com/liyanbing/filter/operations/not"
	_ "github.com/liyanbing/filter/operations/not_in"
	_ "github.com/liyanbing/filter/operations/not_in_ip_range"
	_ "github.com/liyanbing/filter/operations/not_match"
	_ "github.com/liyanbing/filter/operations/version_greater_than"
	_ "github.com/liyanbing/filter/operations/version_greater_than_equal"
	_ "github.com/liyanbing/filter/operations/version_less_than"
	_ "github.com/liyanbing/filter/operations/version_less_than_equal"
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
