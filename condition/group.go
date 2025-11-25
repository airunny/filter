package condition

import (
	"context"
	"errors"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/types"
)

type Group struct {
	logic      Logic
	conditions []Condition
}

func NewGroup(logic Logic) *Group {
	return &Group{
		logic:      logic,
		conditions: make([]Condition, 0),
	}
}

func (s *Group) add(condition Condition) {
	s.conditions = append(s.conditions, condition)
}

func (s *Group) IsConditionOk(ctx context.Context, data interface{}, cache *cache.Cache) (bool, error) {
	for _, condition := range s.conditions {
		ok, err := condition.IsConditionOk(ctx, data, cache)
		if err != nil {
			return false, err
		}

		if ok {
			switch s.logic {
			case LogicNot:
				return false, nil
			case LogicOr:
				return true, nil
			default:
				continue
			}
		} else {
			switch s.logic {
			case LogicNot:
				continue
			case LogicOr:
				continue
			default:
				return false, nil
			}
		}
	}
	return true, nil
}

func BuildGroup(ctx context.Context, items []interface{}, logic Logic) (Condition, error) {
	group := NewGroup(logic)
	for _, item := range items {
		if !types.IsArray(item) {
			return nil, errors.New("condition item is not array")
		}

		subCondition, err := BuildCondition(ctx, item.([]interface{}), LogicAnd)
		if err != nil {
			return nil, err
		}
		group.add(subCondition)
	}
	return group, nil
}
